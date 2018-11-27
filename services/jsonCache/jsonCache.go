package jsonCache

import (
	"encoding/json"
	"fmt"
	"strings"
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
		utils.Info(fmt.Sprintf("JSONCacheSet %v", err))
		return err
	}

	return utils.Info(fmt.Sprintf("JSONCacheSet took: %v key: %v ===============+================", time.Since(s), keys))
}

// CacheGetEx get string cache and convert to json
func CacheGetEx(keys []string, sult interface{}, ttl int) (string, error) {
	s := time.Now()
	str, err := GetEx(joinKeysToSTR(keys...), ttl)
	if err != nil {
		utils.Info(fmt.Sprintf("JSONCacheGetEx Hit miss! %v", err))
		return "", err
	}

	// var result interface{}
	// // fmt.Println(string(str))
	// if err := json.Unmarshal(str, &sult); err != nil {
	// 	utils.Info(fmt.Sprintf("error unmarshal json cache string %v", err))
	// 	fmt.Println("-------------------ssss------------ssss------------ssss------------ssss------------ssss")
	// 	return err
	// }

	// fmt.Println("--------------------xxxx----------------xxxx----------------xxxx----------------xxxx----------------xxxx----------------xxxx----------------xxxx")
	// fmt.Println(result)
	// fmt.Println(sult)
	// fmt.Println("--------------------xxxx----------------xxxx----------------xxxx----------------xxxx----------------xxxx----------------xxxx----------------xxxx")

	return string(str), utils.Info(fmt.Sprintf("JSONCacheGetEx Hit! took: %v keys: %v", time.Since(s), joinKeysToSTR(keys...)))
}

// CacheDelete json cache delete
func CacheDelete(keys []string) error {
	s := time.Now()

	if err := Delete(joinKeysToSTR(keys...)); err != nil {
		return err
	}

	utils.Info(fmt.Sprintf("JSONCacheDelete took %v", time.Since(s)))
	return nil
}

func joinKeysToSTR(keys ...string) string {

	return strings.Join(keys, ":")
}
