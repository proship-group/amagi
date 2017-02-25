package database

import (
	"fmt"
	utils "github.com/b-eee/amagi"
	"github.com/jmcvetta/neoism"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	// CypherDebugDefault cypher debug flag
	CypherDebugDefault = 1

	cypherDebugFlag = "DEBUG_CYPHER"
)

// ExecuteCypherQuery execute cypher query from params
func ExecuteCypherQuery(query neoism.CypherQuery) error {
	s := time.Now()
	err := Neo4jDB.Cypher(&query)
	_ = s

	str := []string{fmt.Sprintf("ExecuteCypher (took: %v) ", time.Since(s))}

	if flag, err := strconv.Atoi(os.Getenv(cypherDebugFlag)); flag == CypherDebugDefault || err != nil {
		str = append(str, fmt.Sprintf("%v", query))
	}
	utils.Info(strings.Join(str, " "))

	// utils.Info(fmt.Sprintf("ExecuteCypherQuery took: %v", time.Since(s)))
	return err
}

// TransacQuery transaction public query exec
func TransactionBegin(qs []*neoism.CypherQuery)  (*neoism.Tx, error)  {

	if flag, err := strconv.Atoi(os.Getenv(cypherDebugFlag)); flag == CypherDebugDefault || err != nil {
		str := []string{fmt.Sprintf("Begin Transactional Cypher ")}
		for key, q := range qs {
			str = append(str, fmt.Sprintf("\n%v >> %v", key, q))
		}
		utils.Info(strings.Join(str, " "))
	}

	return Neo4jDB.Begin(qs)
}

// TransactionErrRlbk display transaction error and send rollback
func TransactionErrRlbk(tx *neoism.Tx) error {
	for _, tErr := range tx.Errors {
		utils.Error(fmt.Sprintf("%v", tErr))
	}

	if err := tx.Rollback(); err != nil {
		utils.Error(fmt.Sprintf("error Rollback %v", err))
		return err
	}

	return nil
}

// TransacQuery transaction public query exec
func TransacQuery(qs []*neoism.CypherQuery) error {
	s := time.Now()

	if flag, err := strconv.Atoi(os.Getenv(cypherDebugFlag)); flag == CypherDebugDefault || err != nil {
		str := []string{fmt.Sprintf("Execute Transactional Cypher ")}
		for key, q := range qs {
			str = append(str, fmt.Sprintf("\n%v >> %v", key, q))
		}
		utils.Info(strings.Join(str, " "))
	}

	tx, err := Neo4jDB.Begin(qs)
	if err != nil {
		TransactionErrRlbk(tx)
		return err
	}

	if err := tx.Commit(); err != nil {
		TransactionErrRlbk(tx)
		return err
	}

	utils.Info(fmt.Sprintf("TransacQuery (took: %v)", time.Since(s)))
	return nil
}
