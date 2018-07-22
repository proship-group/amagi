package actionScript

import (
	"fmt"
	"strings"

	"github.com/b-eee/amagi/services/externalSvc"

	utils "github.com/b-eee/amagi"
)

type (
	// Script script interface
	Script struct {
		PID        string                   `json:"p_id"`
		AID        string                   `json:"a_id"`
		Script     string                   `json:"script"`
		Data       interface{}              `json:"data"`
		Token      string                   `json:"auth_token"`
		ScriptVars []ScriptVariableSettings `json:"script_vars"`
	}

	// ScriptVariableSettings script variable settings
	ScriptVariableSettings struct {
		VarName     string `bson:"var_name" json:"var_name"`
		Description string `bson:"desc" json:"desc"`
		Value       string `bson:"value" json:"value"`
		Enabled     bool   `bson:"enabled" json:"enabled"`
	}
)

// CommonScripts common functions that can use in actionScript
var CommonScripts = `
// common function for linker-api
function callAPI(method, url, params, callback){
	const targetURL = 'http://{HEXA_API_SERVER}/' + url;
	const httpParams = {
		url: targetURL,
		method: method,
		headers: {
			'Authorization': "{HEXA_API_TOKEN}"
		},
		maxRedirects: 1000,
	}
	httpSvc(httpParams, params, function(res) {
		callback(res);
	})
}
`

// ExecuteScript execute action script
func (s *Script) ExecuteScript() error {

	// append common functions into script
	s.Script = fmt.Sprintf("%v%v", s.Script, CommonScripts)

	// replace macro variables
	envVars := map[string]string{
		"{HEXA_API_TOKEN}":  fmt.Sprintf("Bearer %v", s.Token),
		"{HEXA_API_SERVER}": LinkerAPIHost(),
	}
	for _, repVar := range s.ScriptVars {
		if repVar.Enabled {
			varName := fmt.Sprintf("{%v}", repVar.VarName)
			// add replace string if not exists (only use the value that first appeared)
			if _, exists := envVars[varName]; !exists {
				envVars[varName] = repVar.Value
			}
		}
	}
	// utils.Pretty(envVars, "envVars")
	s.ReplaceEnvVars(envVars)
	// utils.Pretty(s.Script, "s.Script")

	req := map[string]interface{}{
		"script": s.Script,
		"data":   s.Data,
	}
	var resp map[string]interface{}
	if err := externalSvc.GenericHTTPRequesterWResp("POST", "http", Host(), "/run", req, &resp); err != nil {
		utils.Error(fmt.Sprintf("error ExecuteScript %v", err))
		return err
	}

	if resp["status"].(float64) != 200 {
		return fmt.Errorf("an error occurred")
	}

	return nil
}

// RunScriptOnUpdate run script on item update[DREPRECATE: use ExecuteScript instead]
func (s *Script) RunScriptOnUpdate() error {
	req := map[string]interface{}{
		"script": s.Script,
		"data":   s.Data,
	}

	var resp map[string]interface{}
	if err := externalSvc.GenericHTTPRequesterWResp("POST", "http", Host(), "/run", req, &resp); err != nil {
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

// ReplaceEnvVars replace env variables values
func (s *Script) ReplaceEnvVars(envVars map[string]string) error {

	for k, v := range envVars {
		// re := regexp.MustCompile(k)
		// s.Script = re.ReplaceAllString(s.Script, v)

		// replace all
		s.Script = strings.Replace(s.Script, k, v, -1)
	}

	return nil
}

// getUserAPIToken get user api token from sql
func getUserAPIToken() error {

	return nil
}

// Host return action script host address
func Host() string {
	return fmt.Sprintf("%s:%s", externalSvc.EnvActionscriptHost, externalSvc.EnvActionscriptPort)
}

// LinkerAPIHost linker api hostname or url
func LinkerAPIHost() string {
	return fmt.Sprintf("%s:%s", externalSvc.EnvLinkerAPIHost, externalSvc.EnvLinkerAPIPort)
}
