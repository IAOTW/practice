package scanner

import (
	"progress/events"
	"time"
)

// HostScanner 主机扫描器
type HostScanner struct {
	eventBus chan<- events.TaskEvent
}

func NewHostScanner(eventBus chan<- events.TaskEvent) *HostScanner {
	return &HostScanner{eventBus: eventBus}
}

func (s *HostScanner) Scan(taskID string, target string) {
	stages := []int{2, 5, 1} // 每个阶段的执行时间
	for idx, duration := range stages {
		s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: idx, IsComplete: false}
		time.Sleep(time.Duration(duration) * time.Second)
		s.eventBus <- events.TaskEvent{TaskID: taskID, StageIndex: idx, IsComplete: true}
	}
}
