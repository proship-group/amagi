package actionScript

import (
	"fmt"

	"github.com/b-eee/amagi/api/httpd"
	"github.com/gin-gonic/gin"

	utils "github.com/b-eee/amagi"
)

// TryActionScript try action script
func TryActionScript(c *gin.Context) {
	var s Script
	if err := httpd.DecodePostRequest(c.Request.Body, &s); err != nil {
		utils.Info(fmt.Sprintf("error TryActionScript DecodePostRequest %v", err))
		return
	}

	s.ReplaceEnvVars(map[string]string{
		"API_TOKEN": "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTQ1MjkwOTAsImlhdCI6MTUxMzkyNDI5MCwic3ViIjoiNThiNTIzODdjYzFjYjQ4NjMyMTAzNjRlIiwidW4iOiJqLnNvbGl2YV9fMzMzMjExMTIyMjIzMTExMjMyMTIyMiJ9.hFqIf_yttewREkt83Q6U-OvcTHh8zPjl1g7ARXmNVyGND4oIn9sPsycsfI3_E6iDWnSjkKTb9x9LD32nsQyLIWaDpbauXYJwttDHwuIR66deIRttYhjKlD7Yr45_AqiMnLwMw0yBJCmVmLVUb8ewMZ580xb2YZ9_4i4c5fRJhW8LAOjYOVPgmdlCuOupHA6_EbfwJgt_NMXA_OH93Woq3GthbvOAUu5JnuzSDOQbA68qmXE7vvlh13n-xp2uaWqYrHQ1r_Wpp63XwBsYTfSwO3j8PmxC0AX3FBYYtKVS6DelZXECtlv7i0S4lZjtFJdU7Y3po40oZ2_MlpBkZN0Wdw",
	})

	if err := s.TryScript(); err != nil {
		httpd.HTTPError(c.Writer, err)
		return
	}

	httpd.RespondNilError(c.Writer)
}
