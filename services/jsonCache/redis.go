package jsonCache

import (
	"fmt"
	"strings"
	"time"

	"github.com/b-eee/amagi/services/database"

	utils "github.com/b-eee/amagi"
	"github.com/garyburd/redigo/redis"
)

var (

	// KeyExpiration max key Time To Live seconds
	KeyExpiration = 10000000

	// DefaultKeySeperator default redis keys seperators for keywords
	DefaultKeySeperator = ":"
)

// Set set key and value string to redis
func Set(key, value string) error {
	if !database.RedisCacheEnabled {
		return fmt.Errorf("redis Set cache disabled")
	}

	s := time.Now()
	c := database.GetRedisConn()
	defer c.Close()

	if _, err := c.Do("SET", key, value); err != nil {
		utils.Error(fmt.Sprintf("erro Redis Set %v", err))
		return err
	}

	return utils.Info(fmt.Sprintf("Set took: %v", time.Since(s)))
}

// SetEx set KV with expire
func SetEx(key, value string, ttl int) error {
	if !database.RedisCacheEnabled {
		return fmt.Errorf("redis SetEx cache disabled")
	}

	s := time.Now()
	c := database.GetRedisConn()
	defer c.Close()

	if err := c.Send("MULTI"); err != nil {
		utils.Error(fmt.Sprintf("error Multi %v", err))
		return err
	}

	if err := c.Send("SET", key, value); err != nil {
		utils.Error(fmt.Sprintf("error Setex %v", err))
		return err
	}
	if err := c.Send("EXPIRE", key, ttl); err != nil {
		utils.Error(fmt.Sprintf("error EXPIRE %v", err))
		return err
	}

	if _, err := c.Do("EXEC"); err != nil {
		utils.Error(fmt.Sprintf("error EXEC %v", err))
		return err
	}

	return utils.Info(fmt.Sprintf("SetEx took: %v", time.Since(s)))
}

// GetEx get kv and update expire
func GetEx(key string, ttl int) ([]uint8, error) {
	if !database.RedisCacheEnabled {
		return []uint8{}, fmt.Errorf("redis SetEx cache disabled")
	}

	s := time.Now()
	c := database.GetRedisConn()
	defer c.Close()

	if err := c.Send("MULTI"); err != nil {
		return []uint8{}, utils.Error(fmt.Sprintf("error Multi %v", err))
	}

	if err := c.Send("EXPIRE", key, ttl); err != nil {
		return []uint8{}, utils.Error(fmt.Sprintf("error EXPIRE %v", err))
	}

	if err := c.Send("GET", key); err != nil {
		return []uint8{}, utils.Error(fmt.Sprintf("error in GET getex %v", err))
	}

	str, err := redis.Values(c.Do("EXEC"))
	if err != nil || str[1] == nil {
		return []uint8{}, utils.Error(fmt.Sprintf("error EXEC in GetEx %v", err))
	}

	return str[1].([]byte), utils.Info(fmt.Sprintf("GetEx took: %v", time.Since(s)))
}

// Get get value from redis
func Get(key string) (string, error) {
	if !database.RedisCacheEnabled {
		return "", fmt.Errorf("redis Get cache disabled")
	}

	c := database.GetRedisConn()
	defer c.Close()

	str, err := redis.String(c.Do("GET", key))
	if err != nil {
		utils.Info(fmt.Sprintf("error Redis Get %v", err))
		return "", err
	}

	return str, nil
}

// Delete delete key w/ values
func Delete(key string) error {
	if !database.RedisCacheEnabled {
		return fmt.Errorf("redis Delete cache disabled")
	}
	c := database.GetRedisConn()
	defer c.Close()

	_, err := c.Do("DEL", key)
	if err != nil {
		utils.Error(fmt.Sprintf("error Delete REdis key %v", err))
		return err
	}

	return nil
}

// DelByPattern delete keys by pattern
func DelByPattern(pattern ...string) error {
	s := time.Now()
	if !database.RedisCacheEnabled {
		return fmt.Errorf("redis delete by pattern disabled")
	}

	c := database.GetRedisConn()
	defer c.Close()

	keys, err := redis.Strings(c.Do("KEYS", JoinKeyWords(pattern...)))
	if err != nil {
		return err
	}

	if err := c.Send("MULTI"); err != nil {
		return utils.Error(fmt.Sprintf("error in DelByPattern MULTI %v", err))
	}

	for _, key := range keys {
		if err := c.Send("DEL", key); err != nil {
			utils.Error(fmt.Sprintf("error EXEC %v", err))
			continue
		}
	}

	if _, err := c.Do("EXEC"); err != nil {
		utils.Error(fmt.Sprintf("error EXEC %v", err))
		return err
	}

	return utils.Info(fmt.Sprintf("DelByPattern took: %v", time.Since(s)))
}

// JoinKeyWords join keywords to produce redis key
func JoinKeyWords(keyword ...string) string {
	return strings.Join(keyword, DefaultKeySeperator)
}
