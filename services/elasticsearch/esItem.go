package elasticsearch

import (
	"fmt"
	"sync"

	"gopkg.in/mgo.v2/bson"

	utils "github.com/b-eee/amagi"
	dbUtils "github.com/b-eee/amagi/services/database"
)

// ESItemSave elasticsearch item save
func ESItemSave(reqItem DistinctItem, fields []map[string]interface{}, wg *sync.WaitGroup) error {
	defer func() {
		wg.Done()
	}()

	var fieldIDs []string
	for _, v := range fields {
		if val, ok := v["search"].(bool); ok {
			if val {
				fieldIDs = append(fieldIDs, v["id"].(string))
			}
		}
	}

	item, err := ESgetItemByID(bson.ObjectIdHex(reqItem.IID), reqItem.DID)
	if err != nil {
		return err
	}

	for _, field := range fieldIDs {
		if value, exists := item[field]; exists {
			req := ESSearchReq{
				IndexName: "datastore",
				Type:      "field",
				BodyJSON: DistinctItem{
					IID:   reqItem.IID,
					PID:   reqItem.PID,
					WID:   reqItem.WID,
					DID:   reqItem.DID,
					FID:   field,
					Value: fmt.Sprintf("%v", value),
				},
			}
			if err := req.ESAddDocument(); err != nil {
				continue
			}
		}
	}

	return nil
}

// ESgetItemByID elastic search get item from mongodb
func ESgetItemByID(itemID bson.ObjectId, DID string) (map[string]interface{}, error) {
	var item map[string]interface{}
	_, sc := dbUtils.BeginMongo()
	c := sc.DB(dbUtils.Db).C(fmt.Sprintf("items_%v", DID))
	defer sc.Close()

	if err := c.Find(bson.M{"_id": itemID}).One(&item); err != nil {
		utils.Error(fmt.Sprintf("error ESgetItemByID %v", err))
		return item, err
	}

	return item, nil
}
