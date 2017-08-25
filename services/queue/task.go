package queue

import (
	"fmt"
	"time"

	utils "github.com/b-eee/amagi"
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
	utils.Info(fmt.Sprintf("executing task for [%v]", t.TaskName))

	// execute specified task
	t.Task()
}
