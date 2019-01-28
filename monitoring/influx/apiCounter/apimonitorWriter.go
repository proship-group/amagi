package apiCounter

import (
	"fmt"
	"time"

	utils "github.com/b-eee/amagi"
	dbUtils "github.com/b-eee/amagi/services/database"
	influx "github.com/influxdata/influxdb1-client/v2"
)

var (
	bp influx.BatchPoints
)

// InitAPIMonitor initialize api monitoring
func InitAPIMonitor() error {
	if err := dbUtils.CreateDatabase(dbUtils.InfluxDB); err != nil {
		return err
	}

	return nil
}

// APICounter api counter and store influxdb
func APICounter(tags map[string]string, fields map[string]interface{}) error {
	if bp == nil {
		initBatchPoints()
	}

	pt, err := influx.NewPoint("api_counter", tags, fields, time.Now())
	if err != nil {
		utils.Error(fmt.Sprintf("error APICounter newpoint %v", err))
		return err
	}
	bp.AddPoint(pt)

	if len(bp.Points()) >= 5 {
		if err := batchInsert(); err == nil {
			initBatchPoints()
			return nil
		}
	}

	return nil
}

func initBatchPoints() {
	bp = createBatchPoints()
}

func batchInsert() error {
	s := time.Now()
	if dbUtils.InfluxClient == nil {
		return nil
	}
	if err := dbUtils.InfluxClient.Write(bp); err != nil {
		utils.Error(fmt.Sprintf("batchInsert err! %v", err))
	}

	utils.Info(fmt.Sprintf("apicounter batchInsert (%v)count success! took %v", len(bp.Points()), time.Since(s)))
	return nil
}

func createBatchPoints() influx.BatchPoints {
	// Create a new point batch
	batchpoints, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  dbUtils.DbNameWHostname(dbUtils.InfluxDB),
		Precision: "s",
	})

	return batchpoints
}
