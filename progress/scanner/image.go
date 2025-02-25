package scanner

import (
	"math/rand"
	"progress/events"
	"time"
)

// ImageScanner 镜像扫描器
type ImageScanner struct {
	eventBus chan<- events.TaskEvent
}

func NewImageScanner(eventBus chan<- events.TaskEvent) *ImageScanner {
	return &ImageScanner{eventBus: eventBus}
}

func (s *ImageScanner) Scan(taskID string, target string) {
	// 初始化阶段
	s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: 0, IsComplete: false}
	time.Sleep(1 * time.Second)
	s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: 0, IsComplete: true}

	// 层扫描阶段
	s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: 1, IsComplete: false}

	// 随机生成层数和每层的文件数
	layers := rand.Intn(5) + 3          // 3-7层
	filesPerLayer := rand.Intn(20) + 10 // 每层10-30个文件
	totalFiles := layers * filesPerLayer
	processedFiles := 0

	for layer := 0; layer < layers; layer++ {
		for file := 0; file < filesPerLayer; file++ {
			// 模拟单个文件处理
			time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
			processedFiles++

			// 计算总体进度百分比
			progress := float64(processedFiles) / float64(totalFiles) * 100

			s.eventBus <- events.TaskEvent{
				TaskID:     taskID,
				StageIndex: 1,
				Progress:   progress,
			}
		}
	}
	s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: 1, IsComplete: true}

	// 结束阶段
	s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: 2, IsComplete: false}
	time.Sleep(1 * time.Second)
	s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: 2, IsComplete: true}
}
