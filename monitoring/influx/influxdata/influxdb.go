package influxdata

import (
	"expvar"
	"fmt"
	"os"
	"runtime"
	"time"

	utils "github.com/b-eee/amagi"
	dbUtils "github.com/b-eee/amagi/services/database"
	influx "github.com/influxdata/influxdb1-client/v2"
)

var (
	reportDelaySec = 60
)

type (
	// MetricGetter mettric getter func
	MetricGetter map[string]func() interface{}
)

// StartInfluxDB influx db monitoring
func StartInfluxDB(appData MetricGetter) {
	if flag := os.Getenv("INFLUX_MONITOR"); flag == "0" {
		fmt.Printf("skipping INFLUX_MONITOR value: %v\n", flag)
		return
	}
	utils.Info(fmt.Sprintf("StartInfluxDB monitoring.."))

	if err := dbUtils.ConnectInfluxDB(); err == nil && os.Getenv("INFLUX_MONITOR") == "1" {
		go MainLoop(appData)
	}
}

// MainLoop for start influx db monitoring
func MainLoop(appData MetricGetter) {
	// Create a new point batch
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  dbUtils.DbNameWHostname(dbUtils.InfluxDB),
		Precision: "s",
	})
	if err != nil {
		utils.Error(fmt.Sprintf("error MainLoop new bp %v", err))
		return
	}

	for t := range time.Tick(time.Duration(reportDelaySec) * time.Second) {
		_ = t
		s := time.Now()
		points := metricPoints()

		// metrics from app
		for k, v := range appData {
			points[k] = v()
		}

		// Create a point and add to batch
		tags := map[string]string{"apicore": "apicore_stats"}
		pt, err := influx.NewPoint("apicore_stats", tags, points, time.Now())
		if err != nil {
			utils.Warn(fmt.Sprintf("error influx NewPoint %v", err))
			continue
		}
		bp.AddPoint(pt)

		// Write the batch
		if err = dbUtils.InfluxClient.Write(bp); err != nil {
			utils.Warn(fmt.Sprintf("error influxdb newPoint write %v", err))
			continue
		}
		utils.Info(fmt.Sprintf("influxWrite success! took: %v", time.Since(s)))
	}
}

// metricPoints apicore metric points collector and builder
func metricPoints() map[string]interface{} {
	memstats := memStats()

	points := map[string]interface{}{
		// goroutines
		"goroutines": runtime.NumGoroutine(),
		"hostname":   HostName(),
		"created_at": time.Now(),

		// HTTP API calls
		// "secure_hits_counter": authentication.GetSecureHitsCounter(),

		// memstats
		"memstats_alloc":       memstats.Alloc,
		"memstats_total_alloc": memstats.TotalAlloc,
		"memstats_sys":         memstats.Sys,
		"memstats_lookups":     memstats.Lookups,
		"memstats_mallocs":     memstats.Mallocs,
		"memstats_frees":       memstats.Frees,

		// GC
		"gc_pause_total_ns": memstats.PauseTotalNs,
		"gc_paus_ns":        memstats.PauseNs,
	}

	return points
}

// memstats get memstats objects from expvar
func memStats() runtime.MemStats {
	memstatsFunc := expvar.Get("memstats").(expvar.Func)
	memstats := memstatsFunc().(runtime.MemStats)

	return memstats
}

// General statistics.
// Alloc      uint64 // bytes allocated and not yet freed
// TotalAlloc uint64 // bytes allocated (even if freed)
// Sys        uint64 // bytes obtained from system (sum of XxxSys below)
// Lookups    uint64 // number of pointer lookups
// Mallocs    uint64 // number of mallocs
// Frees      uint64 // number of frees

// Main allocation heap statistics.
// HeapAlloc    uint64 // bytes allocated and not yet freed (same as Alloc above)
// HeapSys      uint64 // bytes obtained from system
// HeapIdle     uint64 // bytes in idle spans
// HeapInuse    uint64 // bytes in non-idle span
// HeapReleased uint64 // bytes released to the OS
// HeapObjects  uint64 // total number of allocated objects

// Low-level fixed-size structure allocator statistics.
// Inuse is bytes used now.
// Sys is bytes obtained from system.
// StackInuse  uint64 // bytes used by stack allocator
// StackSys    uint64
// MSpanInuse  uint64 // mspan structures
// MSpanSys    uint64
// MCacheInuse uint64 // mcache structures
// MCacheSys   uint64
// BuckHashSys uint64 // profiling bucket hash table
// GCSys       uint64 // GC metadata
// OtherSys    uint64 // other system allocations

// Garbage collector statistics.
// NextGC        uint64 // next collection will happen when HeapAlloc ≥ this amount
// LastGC        uint64 // end time of last collection (nanoseconds since 1970)
// PauseTotalNs  uint64
// PauseNs       [256]uint64 // circular buffer of recent GC pause durations, most recent at [(NumGC+255)%256]
// PauseEnd      [256]uint64 // circular buffer of recent GC pause end times
// NumGC         uint32
// GCCPUFraction float64 // fraction of CPU time used by GC
// EnableGC      bool
// DebugGC       bool

// HostName get app hostname
func HostName() interface{} {
	hostname, err := os.Hostname()
	if err != nil {
		utils.Warn(fmt.Sprintf("can't get hostname %v", err))
		return ""
	}

	return hostname
}
