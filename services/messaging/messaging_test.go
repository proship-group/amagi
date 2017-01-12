package messaging

import (
	"os"
	"testing"

	"github.com/b-eee/amagi/services/messaging/backend"
)

// go test -v -run=TestInitMessaging ./services/messaging
func TestInitMessaging(t *testing.T) {
	os.Setenv("ENV", "local")
	os.Setenv("MESSAGING_BACKEND", "nsq")

	if err := InitMessaging(); err != nil {
		t.Error(err)
	}

	if err := backend.TestConnSeq(); err != nil {
		t.Error(err)
	}
}
