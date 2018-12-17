package externalSvc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	utils "github.com/b-eee/amagi"
)

type (
	// EnvVar env var shortcut
	EnvVar struct {
		Name, Default string
	}
)

var (
	// EnvConfigctlHost env name for configctl host
	EnvConfigctlHost = EnvVar{"BEEE_CONFIGCTL_SERVICE_HOST", "localhost"}
	// EnvConfigctlPort env name for configctl port
	EnvConfigctlPort = EnvVar{"BEEE_CONFIGCTL_SERVICE_PORT", "8083"}

	// EnvApicoreHost env name for apicore host
	EnvApicoreHost = EnvVar{"BEEE_APICORE_SERVICE_HOST", "localhost"}
	// EnvApicorePort env name for apicore port
	EnvApicorePort = EnvVar{"BEEE_APICORE_SERVICE_PORT", "9000"}

	// EnvActionscriptHost env name for actionscript host
	EnvActionscriptHost = EnvVar{"BEEE_ACTIONSCRIPT_SERVICE_HOST", "localhost"}
	// EnvActionscriptPort env name for actionscript port
	EnvActionscriptPort = EnvVar{"BEEE_ACTIONSCRIPT_SERVICE_PORT", "3000"}

	// EnvImporterHost env name for importer host
	EnvImporterHost = EnvVar{"BEEE_IMPORTER_SERVICE_HOST", "localhost"}
	// EnvImporterPort env name for importer port
	EnvImporterPort = EnvVar{"BEEE_IMPORTER_SERVICE_PORT", "8080"}

	// EnvLinkerAPIHost env name for linker-api host
	EnvLinkerAPIHost = EnvVar{"LINKER_API_SERVICE_HOST", "localhost"}
	// EnvLinkerAPIPort env name for linker-api port
	EnvLinkerAPIPort = EnvVar{"LINKER_API_SERVICE_PORT", "7575"}

	// EnvJoblinkerHost env name for joblinker host
	EnvJoblinkerHost = EnvVar{"BEEE_JOBLINKER_SERVICE_HOST", "localhost"}
	// EnvJoblinkerPort env name for joblinker port
	EnvJoblinkerPort = EnvVar{"BEEE_JOBLINKER_SERVICE_PORT", "9010"}

	// EnvNotificatorHost env name for configctl host
	EnvNotificatorHost = EnvVar{"BEEE_NOTIFICATOR_SERVICE_HOST", "localhost"}
	// EnvNotificatorPort env name for configctl port
	EnvNotificatorPort = EnvVar{"BEEE_NOTIFICATOR_SERVICE_PORT", "8081"}

	// EnvLogStockerHost env name for logstocker host
	EnvLogStockerHost = EnvVar{"BEEE_LOGSTOCKER_SERVICE_HOST", "localhost"}
	// EnvLogStockerPort env name for logstocker port
	EnvLogStockerPort = EnvVar{"BEEE_LOGSTOCKER_SERVICE_PORT", "2004"}
)

// String return env value or default
func (v EnvVar) String() string {
	if e := os.Getenv(v.Name); e != "" {
		return e
	}
	return v.Default
}

// GenericHTTPRequester common HTTP requester utils
// TODO DEPRECATE THIS, USE BELOW INSTEAD JP
func GenericHTTPRequester(method, scheme, host, url string, data interface{}) (string, error) {
	s := time.Now()
	hostURL := fmt.Sprintf("%v://%v%v", scheme, host, url)
	mJSON, _ := json.Marshal(data)
	contentReader := bytes.NewReader(mJSON)

	req, _ := http.NewRequest(method, hostURL, contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == 404 {
		utils.Error(fmt.Sprintf("error Creating post to MS %v status: %v", err, resp))
		if resp != nil {
			return resp.Status, err
		}
		return "500", err
	}

	defer resp.Body.Close()

	utils.Info(fmt.Sprintf("GenericHTTPRequester took: %v to: %v", time.Since(s), hostURL))
	return resp.Status, nil
}

// GenericHTTPRequesterWResp generic http requester w/ response return
func GenericHTTPRequesterWResp(method, scheme, host, url string, data interface{}, responseData interface{}) error {
	s := time.Now()
	hostURL := fmt.Sprintf("%v://%v%v", scheme, host, url)
	mJSON, _ := json.Marshal(data)
	contentReader := bytes.NewReader(mJSON)

	req, _ := http.NewRequest(method, hostURL, contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == 404 {
		utils.Error(fmt.Sprintf("error Creating post to MS %v status: %v", err, resp))
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		utils.Error(fmt.Sprintf("error GenericHTTPRequesterWResp decode %v", err))
		return err
	}

	defer resp.Body.Close()

	utils.Info(fmt.Sprintf("GenericHTTPRequester took: %v to: %v", time.Since(s), hostURL))
	return nil
}

// HTTPGetRequest HTTP get Request generic function
func HTTPGetRequest(url string, values url.Values) (*http.Response, error) {
	reqURL := fmt.Sprintf("%v?%v", url, values.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Add("Authorization", "")

	resp, err := client.Do(req)
	if err != nil {
		utils.Error(fmt.Sprintf("error HTTPGetRequest %v", err))
		return resp, err
	}
	return resp, nil
}

// HTTPGetRequestWResponse http generic GETTER with response encoder
func HTTPGetRequestWResponse(url string, values url.Values, result interface{}) error {
	res, err := HTTPGetRequest(url, values)
	if err != nil {
		utils.Error(fmt.Sprintf("error HTTPGetRequestWResponse url=%v %v", url, err))
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		utils.Error(fmt.Sprintf("error HTTPGetRequestWResponse on decode url=%v %v", url, err))
	}

	return nil
}

// APIrequestGetter api configctl credentials getter
func APIrequestGetter(credKey, field string, response interface{}) error {
	v := url.Values{}
	v.Add("credential_key", credKey)
	v.Add("field", field)

	configURL := fmt.Sprintf("http://%v/get_kv/credential/%v", configCtlURL(), os.Getenv("ENV"))

	if err := HTTPGetRequestWResponse(configURL, v, &response); err != nil {
		return err
	}

	return nil
}

func configCtlURL() string {
	return fmt.Sprintf("%s:%s", EnvConfigctlHost, EnvConfigctlPort)
}
