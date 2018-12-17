package database

import (
	"fmt"
	"os"
	"time"

	utils "github.com/b-eee/amagi"
	config "github.com/b-eee/amagi/services/configctl"

	"github.com/garyburd/redigo/redis"
)

var (
	// RedisPool main redis pool
	RedisPool *redis.Pool

	// RedisIPENV redis ip environment name
	RedisIPENV = "REDIS_IP_ENV"
	// RedisPortENV redis port environment name
	RedisPortENV = "REDIS_PORT_ENV"
	// RedisAuthENV redis auth environment name
	RedisAuthENV = "REDIS_AUTH_ENV"

	// EnableCacheEnv enable cache environment variable
	EnableCacheEnv = "ENABLE_REDIS_CACHE"
	// RedisCacheEnabled switch for enable/disable cache
	RedisCacheEnabled = false
)

// StartRedis start connecting to redis
func StartRedis() {
	defer utils.ExceptionDump()
	var env config.Environment
	if redisHost := os.Getenv(RedisIPENV); len(redisHost) != 0 {
		env = config.Environment{
			Host:     redisHost,
			Password: os.Getenv(RedisAuthENV),
			Port:     os.Getenv(RedisPortENV),
		}
	} else {
		env = config.GetDatabaseConf("redis")
	}

	hostString := fmt.Sprintf("%v:%v", env.Host, env.Port)
	RedisPool = newPool(hostString, env.Password)
	RedisCacheEnabled = os.Getenv(EnableCacheEnv) == "1"
}

func newPool(serverConnStr, password string) *redis.Pool {
	utils.Info(fmt.Sprintf("create new redis pool! -->> %v", serverConnStr))
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", serverConnStr)
			if err != nil {
				utils.Error(fmt.Sprintf("error dialing to redis.."))
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}

			}
			return c, err
		},
	}
}

// GetRedisConn get redis connection from pool
func GetRedisConn() redis.Conn {
	if RedisPool == nil {
		StartRedis()
	}
	return RedisPool.Get()
}
