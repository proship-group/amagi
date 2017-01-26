package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"

	utils "github.com/b-eee/amagi"
)

var (
	// ESConn main elasticsearch connection var
	ESConn *elastic.Client
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
		IID       bson.ObjectId `json:"i_id"`
		DID       string        `json:"d_id"`
		IndexName string        `json:"index_name"`
		TypeName  string        `json:"type_name"`
		Item      interface{}   `json:"item"`
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
	termQuery := elastic.NewTermQuery(req.SearchField, req.SearchValues)
	searchResult, err := ESConn.Search().
		Query(termQuery).
		From(0).Size(100).
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
	// TODO DEPRECATE WHEN SESSION UNIFIED -JP
	MongodbStart()

	var resultItems []ResultItem
	var di DistinctItem
	for _, id := range result.Each(reflect.TypeOf(di)) {
		if t, ok := id.(DistinctItem); ok {
			resultItems = append(resultItems, ResultItem{
				IID: bson.ObjectIdHex(t.IID),
				DID: t.DID,
			})
		}
	}

	if err := GetItemsToMongodbByID(&resultItems); err != nil {
		return nil, err
	}

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

			s, sc := BeginMongo()
			defer sc.Close()

			var collectionName string
			switch item.IndexName {
			case "datastore":
				collectionName = fmt.Sprintf("items_%v", item.DID)
			default:
				collectionName = fmt.Sprintf("items_%v", item.DID)
			}

			c := sc.DB(Db).C(collectionName)
			var rItem map[string]interface{}
			if err := c.Find(bson.M{"_id": item.IID}).One(&rItem); err != nil {
				utils.Error(fmt.Sprintf("error find %v", err))
				return
			}

			sItems[index].Item = rItem
			utils.Info(fmt.Sprintf("search ResultItem took: %v", time.Since(s)))
		}(sitem, index)
	}

	wg.Wait()
	return nil
}
