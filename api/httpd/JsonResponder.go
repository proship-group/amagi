package httpd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/unrolled/render"

	utils "github.com/b-eee/amagi"
)

var (
	generatorItemLimit = 6
)

// GenericResponse a generic response from http
type GenericResponse struct {
	Status   int
	Response interface{}
}

// respondToJson write http json resposne
func respondToJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	// start := time.Now()
	r := render.New()
	if err := r.JSON(w, http.StatusOK, data); err != nil {
		fmt.Printf("error respondToJSON %v", err)
		return err
	}
	return nil
	// utils.Info(fmt.Sprintf("respondToJSON took: %v", time.Since(start)))
}

// RespondToJSON public json response writer
func RespondToJSON(w http.ResponseWriter, data interface{}) {
	respondToJSON(w, data)
}

// RespondNilError generic nil error response
func RespondNilError(w http.ResponseWriter) {
	respondToJSON(w, map[string]interface{}{
		"error": nil,
	})
}

// AllowCrossOrigin allow cross origin setup
func AllowCrossOrigin(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// DecodePostRequest decode post request from http
func DecodePostRequest(body io.Reader, target interface{}) error {

	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&target); err != nil {
		return err
	}

	return nil
}

// HTTPError http error
func HTTPError(w http.ResponseWriter, err error) {
	utils.Error(fmt.Sprintf("error %v", err))

	raven.CaptureError(err, nil)
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// ImporterAPIPrefix importer api prefixer
func ImporterAPIPrefix(apiname string) string {
	return fmt.Sprintf("/api/v1/%v", apiname)
}

// ServeContent serve content file for download
func ServeContent(w http.ResponseWriter, r *http.Request, file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		HTTPError(w, err)
		return
	}
	defer os.Remove(file)

	setResponseHeader(w, file)

	http.ServeContent(w, r, fmt.Sprintf(file), time.Now(), bytes.NewReader(data))
}

func setResponseHeader(w http.ResponseWriter, fileName string) error {
	w.Header().Set("Content-Type", fmt.Sprintf("application/octet-stream"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%v.zip\"", fileName))
	w.Header().Set("Content-Transfer-Encoding", "binary")

	return nil
}
