package backend

import (
	"fmt"
	"sync"
	"time"

	"github.com/bitly/go-nsq"

	utils "github.com/b-eee/amagi"
)

// StartNSQ start nsq connection
func StartNSQ(conf MSGBackendConfig) error {
	s := time.Now()
	config := nsq.NewConfig()
	w, err := nsq.NewProducer(conf.Env.Host, config)
	if err != nil {
		utils.Error(fmt.Sprintf("error StartNSQ connection %v", err))
		return err
	}
	defer w.Stop()

	// TEST CODE for Connection
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			e := time.Now()
			chanName := fmt.Sprintf("testing_%v", x)
			msg := []byte(fmt.Sprintf("message_%v", x))

			if err := w.Publish(chanName, msg); err != nil {
				utils.Error(fmt.Sprintf("error StartNSQ Publish %v", err))
				return
			}
			utils.Info(fmt.Sprintf("test publish took: %v chan=%v", time.Since(e), chanName))
		}(i)

	}
	wg.Wait()

	utils.Info(fmt.Sprintf("StartNSQ took: %v", time.Since(s)))
	return nil

}
