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
)

type (
	// MSGBackendConfig backend selector for selected backend
	MSGBackendConfig struct {
		Env     configctl.Environment
		Backend string
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
	}

	return msgConfig
}

// ConnectToMsgBackend connect to msg backend by settings config
func ConnectToMsgBackend(confg MSGBackendConfig) error {
	if (MSGBackendConfig{}) == confg {
		return fmt.Errorf("MSGBackendConfig not set")
	}

	connectionObj := connectionObjLauncher{
		"nsq": StartNSQ,
	}

	utils.Info(fmt.Sprintf(`connecting to
							host=%v
							backend=%v`,
		confg.Env.Host, confg.Backend))

	return connectionObj[confg.Backend](confg)
}
