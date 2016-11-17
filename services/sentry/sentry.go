package sentry

import (
	"fmt"
	"os"

	"github.com/getsentry/raven-go"
)

var (
	// SentryENV sentry logging env name
	SentryENV = "SENTRY_LOGGER"
)

// SendToSentry send to sentry error message
func SendToSentry(msg string) error {

	if os.Getenv(SentryENV) == "1" {
		raven.CaptureError(fmt.Errorf(msg), nil)
	}

	return nil
}
