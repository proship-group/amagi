package database

import (
	"fmt"
	"log"
	"os"

	"github.com/b-eee/amagi/services/configctl"
	"gopkg.in/olivere/elastic.v5"

	utils "github.com/b-eee/amagi"
)

var (
	// ESConn main elasticsearch connection var
	ESConn *elastic.Client
)

// StartElasticSearch start elasticsearch connections
func StartElasticSearch() error {
	env := configctl.GetDBCfgStngWEnvName("elasticsearch", os.Getenv("ENV"))

	esURL := env.Host

	utils.Info(fmt.Sprintf("connecting to esURL=%v", esURL))

	client, err := elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetBasicAuth(env.Username, env.Password))
	if err != nil {

		utils.Fatal(fmt.Sprintf("error StartElasticSearch %v", err))
		return err
	}

	ESConn = client

	utils.Info(fmt.Sprintf("connected to esURL=%v", esURL))
	return nil
}

// ESGetConn get elasticsearch connection
func ESGetConn() *elastic.Client {
	return ESConn
}

// GetESConfigs get elasticsearch configs
func GetESConfigs() configctl.Environment {
	return configctl.GetDBCfgStngWEnvName("elasticsearch", os.Getenv("ENV"))
}
