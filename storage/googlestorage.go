package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	utils "github.com/b-eee/amagi"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

const (
	storageRootURI string = "https://storage.googleapis.com"
)

type googleStorage struct {
	storageService *storage.Service
}

func newGoogleStorage() *googleStorage {
	s := &googleStorage{}
	var err error
	s.storageService, err = storage.New(jwtConfigFromJSON())
	if err != nil {
		utils.Fatal(fmt.Sprintf("error newGoogleStorage %v", err))
		panic(err)
	}

	return s
}

// Save save in cloud storage
func (s *googleStorage) Save(r io.Reader, bucketName, objectName string) (file *CloudStorageFile, err error) {
	aclObject := []*storage.ObjectAccessControl{}
	aclObject = append(aclObject, &storage.ObjectAccessControl{
		Bucket: bucketName,
		Role:   "OWNER",
		Entity: "allAuthenticatedUsers",
		Object: objectName})
	storeData := &storage.Object{Name: objectName, Acl: aclObject}

	res, err := s.storageService.Objects.Insert(bucketName, storeData).Media(r).Do()
	if err != nil {
		utils.Error(fmt.Sprintf("Objects.Insert failed: %v", err))
		return
	}

	uri := fmt.Sprintf("%s/%s/%s", storageRootURI, bucketName, objectName)
	file = &CloudStorageFile{
		FilePath:  res.Name,
		SelfLink:  res.SelfLink,
		MediaLink: res.MediaLink,
		URI:       uri,
		Size:      res.Size,
	}

	return
}

// jwtConfigFromJSON jwt config from json
func jwtConfigFromJSON() *http.Client {
	jwtConfigFilePath := os.Getenv("JWTCONFIG_PATH")
	data, err := ioutil.ReadFile(jwtConfigFilePath)
	if err != nil {
		utils.Fatal(fmt.Sprintf("error jwtConfigFromJSON %v path:%v", err, jwtConfigFilePath))
		panic(err)
	}
	conf, err := google.JWTConfigFromJSON(data, storage.DevstorageFullControlScope)
	if err != nil {
		utils.Fatal(fmt.Sprintf("error jwtConfigFromJSON in google.JWTConfigFromJSON %v", err))
		panic(err)
	}

	client := conf.Client(oauth2.NoContext)
	client.Get("...")
	return client
}
