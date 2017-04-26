package amagi

import (
	"strings"
	// uslack "github.com/b-eee/amagi/api/slack"
	"os"
)

var (
	// AppNamePrefixFromOwner app name prefixer from owner if b-eee or customer
	AppNamePrefixFromOwner = "APPNAME_PREFIX_OWNER"

	// DefaultAppNamePrefix default app name prefix
	DefaultAppNamePrefix = "beee"
)

// GetIndexOfCharInStr retrieve index of char in str
func GetIndexOfCharInStr(str, strPattern string) int {
	return strings.Index(str, strPattern) + 1
}

// AppNamePrefixer app name prefixer from ENV name app owners
func AppNamePrefixer(appName string) string {
	stringSep := "-"

	if prefix := os.Getenv(AppNamePrefixFromOwner); len(prefix) != 0 {
		return strings.Join([]string{prefix, appName}, stringSep)
	}

	return strings.Join([]string{DefaultAppNamePrefix, appName}, stringSep)
}
