package elasticsearch

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/b-eee/amagi/services/database"
	"github.com/jmcvetta/neoism"
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
		IndexName string
		Type      string
		Context   context.Context
		BodyJSON  interface{}

		SearchName   string
		SearchField  string
		SearchValues interface{}

		SortField string
		SortAsc   bool

		UserID        string
		UserBasicInfo UserBasicInfo
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

	// ItemMofifiersProc item modifier processors
	ItemMofifiersProc func(ItemModifier) error

	// ItemMap generic item map
	ItemMap map[string]interface{}
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
	s := time.Now()
	if _, err := database.ESGetConn().Index().
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
	hl := elastic.NewHighlight().
		Fields(elastic.NewHighlighterField(req.SearchField)).
		PreTags("<em class='searched_em'>").PostTags("</em")

	joinedText := buildRegexpString(req.SearchValues)
	regexpQuery := elastic.NewRegexpQuery(req.SearchField, joinedText).
		Boost(1.2).
		Flags("INTERSECTION|COMPLEMENT|EMPTY")

	searchResult, err := database.ESGetConn().Search().
		Highlight(hl).
		Query(regexpQuery).
		From(0).Size(50).
		Do(CreateContext())
	if err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("ESTermQuery took: %v ms hits: %v", searchResult.TookInMillis, searchResult.Hits.TotalHits))
	return searchResult, nil
}

func buildRegexpString(str interface{}) string {
	var st []string
	for _, t := range strings.Split(fmt.Sprintf("%v", str), " ") {

		// regexps
		st = append(st, fmt.Sprintf("%v.*?+", t))
		st = append(st, fmt.Sprintf("*.%v", t))
		st = append(st, fmt.Sprintf("%v", t))
		st = append(st, fmt.Sprintf("[%v]", t))
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

		if t, ok := id.(DistinctItem); ok {
			resultItem := ResultItem{
				IID:      bson.ObjectIdHex(t.IID),
				DID:      t.DID,
				FID:      t.FID,
				Order:    index,
				MaxScore: (*searchItem.Score),
			}

			if len(result.Hits.Hits[index].Highlight["value"]) != 0 {
				resultItem.Value = result.Hits.Hits[index].Highlight["value"][0]
			}

			resultItems = append(resultItems, resultItem)
		}
	}

	if err := GetItemsToMongodbByID(&resultItems, esSearchReq); err != nil {
		return nil, err
	}

	utils.Info(fmt.Sprintf("ESSearchItems took: %v", time.Since(s)))
	return resultItems, nil
}

// GetItemsToMongodbByID get items to mongodb by ID
func GetItemsToMongodbByID(searchedItems *[]ResultItem, esSearchReq ESSearchReq) error {

	sItems := (*searchedItems)
	var wg sync.WaitGroup
	for index, sitem := range sItems {
		wg.Add(1)

		go func(item ResultItem, i int) {
			defer wg.Done()
			itemModifier := ItemModifier{
				SItems:        searchedItems,
				Item:          &item,
				Index:         i,
				UserBasicInfo: esSearchReq.UserBasicInfo,
			}

			switch itemModifier.Item.IndexName {
			case "datastore":
				itemModifier.CollectionName = fmt.Sprintf("items_%v", itemModifier.Item.DID)
			default:
				itemModifier.CollectionName = fmt.Sprintf("items_%v", itemModifier.Item.DID)
			}

			var pg sync.WaitGroup
			procFinders := []ItemMofifiersProc{
				findItemByID,
				findDatastoreByID,
			}
			for _, pf := range procFinders {
				pg.Add(1)
				go func(procFinder ItemMofifiersProc) {
					defer pg.Done()

					if err := procFinder(itemModifier); err != nil {
						return
					}
				}(pf)
			}
			pg.Wait()
		}(sitem, index)
	}

	wg.Wait()
	return nil
}

func findDatastoreByID(itemModifier ItemModifier) error {
	_, sc := database.BeginMongo()
	c := sc.DB(database.Db).C(DatastoreCollection)
	defer sc.Close()

	var ds Datastore
	if err := c.Find(bson.M{"d_id": itemModifier.Item.DID}).One(&ds); err != nil {
		utils.Error(fmt.Sprintf("error findDatastoreByID %v", err))
		return err
	}

	(*itemModifier.SItems)[itemModifier.Index].PID = ds.ProjectID
	(*itemModifier.SItems)[itemModifier.Index].DID = itemModifier.Item.DID
	return nil
}

func findItemByID(itemModifier ItemModifier) error {
	s, sc := database.BeginMongo()
	colSettings := make(chan []map[string]interface{})
	defer close(colSettings)
	defer sc.Close()

	go GetDatastoreColSettings(itemModifier.UserBasicInfo.ID, itemModifier.Item.DID, colSettings)

	c := sc.DB(database.Db).C(itemModifier.CollectionName)
	var rItem map[string]interface{}
	if err := c.Find(bson.M{"_id": itemModifier.Item.IID, "access_keys": bson.M{"$in": itemModifier.UserBasicInfo.AccessKeys}}).One(&rItem); err != nil {
		utils.Error(fmt.Sprintf("error find %v", err))
		return nil
	}

	(*itemModifier.SItems)[itemModifier.Index].Item = rItem
	(*itemModifier.SItems)[itemModifier.Index].Value = itemModifier.Item.Value

	cols := <-colSettings
	(*itemModifier.SItems)[itemModifier.Index].Title = ItemTitle(rItem, cols...)
	(*itemModifier.SItems)[itemModifier.Index].ColSettings = colSettingsMap(cols)
	_ = s
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