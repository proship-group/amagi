package helpers

import (
	"os"
	"strconv"
)

// GetEnvIntValue get environment integer value with default value
func GetEnvIntValue(envName string, defaultVal int) int {
	if value := os.Getenv(envName); len(value) != 0 {
		i, err := strconv.Atoi(value)
		if err != nil {
			return defaultVal
		}

		return i
	}

	return defaultVal
}
