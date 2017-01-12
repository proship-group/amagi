package messaging

import (
	"fmt"

	"github.com/b-eee/amagi/services/messaging/backend"

	utils "github.com/b-eee/amagi"
)

var (
	// CurrentMSGBackend current messaging backend config
	CurrentMSGBackend BackendConfig
)

type (
	// BackendConfig backend config
	BackendConfig struct {
		ConfigEnv backend.MSGBackendConfig
	}

	// SubscribeReq subscribe request
	SubscribeReq struct {
		Topic   string
		Channel string
	}

	// PublishReq publish request to backend
	PublishReq struct {
		Topic string
		Body  []byte
	}
)

// InitMessaging initialize messaging with backend and return msg backend config
func InitMessaging() *BackendConfig {
	confBackend := backend.SetAndInitBackend()

	if err := backend.ConnectToMsgBackend(SetCurrentMSGBackend(confBackend).ConfigEnv); err != nil {
		panic(fmt.Errorf("can't set messaging backind %v", err))
		// return &BackendConfig{}
	}

	backend := GetCurrentMSGBackend()
	return &backend
}

// SetCurrentMSGBackend set current messaging backend and return MSGBackendConfig
func SetCurrentMSGBackend(config backend.MSGBackendConfig) BackendConfig {
	newBackend := BackendConfig{
		ConfigEnv: config,
	}

	CurrentMSGBackend = newBackend
	return GetCurrentMSGBackend()
}

// GetCurrentMSGBackend get current messaging backend
func GetCurrentMSGBackend() BackendConfig {
	return CurrentMSGBackend
}

// Subscribe subscribe request to messaging
func (msg *BackendConfig) Subscribe(req SubscribeReq) error {
	var r backend.MSGBackendSubscReq

	if len(req.Topic) != 0 && len(req.Channel) != 0 {
		r.Topic = req.Topic
		r.Channel = req.Channel
	}

	utils.Info(fmt.Sprintf("subscribing to chan=%v topic=%v", r.Topic, r.Channel))
	if err := backend.SubscribeToBackend(GetCurrentMSGBackend().ConfigEnv, r); err != nil {
		return err
	}

	return nil
}

// Publish Publish request to messaging
func (msg *BackendConfig) Publish(req PublishReq) error {
	var r backend.MSGBackendPubReq

	if len(req.Topic) != 0 && len(req.Body) != 0 {
		r.Topic = req.Topic
		r.Body = req.Body
	}

	if err := backend.PublishToBackend(GetCurrentMSGBackend().ConfigEnv, r); err != nil {
		return err
	}

	return nil
}
