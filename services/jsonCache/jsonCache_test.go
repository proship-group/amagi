package jsonCache

import (
	"fmt"
	"os"
	"testing"

	"github.com/b-eee/amagi/services/database"
)

func init() {
	os.Setenv("ENV", "local")
	database.StartRedis()

	os.Setenv("ENABLE_REDIS_CACHE", "1")
}

type (
	// ApplicationDatastore application datastore struct
	ApplicationDatastore struct {
		Application struct {
			ApplicationID string `json:"application_id"`
			DisplayID     string `json:"display_id"`
			DisplayOrder  int    `json:"display_order"`
			// Datastores    []linkermodels.Datastore `json:"datastores"`
		} `json:"application"`
	}
)

func TestCacheGetEx(t *testing.T) {
	keys := []string{"applications", "list", "5bdabd6d8726e333886baaca"}
	var res interface{}
	err := CacheGetEx(keys, &res, 12000)
	if err != nil {
		t.Error(err)
	}

	// r := []struct {
	// 	Application linkermodels.Application `json:"application"`
	// }{}
	fmt.Println(res)
}

func TestDelByPattern(t *testing.T) {
	database.StartRedis()

	keys := []string{"get_paginate_items_with_search", "*", "5bdabd6d8726e333886baaca"}

	if err := DelByPattern(keys...); err != nil {
		t.Error(err)
	}

}
