package mailerBackend

import (
	"os"
)

var (
	// CurrentBackendConf current email backend config or flag
	CurrentBackendConf string

	// BackendFlagENV backend email flag environment variable
	BackendFlagENV = "EMAIL_BACKEND_FLAG"
)

type (
	// Request email request
	Request struct {
		SGTemplateID string
		Subject      string
		Receivers    []string
		Sender       string
		TemplateData map[string]interface{}
	}
)

// SetEmailBackend set email backend config
func SetEmailBackend() string {
	var backend string
	switch os.Getenv(BackendFlagENV) {
	case "postfix":
		backend = "postfix"
	case "smtp.google":
		backend = "smtp.google"
	}

	return backend
}

// Send send email from backend
func Send(r *Request) error {
	var err error
	switch SetEmailBackend() {
	case "postfix":
		err = PostfixSendEmail()
	case "smtp.google":
		err = SMTPGoogleSendEmail(*r)
	}

	return err
}
