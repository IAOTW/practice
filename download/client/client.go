package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	ChunkSize = 1024 * 1024 // 每个分片大小为 1MB
)

// PathExists 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateDirIfNotExists 如果目录不存在则创建
func CreateDirIfNotExists(dir string) error {
	exists, err := PathExists(dir)
	if err != nil {
		return err
	}
	if !exists {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}
	return nil
}

func downloadChunk(client *http.Client, url string, start, end int64, filePath string, fileMode os.FileMode) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发起请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}
	fmt.Println("获取到分片，", fmt.Sprintf("%d - %d", start, end))
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, fileMode)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	_, err = file.Seek(start, os.SEEK_SET)
	if err != nil {
		return fmt.Errorf("设置文件写入位置失败: %w", err)
	}

	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return writer.Flush()
}

// downloadFileInChunks 对应http.ServeContent的分片下载
func downloadFileInChunks(client *http.Client, url string, filePath string, contentLength int64, fileMode os.FileMode) error {
	var errors []error

	for start := int64(0); start < contentLength; start += ChunkSize {
		end := start + ChunkSize - 1
		if end >= contentLength {
			end = contentLength - 1
		}
		err := downloadChunk(client, url, start, end, filePath, fileMode)
		if err != nil {
			return fmt.Errorf("分片下载过程中出现错误: %v", errors)
		}
	}

	return nil
}

func downloadFullFile(client *http.Client, url string, filePath string, fileMode os.FileMode) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发起请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

func GenericHTTPFileDownload(client *http.Client, url string, filePath string, fileMode os.FileMode) error {
	if client == nil {
		client = http.DefaultClient
	}
	// 获取文件所在目录
	dir := filepath.Dir(filePath)
	// 检查并创建目录
	if err := CreateDirIfNotExists(dir); err != nil {
		return err
	}

	// 检查文件是否存在，如果存在则删除l
	if exists, _ := PathExists(filePath); exists {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("删除已存在的文件失败: %w", err)
		}
	}
	resp, err := client.Head(url)
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %w", err)
	}
	defer resp.Body.Close()

	contentLengthStr := resp.Header.Get("Content-Length")
	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil {
		return fmt.Errorf("解析文件大小失败: %w", err)
	}

	if resp.Header.Get("Accept-Ranges") == "bytes" {
		// 分片下载
		return downloadFileInChunks(client, url, filePath, contentLength, fileMode)
	}

	// 不支持分片下载，直接下载
	return downloadFullFile(client, url, filePath, fileMode)
}

func main() {
	err := GenericHTTPFileDownload(nil, "http://localhost:8082/download/file2.iso", "./downloads/file2.iso", os.ModePerm)
	if err != nil {
		fmt.Println("规则库文件下载失败，err:", err)
	}
}
