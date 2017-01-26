package database

import (
	"apicore/lib/database"
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"

	"strconv"

	utils "github.com/b-eee/amagi"
	"github.com/jmcvetta/neoism"
)

var (
	// ESConn main elasticsearch connection var
	ESConn *elastic.Client

	// DatastoreNode datastore node labels
	DatastoreNode = "ds:Datastore"
	// DatastoreFields datastore fields
	DatastoreFields = "fs:Fields"
	// FieldNodeLabel field node label
	FieldNodeLabel = "f:Field"
	// RoleNodeLabel role node label name
	RoleNodeLabel = "r:Role"
	// UserNodeLabel user node name label
	UserNodeLabel = "u:User"
)

type (
	// ESSearchReq elasticsearch search request
	ESSearchReq struct {
		IndexName string
		Type      string
		Context   context.Context
		BodyJSON  interface{}

		SearchName   string
		SearchField  string
		SearchValues interface{}

		SortField string
		SortAsc   bool

		UserID string
	}

	// Testing test struct
	Testing struct {
		User    string
		Message string
	}

	// DistinctItems distincted item values for elasticsearch
	DistinctItems struct {
		ID           string       `bson:"_id" json:"_id"`
		DistinctItem DistinctItem `bson:"distinct_item" json:"distinct_item"`
	}

	// DistinctItem unwinded item
	DistinctItem struct {
		IID   string `bson:"i_id" json:"i_id"`
		DID   string `bson:"d_id" json:"d_id"`
		FID   string `bson:"f_id" json:"f_id"`
		Index string `json:"index"`
		Value string `bson:"value" json:"value"`
	}

	// ResultItem result item from mongodb
	ResultItem struct {
		IID         bson.ObjectId          `json:"i_id"`
		DID         string                 `json:"d_id"`
		FID         string                 `json:"f_id"`
		MaxScore    float64                `json:"max_score"`
		Order       int                    `json:"order"`
		IndexName   string                 `json:"index_name"`
		TypeName    string                 `json:"type_name"`
		ColSettings map[string]interface{} `json:"col_settings"`
		Value       string                 `json:"value"`
		Item        interface{}            `json:"item"`
		Title       string                 `json:"title"`
	}

	// ItemMap generic item map
	ItemMap map[string]interface{}
)

// StartElasticSearch start elasticsearch connections
func StartElasticSearch() error {
	esURL := "http://104.198.115.53:9200"

	utils.Info(fmt.Sprintf("connecting to esURL=%v", esURL))

	client, err := elastic.NewClient(elastic.SetURL(esURL), elastic.SetSniff(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {

		utils.Fatal(fmt.Sprintf("error StartElasticSearch %v", err))
		return err
	}

	ESConn = client

	return nil
}

// ESCreateIndex elasticsearch create index
func ESCreateIndex(indexName string) error {
	ctx := context.Background()
	// Create an index
	if _, err := ESConn.CreateIndex(indexName).Do(ctx); err != nil {
		// Handle error
		panic(err)
	}

	return nil
}

// ESAddDocument add document to the index
func (req *ESSearchReq) ESAddDocument() error {
	s := time.Now()
	if _, err := ESConn.Index().
		Index(req.IndexName).
		Type(req.Type).
		Id("1").
		BodyJson(req.BodyJSON).
		Refresh("true").
		Do(CreateContext()); err != nil {
		return err
	}

	utils.Info(fmt.Sprintf("ESaddDocument took: %v", time.Since(s)))
	return nil
}

// ESTermQuery new term query
func (req *ESSearchReq) ESTermQuery(result *elastic.SearchResult) (*elastic.SearchResult, error) {
	// termQuery := elastic.NewTermQuery(req.SearchField, req.SearchValues)
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField(req.SearchField))
	hl = hl.PreTags("<em class='searched_em'>").PostTags("</em")
	// phrase := elastic.NewMatchQuery(req.SearchField, req.SearchValues)

	q := elastic.NewQueryStringQuery(fmt.Sprintf("%v", req.SearchValues))
	searchResult, err := ESConn.Search().
		Highlight(hl).
		Query(q).
		Do(CreateContext())
	if err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("ESTermQuery took: %v ms hits: %v", searchResult.TookInMillis, searchResult.Hits.TotalHits))
	return searchResult, nil
}

// CreateContext create context
func CreateContext() context.Context {
	return context.Background()
}

// ESSearchItems search items to mongodb by ID
func ESSearchItems(result elastic.SearchResult) ([]ResultItem, error) {
	s := time.Now()
	// TODO DEPRECATE WHEN SESSION UNIFIED -JP
	MongodbStart()
	defer MongodbSession.Close()

	var resultItems []ResultItem
	var di DistinctItem

	for index, id := range result.Each(reflect.TypeOf(di)) {
		// original searched item before reflect
		searchItem := result.Hits.Hits[index]
		if t, ok := id.(DistinctItem); ok {
			resultItems = append(resultItems, ResultItem{
				IID: bson.ObjectIdHex(t.IID),
				DID: t.DID,
				FID: t.FID,
				// Value:    result.Hits.Hits[index].Highlight["value"][0],
				Order:    index,
				MaxScore: (*searchItem.Score),
			})
		}
	}

	if err := GetItemsToMongodbByID(&resultItems); err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("ESSearchItems took: %v", time.Since(s)))
	return resultItems, nil
}

// GetItemsToMongodbByID get items to mongodb by ID
func GetItemsToMongodbByID(searchedItems *[]ResultItem) error {

	sItems := (*searchedItems)
	var wg sync.WaitGroup
	for index, sitem := range sItems {
		wg.Add(1)

		go func(item ResultItem, i int) {
			defer wg.Done()

			if err := findItem(searchedItems, &item, i); err != nil {
				return
			}
		}(sitem, index)
	}

	wg.Wait()
	return nil
}

func findItem(sItems *[]ResultItem, item *ResultItem, i int) error {
	s, sc := BeginMongo()
	colSettings := make(chan []map[string]interface{})
	defer close(colSettings)
	defer sc.Close()

	go GetDatastoreColSettings("588028a66aeb5890349a6c98", item.DID, colSettings)

	var collectionName string
	switch item.IndexName {
	case "datastore":
		collectionName = fmt.Sprintf("items_%v", item.DID)
	default:
		collectionName = fmt.Sprintf("items_%v", item.DID)
	}

	c := sc.DB(Db).C(collectionName)
	fmt.Println("item collection", item.DID)
	var rItem map[string]interface{}
	if err := c.Find(bson.M{"_id": item.IID}).One(&rItem); err != nil {
		utils.Error(fmt.Sprintf("error find %v", err))
		return nil
	}

	(*sItems)[i].Item = rItem
	(*sItems)[i].Value = item.Value

	cols := <-colSettings
	(*sItems)[i].Title = ItemTitle(rItem, cols...)
	(*sItems)[i].ColSettings = colSettingsMap(cols)
	_ = s
	// utils.Info(fmt.Sprintf("search findItem ResultItem took: %v", time.Since(s)))
	return nil
}

func colSettingsMap(colSettings []map[string]interface{}) map[string]interface{} {
	colMap := make(map[string]interface{})

	for _, col := range colSettings {
		k := col["id"].(string)
		colMap[k] = col
	}

	return colMap
}

// GetDatastoreColSettings get datastore column settings
func GetDatastoreColSettings(userID, datastoreID string, colSets chan []map[string]interface{}) error {
	s := time.Now()
	var colSettings []map[string]interface{}

	res := []struct {
		ColumnSettings struct {
			Data map[string]interface{} `json:"data"`
		} `json:"column_settings"`
	}{}

	cypher := neoism.CypherQuery{
		Statement: fmt.Sprintf(`
			MATCH (%v {d_id:{datastoreID}})-[]->(%v)-[]->(%v)<-[:CAN_USE]-(%v)<-[:HAS]-(%v {u_id:{userID}})
            RETURN DISTINCT f as column_settings
			UNION
			MATCH (%v {d_id:{datastoreID}})-[]->(%v)-[]->(%v)<-[:CAN_USE]-(gr:Role)<-[:HAS]-(pg:Group)-[:PARENT*0..]->(g:Group)<-[:HAS]-(%v {u_id:{userID}})
            RETURN DISTINCT f as column_settings
        `, DatastoreNode, DatastoreFields, FieldNodeLabel, RoleNodeLabel, UserNodeLabel,
			DatastoreNode, DatastoreFields, FieldNodeLabel, UserNodeLabel),
		Parameters: neoism.Props{"datastoreID": datastoreID, "userID": userID},
		Result:     &res,
	}
	if err := database.ExecuteCypherQuery(cypher); err != nil {
		utils.Error(fmt.Sprintf("error GetDatastoreColSettings %v", err))
		colSets <- colSettings
		return err
	}

	for _, i := range res {
		colSettings = append(colSettings, i.ColumnSettings.Data)
	}

	colSets <- colSettings
	utils.Info(fmt.Sprintf("GetDatastoreColSettings took: %v", time.Since(s)))
	return nil
}

// ItemTitle extract item ttiles
func ItemTitle(item map[string]interface{}, columnSettings ...map[string]interface{}) string {
	var titles []map[string]interface{}
	for _, item := range columnSettings {
		if value, exists := item["as_title"].(bool); exists && value {
			titles = append(titles, map[string]interface{}{
				"columnID":           item["id"].(string),
				"title_order_number": item["title_order"].(string),
			})
		}
	}

	titleBuilt := make([]string, len(titles))
	for _, t := range titles {
		val := item[t["columnID"].(string)]
		if index, err := strconv.Atoi(t["title_order_number"].(string)); err == nil {
			// TODO ADD ERR HANDLER IF OUT OF RANGE -JP
			titleBuilt[index-1] = fmt.Sprintf("%v", val)
		}
	}
	return strings.Join(titleBuilt, " ")
}
