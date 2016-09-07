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
)

type (
	// Host host that is using this logger
	Host struct {
		Hostname  func() string
		Env       string
		TokenID   string
		ChannelID string
	}
)

// Init initialize slack API token
func Init(host Host) error {
	SlackToken = slack.New(host.TokenID)
	LogChannel = host.ChannelID
	tokenID = host.TokenID

	return nil
}

// Send error log to slack channel
func Send(errmsg interface{}, errStr string) error {
	if ok, err := validateCanSend(); !ok && err != nil {
		return err
	}
	params := slack.PostMessageParameters{}

	chanID, _, err := SlackToken.PostMessage(LogChannel,
		fmt.Sprintf("```error %v\n ====\n %v ```", errmsg, errStr),
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
