package configctl

import (
	"fmt"
	"net/url"
	"os"

	"github.com/b-eee/amagi/services/externalSvc"
)

func configCtlURL() string {
	return fmt.Sprintf("%s:%s", externalSvc.EnvConfigctlHost, externalSvc.EnvConfigctlPort)
}

// APIrequestGetter api configctl credentials getter
// TODO deprecate moved to amagi -JP
func APIrequestGetter(credKey, field string) (map[string]string, error) {
	v := url.Values{}
	v.Add("credential_key", credKey)
	v.Add("field", field)

	configURL := fmt.Sprintf("http://%v/get_kv/credential/%v", configCtlURL(), os.Getenv("ENV"))

	var resp map[string]string
	if err := externalSvc.HTTPGetRequestWResponse(configURL, v, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
