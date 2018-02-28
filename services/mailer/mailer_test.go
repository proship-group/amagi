package mailer

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("SENDGRID_API_KEY", "")
}

// TestSendEmail test sending email
func TestSendEmail(t *testing.T) {
	e := Email{
		Subject: "test title",
	}
	if err := e.SendEmail(); err != nil {
		t.Error(err.Error())
	}
}
