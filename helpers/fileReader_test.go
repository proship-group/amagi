package helpers

import (
	"fmt"
	"testing"
)

func TestLoadYamlWStruct(t *testing.T) {
	st := struct {
		Name string `yaml:"name"`
	}{}

	byteSlice, err := LoadFileFromGOpath("/github.com/b-eee/amagi/helpers/test.yaml")
	if err != nil {
		t.Error(err)
	}
	if err := LoadYamlWStruct(byteSlice, &st); err != nil {
		t.Error(err)
	}

	fmt.Println(st)
}
