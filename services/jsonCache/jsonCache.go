package jsonCache

import (
	"encoding/json"
	"fmt"
	"time"

	utils "github.com/b-eee/amagi"
)

var (
	separator = ":"
)

// JSONStringify stringify json data
func JSONStringify(data interface{}) (string, error) {
	str, err := json.Marshal(data)
	if err != nil {
		utils.Info(fmt.Sprintf("error JSONStringify %v", err))
		return "", err
	}

	return string(str), nil
}

// CacheSetEx simple cache json data to string value, ttl in seconds
func CacheSetEx(value interface{}, keys []string, ttl int) error {
	s := time.Now()
	strVal, err := JSONStringify(value)
	if err != nil {
		return err
	}

	if err := SetEx(joinKeysToSTR(keys...), strVal, ttl); err != nil {
		utils.Info(fmt.Sprintf("CacheSetEx %v", err))
		return err
	}

	return utils.Info(fmt.Sprintf("CacheSetEx took: %v key: %v ", time.Since(s), keys))
}

// CacheGetEx get string cache and convert to json
func CacheGetEx(keys []string, target interface{}, ttl int) error {
	s := time.Now()
	str, err := GetEx(joinKeysToSTR(keys...), ttl)
	if err != nil {
		utils.Info(fmt.Sprintf("CacheGetEx Hit miss! %v ", err))
		return err
	}

	if err := json.Unmarshal(str, target); err != nil {
		utils.Info(fmt.Sprintf("error in CacheGetEx json Unmarshal %v", err))
		return err
	}

	// return string(str), utils.Info(fmt.Sprintf("CacheGetEx Hit! took: %v keys: %v", time.Since(s), joinKeysToSTR(keys...)))
	return utils.Info(fmt.Sprintf("CacheGetEx Hit! took: %v keys: %v", time.Since(s), joinKeysToSTR(keys...)))
}

// CacheDelete json cache delete
func CacheDelete(keys []string) error {
	s := time.Now()

	if err := Delete(joinKeysToSTR(keys...)); err != nil {
		return err
	}

	utils.Info(fmt.Sprintf("CacheDelete took %v", time.Since(s)))
	return nil
}

// StringToStruct util for unmarshalling string to struct
func StringToStruct(result string, target interface{}) error {
	if err := json.Unmarshal([]byte(result), &target); err != nil {
		return utils.Error(fmt.Sprintf("error in StringToStruct %v", err))
	}

	return nil
}

func joinKeysToSTR(keys ...string) string {

	return JoinKeyWords(keys...)
}
