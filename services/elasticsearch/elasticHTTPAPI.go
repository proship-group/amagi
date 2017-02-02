package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/b-eee/amagi/services/configctl"

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

func ESReqHTTPPost(index, apiname string, query []byte) error {
	esUrl := ESURL(index, apiname)

	utils.Info(fmt.Sprintf("sending post to %v", esUrl))

	req, err := http.NewRequest("POST", esUrl, bytes.NewBuffer(query))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	r := struct {
		Updated float64 `json:"updated"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		panic(err)
	}

	if r.Updated == 0 {
		return fmt.Errorf("updated failed updated=%v", r.Updated)
	}

	return nil
}

func ESURL(index, apiname string) string {
	env := configctl.GetDBCfgStngWEnvName("elasticsearch", os.Getenv("ENV"))

	return fmt.Sprintf("%v/%v/%v", env.Host, index, apiname)
}
