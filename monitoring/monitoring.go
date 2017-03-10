package monitoring

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	utils "github.com/b-eee/amagi"
)

var (
	reportDelaySec = 15
)

type (
	// MonitorTasks monitor task interface
	MonitorTasks struct {
		TaskName string
		Task     func(*sync.WaitGroup) error
	}
)

// ReportGoRoutines print go routines count to console
// TODO Move this routine to seperate file -JP
func ReportGoRoutines() {
	utils.Info(fmt.Sprintf("\tGOMAXPROCS/logicalCPU -->> %v", runtime.NumCPU()))
	c := time.Tick(time.Duration(reportDelaySec) * time.Second)
	for now := range c {
		_ = now
		utils.Info(fmt.Sprintf("\tcurrently have goroutines -->> %v", runtime.NumGoroutine()))
	}
}

// InitAppMonit initialize and configure app monitoring services
func InitAppMonit() error {
	// add tasks as list
	tasks := []MonitorTasks{
		MonitorTasks{"sentry", SentryService},
	}

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)

		go func(t MonitorTasks) {
			utils.Info(fmt.Sprintf("starting monitoring task -->> %v", t.TaskName))
			if err := t.Task(&wg); err != nil {
				return
			}

		}(task)
	}
	wg.Wait()

	return nil
}
