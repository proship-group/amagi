package messaging

import (
	"os"
	"testing"
)

// go test -v -run=TestInitMessaging ./services/messaging
func TestInitMessaging(t *testing.T) {
	os.Setenv("ENV", "local")
	os.Setenv("MESSAGING_BACKEND", "nsq")

	if err := InitMessaging(); err != nil {
		t.Error(err)
	}
}
