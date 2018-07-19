package queue

import (
	"fmt"
	"os"
	"strconv"
	"time"

	utils "github.com/b-eee/amagi"
)

const (
	// DequeuerSleepDurationEnv the env var name of the sleep duration
	DequeuerSleepDurationEnv = "QUEUE_DEQUEUER_INTERVAL_MS"

	defaultSleepDuration     = (1 * time.Second)
	defaultMaxConcurrentExec = 1
)

// Dequeue loop process for dequeuing the queue
func Dequeue(itemPtr Executor) {
	sleepDuration := getSleepDuration()
	queueItem := Queue{}
	queueItem.ItemData = itemPtr

	utils.Info(fmt.Sprintf("Dequeuer started with %v sleeping time...", sleepDuration))

	for {
		// TODO: add concurrency settings? like how many max concurrent execution at the same time
		func() {
			if err := queueItem.Dequeue(); err != nil {
				time.Sleep(sleepDuration)
				return
			}
			defer queueItem.CleanUp()

			itemString := fmt.Sprintf("queue `%v` with Identity `%v`",
				queueItem.ID.Hex(),
				queueItem.ItemData.Identity(),
			)
			utils.Info(fmt.Sprintf("[Amagi-Queue] Starting process for %s", itemString))
			procStart := time.Now()
			if err := queueItem.ItemData.Execute(); err != nil {
				utils.Error(fmt.Sprintf("[Amagi-Queue] error queueItem.Execute for %s: %v", itemString, err))
				defer queueItem.Fail()
				return
			}
			queueItem.Success()
			utils.Info(fmt.Sprintf("[Amagi-Queue] Queued %s is done, took: %v",
				itemString,
				time.Since(procStart),
			))
		}()
	}
}

func getSleepDuration() time.Duration {
	if durationEnv := os.Getenv(DequeuerSleepDurationEnv); durationEnv != "" {
		duration, err := strconv.Atoi(durationEnv)
		if err != nil {
			utils.Error(fmt.Sprintf("[Amagi-Queue] Invalid dequeuer sleep duration value: %v", err))
			utils.Warn(fmt.Sprintf("[Amagi-Queue] Using default sleep duration: %v", defaultSleepDuration))
			return defaultSleepDuration
		}
		return (time.Duration(duration) * time.Millisecond)
	}
	return defaultSleepDuration
}
