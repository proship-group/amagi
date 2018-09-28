package slack

import (
	"fmt"

	"github.com/nlopes/slack"
)

var (
	// SlackToken slack token for API access
	SlackToken *slack.Client

	// LogChannel channel ID string
	LogChannel string

	// AppName app name
	AppName string

	tokenID string

	// CurrentHost current host configured
	CurrentHost *Host
)

type (
	// Host host that is using this logger
	Host struct {
		Hostname  func() string
		Env       string
		TokenID   string
		ChannelID string

		PublishKey   string
		SubscribeKey string
		SecretKey    string
		MicroAppName string
		Color        string
	}
)

// Init initialize slack API token
func Init(host Host) error {
	SlackToken = slack.New(host.TokenID)
	LogChannel = host.ChannelID
	tokenID = host.TokenID
	CurrentHost = &host

	fmt.Printf("Slack token resp: %v\n", testSlackAuth())

	if err := setHostColor(CurrentHost); err != nil {
		return err
	}

	printHostConnections()
	return nil
}

// Send error log to slack channel
func Send(errmsg interface{}, errStr string) error {
	if ok, err := validateCanSend(); !ok && err != nil {
		return err
	}
	params := slack.PostMessageParameters{}

	chanID, _, err := SlackToken.PostMessage(LogChannel,
		fmt.Sprintf("```host=%v error %v\n ====\n %v ```", CurrentHost.Hostname(), errmsg, errStr),
		params)
	if err != nil {
		errMsg := fmt.Errorf("error sending to slack %v", err)
		fmt.Println(errMsg)
		return errMsg
	}

	fmt.Printf("Message Sent to channel %v", chanID)
	return nil
}

func validateCanSend() (bool, error) {
	if _, err := SlackToken.AuthTest(); err != nil {
		return false, err
	}

	return true, nil
}

// GetCurrentConfiguredHost get current configured host
func GetCurrentConfiguredHost() *Host {
	return CurrentHost
}

// HostName return hostname
func HostName() string {
	if CurrentHost == nil {
		return ""
	}

	return CurrentHost.Hostname()
}

// GetMicroAppName get micro service app name
func GetMicroAppName() string {

	return CurrentHost.MicroAppName
}

func testSlackAuth() bool {
	resp, err := SlackToken.AuthTest()
	if err != nil {
		fmt.Printf("error Slack not auth! %v\n", err)
		return false
	}

	fmt.Printf("slack AuthResponse: %v\n", resp)
	return true
}

func printHostConnections() {
	// Hostname  func() string
	// Env       string
	// TokenID   string
	// ChannelID string

	// PublishKey   string
	// SubscribeKey string
	// SecretKey    string
	// MicroAppName string
	// Color        string
	fmt.Println("printHostConnections")
	str := fmt.Sprintf(`
		HostName: %v
		Env: %v
		ChannelID: %v

		PublishKey: %v
		SubscribeKey: %v
		MicroAppName: %v
		Color: %v

	`, CurrentHost.Hostname(),
		CurrentHost.Env,
		CurrentHost.ChannelID,
		CurrentHost.PublishKey,
		CurrentHost.SubscribeKey,
		CurrentHost.MicroAppName,
		CurrentHost.Color)
	fmt.Println(str)
}
