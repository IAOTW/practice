package scanner

// Scanner 定义扫描器接口
type Scanner interface {
	Scan(taskID string, target string)
}
