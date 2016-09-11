package pubnub

import (
	"fmt"
	"github.com/pubnub/go/messaging"
	"sync"

	// utils "github.com/b-eee/amagi"
	"github.com/b-eee/amagi/api/slack"
)

var (
	// PubNubConn pubnub connection
	PubNubConn *messaging.Pubnub
)

type (
	// Credentials pub nub credentials key struct
	Credentials struct {
		PublishKey   string
		SubscribeKey string
		SecretKey    string
	}
)

// SetPubNubConnection set pubnub connection
func SetPubNubConnection() {
	pubNub, err := CreateCredentials()
	if err != nil {
		PubNubConn = nil
	}

	PubNubConn = pubNub
}

// GetPubNubCredentials get pubnub credentials
func GetPubNubCredentials() (*slack.Host, error) {

	return slack.GetCurrentConfiguredHost(), nil
}

// CreateCredentials create pubnub credentials
func CreateCredentials() (*messaging.Pubnub, error) {
	credentials, err := GetPubNubCredentials()
	if err != nil {
		fmt.Printf("error GetPubNubCredentials %v", err)
		return &messaging.Pubnub{}, err
	}

	pubnub := messaging.NewPubnub(credentials.PublishKey, credentials.SubscribeKey, credentials.SecretKey, "", true, "")
	return pubnub, nil
}

// Publish publish to channel
func Publish(channel, message string, wg *sync.WaitGroup) error {
	publishSuccessChannel := make(chan []byte)
	publishErrorChannel := make(chan []byte)

	go PubNubConn.Publish(channel, message,
		publishSuccessChannel, publishErrorChannel)

	go handleResult(publishSuccessChannel, publishErrorChannel, messaging.GetNonSubscribeTimeout(), "publishing", wg)

	defer func() {
		PubNubConn.CloseExistingConnection()
	}()
	return nil
}
