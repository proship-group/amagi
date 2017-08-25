package queue

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestDoSchedules(t *testing.T) {
	q := make(chan int)
	sampleTask := Task{
		TaskName: "test",
		Task:     testingDo,
	}

	os.Setenv("TestHour", "20")
	os.Setenv("TestMin", "15")

	Exec().SetHourMinute(os.Getenv("TestHour"), os.Getenv("TestMin")).SchedulerDuration(TaskTimeGen).Tasks(
		sampleTask,
	).DoSchedules()

	go func() {
		time.Sleep(5 * time.Minute)
		q <- 1
	}()

	<-q
}

func testingDo() {
	fmt.Println("testingDo!")
}

func calcTime(sc *Scheduler) time.Duration {
	s := time.Now()
	hour, _ := strconv.Atoi("19")
	min, _ := strconv.Atoi("46")
	target := time.Date(s.Year(), s.Month(), s.Day(), hour, min, 0, 0, time.FixedZone("Asia/Tokyo", 9*60*60))

	if sc.LastExecution != (time.Time{}) {
		s = sc.LastExecution.Add(Hourly)
		hour = s.Hour()
		min = s.Minute()
		target = s
	}

	sc.LastExecution = s
	return time.Until(target)
}
