package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	utils "github.com/b-eee/amagi"
	config "github.com/b-eee/amagi/services/configctl"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/denisenkom/go-mssqldb"
)

var (
	// SQLDB sql db connection
	SQLDB *gorm.DB

	mssqlMAXIdleConns          = 100
	mssqlMaxOpenConns          = 1000
	mssqlConnMaxLifeTime int64 = 120
	reportDelaySec             = 30
)

// OpenMssqlConn open myssql connection
func OpenMssqlConn() *gorm.DB {
	defer utils.ExceptionDump()

	s := time.Now()
	utils.Info(fmt.Sprintf("connecting to MySQL.. %v", ConnStr()))
	db, err := gorm.Open("mysql", ConnStr())
	if err != nil {
		utils.Error(fmt.Sprintf("failed to connect database"))
		// TODO panic if connection failed!!
		// HANDLE THIS!
		// panic(err)
		log.Fatal(err)
	}

	db.DB().SetMaxIdleConns(mssqlMAXIdleConns)
	db.DB().SetMaxOpenConns(mssqlMaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(120 * time.Second))

	if os.Getenv("DEBUG_SQL") == "true" {
		db.LogMode(true)
		utils.Info(fmt.Sprintf("Entering DEBUG_SQL %v", os.Getenv("DEBUG_SQL")))
	}

	// SQLDB = db
	utils.Info(fmt.Sprintf("connected to MySQL took: %v to %v", time.Since(s), mySQLServer))
	return db
}

// StartGormMssql start gorm mysql connection
// TODO RENAME TO SQL -JP
func StartGormMssql(initializeTable func(sqlDB *gorm.DB)) {
	defer utils.ExceptionDump()

	SQLDB = OpenMssqlConn()

	go keepAliveConn(SQLDB)
	initializeTable(SQLDB)
}

func keepAliveConn(db *gorm.DB) {
	for t := range time.Tick(time.Duration(reportDelaySec) * time.Second) {
		_ = t
		if err := utils.KeepAlive(db); err != nil {
			utils.Error(fmt.Sprintf("===keepalive mysql Failed=== connections: %v", countOpenConnections()))
			SQLDB = OpenMssqlConn()
			continue
		}

		utils.Info(fmt.Sprintf("===keepalive mysql Success=== open_connections: %v", countOpenConnections()))
	}
}

func countOpenConnections() int {
	return SQLDB.DB().Stats().OpenConnections
}

// MssqlCopyConn use for concurrent connections. check if connec tion is nil, else create new
// TODO RENAME - JP
func MssqlCopyConn() *gorm.DB {
	if err := SQLDB.DB().Ping(); err != nil {
		SQLDB = OpenMssqlConn()
		return SQLDB
	}

	return SQLDB
}

var (
	mySQLServer string
	keepAlive   = 10
)

func startMysql() {
	utils.Info("start mysql")
}

// ConnStr construct sql connection string
func ConnStr() string {
	return buildMySQLConnStr()
}

type con struct {
	k string
	v interface{}
}

func buildMySQLConnStr() string {
	env, mySQLServerStr := setMySQLlHost()
	mySQLServer = mySQLServerStr

	var connStrOptions []string
	str := []con{
		con{k: "charset", v: "utf8"},
		con{k: "parseTime", v: "True"},
		con{k: "loc", v: "Local"},
		// con{k: "encrypt", v: true},
		// con{k: "TrustServerCertificate", v: false},
	}

	for _, st := range str {
		connStrOptions = append(connStrOptions, fmt.Sprintf("%v=%v", st.k, st.v))
	}

	return mysqlHostString(env, strings.Join(connStrOptions, "&"))
}

func setMySQLlHost() (config.Environment, string) {
	env := config.GetDatabaseConf("mysql")

	return env, fmt.Sprintf("%v:%v", env.Host, env.Port)
}

func mysqlHostString(env config.Environment, optionStrings string) string {
	host := fmt.Sprintf("%v:%v", env.Host, env.Port)
	protocol := "tcp"
	if env.Port == "" {
		host = fmt.Sprintf("%v", env.Host)
		protocol = "ip"
	}

	return fmt.Sprintf("%v:%v@%v(%v)/%v?%v", env.Username, env.Password, protocol, host, env.Database, optionStrings)
}
