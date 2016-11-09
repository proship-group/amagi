package monitoring

import (
	utils "amagi"
	"fmt"
	"runtime"
	"time"
)

var (
	reportDelaySec = 15
)

// reportGoRoutines print go routines count to console
func ReportGoRoutines() {
	c := time.Tick(time.Duration(reportDelaySec) * time.Second)
	for now := range c {
		_ = now
		utils.Info(fmt.Sprintf("currently have goroutines -->> %v", runtime.NumGoroutine()))
		utils.Info(fmt.Sprintf("GOMAXPROCS/logicalCPU -->> %v", runtime.NumCPU()))
	}
}
