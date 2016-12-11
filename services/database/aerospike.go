package database

import (
	"fmt"
	"os"
	"time"

	as "github.com/aerospike/aerospike-client-go"
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
		Objects   struct{}
		Bins      []string
	}
)

// StartAerospike start aerospike connection
func StartAerospike() error {
	asHost := getAerospikeDBHost()
	utils.Info(fmt.Sprintf("aerospike connection to -->> %v", asHost))
	client, err := as.NewClient(asHost, 3000)
	if err != nil {
		utils.Error(fmt.Sprintf("error StartAerospike %v", err))
		return err
	}

	ASClient = client
	utils.Info(fmt.Sprintf("aerospike connected to %v", asHost))
	return nil
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

	utils.Info(fmt.Sprintf("ASWriteObject took: %v", time.Since(s)))
	return nil
}

// ASReadKey aerospike read single or multiple bins from key
func ASReadKey(asQuery ASQuery) error {
	record, err := ASClient.Get(nil, asQuery.ASKey, asQuery.Bins...)
	if err != nil {
		return err
	}

	fmt.Println(record)

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

// GetASClient get aerospike client
func GetASClient() (*as.Client, time.Time) {
	return ASClient, time.Now()
}
