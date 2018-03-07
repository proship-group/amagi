package mailer

import (
	"os"
	"testing"
)

func init() {
}

// TestSendEmail test sending email
func TestSendEmail(t *testing.T) {
	e := Email{
		From:          "tester@b-ee.com",
		To:            os.Getenv("TARGET"),
		Subject:       "test title",
		PlainTextBody: "weee",
	}
	if err := e.SendEmail(); err != nil {
		t.Error(err.Error())
	}
}
