package uptime

import (
	netUtils "github.com/b-eee/amagi/api/net"
)

type (
	// HealthCheckResp health check response
	HealthCheckResp struct {
		StatusCode int    `json:"status_code"`
		HostName   string `json:"host_name"`
	}
)

// Healthy healthy response
func Healthy() HealthCheckResp {
	return HealthCheckResp{
		StatusCode: 200,
		HostName:   netUtils.AppHostName(),
	}
}
