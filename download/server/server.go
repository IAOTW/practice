package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// 文件分片下载并发处理方案
var (
	FileManagerOnce sync.Once
	fileManager     *FileManager
)

func getFileManager() *FileManager {
	FileManagerOnce.Do(func() {
		fileManager = NewFileManager(10 * time.Second)
	})
	return fileManager
}

type FileManager struct {
	files    sync.Map
	lifeTime time.Duration
}

type FileMeta struct {
	RefCount      int
	LastAccessed  time.Time
	ZeroTime      time.Time
	CancelCleanup func()
	mu            sync.Mutex
}

// 创建文件，可以允许多参数，要求第一个参数是必须是文件路径，返回error
type FileGenerator func(filePath string, args ...interface{}) error

func NewFileManager(lifeTime time.Duration) *FileManager {
	return &FileManager{
		lifeTime: lifeTime,
	}
}

func (m *FileManager) Acquire(filePath string, generate FileGenerator, args ...interface{}) (*FileMeta, error) {
	fileMeta, _ := m.files.LoadOrStore(filePath, &FileMeta{})
	fm := fileMeta.(*FileMeta)
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.RefCount == 0 {
		// 传递参数到生成函数
		if err := generate(filePath, args...); err != nil {
			m.files.Delete(filePath)
			return nil, err
		}
	} else if fm.RefCount > 0 && fm.CancelCleanup != nil {
		fm.CancelCleanup()
		fm.CancelCleanup = nil
	}

	fm.RefCount++
	fm.LastAccessed = time.Now()
	return fm, nil
}

func (m *FileManager) Release(filePath string) {
	fileMeta, exists := m.files.Load(filePath)
	if !exists {
		return
	}

	fm := fileMeta.(*FileMeta)
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.RefCount--

	fmt.Printf("文件引用计数 [%s] 减少至: %d\n", filePath, fm.RefCount)

	if fm.RefCount == 0 {
		fm.ZeroTime = time.Now()
		// 异步触发清理检查
		go m.delayedCleanup(filePath)
	}
}

func (m *FileManager) delayedCleanup(filePath string) {
	meta, ok := m.files.Load(filePath)
	if !ok {
		return
	}

	fm := meta.(*FileMeta)
	fm.mu.Lock()

	// 取消之前的清理任务
	if fm.CancelCleanup != nil {
		fm.CancelCleanup()
	}

	ctx, cancel := context.WithCancel(context.Background())
	fm.CancelCleanup = cancel
	fm.mu.Unlock() // 立即释放锁

	go func() {
		select {
		case <-time.After(m.lifeTime):
			fm.mu.Lock()
			defer fm.mu.Unlock()

			// 双重检查
			if fm.RefCount == 0 && time.Since(fm.ZeroTime) >= m.lifeTime {
				if err := os.Remove(filePath); err == nil {
					m.files.Delete(filePath)
					fmt.Printf("成功清理文件: %s\n", filePath)
				}
				fm.CancelCleanup = nil
			}

		case <-ctx.Done():
			fmt.Printf("清理任务取消: %s\n", filePath)
		}
	}()
}

func serveFile(w http.ResponseWriter, r *http.Request, path string) {
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "文件打开失败", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func prepareDownloadFile(filePath string) error {
	fm := getFileManager()
	_, err := fm.Acquire(filePath, generateTempFile)
	return err
}

func main() {
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		filePath := r.URL.Path[len("/download/"):]
		if filePath == "" {
			http.Error(w, "Missing file ID", http.StatusBadRequest)
			return
		}

		if err := prepareDownloadFile(filePath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		serveFile(w, r, filePath)
		getFileManager().Release(filePath)
	})

	fmt.Println("Server running at :8080")
	_ = http.ListenAndServe(":8080", nil)
}

func generateTempFile(path string, args ...interface{}) error {
	tmpPath := path + ".tmp"
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer file.Close()

	size := 100 * 1024 * 1024 // 100MB
	if _, err := io.CopyN(file, rand.Reader, int64(size)); err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}
