package auth

import (
	"fmt"
	"os"
	"strconv"
)

var (
	// SessionExpireOffSetENV session expire offset environment variable name
	SessionExpireOffSetENV = "SESSION_EXPIRE_OFFSET"
	// expireOffset session token in seconds
	// default 4 days
	expireOffset = 345600
)

// sessionExpireOffSet get session expire offset from ENV
func sessionExpireOffSet() int64 {
	envOffset := os.Getenv(SessionExpireOffSetENV)
	if len(envOffset) == 0 {
		return int64(expireOffset)
	}

	i, err := strconv.ParseInt(envOffset, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("can't convert str %v to int64", envOffset))
	}

	return i
}
