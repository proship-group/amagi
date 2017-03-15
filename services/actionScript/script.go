package actionScript

import (
	"fmt"
	"os"

	"github.com/b-eee/amagi/services/externalSvc"

	utils "github.com/b-eee/amagi"
)

type (
	// Script script interface
	Script struct {
		Script string      `json:"script"`
		Data   interface{} `json:"data"`
	}
)

// TryScript try script
func (s *Script) TryScript() error {
	req := map[string]interface{}{
		"script": s.Script,
		"data":   s.Data,
	}

	fmt.Println("sending req")

	var resp map[string]interface{}
	if err := externalSvc.GenericHTTPRequesterWResp("POST", "http", actionScriptHost(), "/try", req, &resp); err != nil {
		utils.Error(fmt.Sprintf("error TryScript %v", err))

		return err
	}

	if resp["status"].(float64) != 200 {
		return fmt.Errorf("an error occurred")
	}

	return nil
}

// RunScriptOnUpdate run script on item update
func (s *Script) RunScriptOnUpdate() error {

	req := map[string]interface{}{
		"script": s.Script,
		"data":   s.Data,
	}

	var resp map[string]interface{}
	if err := externalSvc.GenericHTTPRequesterWResp("POST", "http", actionScriptHost(), "/run", req, &resp); err != nil {
		utils.Error(fmt.Sprintf("error TryScript %v", err))

		// skip error handler as actionScriptHost may not exists
		// TODO -JP
		return nil
	}

	if resp["status"].(float64) != 200 {
		return fmt.Errorf("an error occurred")
	}

	return nil
}

func actionScriptHost() string {
	host := "localhost:3000"
	switch os.Getenv("ENV") {
	case "dev", "stg", "prod":
		host = "beee-actionscript:3000"
	}

	return host
}
