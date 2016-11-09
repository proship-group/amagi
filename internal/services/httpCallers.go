package services

import (
	"fmt"
	"net/http"
	"net/url"

	utils "amagi"
)

// HTTPGetRequest HTTP get Request generic function
func HTTPGetRequest(url string, values url.Values) (*http.Response, error) {

	reqURL := fmt.Sprintf("%v?%v", url, values.Encode())
	resp, err := http.Get(reqURL)
	if err != nil {
		utils.Error(fmt.Sprintf("error HTTPGetRequest %v", err))
		return resp, err
	}
	return resp, nil
}
