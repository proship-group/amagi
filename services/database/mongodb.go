package database

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	config "github.com/b-eee/amagi/services/configctl"

	utils "github.com/b-eee/amagi"
	mongodb "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// MongodbSession main mongodb cluster session
	MongodbSession *mongodb.Session

	// Db current app database
	// set this in config.yml!
	Db                string
	mongodbClusterKey = "mongodb_cluster1"

	verboseMongodb = false

	// ensureMaxWrite counter for max mongodbhost
	ensureMaxWrite = 1

	maxSyncTimeout time.Duration = 1

	// JournalPatternDB journal pattern database prefix name
	JournalPatternDB = "journal"

	// SumJournalCol summed journal collection name
	SumJournalCol = "sum_journal"

	JournalSource = "journal_source"
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

	setDatabaseName(cfe)

	// MongodbSession.SetMode(mongodb.Monotonic, true)
	// PRINT MONGODB CONNECTED/LIVE SERVERS
	printLiveServers(session)
}

func buildMongodBconn(cfe config.Environment, hosts string) mongodb.DialInfo {
	mongodbHosts := splitMongodbInstances(hosts)
	ensureMaxWrite = len(mongodbHosts)
	conn := mongodb.DialInfo{
		Addrs:    mongodbHosts,
		Timeout:  10 * time.Second,
		Source:   "admin",
		Database: cfe.Database,
		Username: cfe.Username,
		Password: cfe.Password,

		Direct: false,
	}

	if os.Getenv("HOST") != "" {
		conn.Addrs = []string{os.Getenv("HOST")}
	}

	// if os.Getenv("SECURE") == "true" {
	// 	conn.Username = cfe.Username
	// 	conn.Password = cfe.Password
	// }

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
	if os.Getenv("HOST") != "" {
		env.Host = os.Getenv("HOST")
		// env.Port = "27017"
	}

	return env, fmt.Sprintf("%v", env.Host)
}

// SessionCopy make copy of a mongodb session
func SessionCopy() *mongodb.Session {
	MongodbSession.Ping()

	sc := MongodbSession.Copy()

	// setMode ref: https://godoc.org/labix.org/v2/mgo#Session.SetMode
	sc.SetMode(mongodb.PrimaryPreferred, true)

	// https://godoc.org/labix.org/v2/mgo#Session.SetSafe
	// W: 2  atleast two instances confirm of writes
	sc.SetSafe(&mongodb.Safe{WMode: "majority", W: ensureMaxWrite, FSync: true})
	sc.SetSyncTimeout(maxSyncTimeout * time.Second)
	sc.SetSocketTimeout(maxSyncTimeout * time.Hour)
	return sc
}

func printLiveServers(session *mongodb.Session) {
	utils.Info(fmt.Sprintf("liveServers %v", session.LiveServers()))
}

// BeginMongo begin mongodb session with time now
func BeginMongo() (time.Time, *mongodb.Session) {
	return time.Now(), SessionCopy()

}

// setDatabaseName set database name
func setDatabaseName(env config.Environment) error {
	if env.Database != "" {
		Db = env.Database
		return nil
	}

	if dbFromEnv := os.Getenv("APP_MONGODB"); len(dbFromEnv) != 0 {
		Db = dbFromEnv
		utils.Info(fmt.Sprintf("APP_MONGODB set to=%v", Db))
	}

	return nil
}

// MongoInsert can be used as generic slice data collection insert to mongodb
func MongoInsert(collection string, data ...interface{}) error {
	s, sc := BeginMongo()
	c := sc.DB(Db).C(collection)
	defer sc.Close()

	if err := c.Insert(data...); err != nil {
		utils.Error(fmt.Sprintf("error MongoInsert %v collection: %v", err, collection))
		return err
	}

	utils.Info(fmt.Sprintf("MongoInsert took: %v collection: %v len(%v)", time.Since(s), collection, len(data)))
	return nil
}

// MongoUpdate update data in mongodb
func MongoUpdate(collection string, selector interface{}, update interface{}) error {
	s, sc := BeginMongo()
	c := sc.DB(Db).C(collection)
	defer sc.Close()

	if err := c.Update(selector, update); err != nil {
		utils.Error(fmt.Sprintf("error MongoUpdate %v collection: %v", err, collection))
		return err
	}

	utils.Info(fmt.Sprintf("MongoUpdate took: %v collection: %v", time.Since(s), collection))
	return nil
}

// MongoRemove mongo db remove collection
func MongoRemove(collection string, selector interface{}) error {
	s, sc := BeginMongo()
	c := sc.DB(Db).C(collection)
	defer sc.Close()

	if err := c.Remove(selector); err != nil {
		utils.Error(fmt.Sprintf("error MongoRemove on selector %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("MongoRemove took: %v", time.Since(s)))
	return nil
}

// MongoCreateCollection mongodb create collection with info
func MongoCreateCollection(collection string, info *mongodb.CollectionInfo) error {
	s, sc := BeginMongo()
	c := sc.DB(Db).C(collection)
	defer sc.Close()

	// info := mgo.CollectionInfo{ForceIdIndex: false, DisableIdIndex: true}
	if err := c.Create(info); err != nil {
		utils.Error(fmt.Sprintf("error MongoCreateCollection %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("createDatastoreCollection took: %v collectionName=%v", time.Since(s), collection))
	return nil

}

// CountCollection count items from a collection
func CountCollection(collection string, query bson.M) int {
	s, sc := BeginMongo()
	c := sc.DB(Db).C(collection)
	defer sc.Close()

	count, err := c.Find(query).Count()
	if err != nil {
		utils.Error(fmt.Sprintf("error CountCollection %v", err))
		return 0
	}

	utils.Info(fmt.Sprintf("CountCollection took: %v items(%v)", time.Since(s), count))
	return count
}

// MongoEnsureIndex ensure index in collection
func MongoEnsureIndex(collection string, index mongodb.Index) error {
	s, sc := BeginMongo()
	c := sc.DB(Db).C(collection)
	defer sc.Close()

	if err := c.EnsureIndex(index); err != nil {
		utils.Error(fmt.Sprintf("error MongoEnsureIndex %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("MongoEnsureIndex took: %v", time.Since(s)))
	return nil
}

// CollectionName collection name contstructor
// TODO Deprecate
func CollectionName(dtID string) string {
	return fmt.Sprintf("%v_%v", JournalPatternDB, dtID)
}

// ConstructColName contruct collection name with pattern
func ConstructColName(dtID, colPattern string) string {
	return fmt.Sprintf("%v_%v", colPattern, dtID)
}
