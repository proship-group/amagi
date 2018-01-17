package server

import "os"

var (
	// ServerportENV server port environment name
	ServerportENV = "SERVER_PORT_ENV"
)

// GetServerPortFrmEnv get server port from env if defined or return default from params
func GetServerPortFrmEnv(defaultServerPort string) string {
	if portEnv := os.Getenv(ServerportENV); len(portEnv) != 0 {
		return portEnv
	}

	return defaultServerPort
}
