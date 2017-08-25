package database

import (
	"os"
	"testing"
)

func TestStartElasticSearch(t *testing.T) {
	os.Setenv("ENV", "local")

	if err := StartElasticSearch(); err != nil {
		t.Error(err)
	}
}
