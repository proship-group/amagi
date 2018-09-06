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
	ensureMaxWrite     = 1
	defaultMongodbPort = "27017"

	maxSyncTimeout time.Duration = 1

	// JournalPatternDB journal pattern database prefix name
	JournalPatternDB = "journal"

	// SumJournalCol summed journal collection name
	SumJournalCol = "sum_journal"

	JournalSource = "journal_source"

	// MongodbHostsWithPortENV mongodb host with port from env
	MongodbHostsWithPortENV = "MONGODB_HOST_ENV"

	// MongodbHostUserENV mongodb host user Env name
	MongodbHostUserENV = "MONGODB_USER"

	// MongodbPassENV mongodb password env
	MongodbPassENV = "MONGODB_PASS"
)

// MongodbStart start connecting to mongodb
func MongodbStart() {
	defer utils.ExceptionDump()

	cfe, host := setMongodbHost()

	mongoDBDialInfo := buildMongodBconn(cfe, host)
	session, err := mongodb.DialWithInfo(&mongoDBDialInfo)
	if err != nil {
		panic(fmt.Sprintf("failed to connect mongodb:%v", err))
	}
	MongodbSession = session

	if os.Getenv("MONGODB_DEBUG") == "1" {
		mongodb.SetDebug(true)
		var aLogger *log.Logger
		aLogger = log.New(os.Stderr, "", log.LstdFlags)
		mongodb.SetLogger(aLogger)
	}

	setDatabaseName(cfe)
	utils.Info(fmt.Sprintf("connected to mongodb... %v(db:%s)", host, Db))

	// MongodbSession.SetMode(mongodb.Monotonic, true)
	// PRINT MONGODB CONNECTED/LIVE SERVERS
	printLiveServers(session)
}

func buildMongodBconn(cfe config.Environment, hosts string) mongodb.DialInfo {
	mongodbHosts := splitMongodbInstances(hosts)
	ensureMaxWrite = len(mongodbHosts)
	for i, host := range mongodbHosts {
		// if the host doesnt have a port specified, add one from cfe.Port
		if !strings.Contains(host, ":") {
			mongodbHosts[i] = fmt.Sprintf("%s:%s", host, cfe.Port)
		}
	}
	utils.Info(fmt.Sprintf("connecting to mongodb.. %s", mongodbHosts))
	conn := mongodb.DialInfo{
		Addrs:    mongodbHosts,
		Timeout:  10 * time.Second,
		Source:   "admin",
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

		// conn.Addrs = []string{cfe.Host}
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
	var env config.Environment
	if len(os.Getenv(MongodbHostsWithPortENV)) == 0 {
		env = config.GetDatabaseConf("mongodb")
		if os.Getenv("HOST") != "" {
			env.Host = os.Getenv("HOST")
			// env.Port = "27017"
		}
		if env.Port == "" {
			// set default port
			env.Port = defaultMongodbPort
		}
	}

	// override mongodb env from configctl if set from environments
	if fromHost, err := setHostFromENVS(&env); fromHost && err == nil {
		return env, fmt.Sprintf("%v", env.Host)
	}

	return env, fmt.Sprintf("%v", env.Host)
}

func setHostFromENVS(env *config.Environment) (bool, error) {
	if host := os.Getenv(MongodbHostsWithPortENV); len(host) != 0 {
		env.Host = host
		env.Password = os.Getenv(MongodbPassENV)
		env.Username = os.Getenv(MongodbHostUserENV)
		return true, utils.Info(fmt.Sprintf("setHostFomrENVS is true"))
	}

	return false, utils.Info(fmt.Sprintf("setHostFomrENVS is false"))
}

// IsConnected Check connected
func IsConnected() bool {
	connected := true

	if MongodbSession == nil {
		connected = false
	}

	return connected
}

// SessionCopy make copy of a mongodb session
func SessionCopy() *mongodb.Session {
	// TRIAL prevent connection drop on master reschedule(on replica) -JP
	MongodbSession.Refresh()

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
	utils.Info(fmt.Sprintf("mongodb liveServers=%v", session.LiveServers()))
}

// BeginMongo begin mongodb session with time now
func BeginMongo() (time.Time, *mongodb.Session) {
	return time.Now(), SessionCopy()
}

// BeginMongoConn begin mongodb connection with session and collection
type BeginMongoConn struct {
	Time time.Time
	Conn *mongodb.Session
	Col  *mongodb.Collection
}

// BeginMongoWCol begin mongodb session with collection initialized
func BeginMongoWCol() func(string) BeginMongoConn {
	return func(cname string) BeginMongoConn {
		conn := SessionCopy()
		return BeginMongoConn{time.Now(), conn, conn.DB(Db).C(cname)}
	}
}

// setDatabaseName set database name
func setDatabaseName(env config.Environment) error {
	fmt.Println(os.Getenv("APP_MONGODB"), "APP_MONGODB")

	if dbFromEnv := os.Getenv("APP_MONGODB"); len(dbFromEnv) != 0 {
		Db = dbFromEnv
		utils.Info(fmt.Sprintf("APP_MONGODB set to=%v", Db))
		return nil
	}

	if env.Database != "" {
		Db = env.Database
		return nil
	}

	panic("mongodb database name is not set")
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
