package queue

import (
	"fmt"
	"opsManager/lib/task"
	"os"
	"strconv"
	"time"

	utils "github.com/b-eee/amagi"
)

type (
	// Scheduler model for multiple tasking
	Scheduler struct {
		IntervalDuration time.Duration
		MainTaskName     string
		TaskHandlers     []Task
		Quit             chan int
	}
)

// Duration the interval duration for the task to execute
func (s *Scheduler) Duration(duration time.Duration) *Scheduler {
	s.IntervalDuration = duration

	return s
}

// Tasks set tasks for the pipeline
func (s *Scheduler) Tasks(task ...Task) *Scheduler {
	s.TaskHandlers = append(s.TaskHandlers, task...)

	return s
}

// LoopDuration loop duration getter from ENV
func LoopDuration(envName string) time.Duration {
	defaultLoopDuration := 60 * task.MinutesMultiplier
	if val, err := strconv.Atoi(os.Getenv(envName)); err == nil {
		return time.Duration(val) * task.MinutesMultiplier
	}

	return defaultLoopDuration
}

// Do do scheduler task
func (s *Scheduler) Do() *Scheduler {
	utils.Info(fmt.Sprintf("task %v started..  interval=%v", s.MainTaskName, s.IntervalDuration))
	ticker := time.NewTicker(s.IntervalDuration)

	go func() {
		for {
			select {
			case <-ticker.C:
				// run all tasks from main task in sync
				// TODO TO RUN ALL TASKS FROM MAIN TASK IN PARALLEL? -JP
				for _, t := range s.TaskHandlers {
					t.Exec()
				}
			case <-s.Quit:
				ticker.Stop()
				return
			}
		}
	}()

	return s
}
