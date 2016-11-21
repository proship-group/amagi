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
		fmt.Printf("can't set pubnub %v\n", err)
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
func Publish(message string) error {
	if cantPublish() {
		return fmt.Errorf("can't publish")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	publishSuccessChannel := make(chan []byte)
	publishErrorChannel := make(chan []byte)

	channel, message := buildMsgAndChan(message)

	go PubNubConn.Publish(channel, message,
		publishSuccessChannel, publishErrorChannel)

	// DISABLED FOR THE MEANTIME --JP
	// go handleResult(publishSuccessChannel, publishErrorChannel, messaging.GetNonSubscribeTimeout(), "publishing -->>", channel, &wg)

	defer func() {
		PubNubConn.CloseExistingConnection()
	}()

	wg.Wait()
	return nil
}

func buildMsgAndChan(msg string) (string, string) {
	channel := ChanName([]string{"log", "stream"}...)
	message := formatHostName(msg, slack.GetMicroAppName(), slack.GetCurrentConfiguredHost())

	// pubnub.Publish(channel, message, &wg)
	return channel, message
}

func cantPublish() bool {
	return slack.CurrentHost.PublishKey == "" &&
		slack.CurrentHost.SubscribeKey == "" &&
		slack.CurrentHost.SecretKey == ""
}
