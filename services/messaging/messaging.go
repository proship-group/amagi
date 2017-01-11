package messaging

import (
	"github.com/b-eee/amagi/services/messaging/backend"
)

var (
	// BackendConfig backend config
	BackendConfig struct {
		ConfigEnv backend.MSGBackendConfig
	}
)

// InitMessaging initialize messaging and backend
func InitMessaging() error {
	confBackend := backend.SetAndInitBackend()

	if err := backend.ConnectToMsgBackend(confBackend); err != nil {
		return err
	}

	return nil
}
