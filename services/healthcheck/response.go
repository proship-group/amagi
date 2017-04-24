package healthcheck

import (
	"fmt"
	"os"
	"time"

	utils "github.com/b-eee/amagi"
	netUtils "github.com/b-eee/amagi/api/net"
)

type (
	// RespHealthCheck health check response
	RespHealthCheck struct {
		StatusCode int    `bson:"status_code" json:"status_code"`
		HostName   string `bson:"host_name" json:"host_name"`
		AppName    string `bson:"app_name" json:"app_name"`

		ResponseTime int           `bson:"response_time" json:"response_time,omitempty"`
		ResponseDate time.Duration `bson:"-" json:"response_date,omitempty"`
		Online       string        `bson:"online" json:"online,omitempty"`
	}
)

// Healthy healthy response
func Healthy(msg string, callback func(string, time.Time)) RespHealthCheck {
	callback(msg, time.Now())
	return RespHealthCheck{
		StatusCode: 200,
		HostName:   netUtils.AppHostName(),
		AppName:    os.Getenv("APP_NAME"),
	}
}

// ResponseLogMessage response logger message
func ResponseLogMessage(msg string, start time.Time) {
	if responseLog := os.Getenv("HEALTH_CHECK_LOG"); responseLog == "true" {
		utils.Info(fmt.Sprintf("healthCheck success took=%v", time.Since(start)))
	}
}
