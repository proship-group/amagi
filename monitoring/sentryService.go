package monitoring

import (
	"os"
	"sync"

	"github.com/getsentry/raven-go"
)

var (
	// SentryDSNEnv sentry dsn env name
	SentryDSNEnv = "SENTRY_DSN"
)

// SentryService sentry service reporting
// set and initialize sentry api settings
func SentryService(wg *sync.WaitGroup) error {
	defer wg.Done()

	if dsn := os.Getenv(SentryDSNEnv); len(dsn) != 0 {
		raven.SetDSN(dsn)
	}

	return nil
}
