package amagi

import (
	"strings"
)

// GetIndexOfCharInStr retrieve index of char in str
func GetIndexOfCharInStr(str, strPattern string) int {
	return strings.Index(str, strPattern) + 1
}
