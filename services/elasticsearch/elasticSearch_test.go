package elasticsearch

import (
	"os"
	"testing"

	dbUtils "github.com/b-eee/amagi/services/database"
)

func init() {
	os.Setenv("ENV", "local")
	dbUtils.StartElasticSearch()
}

func TestESDeleteDocument(t *testing.T) {
	req := ESSearchReq{
		Type:      "field",
		IndexName: "datastore,histories",
		BodyJSON: DistinctItem{
			IID: "589195b26aeb58d98c5dea9f",
		},
	}

	if err := req.ESDeleteDocument(); err != nil {
		t.Error(err)
	}
}

func TestESUpdateDocument(t *testing.T) {
	req := ESSearchReq{
		Type:      "field",
		IndexName: "datastore",
		BodyJSON: DistinctItem{
			IID:   "5891a3c26aeb58d3d0b73c0a",
			FID:   "de18fb55-884f-44dd-8b5e-6d91f9392ea1",
			Value: "testing!",
		},
		UpdateChanges: map[string]interface{}{
			"de18fb55-884f-44dd-8b5e-6d91f9392ea1": "testing!",
		},
	}

	if err := req.ESUpdateDocument(); err != nil {
		t.Error(err)
	}
}
