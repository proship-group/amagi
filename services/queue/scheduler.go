package queue

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/b-eee/amagi/helpers"

	utils "github.com/b-eee/amagi"
)

var (
	// Daily schedule incremen
	Daily = 24 * time.Hour

	// Minute minute schedule increment
	Minute = 60 * time.Second

	// Hourly hourly schedule increment
	Hourly = 60 * time.Minute
)

type (
	// Scheduler model for multiple tasking
	Scheduler struct {
		IntervalDuration  time.Duration
		ScheduledDuration func(*Scheduler) time.Duration
		LastExecution     time.Time
		MainTaskName      string
		TaskHandlers      []Task

		SchedulerTimeHour   string
		SchedulerTimeMinute string
		SchedulerIncrement  time.Duration

		Quit chan int
	}
)

// Duration the interval duration for the task to execute
func (s *Scheduler) Duration(duration time.Duration) *Scheduler {
	s.IntervalDuration = duration

	return s
}

// SchedulerDuration schedule a duration from re-calc time
func (s *Scheduler) SchedulerDuration(schedule func(*Scheduler) time.Duration) *Scheduler {
	s.ScheduledDuration = schedule

	return s
}

// SetSchedulerIncrement set scheduler increment
func (s *Scheduler) SetSchedulerIncrement(inc time.Duration) *Scheduler {
	s.SchedulerIncrement = inc
	return s
}

// SetHourMinute set hour and minute for main task scheduler
func (s *Scheduler) SetHourMinute(hour, minute string) *Scheduler {
	utils.Info(fmt.Sprintf("Task set for hour/minute %v/%v", hour, minute))

	s.SchedulerTimeHour = hour
	s.SchedulerTimeMinute = minute

	return s
}

// Tasks set tasks for the pipeline
func (s *Scheduler) Tasks(task ...Task) *Scheduler {
	s.TaskHandlers = append(s.TaskHandlers, task...)

	return s
}

// LoopDuration loop duration getter from ENV
func LoopDuration(envName string) time.Duration {
	defaultLoopDuration := 60 * SecondsMultiplier
	if val, err := strconv.Atoi(os.Getenv(envName)); err == nil {
		return time.Duration(val) * MinutesMultiplier
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

// DoSchedules do task from a specified main schedule
func (s *Scheduler) DoSchedules() *Scheduler {
	for _, t := range s.TaskHandlers {
		go func(task Task) {
			for {

				sleepTime := s.ScheduledDuration(s)
				// ss := time.Until(sleepTime)
				if s.LastExecution != (time.Time{}) {
					fmt.Printf("next execution for task is %v or %v========\n", helpers.TimeToStr(s.LastExecution), sleepTime)
				}
				time.Sleep(sleepTime)
				task.Exec()
			}
		}(t)
	}

	return s
}

// TaskTimeGen generate increment timer
func TaskTimeGen(sc *Scheduler) time.Duration {
	s := time.Now()
	hour, _ := strconv.Atoi(sc.SchedulerTimeHour)
	min, _ := strconv.Atoi(sc.SchedulerTimeMinute)
	target := time.Date(s.Year(), s.Month(), s.Day(), hour, min, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))

	// TODO validate if time exceeded -JP
	// if time.Now() > target. {

	// }

	if sc.LastExecution != (time.Time{}) {
		s = sc.LastExecution.Add(sc.SchedulerIncrement)
		hour = s.Hour()
		min = s.Minute()
		target = s
	}

	sc.LastExecution = s
	return time.Until(target)
}
