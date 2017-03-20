package elasticsearch

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/b-eee/amagi/services/database"

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

	// IndexNameGlobalSearch elasticsearch index name
	IndexNameGlobalSearch  = "global_search"

	// TypeNameFullTextSearch elasticsearch index type name of full text search
	TypeNameFullTextSearch = "fulltext_search"
	// TypeNameFileSearch elasticsearch index type name of file search
	TypeNameFileSearch     = "file_search"

	FieldNameCategory = "category"
	FieldNameKeyType = "key_type"

	// IndexNameItem index category name for items
	IndexNameItems = "items"
	// IndexNameItem index category name for items
	IndexNameFieldValues = "field_values"
	// IndexNameQueries index category name for queries
	IndexNameQueries = "queries"
	// IndexNameNewActions index category name for new actions
	IndexNameNewActions = "newactions"
	// IndexNameHistories index category name for histories
	IndexNameHistories = "histories"
	// IndexNameFiles index category name for files
	IndexNameFiles = "files"
	// IndexNameItem index category name for datastores (use for delete index request)
	IndexNameDatastores = "datastores"
	// IndexNameItem index category name for projects (use for delete index request)
	IndexNameProjects = "projects"

	//IndexTypeNameMenu = "shortcut_menu"
	//IndexTypeNameComment = "item_comment"
	//IndexTypeNameFileContent = "file_content"
	//IndexTypeNameFieldValue = "field_value"

)

type (
	// ESSearchReq elasticsearch search request
	ESSearchReq struct {
		IndexName  string
		Type       string
		Context    context.Context
		BodyJSON   DistinctItem
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
		UserID      string
		WorkspaceID string
		AccessKeys  []string
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
		WID string `bson:"w_id" json:"w_id"`
		PID string `bson:"p_id" json:"p_id"`
		DID string `bson:"d_id" json:"d_id"`
		QID string `bson:"q_id" json:"q_id"`
		IID string `bson:"i_id" json:"i_id"`
		FID string `bson:"f_id" json:"f_id"`
		AID string `bson:"a_id" json:"a_id"`
		HID string `bson:"h_id" json:"h_id"`

		Category string `json:"category"`
		//KeyType  string `json:"key_type"`

		Title string `bson:"title" json:"title"`
		Value string `bson:"value" json:"value"`
		Keys  string `bson:"keys" json:"keys"`

		FileID     string `bson:"file_id" json:"file_id"`
		Attachment struct {
			Content interface{} `json:"content,omitempty"`
		} `json:"attachment,omitempty"`
	}
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

	// USE GLOBAL COMMON INDEX !! TODO: Refactor code ,HI
	req.IndexName = IndexNameGlobalSearch
	req.Type = TypeNameFullTextSearch

	// indexname should be lowercase
	indexName := strings.ToLower(req.IndexName)

	if exists, err := database.ESGetConn().IndexExists(indexName).Do(CreateContext()); !exists || err != nil {
		utils.Warn(fmt.Sprintf("index does not exists! err=%v creating index.. %v", err, indexName))
		if err := ESCreateIndex(strings.ToLower(indexName)); err != nil {
			return err
		}
	}

	_, err := database.ESGetConn().Index().
		Index(indexName).
		Type(req.Type).
		BodyJson(req.BodyJSON).
		Refresh("true").
		Do(CreateContext())
	if err != nil{
		utils.Error(fmt.Sprintf("error ESAddDocument %v", err))
		return err
	}

	//utils.Pretty(req.BodyJSON, "ES create document")
	//utils.Pretty(res,"add es document response")

	utils.Info(fmt.Sprintf("ESaddDocument took: %v [category=%v]",
		time.Since(s), req.BodyJSON.Category))
	return nil
}

// ESDeleteDocument delete document
func (req *ESSearchReq) ESDeleteDocument() error {
	s := time.Now()

	// USE GLOBAL COMMON INDEX !!  TODO: Refactor ,HI
	req.IndexName = IndexNameGlobalSearch

	// set delete query key-value
	var key, value string
	switch req.BodyJSON.Category {
	case IndexNameItems:
		key = "i_id"
		value = req.BodyJSON.IID
	case IndexNameFieldValues:
		key = "f_id"
		value = req.BodyJSON.FID
	case IndexNameNewActions:
		key = "a_id"
		value = req.BodyJSON.AID
	case IndexNameQueries:
		key = "q_id"
		value = req.BodyJSON.QID
	case IndexNameHistories:
		key = "h_id"
		value = req.BodyJSON.HID
	case IndexNameFiles:
		key = "file_id"
		value = req.BodyJSON.FileID
	case IndexNameDatastores:
		key = "d_id"
		value = req.BodyJSON.DID
	case IndexNameProjects:
		key = "p_id"
		value = req.BodyJSON.PID
	default:
		return fmt.Errorf("Invalid category [ %v ]", req.BodyJSON.Category )
	}

	res, err := elastic.NewDeleteByQueryService(database.ESGetConn()).
		// for multiple index search query, pass in slice of string
		Index(strings.Split(req.IndexName, ",")...).
		Query(elastic.NewBoolQuery().
			Must(
				elastic.NewMatchQuery(FieldNameCategory, req.BodyJSON.Category),
				elastic.NewMatchQuery(key, value),
			)).
		Do(CreateContext())
	if err != nil {
		utils.Error(fmt.Sprintf("error ESDeleteDocument %v", err))
		return err
	}

	//utils.Pretty(req,"delete documents")

	if res.Deleted == 0 {
		return fmt.Errorf("deleted: %v", res.Deleted)
	}

	utils.Info(fmt.Sprintf("ESDeleteDocument took: %v [category=%v]",
		time.Since(s), req.BodyJSON.Category))
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
func (req *ESSearchReq) ESTermQuery(result *elastic.SearchResult) (*elastic.SearchResult, error) {

	query := elastic.NewSimpleQueryStringQuery(fmt.Sprintf("%v", req.SearchValues)).
		Field("_all").
		Field("title").
		Field("value").
		Field("keys").
		Field("attachment.content").
		DefaultOperator("AND").
		AnalyzeWildcard(true)

	searchResult, err := database.ESGetConn().Search().
		Index(IndexNameGlobalSearch).
		Highlight(ResultHighlighter()).
		Query(query).
		PostFilter(elastic.NewMatchQuery("w_id", req.UserBasicInfo.WorkspaceID)).
		From(0).
		Size(200).
		Do(CreateContext())
	if err != nil {
		return nil, err
	}

	//utils.Pretty(searchResult, "Result of ES")

	utils.Info(fmt.Sprintf("ESTermQuery took: %v ms hits: %v", searchResult.TookInMillis, searchResult.Hits.TotalHits))
	return searchResult, nil
}

// ResultHighlighter create result highlighter
func ResultHighlighter() *elastic.Highlight {
	return elastic.NewHighlight().
		Fields(elastic.NewHighlighterField("title"),
			elastic.NewHighlighterField("value"),
			elastic.NewHighlighterField("keys"),
			elastic.NewHighlighterField("attachment.content")).
		PreTags("<em class='searched_em'>").PostTags("</em>")
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
			utils.Error(fmt.Sprintf("ESAddDocument error %v", err))
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
