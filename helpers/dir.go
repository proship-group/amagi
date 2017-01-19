package helpers

import (
	"fmt"
	"os"
)

// GetCwd get current working directory
func GetCwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}
