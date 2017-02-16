package queue

import (
	"fmt"
	"time"
)

var (
	// SecondsMultiplier seconds multiplier for duration
	SecondsMultiplier = 1 * time.Second

	// MinutesMultiplier minute multiplier for duration
	MinutesMultiplier = 1 * time.Minute
)

type (
	// Task task interface
	Task struct {
		TaskName string
		Task     func()

		Quit chan int
	}
)

// Exec task execution
func (t *Task) Exec() {
	fmt.Printf("executing task for %v\n", t.TaskName)

	// execute specified task
	t.Task()
}
