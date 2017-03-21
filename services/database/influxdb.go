package database

import (
	"fmt"
	"os"
	"time"

	"github.com/b-eee/amagi/services/configctl"

	utils "github.com/b-eee/amagi"
	influx "github.com/influxdata/influxdb/client/v2"
)

var (
	// InfluxDB database name
	InfluxDB = "apicore_db"

	// APIMonitorDB apimonitoring db name
	APIMonitorDB = "apicore_api_access"

	// InfluxClient influx db client
	InfluxClient influx.Client
)

// ConnectInfluxDB start connecting to influxdb
func ConnectInfluxDB() error {
	host, conf := influxdbConf()
	if os.Getenv("INFLUX_MONITOR") == "0" {
		return nil
	}

	utils.Info(fmt.Sprintf("connecting to influxdb"))
	client, err := influx.NewHTTPClient(conf)
	if err != nil {
		utils.Error(fmt.Sprintf("error connectInfluxDB %v", err))
		return err
	}

	InfluxClient = client
	utils.Info(fmt.Sprintf("connected to influxDB %v", host))

	// initialize database
	for _, db := range []string{InfluxDB, APIMonitorDB} {
		dbHostname, err := dbNameConstruct(db)
		if err != nil {
			continue
		}
		if err := createDatabaseIfNotExists(dbHostname, InfluxClient); err != nil {
			continue
		}
	}

	return nil
}

func influxdbConf() (string, influx.HTTPConfig) {
	conf := configctl.GetDatabaseConf("influxdb")
	if len(conf.Host) == 0 {
		// disable INFLUX_MONITOR
		os.Setenv("INFLUX_MONITOR", "0")
		utils.Warn(fmt.Sprintf("INFLUX_MONITOR status=disabled"))
		return "", influx.HTTPConfig{}
	}

	strHost := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
	httpConf := influx.HTTPConfig{
		Addr:    strHost,
		Timeout: time.Duration(2 * time.Second),
	}

	if len(conf.Username) != 0 && len(conf.Password) != 0 {
		httpConf.Username = conf.Username
		httpConf.Password = conf.Password
	}

	return strHost, httpConf
}

func dbNameConstruct(dbName string) (string, error) {
	// TODO REFACTOR TO USE REGEX INSTEAD OF REPLACE!
	return fmt.Sprintf("%v_%v", dbName, os.Getenv("ENV")), nil
}

// DbNameWHostname database name with hostname
func DbNameWHostname(dbName string) string {
	db, err := dbNameConstruct(dbName)
	if err != nil {
		return ""
	}
	return db
}

// CreateDatabase create database public
func CreateDatabase(dbName string) error {
	if err := createDatabaseIfNotExists(dbName, InfluxClient); err != nil {
		return err
	}
	return nil
}

func createDatabaseIfNotExists(dbName string, influxClient influx.Client) error {
	s := time.Now()
	q := influx.NewQuery(fmt.Sprintf("CREATE DATABASE %v", dbName), "", "")
	if _, err := influxClient.Query(q); err != nil {
		utils.Error(fmt.Sprintf("createDatabaseIfNotExists error %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("createDatabaseIfNotExists took: %v dbName=%v", time.Since(s), dbName))
	return nil
}
