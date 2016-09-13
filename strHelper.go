package amagi

import (
	// "fmt"
	"strings"
	// uslack "github.com/b-eee/amagi/api/slack"
)

// GetIndexOfCharInStr retrieve index of char in str
func GetIndexOfCharInStr(str, strPattern string) int {
	return strings.Index(str, strPattern) + 1
}
