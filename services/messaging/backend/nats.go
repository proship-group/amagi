package backend

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats"

	utils "github.com/b-eee/amagi"
)

var (
	// NC main nats connection
	NC *nats.Conn
)

// StartNATS start nats connection and settings
func StartNATS(conf MSGBackendConfig) error {
	fmt.Println(conf)

	hosts := natsHosts(conf)
	nc, err := nats.Connect(hosts)
	if err != nil {
		utils.Error(fmt.Sprintf("error StartNATS %v", err))
		return err
	}

	setNATSConn(nc)

	return nil
}

func setNATSConn(nc *nats.Conn) {
	NC = nc
}

func getNATSConn() *nats.Conn {
	return NC
}

func natsHosts(conf MSGBackendConfig) string {
	var splitHost []string
	for _, host := range strings.Split(conf.Env.Host, ",") {
		protocol := "nats://"
		if conf.Env.Username != "" && conf.Env.Password != "" {
			protocol = fmt.Sprintf("%v%v:%v", protocol, conf.Env.Username, conf.Env.Password)
		}

		splitHost = append(splitHost, fmt.Sprintf("%v@%v", protocol, host))
	}

	return strings.Join(splitHost, ", ")

}

// NATSPublish nats publish interface
func NATSPublish(req NSQPubReq) error {
	s := time.Now()
	if err := getNATSConn().Publish(req.Topic, req.Body); err != nil {
		utils.Error(fmt.Sprintf("error NATSPublish %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("NATSPublish took: %v chan=%v", time.Since(s), req.Topic))
	return nil
}
