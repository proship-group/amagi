package backend

import (
	"fmt"
	"sync"
	"time"

	"github.com/bitly/go-nsq"

	utils "github.com/b-eee/amagi"
)

var (
	// NSQProducer current nsq producer
	NSQProducer *nsq.Producer
)

type (
	// NSQPubReq nsq publish request
	NSQPubReq struct {
		Topic string
		Body  []byte
	}
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
	NSQProducer = w
	// defer w.Stop()
	TestConn()

	utils.Info(fmt.Sprintf("StartNSQ took: %v", time.Since(s)))
	return nil

}

// TestConn test nsq connection and publish
func TestConn() {
	// TEST CODE for Connection
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			e := time.Now()
			chanName := fmt.Sprintf("testing_%v", x)
			msg := []byte(fmt.Sprintf("message_%v", x))

			req := NSQPubReq{
				Topic: chanName,
				Body:  msg,
			}
			if err := NSQPublish(req); err != nil {
				return
			}

			utils.Info(fmt.Sprintf("test publish took: %v chan=%v", time.Since(e), chanName))
		}(i)

	}
	wg.Wait()

}

// NSQPublish nsq publish from nsq producer
func NSQPublish(req NSQPubReq) error {
	if err := NSQProducer.Publish(req.Topic, req.Body); err != nil {
		utils.Error(fmt.Sprintf("error NSQPublish Publish %v", err))
		return err
	}

	return nil
}
