package server

import "os"

var (
	// ServerportENV server port environment name
	ServerportENV     = "SERVER_PORT_ENV"
	defaultServerPort = ":9000"
)

// GetServerPortEnv get server port from env
func GetServerPortEnv() string {
	if portEnv := os.Getenv(ServerportENV); len(portEnv) != 0 {
		return portEnv
	}

	return defaultServerPort
}
