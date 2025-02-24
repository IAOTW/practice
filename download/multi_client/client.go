package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	baseURL     = "http://localhost:8080/download/"
	outputDir   = "./downloads"
	chunkSize   = 5 * 1024 * 1024 // 5MB
	concurrency = 4               // 每个用户的并发线程数
	maxRetries  = 3               // 分片下载重试次数
)

// 模拟独立用户
type User struct {
	ID        string
	UserAgent string
}

// 定义要下载的两个文件
var (
	fileIDs = []string{"file1.zip", "file2.iso"} // 两个不同文件
)

func main() {
	os.MkdirAll(outputDir, 0755)

	// 模拟 6 个独立用户：3个下载file1，3个下载file2
	users := []*User{
		// 下载 file1.zip 的用户组
		{ID: "user1-file1", UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"},
		{ID: "user2-file1", UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"},
		{ID: "user3-file1", UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)"},

		// 下载 file2.iso 的用户组
		{ID: "user1-file2", UserAgent: "Mozilla/5.0 (X11; Linux x86_64)"},
		{ID: "user2-file2", UserAgent: "Mozilla/5.0 (Android 12; Mobile; rv:68.0)"},
		{ID: "user3-file2", UserAgent: "Mozilla/5.0 (iPad; CPU OS 15_0 like Mac OS X)"},
	}

	var wg sync.WaitGroup
	wg.Add(len(users))

	for _, user := range users {
		go func(u *User) {
			defer wg.Done()

			// 修复后的文件分配逻辑
			var fileID string
			switch {
			case strings.HasSuffix(u.ID, "-file1"):
				fileID = fileIDs[0] // file1.zip
			case strings.HasSuffix(u.ID, "-file2"):
				fileID = fileIDs[1] // file2.iso
			default:
				fmt.Printf("[%s] 无效用户配置\n", u.ID)
				return
			}

			downloadAsUser(u, fileID)
		}(user)
	}

	wg.Wait()
	fmt.Println("所有用户下载完成!")
}

// ----------------- 核心下载逻辑 -----------------
func downloadAsUser(user *User, fileID string) {
	// 生成用户专属文件名（user1-file1_file1.zip）
	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s", user.ID, fileID))

	// 获取文件大小
	size, err := getFileSizeWithUser(fileID, user)
	if err != nil {
		fmt.Printf("[%s] 获取文件大小失败: %v\n", user.ID, err)
		return
	}

	// 创建用户专属文件
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("[%s] 创建文件失败: %v\n", user.ID, err)
		return
	}
	defer file.Close()

	// 预分配空间
	if err := file.Truncate(size); err != nil {
		fmt.Printf("[%s] 预分配空间失败: %v\n", user.ID, err)
		return
	}

	// 分片下载
	chunks := calculateChunks(size)
	var dlWg sync.WaitGroup
	chunkChan := make(chan chunk, len(chunks))

	// 启动 Worker
	for i := 0; i < concurrency; i++ {
		dlWg.Add(1)
		go userDownloadWorker(user, fileID, chunkChan, file, &dlWg)
	}

	// 分发任务
	for _, c := range chunks {
		chunkChan <- c
	}
	close(chunkChan)
	dlWg.Wait()

	fmt.Printf("[%s] 下载完成: %s\n", user.ID, outputPath)
}

// ----------------- 工具类型和函数 -----------------
type chunk struct {
	start int64
	end   int64
}

func calculateChunks(totalSize int64) []chunk {
	var chunks []chunk
	for start := int64(0); start < totalSize; start += chunkSize {
		end := start + chunkSize - 1
		if end >= totalSize {
			end = totalSize - 1
		}
		chunks = append(chunks, chunk{start, end})
	}
	return chunks
}

type offsetWriter struct {
	file   *os.File
	offset int64
}

func (w *offsetWriter) Write(p []byte) (n int, err error) {
	n, err = w.file.WriteAt(p, w.offset)
	w.offset += int64(n)
	return
}

// ----------------- 工具函数 (带用户模拟) -----------------
func getFileSizeWithUser(fileID string, user *User) (int64, error) {
	req, _ := http.NewRequest("HEAD", baseURL+fileID, nil)
	setUserHeaders(req, user) // 设置用户特征

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("状态码 %d", resp.StatusCode)
	}

	return strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
}

func userDownloadWorker(user *User, fileID string, chunkChan <-chan chunk, file *os.File, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{Timeout: 30 * time.Second}

	for c := range chunkChan {
		for retry := 0; retry < maxRetries; retry++ {
			err := downloadChunkWithUser(client, user, fileID, file, c.start, c.end)
			if err == nil {
				break
			}
			fmt.Printf("[%s] 分片 %d-%d 下载失败 (重试 %d): %v\n",
				user.ID, c.start, c.end, retry+1, err)
		}
	}
}

func downloadChunkWithUser(client *http.Client, user *User, fileID string, file *os.File, start, end int64) error {
	req, _ := http.NewRequest("GET", baseURL+fileID, nil)
	setUserHeaders(req, user)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("异常状态码: %d", resp.StatusCode)
	}

	writer := &offsetWriter{file: file, offset: start}
	_, err = io.CopyBuffer(writer, resp.Body, make([]byte, 512*1024))
	return err
}

// 设置用户特征
func setUserHeaders(req *http.Request, user *User) {
	req.Header.Set("User-Agent", user.UserAgent)
	req.Header.Set("X-User-ID", user.ID)
	// 可添加更多模拟头，如 Accept-Language 等
}
