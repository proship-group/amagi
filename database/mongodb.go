package database

import (
	"amagi/config"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	utils "amagi"

	mongodb "gopkg.in/mgo.v2"
)

var (
	// MongodbSession main mongodb cluster session
	MongodbSession *mongodb.Session

	// Db current app database
	// set this in config.yml!
	Db                = "beee-dev"
	mongodbClusterKey = "mongodb_cluster1"

	verboseMongodb = false

	// ensureMaxWrite counter for max mongodbhost
	ensureMaxWrite = 1

	maxSyncTimeout time.Duration = 60
)

// MongodbStart start connecting to mongodb
func MongodbStart() {
	defer utils.ExceptionDump()

	cfe, host := setMongodbHost()
	utils.Info(fmt.Sprintf("connecting to mongodb.. %s", host))

	mongoDBDialInfo := buildMongodBconn(cfe, host)
	session, err := mongodb.DialWithInfo(&mongoDBDialInfo)
	if err != nil {
		// fmt.Println(err)
		panic(err)
	}
	utils.Info(fmt.Sprintf("connected to mongodb... %v", cfe.Host))
	MongodbSession = session

	if os.Getenv("MONGODB_DEBUG") == "1" {
		mongodb.SetDebug(true)
		var aLogger *log.Logger
		aLogger = log.New(os.Stderr, "", log.LstdFlags)
		mongodb.SetLogger(aLogger)
	}
	// MongodbSession.SetMode(mongodb.Monotonic, true)
	// PRINT MONGODB CONNECTED/LIVE SERVERS
	// printLiveServers(session)
}

func buildMongodBconn(cfe config.Environment, hosts string) mongodb.DialInfo {
	mongodbHosts := splitMongodbInstances(hosts)
	ensureMaxWrite = len(mongodbHosts)
	conn := mongodb.DialInfo{
		Addrs:    mongodbHosts,
		Timeout:  10 * time.Second,
		Username: cfe.Username,
		Password: cfe.Password,
		Direct:   false,
	}

	if cfe.TLS == "true" {
		conn.DialServer = func(addr *mongodb.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}

		conn.Addrs = []string{cfe.Host}
		conn.Direct = false
	} else {
		return conn
	}

	// SET CONFIGS UNLESS DEFINED
	if cfe.ReplicaSetName != "" {
		conn.ReplicaSetName = cfe.ReplicaSetName
	}

	if cfe.Source != "" {
		conn.Source = cfe.Source
	}

	return conn
}

// splitMongodbInstances split mongodb instances string to slices
func splitMongodbInstances(instances string) []string {
	var hosts []string
	for _, host := range strings.Split(instances, ",") {
		hosts = append(hosts, host)
	}

	return hosts
}

func setMongodbHost() (config.Environment, string) {
	env := config.GetDatabaseConf("mongodb")

	return env, fmt.Sprintf("%v", env.Host)
}

// SessionCopy make copy of a mongodb session
func SessionCopy() *mongodb.Session {
	sc := MongodbSession.Copy()

	// setMode ref: https://godoc.org/labix.org/v2/mgo#Session.SetMode
	sc.SetMode(mongodb.Eventual, true)

	// https://godoc.org/labix.org/v2/mgo#Session.SetSafe
	// W: 2  atleast two instances confirm of writes
	sc.SetSafe(&mongodb.Safe{WMode: "majority", W: ensureMaxWrite, FSync: true})

	sc.SetSyncTimeout(maxSyncTimeout * time.Second)
	sc.SetSocketTimeout(maxSyncTimeout * time.Second)
	return sc
}

// func printLiveServers(session *mongodb.Session) {
// 	utils.Info(fmt.Sprintf("liveServers %v", session.LiveServers()))
// }
