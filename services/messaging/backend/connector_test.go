package backend

import (
	"os"
	"testing"
)

func TestSetAndInitBackend(t *testing.T) {
	os.Setenv("MESSAGING_BACKEND", "nats")
	os.Setenv("ENV", "local")

	if conf := SetAndInitBackend(); (MSGBackendConfig{}) == conf {
		t.Errorf("MSGBackendConfig not set")
	}
}
