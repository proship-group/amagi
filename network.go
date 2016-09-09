package amagi

import "os"

var (
	localDevEnv = "LOCAL_DEV"
)

// IsLocal detect if server is in local dev
func IsLocal() bool {
	local := os.Getenv(localDevEnv)
	if local == "" {
		return false
	}
	return true
}
