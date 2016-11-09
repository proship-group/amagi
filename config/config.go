package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"amagi/internal/services"
)

// Environment struct for environment val
type Environment struct {
	Host           string
	Port           string
	Username       string
	Password       string
	Database       string
	ReplicaSetName string
	Source         string
	TLS            string `yaml:"TLS"`
}

var (
	// CurrentSetEnv current settted env
	CurrentSetEnv string
)

// InitConfig initialize app databases configs
func InitConfig(dbCallback func()) {
	dbCallback()
}

// GetDatabaseConf get database config for db storage name
func GetDatabaseConf(dbStorageName string) Environment {
	return GetDBConfigSettings(dbStorageName)
}

// GetDBConfigSettings get DBconfig settings
func GetDBConfigSettings(dbName string) Environment {
	var result Environment
	configCtlURL := fmt.Sprintf("%v/db/get/%v/%v", ConfCtLHost(), dbName, os.Getenv("ENV"))
	res, err := services.HTTPGetRequest(configCtlURL, url.Values{})
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		panic(err)
	}

	return result
}
