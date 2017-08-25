package messaging

import (
	"encoding/json"
	// utils "github.com/b-eee/amagi"
)

// GetBytes get arbitrary interface bytes
func GetBytes(data interface{}) ([]byte, error) {
	// defer utils.ExceptionDump()

	m, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return m, nil
}
