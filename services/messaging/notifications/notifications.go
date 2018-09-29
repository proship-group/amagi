package notifications

import (
	"fmt"
	"time"

	utils "github.com/b-eee/amagi"
	"github.com/b-eee/amagi/services/externalSvc"
)

var (
	// NotificatorURL notificator service URL
	NotificatorURL string
)

// Init initialize
func Init() {
}

// BuildNotificatorURL build notificator URL
func BuildNotificatorURL() string {
	NotificatorURL = fmt.Sprintf("%s:%s", externalSvc.EnvNotificatorHost, externalSvc.EnvNotificatorPort)
	return NotificatorURL
}

// Publish publish message to pusher
func Publish(message interface{}, channel, event string) error {
	s := time.Now()
	url := "/notification/api/push"

	data := map[string]interface{}{
		"message": message,
		"channel": channel,
		"event":   event,
	}

	if _, err := externalSvc.GenericHTTPRequester("POST", "http", BuildNotificatorURL(), url, data); err != nil {
		utils.Error(fmt.Sprintf("error Publish %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("Publish took: %v channel=%v message:%v", time.Since(s), channel, message))
	return nil
}
