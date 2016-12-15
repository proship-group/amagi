package monitoring

import (
	"fmt"
	"runtime"
	"time"

	utils "github.com/b-eee/amagi"
)

var (
	reportDelaySec = 15
)

// ReportGoRoutines print go routines count to console
func ReportGoRoutines() {
	utils.Info(fmt.Sprintf("\tGOMAXPROCS/logicalCPU -->> %v", runtime.NumCPU()))
	c := time.Tick(time.Duration(reportDelaySec) * time.Second)
	for now := range c {
		_ = now
		utils.Info(fmt.Sprintf("\tcurrently have goroutines -->> %v", runtime.NumGoroutine()))
	}
}
