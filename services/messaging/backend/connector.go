package backend

import (
	"fmt"
	"os"

	"github.com/b-eee/amagi/services/configctl"

	utils "github.com/b-eee/amagi"
)

var (
	// NSQHOSTBackendENV nsq host backend URL
	NSQHOSTBackendENV = "NSQ_HOST"

	// MessagingBackendENV  messaging backend select
	MessagingBackendENV = "MESSAGING_BACKEND"

	// CurrentMSGBackendConfig current messaging backend config
	CurrentMSGBackendConfig MSGBackendConfig
)

type (
	// MSGBackendConfig backend selector for selected backend
	MSGBackendConfig struct {
		Env     configctl.Environment
		Backend string
	}

	// MSGBackendSubscReq messaging backend subscribe request
	MSGBackendSubscReq struct {
		Topic   string
		Channel string
	}

	// MSGBackendPubReq messaging backend subscribe request
	MSGBackendPubReq struct {
		Topic string
		Body  []byte
	}

	connectionObjLauncher map[string]func(MSGBackendConfig) error
)

// SetAndInitBackend set and initialize backend
func SetAndInitBackend() MSGBackendConfig {
	msgConfig := MSGBackendConfig{}

	switch os.Getenv(MessagingBackendENV) {
	case "nsq":
		msgConfig.Backend = "nsq"
		msgConfig.Env = configctl.GetDBCfgStngWEnvName("nsq", os.Getenv("ENV"))
	case "nats":
		msgConfig.Backend = "nats"
		msgConfig.Env = configctl.GetDBCfgStngWEnvName("nats", os.Getenv("ENV"))
	}

	return SetMSGBackendConfig(msgConfig)
}

// SetMSGBackendConfig set current MSGBackendConfig and return CurrentMSGBackendConfig
func SetMSGBackendConfig(backendConf MSGBackendConfig) MSGBackendConfig {
	CurrentMSGBackendConfig = backendConf

	return CurrentMSGBackendConfig
}

// GetMSGBackendConfig get current MSGBackendConfig
func GetMSGBackendConfig() MSGBackendConfig {
	return CurrentMSGBackendConfig
}

// ConnectToMsgBackend connect to msg backend by settings config
func ConnectToMsgBackend(confg MSGBackendConfig) error {
	if (MSGBackendConfig{}) == confg {
		return fmt.Errorf("MSGBackendConfig not set")
	}

	connectionObj := connectionObjLauncher{
		"nsq":  StartNSQ,
		"nats": StartNATS,
	}

	utils.Info(fmt.Sprintf(`connecting to
							host=%v
							backend=%v`,
		confg.Env.Host, confg.Backend))

	return connectionObj[confg.Backend](confg)
}

// SubscribeToBackend subscribe to messaging backend
func SubscribeToBackend(confg MSGBackendConfig, req MSGBackendSubscReq) error {
	switch confg.Backend {
	case "nsq":
		nsqReq := NSQConsumerReq{
			Topic:   req.Topic,
			Channel: req.Channel,
		}

		if err := NSQCreateConsumer(confg, nsqReq); err != nil {
			return err
		}
	case "nats":
	}

	return nil
}

// PublishToBackend publish to messaging backend
func PublishToBackend(confg MSGBackendConfig, req MSGBackendPubReq) error {
	nsqReq := NSQPubReq{
		Topic: req.Topic,
		Body:  req.Body,
	}

	switch confg.Backend {
	case "nsq":
		if err := NSQPublish(nsqReq); err != nil {
			return err
		}
	case "nats":
		if err := NATSPublish(nsqReq); err != nil {
			return err
		}
	}

	return nil
}
