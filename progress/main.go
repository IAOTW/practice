package main

import (
	"fmt"
	"progress/task"
	"time"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J") // 清除屏幕并重置光标位置
}

func printTaskProgress(tm *task.TaskManager, taskID string) {
	progress, stages := tm.GetProgress(taskID)
	fmt.Printf("[%s] %.1f%% Complete\n", taskID, progress)
	for _, stage := range stages {
		fmt.Printf("  - %s: ", stage.Name)
		if stage.Total > 0 {
			fmt.Printf("%d/%d (%.1f%%)\n", stage.Current, stage.Total,
				float64(stage.Current)/float64(stage.Total)*100)
		} else {
			fmt.Printf("%s\n", map[bool]string{true: "完成", false: "进行中"}[stage.Complete])
		}
	}
}

func main() {
	tm := task.NewTaskManager()

	tm.StartTask("img-1", "image")
	tm.StartTask("host-1", "host")

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		clearScreen()
		fmt.Println("=== Current Progress ===")
		printTaskProgress(tm, "img-1")
		printTaskProgress(tm, "host-1")
	}
}
