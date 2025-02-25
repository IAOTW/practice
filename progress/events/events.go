package events

// TaskEvent 统一的事件结构
type TaskEvent struct {
	TaskID     string
	StageIndex int
	IsComplete bool
	Progress   float64 // 修改为float64，表示百分比进度(0-100)
}

// StageConfig 定义阶段配置
type StageConfig struct {
	Name   string
	Weight float64
	Total  int // 该阶段需要处理的总数，如果是0则表示简单阶段
}
