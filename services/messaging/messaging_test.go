package messaging

import (
	"os"
	"testing"
	// "github.com/b-eee/amagi/services/messaging/backend"
)

// go test -v -run=TestInitMessaging ./services/messaging
func TestInitMessaging(t *testing.T) {
	os.Setenv("ENV", "local")
	os.Setenv("MESSAGING_BACKEND", "nats")

	// if err := InitMessaging(); err != nil {
	// 	t.Error(err)
	// }

	b := InitMessaging()
	b.Publish(PublishReq{Topic: "test", Body: []byte("test")})

	// if err := backend.TestConnSeq(); err != nil {
	// 	t.Error(err)
	// }
}
