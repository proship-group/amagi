package apiCounter

import (
	"fmt"

	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// IncrAPIAccess increment user API access per minute
func IncrAPIAccess(token *jwt.Token, path string, duration time.Duration) error {

	claims := token.Claims.(jwt.MapClaims)

	tags := map[string]string{
		"api_name":    path,
		"consumer_id": claims["sub"].(string),
	}
	fields := map[string]interface{}{
		"consumer_id":   claims["sub"],
		"created_at":    time.Now(),
		"api_name":      path,
		"response_time": duration.Nanoseconds() / 1000000,
	}

	// TODO DON'T RETURN ERROR ON DISABLE -JP
	if os.Getenv("INFLUX_MONITOR") == "0" {
		return fmt.Errorf("INFLUX_MONITOR DISABLED %v", os.Getenv("INFLUX_MONITOR"))
	}

	return APICounter(tags, fields)
}
