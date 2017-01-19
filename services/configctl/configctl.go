package configctl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/b-eee/amagi/services/externalSvc"

	utils "github.com/b-eee/amagi"
	yaml "gopkg.in/yaml.v2"
)

// DatabaseConfig struct holder for dastabase config
type DatabaseConfig map[string]struct {
	Local      Environment `json:"local"`
	Dev        Environment `json:"dev"`
	Production Environment `json:"production"`
	Test       Environment `json:"test"`
}

// Environment struct for environment val
type Environment struct {
	Database       string `yaml:"database" json:"database"`
	Host           string `yaml:"host" json:"host"`
	Port           string `yaml:"port" json:"port"`
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password" json:"password"`
	ReplicaSetName string `yaml:"ReplicaSetName" json:"ReplicaSetName"`
	Source         string `yaml:"source" json:"source"`
	TLS            string `yaml:"TLS" json:"tls"`
	ExternalIP     string `yaml:"externalIP" json:"external_ip"`

	// MINIO
	MIOAccessKeyID     string `yaml:"minio_accesskey_id" json:"minio_accesskey_id"`
	MIOSecretAccessKey string `yaml:"minio_secretaccess_key" json:"minio_secretaccess_key"`
}

var (
	databaseFileName = "database.yml"

	// EnvDevStr env development string
	EnvDevStr = "dev"

	// EnvLocalStr env local string
	EnvLocalStr = "local"

	// EnvProductionStr env production string
	EnvProductionStr = "production"

	// EnvTestStr env production string
	EnvTestStr = "test"

	// EnvDockerStr	env docker string
	EnvDockerStr = "docker"

	// CurrentSetEnv current settted env
	CurrentSetEnv string

	// DBconfig app databases configs
	DBconfig DatabaseConfig
)

// InitConfig initialize app databases configs
func InitConfig(dbCallback func()) {
	dbConfig, err := ReadDb("")
	if err != nil {
		fmt.Println("error Init Config ", err)
		os.Exit(1)
	}

	DBconfig = dbConfig
	dbCallback()
}

// ReadDb read database yaml
func ReadDb(path string) (DatabaseConfig, error) {
	if path == "" {
		path = fmt.Sprintf("./config/%s", databaseFileName)
	}
	var db DatabaseConfig
	filename, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("error reading file path ", err)
	}
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		utils.Error(fmt.Sprintf("error reading database yaml %v", err))
		return db, err
	}

	err = yaml.Unmarshal(yamlFile, &db)
	if err != nil {
		utils.Error(fmt.Sprintf("error unmarshalling yamlfile %v", err))
		return db, err
	}

	return db, nil
}

// EnvironmentVar get environment variable
func EnvironmentVar(env string) string {
	CurrentSetEnv = os.Getenv(env)
	return CurrentSetEnv
}

// GetDatabaseConf get database config for db storage name
func GetDatabaseConf(dbStorageName string) Environment {

	return GetDBConfigSettings(dbStorageName)
}

// GetDBConfigSettings get DBconfig settings
func GetDBConfigSettings(dbName string) Environment {
	var result Environment
	configCtlURL := fmt.Sprintf("%v/db/get/%v/%v", ConfCtLHost(), dbName, os.Getenv("ENV"))
	res, err := externalSvc.HTTPGetRequest(configCtlURL, url.Values{})
	if err != nil || res.StatusCode == 404 || res.StatusCode == 500 {
		panic(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		panic(err)
	}

	return result
}

// GetDBCfgStngWEnvName get DBconfig settings with environment name
func GetDBCfgStngWEnvName(dbName, envName string) Environment {
	var result Environment
	configCtlURL := fmt.Sprintf("%v/db/get/%v/%v", ConfCtLHost(), dbName, envName)
	res, err := externalSvc.HTTPGetRequest(configCtlURL, url.Values{})
	if err != nil || res.StatusCode == 404 || res.StatusCode == 500 {
		panic(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		panic(err)
	}

	return result
}
