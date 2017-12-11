package helpers

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// LoadYamlWStruct load yaml file from with target struct/type
func LoadYamlWStruct(yamlData []byte, data interface{}) error {
	if err := yaml.Unmarshal(yamlData, data); err != nil {
		return err
	}

	return nil
}

// LoadFileFromGOpath load file from gopath
func LoadFileFromGOpath(filePath string) ([]byte, error) {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), filePath))
	if err != nil {
		return nil, err
	}

	return buf, nil
}
