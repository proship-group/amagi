package database

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	as "github.com/aerospike/aerospike-client-go"
	al "github.com/aerospike/aerospike-client-go/logger"
	utils "github.com/b-eee/amagi"
)

var (
	// ASClient aerospike client
	ASClient *as.Client
)

type (
	// ASQuery aerospike query request
	ASQuery struct {
		Namespace string
		SetName   string
		Key       interface{}
		ASKey     *as.Key
		Objects   interface{}
		Bins      []string

		Record *as.Record
	}
)

// StartAerospike start aerospike connection
func StartAerospike() {
	asHost := getAerospikeDBHost()
	utils.Info(fmt.Sprintf("aerospike connection to -->> %v", asHost))
	client, err := as.NewClient(asHost, 3000)
	if err != nil {
		utils.Error(fmt.Sprintf("error StartAerospike connection %v", err))
		return
	}

	var buf bytes.Buffer
	lgr := log.New(&buf, "logger: ", log.Lshortfile)
	al.Logger.SetLogger(lgr)
	al.Logger.SetLevel(al.DEBUG)

	ASClient = client
	utils.Info(fmt.Sprintf("aerospike connected to %v", asHost))
	return
}

func getAerospikeDBHost() string {
	host := os.Getenv("AEROSPIKE")
	if len(host) == 0 {
		panic(fmt.Errorf("aerospike not set"))
	}

	return host
}

// ASWriteObject aerospike write KV
func ASWriteObject(asQuery ASQuery) error {
	s := time.Now()
	key, err := as.NewKey(asQuery.Namespace, asQuery.SetName, asQuery.Key)
	if err != nil {
		return err
	}

	if err := ASClient.PutObject(nil, key, &asQuery.Objects); err != nil {
		return err
	}

	return utils.Info(fmt.Sprintf("ASWriteObject took: %v", time.Since(s)))
}

// ASReadKey aerospike read single or multiple bins from key
func ASReadKey(asQuery ASQuery) error {
	record, err := ASClient.Get(nil, asQuery.ASKey, asQuery.Bins...)
	if err != nil {
		return err
	}

	asQuery.Record = record

	return nil
}

// ASReadFromKey aerospike read value from key with return value reference
func ASReadFromKey(asQuery *ASQuery) error {

	return nil
}

// ASCreateKey aerospike create key
func ASCreateKey(asQuery ASQuery) (*as.Key, error) {
	key, err := as.NewKey(asQuery.Namespace, asQuery.SetName, asQuery.Key)
	if err != nil {
		return key, err
	}

	return key, nil
}

// ASReadRecordBins aerospike read record with bins
func ASReadRecordBins(asQuery ASQuery, target *as.Record) error {
	record, err := ASClient.Get(nil, asQuery.ASKey, asQuery.Bins...)
	if err != nil {
		return err
	}

	(*target) = *record
	return nil
}

// ASReadObject aerospike read object
func ASReadObject(asQuery ASQuery, target interface{}) error {
	if err := ASClient.GetObject(nil, asQuery.ASKey, &target); err != nil {
		utils.Error(fmt.Sprintf("error ASReadObject %v", err))
		return err
	}

	return nil
}

// GetASClient get aerospike client
func GetASClient() (*as.Client, time.Time) {
	return ASClient, time.Now()
}
