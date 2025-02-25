package task

import (
	"progress/events"
	"progress/scanner"
	"sync"
)

// Progress 记录进度信息
type Progress struct {
	Name     string
	Weight   float64
	Current  int
	Total    int
	Complete bool
}

// Task 任务信息
type Task struct {
	ID     string
	Type   string
	Stages []Progress
	mu     sync.Mutex
}

// TaskManager 任务管理器
type TaskManager struct {
	tasks    sync.Map
	eventBus chan events.TaskEvent
}

func NewTaskManager() *TaskManager {
	tm := &TaskManager{
		eventBus: make(chan events.TaskEvent, 100),
	}
	go tm.handleEvents()
	return tm
}

func (tm *TaskManager) handleEvents() {
	for event := range tm.eventBus {
		if task, ok := tm.tasks.Load(event.TaskID); ok {
			t := task.(*Task)
			t.mu.Lock()
			if event.StageIndex < len(t.Stages) {
				stage := &t.Stages[event.StageIndex]
				if event.IsComplete {
					stage.Complete = true
					stage.Current = stage.Total
				} else if event.Progress > 0 {
					// 根据百分比计算当前进度值
					stage.Current = int(float64(stage.Total) * event.Progress / 100)
					if stage.Current > stage.Total {
						stage.Current = stage.Total
					}
				}
			}
			t.mu.Unlock()
		}
	}
}

func (tm *TaskManager) StartTask(taskID, taskType string) {
	var stages []events.StageConfig
	switch taskType {
	case "image":
		stages = []events.StageConfig{
			{"初始化", 5, 0},
			{"扫描层", 90, 45},
			{"完成", 5, 0},
		}
	case "host":
		stages = []events.StageConfig{
			{"发现", 20, 0},
			{"检查", 70, 0},
			{"报告", 10, 0},
		}
	}

	task := &Task{
		ID:     taskID,
		Type:   taskType,
		Stages: make([]Progress, len(stages)),
	}

	for i, cfg := range stages {
		task.Stages[i] = Progress{
			Name:   cfg.Name,
			Weight: cfg.Weight,
			Total:  cfg.Total,
		}
	}

	tm.tasks.Store(taskID, task)
	go tm.runTask(task)
}

func (tm *TaskManager) GetProgress(taskID string) (float64, []Progress) {
	if taskData, ok := tm.tasks.Load(taskID); ok {
		task := taskData.(*Task)
		task.mu.Lock()
		defer task.mu.Unlock()

		total := 0.0
		for _, stage := range task.Stages {
			if stage.Total > 0 {
				total += float64(stage.Current) / float64(stage.Total) * stage.Weight
			} else if stage.Complete {
				total += stage.Weight
			}
		}
		return total, task.Stages
	}
	return 0, nil
}

func (tm *TaskManager) runTask(task *Task) {
	var s scanner.Scanner
	switch task.Type {
	case "image":
		s = scanner.NewImageScanner(tm.eventBus)
	case "host":
		s = scanner.NewHostScanner(tm.eventBus)
	default:
		return
	}

	// 执行扫描任务
	s.Scan(task.ID, "")
}
