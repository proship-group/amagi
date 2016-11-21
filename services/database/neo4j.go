package database

import (
	"fmt"

	utils "github.com/b-eee/amagi"
	config "github.com/b-eee/amagi/services/configctl"
	"github.com/jmcvetta/neoism"
)

var (
	neo4jConnection *neoism.Database

	// Neo4jDB public neo4j connection
	Neo4jDB *neoism.Database

	neo4jHost = "http://neo4j:jinpol0405@localhost:7474/db/data"
)

// ProjectIssueRelation project relationsship query result
type ProjectIssueRelation struct {
	// I neoism.Node
	I  string `json:"i.object_id"`
	RP string `json:"rp.object_id"`
}

// StartNeo4j make connection to neo4j
func StartNeo4j() {
	defer utils.ExceptionDump()

	utils.Info(fmt.Sprintf("connecting to neo4j.. %v", setNeo4jHost()))

	db, err := neoism.Connect(setNeo4jHost())
	if err != nil {
		panic(err)
	}

	neo4jConnection = db
	Neo4jDB = db

	utils.Info(fmt.Sprintf("connected to neo4j.."))
}

// setNeo4jHost set neo4j host
func setNeo4jHost() string {
	env := config.GetDatabaseConf("neo4j")
	var host string
	if env.Username != "" {
		host = fmt.Sprintf("%v:%v@", env.Username, env.Password)
	}

	return fmt.Sprintf("http://%v%v:%v/db/data", host, env.Host, env.Port)
}
