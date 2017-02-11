package elasticsearch

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/b-eee/amagi/services/database"
	"gopkg.in/mgo.v2/bson"

	utils "github.com/b-eee/amagi"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	// ESConn main elasticsearch connection var
	ESConn *elastic.Client

	// DatastoreNode datastore node labels
	DatastoreNode = "ds:Datastore"
	// ProjectNode project node name
	ProjectNode = "pj:Project"
	// WorkspaceNode workspace nodename
	WorkspaceNode = "ws:Workspace"
	// DatastoreFields datastore fields
	DatastoreFields = "fs:Fields"
	// FieldNodeLabel field node label
	FieldNodeLabel = "f:Field"
	// RoleNodeLabel role node label name
	RoleNodeLabel = "r:Role"
	// UserNodeLabel user node name label
	UserNodeLabel = "u:User"

	// DatastoreCollection collection name
	DatastoreCollection = "data_stores"
)

type (
	// ESSearchReq elasticsearch search request
	ESSearchReq struct {
		IndexName  string
		Type       string
		Context    context.Context
		BodyJSON   interface{}
		FileBase64 string

		SearchName   string
		SearchField  string
		SearchValues interface{}

		SortField string
		SortAsc   bool

		UserID        string
		UserBasicInfo UserBasicInfo

		UpdateChanges map[string]interface{}
	}

	// UserBasicInfo user basic info
	UserBasicInfo struct {
		ID         string
		AccessKeys []string
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
		QID   string `bson:"q_id" json:"q_id"`
		IID   string `bson:"i_id" json:"i_id"`
		DID   string `bson:"d_id" json:"d_id"`
		FID   string `bson:"f_id" json:"f_id"`
		PID   string `bson:"p_id" json:"p_id"`
		WID   string `bson:"w_id" json:"w_id"`
		Index string `json:"index"`
		Value string `bson:"value" json:"value"`

		Attachment struct {
			Content interface{} `json:"content,omitempty"`
		} `json:"attachment,omitempty"`
	}

	// ResultItem result item from mongodb
	ResultItem struct {
		QID         string                 `json:"q_id"`
		IID         bson.ObjectId          `json:"i_id"`
		DID         string                 `json:"d_id"`
		FID         string                 `json:"f_id"`
		PID         string                 `json:"p_id"`
		WID         string                 `json:"w_id"`
		MaxScore    float64                `json:"max_score"`
		Order       int                    `json:"order"`
		IndexName   string                 `json:"index_name"`
		TypeName    string                 `json:"type_name"`
		ColSettings map[string]interface{} `json:"col_settings"`
		Value       string                 `json:"value"`
		Item        interface{}            `json:"item"`
		Title       string                 `json:"title"`
	}

	// Datastore datastore struct for mgo
	Datastore struct {
		ID               string    `bson:"_id" json:"_id"`
		DID              string    `bson:"d_id" json:"d_id"`
		DatastoreID      string    `bson:"datastoreID" json:"datastoreID"`
		CreatedAt        time.Time `bson:"created_at"`
		Deleted          bool      `bson:"deleted" json:"deleted"`
		Encoding         string    `json:"encoding"`
		Failed           bool      `json:"failed"`
		Imported         bool      `json:"imported"`
		Name             string    `json:"name"`
		NoStatus         bool      `json:"no_status"`
		Progress         int       `json:"progress"`
		ProjectID        string    `bson:"project_id" json:"project_id"`
		ShowInMenu       bool      `json:"show_in_menu"`
		Status           int       `json:"status"`
		StatusFieldIndex int       `json:"status_field_index"`
		Uploading        bool      `json:"uploading"`
	}

	// ItemModifier item modifier request
	ItemModifier struct {
		SItems         *[]ResultItem
		Item           *ResultItem
		Index          int
		UserID         string
		UserBasicInfo  UserBasicInfo
		CollectionName string
	}

	// GenericESItem generic elasticsearch item for save
	GenericESItem struct {
		QID       string `bson:"q_id" json:"q_id"`
		IID       string `bson:"i_id" json:"i_id"`
		DID       string `bson:"d_id" json:"d_id"`
		FID       string `bson:"f_id" json:"f_id"`
		PID       string `bson:"p_id" json:"p_id"`
		IndexName string `bson:"index_name" json:"index_name"`
		TypeName  string `bson:"type_name" json:"type_name"`
		Value     string `bson:"value" json:"value"`
	}

	// ItemMofifiersProc item modifier processors
	ItemMofifiersProc func(ItemModifier) error

	// ItemMap generic item map
	ItemMap map[string]interface{}

	// ItemFieldSettings item field settings
	ItemFieldSettings map[string]map[string]interface{}
)

// ESCreateIndex elasticsearch create index
func ESCreateIndex(indexName string) error {
	ctx := context.Background()
	// Create an index
	if _, err := database.ESGetConn().CreateIndex(indexName).Do(ctx); err != nil {
		// Handle error
		panic(err)
	}

	return nil
}

// ESAddDocument add document to the index
func (req *ESSearchReq) ESAddDocument() error {
	// indexname should be lowercase
	indexName := strings.ToLower(req.IndexName)

	if exists, err := database.ESGetConn().IndexExists(indexName).Do(CreateContext()); !exists || err != nil {
		utils.Error(fmt.Sprintf("index does not exists! err=%v creating index.. %v", err, indexName))
		if err := ESCreateIndex(strings.ToLower(indexName)); err != nil {
			return err
		}
	}

	s := time.Now()
	if _, err := database.ESGetConn().Index().
		Index(indexName).
		Type(req.Type).
		BodyJson(req.BodyJSON).
		Refresh("true").
		Do(CreateContext()); err != nil {

		utils.Error(fmt.Sprintf("error ESAddDocument %v", err))
		return err
	}

	utils.Info(fmt.Sprintf("ESaddDocument took: %v index=%v", time.Since(s), indexName))
	return nil
}

// ESDeleteDocument delete document
func (req *ESSearchReq) ESDeleteDocument() error {
	del := elastic.NewMatchQuery("i_id", req.BodyJSON.(DistinctItem).IID)

	res, err := elastic.NewDeleteByQueryService(database.ESGetConn()).
		// for multiple index search query, pass in slice of string
		Index(strings.Split(req.IndexName, ",")...).
		Query(del).
		Do(CreateContext())
	if err != nil {
		utils.Error(fmt.Sprintf("error ESDeleteDocument %v", err))
		return err
	}

	fmt.Println("deleted: ", res.Deleted)
	if res.Deleted == 0 {
		return fmt.Errorf("deleted: %v", res.Deleted)
	}

	return nil
}

// ESUpdateDocument update elasticSearch document
func (req *ESSearchReq) ESUpdateDocument() error {

	// TODO handle in switch instead
	// req.ESHTTPItemUpdate()
	return nil

}

type concurrentSearch struct {
	Query elastic.Query
	Field string
}

// ESTermQuery new term query
// manual settings for setting default tokenizer for kuromoji
// $ curl -u elastic -XPOST 'http://104.198.115.53:9400/datastore/_close'
// $ curl -u elastic -XPUT 'http://104.198.115.53:9400/datastore/_settings?preserve_existing=true' -d '{   "index.analysis.analyzer.default.tokenizer" : "kuromoji",   "index.analysis.analyzer.default.type" : "custom" }'
// $ curl -u elastic -XPOST 'http://104.198.115.53:9400/datastore/_open'
func (req *ESSearchReq) ESTermQuery(result *elastic.SearchResult) (*elastic.SearchResult, error) {
	// joinedText := buildRegexpString(req.SearchValues)
	// query := elastic.NewRegexpQuery(req.SearchField, joinedText).
	// 	Boost(1.2).Analyzer("analyzer")

	query := elastic.NewSimpleQueryStringQuery(fmt.Sprintf("%v", req.SearchValues)).
		Analyzer("kuromoji").
		Flags("OR|AND|PREFIX")

	searchResult, err := database.ESGetConn().Search().
		Highlight(ResultHighlighter(req.SearchField)).
		Query(query).
		From(0).
		Do(CreateContext())
	if err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("ESTermQuery took: %v ms hits: %v", searchResult.TookInMillis, searchResult.Hits.TotalHits))
	return searchResult, nil
}

// ResultHighlighter create result highlighter
func ResultHighlighter(field string) *elastic.Highlight {
	return elastic.NewHighlight().
		Fields(elastic.NewHighlighterField(field)).
		PreTags("<em class='searched_em'>").PostTags("</em")

}

// ESBulkDeleteDocuments bulk delete elasticsearch document
func ESBulkDeleteDocuments(requests ...ESSearchReq) error {
	for _, req := range requests {
		if err := req.ESDeleteDocument(); err != nil {
			continue
		}
	}

	return nil
}

// ESBulkAddDocuments bulk delete elasticsearch document
func ESBulkAddDocuments(requests ...ESSearchReq) error {
	for _, req := range requests {
		if err := req.ESAddDocument(); err != nil {
			continue
		}
	}

	return nil
}

func buildRegexpString(str interface{}) string {
	var st []string
	for _, t := range strings.Split(fmt.Sprintf("%v", str), " ") {

		// regexps
		st = append(st, fmt.Sprintf("%v", t))
		st = append(st, fmt.Sprintf("%v.*", t))
		st = append(st, fmt.Sprintf("*.%v", t))
		st = append(st, fmt.Sprintf("(%v)", t))
		// st = append(st, fmt.Sprintf("%v", t))
		// st = append(st, fmt.Sprintf("[%v]", t))
	}

	return strings.Join(st, "|")
}

// CreateContext create context
func CreateContext() context.Context {
	return context.Background()
}

// ESSearchItems search items to mongodb by ID
func ESSearchItems(result elastic.SearchResult, esSearchReq ESSearchReq) ([]ResultItem, error) {
	s := time.Now()
	// TODO DEPRECATE WHEN SESSION UNIFIED -JP
	database.MongodbStart()
	database.StartNeo4j()
	defer func() {
		database.MongodbSession.Close()
	}()

	var resultItems []ResultItem
	var di DistinctItem

	for index, id := range result.Each(reflect.TypeOf(di)) {
		// original searched item before reflect
		searchItem := result.Hits.Hits[index]
		indexName := result.Hits.Hits[index].Index
		typeName := result.Hits.Hits[index].Type

		if t, ok := id.(DistinctItem); ok {
			resultItem := ResultItem{
				DID:       t.DID,
				FID:       t.FID,
				PID:       t.PID,
				QID:       t.QID,
				Order:     index,
				IndexName: indexName,
				TypeName:  typeName,
				MaxScore:  (*searchItem.Score),
			}

			if len(t.IID) != 0 {
				resultItem.IID = bson.ObjectIdHex(t.IID)
			}

			if len(result.Hits.Hits[index].Highlight["value"]) != 0 {
				resultItem.Value = result.Hits.Hits[index].Highlight["value"][0]
			}

			resultItems = append(resultItems, resultItem)
		}
	}

	if err := GetItemsByCollections(&resultItems, esSearchReq); err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("ESSearchItems took: %v", time.Since(s)))
	return resultItems, nil
}

type protectedObjectResults struct {
	sync.RWMutex
	ObjectMapRes map[bson.ObjectId]ResultItem
}

func (po *protectedObjectResults) Get(key bson.ObjectId) ResultItem {
	po.RLock()
	defer po.RUnlock()

	if value, exists := po.ObjectMapRes[key]; exists {
		return value
	}

	return ResultItem{}
}

func (po *protectedObjectResults) Set(key bson.ObjectId, resultItem ResultItem) {
	po.Lock()
	po.ObjectMapRes[key] = resultItem
	po.Unlock()
}

// GetItemsByCollections get items by group of collections
func GetItemsByCollections(searchedItems *[]ResultItem, esSearchReq ESSearchReq) error {
	s := time.Now()

	sItems := (*searchedItems)
	groupedItems := make(map[string][]ItemModifier)

	for index := range sItems {
		itemModifier := ItemModifier{
			SItems:        searchedItems,
			Item:          &sItems[index],
			Index:         index,
			UserBasicInfo: esSearchReq.UserBasicInfo,
		}

		switch itemModifier.Item.IndexName {

		// condition block will default to search datastore collections
		case "datastore", "histories":
			collection := fmt.Sprintf("items_%v", itemModifier.Item.DID)

			if len(groupedItems[collection]) == 0 {
				groupedItems[collection] = []ItemModifier{}
			}
			groupedItems[collection] = append(groupedItems[collection], itemModifier)
		case "queries":
			collection := "queries"
			itemModifier.CollectionName = collection
			groupedItems[collection] = append(groupedItems[collection], itemModifier)
		}
	}

	// fix fatal error: concurrent map writes
	protectedItem := protectedObjectResults{ObjectMapRes: make(map[bson.ObjectId]ResultItem)}
	var newResults []ResultItem
	var wg sync.WaitGroup
	for key := range groupedItems {
		wg.Add(1)

		go func(k string) {
			defer wg.Done()

			var ids []bson.ObjectId
			for _, s := range groupedItems[k] {
				if len(s.Item.IID) != 0 && s.Item.IndexName != "queries" {
					ids = append(ids, s.Item.IID)
				}
				switch s.Item.IndexName {
				case "queries":
					ids = append(ids, bson.ObjectIdHex(s.Item.QID))
					protectedItem.Set(bson.ObjectIdHex(s.Item.QID), (*s.Item))
				default:
					protectedItem.Set(s.Item.IID, (*s.Item))
				}

			}

			// for additional data
			results, err := FindItemsInCollectionByIDS(k, ids...)
			if err != nil {
				return
			}
			for _, r := range results {
				id := r["_id"].(bson.ObjectId)
				updateObj := ResultItem{
					IID:       id,
					Value:     protectedItem.Get(id).Value,
					DID:       protectedItem.Get(id).DID,
					PID:       protectedItem.Get(id).PID,
					IndexName: protectedItem.Get(id).IndexName,
					TypeName:  protectedItem.Get(id).TypeName,
				}
				if title, exists := r["title"].(string); exists {
					updateObj.Title = title
				}

				if pid, exists := r["p_id"].(string); exists {
					updateObj.PID = pid
				}

				if pid, exists := r["project_id"].(string); exists {
					updateObj.PID = pid
				}

				if qid, exists := r["q_id"].(string); exists {
					updateObj.QID = qid
				}

				protectedItem.Set(id, updateObj)

			}

		}(key)
	}
	wg.Wait()

	for _, item := range protectedItem.ObjectMapRes {
		newResults = append(newResults, item)
	}

	(*searchedItems) = newResults
	utils.Info(fmt.Sprintf("GetItemsByCollections took: %v", time.Since(s)))
	return nil
}

// FindItemsInCollectionByIDS find items in collection by item ids
func FindItemsInCollectionByIDS(collectionName string, IDS ...bson.ObjectId) ([]map[string]interface{}, error) {
	if len(IDS) == 0 {
		return []map[string]interface{}{}, nil
	}

	var res []map[string]interface{}
	_, sc := database.BeginMongo()
	c := sc.DB(database.Db).C(collectionName)
	defer sc.Close()

	if err := c.Find(bson.M{"_id": bson.M{"$in": IDS}}).All(&res); err != nil {
		utils.Error(fmt.Sprintf("error in FindItemsInCollectionByIDS %v", err))
		return res, err
	}
	return res, nil
}
