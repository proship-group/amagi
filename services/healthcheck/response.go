package healthcheck

import (
	"time"

	"os"

	netUtils "github.com/b-eee/amagi/api/net"
)

type (
	// RespHealthCheck health check response
	RespHealthCheck struct {
		StatusCode int    `bson:"status_code" json:"status_code"`
		HostName   string `bson:"host_name" json:"host_name"`
		AppName    string `bson:"app_name" json:"app_name"`

		ResponseTime int           `bson:"response_time" json:"response_time,omitempty"`
		ResponseDate time.Duration `bson:"response_date" json:"response_date,omitempty"`
		Online       bool          `bson:"online" json:"online,omitempty"`
	}
)

// Healthy healthy response
func Healthy() RespHealthCheck {
	return RespHealthCheck{
		StatusCode: 200,
		HostName:   netUtils.AppHostName(),
		AppName:    os.Getenv("APP_NAME"),
	}
}
