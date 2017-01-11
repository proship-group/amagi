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

	backend.ConnectToMsgBackend(confBackend)
	return nil
}
