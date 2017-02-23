package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/b-eee/amagi/services/configctl"
	"github.com/b-eee/amagi/services/database"

	utils "github.com/b-eee/amagi"
)

// ESHTTPItemUpdate elasticsearch update API
func (req *ESSearchReq) ESHTTPItemUpdate() error {
	itemReq := req.BodyJSON.(DistinctItem)
	query := fmt.Sprintf(`
		{
			"query":{
				"bool":{
					"must":[
						{
						"match":{
							"i_id":"%v"
						    }
						},
						{
						"match":{
							"f_id":"%v"
						    }
						}
					]
				}
			},
			"script":{
				"inline":"ctx._source.value = params.value",
				"lang":"painless",
				"params":{
					"value":"%v"
				}
			}
		}
	`, itemReq.IID, itemReq.FID, itemReq.Value)

	if err := ESReqHTTPPost(req.IndexName, "_update_by_query", []byte(query)); err != nil {
		utils.Error(fmt.Sprintf("error ESHTTPUpdate %v", err))

		// Create new index record instead
		if err := req.ESAddDocument(); err != nil {
			return err
		}

		return err
	}

	return nil
}

// ESDeleteHTTPByQuery es delete http api
func (req *ESSearchReq) ESDeleteHTTPByQuery(query map[string]interface{}) error {

	str, err := json.Marshal(query)
	if err != nil {
		utils.Error(fmt.Sprintf("error ESDeleteHTTPByQuery json marshal %v", err))
		return err
	}

	if err := ESReqHTTPPost(req.IndexName, "_delete_by_query", str); err != nil {
		// TODO error handling -jp
	}

	return nil
}

// ESFileAttachIndex file attach index using elasticsearch HTTP API instead
// http://stackoverflow.com/a/40334033/1175415
func (req *ESSearchReq) ESFileAttachIndex() error {
	s := time.Now()
	if err := createIngestPipeline(); err != nil {
		return err
	}

	if err := putFileIngestAttachment(req.BodyJSON.(DistinctItem), req.FileBase64); err != nil {
		return err
	}

	utils.Info(fmt.Sprintf("ESFileAttachIndex took: %v", time.Since(s)))
	return nil
}

func putFileIngestAttachment(item DistinctItem, fileBase64 string) error {
	query := fmt.Sprintf(`
		{
			"data": "%v",
			"i_id": "%v",
			"d_id": "%v",
			"w_id": "%v",
			"p_id": "%v"
		}	
	`, fileBase64, item.IID, item.DID, item.WID, item.PID)

	if err := ESReqHTTPPut(fmt.Sprintf("datastore/file/%v?pipeline=attachment", item.IID), []byte(query)); err != nil {
		return err
	}

	return nil
}

func createIngestPipeline() error {
	query := fmt.Sprintf(`
		{
			"description" : "Extract attachment information",
			"processors" : [
					{
						"attachment" : 
							{
								"field" : "data",
								"indexed_chars": -1
						}
					}
			]
		}	
	`)

	if err := ESReqHTTPPut("_ingest/pipeline/attachment", []byte(query)); err != nil {
		return err
	}

	return nil
}

// ESReqHTTPPut basic http put
func ESReqHTTPPut(api string, query []byte) error {
	esURL := esURLWoIndex(api)
	utils.Info(fmt.Sprintf("ESreqHTTPPut url %v", esURL))
	client := &http.Client{
		Timeout: time.Duration(60 * time.Second),
	}
	req, err := http.NewRequest("PUT", esURL, bytes.NewBuffer(query))
	if err != nil {
		utils.Error(fmt.Sprintf("error ESReqHTTPPut %v", err))
		return err
	}

	conf := database.GetESConfigs()
	req.SetBasicAuth(conf.Username, conf.Password)
	resp, err := client.Do(req)
	if err != nil {
		utils.Error(fmt.Sprintf("error on client.Do ESReqHTTPPut %v", err))
		return err
	}
	defer resp.Body.Close()
	utils.Info(fmt.Sprintf("ESREqHTTPPut for %v code=%v", api, resp.StatusCode))
	if resp.StatusCode == 400 {
		var r interface{}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {

		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}
	return nil
}

// ESReqHTTPPost basic es request http post
func ESReqHTTPPost(index, apiname string, query []byte) error {
	s := time.Now()
	esURL := esURL(index, apiname)

	utils.Info(fmt.Sprintf("sending post to %v", esURL))

	req, err := http.NewRequest("POST", esURL, bytes.NewBuffer(query))
	req.Header.Set("Content-Type", "application/json")

	conf := database.GetESConfigs()
	req.SetBasicAuth(conf.Username, conf.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	r := struct {
		Updated float64 `json:"updated"`
		Deleted float64 `json:"deleted"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		panic(err)
	}

	if r.Updated == 0 {
		return fmt.Errorf("updated failed updated=%v", r.Updated)
	}

	utils.Info(fmt.Sprintf("ESReqHTTPPost took: %v updated: %v deleted: %v", time.Since(s), r.Updated, r.Deleted))
	return nil
}

// ESReqHTTPDelete delete item by path and ID
func ESReqHTTPDelete(path string) error {
	esURL := esURL(path, "")

	utils.Info(fmt.Sprintf("sending post to %v", esURL))

	req, err := http.NewRequest("DELETE", esURL, bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")

	conf := database.GetESConfigs()
	req.SetBasicAuth(conf.Username, conf.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	utils.Info(fmt.Sprintf("ESReqHTTPDelete status=%v", resp.StatusCode))
	return nil
}

func esURLWoIndex(apiname string) string {
	env := configctl.GetDBCfgStngWEnvName("elasticsearch", os.Getenv("ENV"))
	return fmt.Sprintf("%v/%v", env.Host, apiname)
}

func esURL(index, apiname string) string {
	env := configctl.GetDBCfgStngWEnvName("elasticsearch", os.Getenv("ENV"))

	return fmt.Sprintf("%v/%v/%v", env.Host, index, apiname)
}
