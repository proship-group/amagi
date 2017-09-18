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
	var configURL string
	switch os.Getenv("ENV") {
	case "local":
		configURL = "localhost:8083"
	default:
		configURL = "beee-configctl:8083"
	}

	return configURL
}
