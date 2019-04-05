package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/b-eee/amagi/services/configctl"
	"gopkg.in/olivere/elastic.v5"

	utils "github.com/b-eee/amagi"
)

var (
	// ESConn main elasticsearch connection var
	ESConn *elastic.Client
)

// StartElasticSearch start elasticsearch connections
func StartElasticSearch() {
	env := configctl.GetDBCfgStngWEnvName("elasticsearch", os.Getenv("ENV"))
	esURL := env.Host
	utils.Info(fmt.Sprintf("connecting to elasticsearch.. %v", esURL))

	var client *elastic.Client
	var err error
	if debug, e := strconv.ParseBool(os.Getenv("DEBUG_ELASTIC")); e == nil && debug {
		client, err = elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(false),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
			elastic.SetBasicAuth(env.Username, env.Password))
	} else {
		client, err = elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(false),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
			elastic.SetBasicAuth(env.Username, env.Password))
	}
	if err != nil {
		utils.Fatal(fmt.Sprintf("error StartElasticSearch %v", err))
		// panic(err)
	}

	ESConn = client

	utils.Info(fmt.Sprintf("connected to elasticsearch... %v", esURL))
}

// ESClientExists client exists check
func ESClientExists() bool {
	return ESConn != nil
}

// ESGetConn get elasticsearch connection
func ESGetConn() *elastic.Client {
	// fmt.Println("client %v", ESConn == nil)
	return ESConn
}

// GetESConfigs get elasticsearch configs
func GetESConfigs() configctl.Environment {
	return configctl.GetDBCfgStngWEnvName("elasticsearch", os.Getenv("ENV"))
}
