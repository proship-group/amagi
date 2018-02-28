package mailer

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("SENDGRID_API_KEY", "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")
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
