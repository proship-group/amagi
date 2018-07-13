package gcs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	utils "github.com/b-eee/amagi"
	"github.com/b-eee/amagi/helpers"
	"github.com/b-eee/amagi/services/storage"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	gcs "google.golang.org/api/storage/v1"
)

type (
	// Service gcs service struct
	Service struct {
		Endpoint      string
		Project       string
		ServieAccount string
		Region        string
		BucketName    string
		PublicPath    string

		// helper variables / cache variables
		cacheVars struct {
			client *gcs.Service
		}
	}
)

var (
	// GCPKeyFileNAme key file name for api key access
	GCPKeyFileNAme = "GCP_KEYFILE_NAME"
)

// GetBucketName return bucket name
func (svc *Service) GetBucketName() string {
	return svc.BucketName
}

// GetPublicPath return GetPublicPath
func (svc *Service) GetPublicPath() string {
	return svc.PublicPath
}

// GetRegion return GetRegion
func (svc *Service) GetRegion() string {
	return svc.Region
}

// GetEndpoint return GetEndpoint
func (svc *Service) GetEndpoint() string {
	return svc.Endpoint
}

// Initialize initialize service
func (svc *Service) Initialize() error {
	// check if client has already been created
	if svc.cacheVars.client != nil {
		return nil
	}
	// if not yet created...
	// check if all fields are existing
	if ok := (svc.ServieAccount != "" &&
		svc.Project != "" &&
		svc.Region != "" &&
		svc.BucketName != "" &&
		svc.PublicPath != ""); !ok {
		return fmt.Errorf("Invalid service struct: %v", svc)
	}
	client, err := gcs.New(ServiceAccountClient(svc.ServieAccount))
	if err != nil {
		utils.Error(fmt.Sprintf("Unable to create storage service: %v", err))
		return err
	}
	svc.cacheVars.client = client

	// Ensure bucket
	if _, err := client.Buckets.Insert(svc.Project, &gcs.Bucket{Name: svc.BucketName}).Do(); err != nil {
		if !googleapi.IsNotModified(err) && err.(*googleapi.Error).Code != 409 {
			return err
		}
	}

	return nil
}

// NewObject base method in creating objects
func (svc *Service) NewObject(bucket, objectName string, file io.Reader, contentType string, acls ...*gcs.ObjectAccessControl) (*storage.ObjectInfo, error) {
	res, err := svc.cacheVars.client.Objects.Insert(
		bucket,
		&gcs.Object{
			Name: objectName,
			Acl:  acls,
		},
	).Media(file).Do()
	if err != nil {
		utils.Error(fmt.Sprintf("Objects.Insert failed: %v", err))
		return nil, err
	}
	updated, err := time.Parse(timeLayout, res.Updated)
	if err != nil {
		updated = time.Now()
	}

	mediaLink := res.SelfLink
	mediaLinkURL, err := url.Parse(mediaLink)
	if err != nil {
		utils.Error(fmt.Sprintf("error at parse url %v", err))
	}
	mediaLink = fmt.Sprintf("/storage/%s", mediaLinkURL.Path[1:])
	return &storage.ObjectInfo{
		Name:         res.Name,
		SelfLink:     res.SelfLink,
		MediaLink:    mediaLink,
		ContentType:  res.ContentType,
		Size:         res.Size,
		ETag:         res.Etag,
		LastModified: updated,
	}, nil
}

// CreateObject create object
func (svc *Service) CreateObject(objectName string, file io.Reader, contentType string) (*storage.ObjectInfo, error) {
	return svc.NewObject(svc.BucketName, objectName, file, contentType)
}

// SaveObject saves the data with a unique random name
func (svc *Service) SaveObject(file io.Reader) (*storage.ObjectInfo, error) {
	objectName := helpers.RandString6(storage.RandObjNameLength)
	return svc.NewObject(svc.BucketName, objectName, file, "application/octet-stream")
}

// CreatePublicObject creates the object in the public storage/folder
func (svc *Service) CreatePublicObject(objectName string, file io.Reader, contentType string) (*storage.ObjectInfo, error) {
	objectName = fmt.Sprintf("%s/%s", svc.PublicPath, objectName)
	return svc.NewObject(
		svc.BucketName,
		objectName,
		file, contentType,
		svc.publicACL(objectName))
}

// SavePublicObject save the object in public storage
func (svc *Service) SavePublicObject(file io.Reader) (*storage.ObjectInfo, error) {
	objectName := fmt.Sprintf("%s/%s", svc.PublicPath, helpers.RandString6(storage.RandObjNameLength))
	return svc.NewObject(
		svc.BucketName,
		objectName,
		file,
		"application/octet-stream",
		svc.publicACL(objectName))
}

// DownloadObjectDest save object to file system
func (svc *Service) DownloadObjectDest(objectName, destFilename string) (localPath string, err error) {
	localPath = storage.NewTempFile(destFilename)
	output, err := os.Create(localPath)
	if err != nil {
		utils.Error(fmt.Sprintf("Error SavePublicObjectAsFile While Creating file %v error=%v", localPath, err))
		return "", err
	}

	file, err := svc.GetObject(objectName)
	if err != nil {
		utils.Warn(fmt.Sprintf("error SavePublicObjectAsFile url=%v err=%v", objectName, err))
		return "", err
	}
	defer file.Close()
	defer output.Close()

	n, err := io.Copy(output, file)
	utils.Info(fmt.Sprintf("%v bytes downloaded for %v ", n, objectName))
	return localPath, nil
}

// DownloadObject save object to file system with a uniques name
func (svc *Service) DownloadObject(objectName string) (localPath string, err error) {
	destFilename := helpers.RandString6(storage.RandObjNameLength)
	return svc.DownloadObjectDest(objectName, destFilename)
}

// GetObject get object from storage
func (svc *Service) GetObject(objectName string) (io.ReadCloser, error) {
	resp, err := svc.cacheVars.client.Objects.Get(svc.BucketName, objectName).Download()
	if err != nil {
		utils.Warn(fmt.Sprintf("error GetObject url=%v err=%v", objectName, err))
		return nil, err
	}
	return resp.Body, nil
}

// DeleteObject base function for deleting objects from storage
func (svc *Service) DeleteObject(objectName string) error {
	if err := svc.cacheVars.client.Objects.Delete(svc.BucketName, objectName).Do(); err != nil {
		utils.Warn(fmt.Sprintf("error DeleteObject: %v[%v/%v]", err, svc.BucketName, objectName))
		return err
	}
	return nil
}

// GetObjectInfo get object info
func (svc *Service) GetObjectInfo(objectName string) (*storage.ObjectInfo, error) {
	obj, err := svc.cacheVars.client.Objects.Get(svc.BucketName, objectName).Do()
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(timeLayout, obj.Updated)
	if err != nil {
		updated = time.Now()
	}
	mediaLink := obj.MediaLink
	mediaLinkURL, err := url.Parse(mediaLink)
	if err != nil {
		utils.Error(fmt.Sprintf("error at parse url %v", err))
		mediaLink = fmt.Sprintf("/storage/%s", mediaLinkURL.Path[1:])
	}
	return &storage.ObjectInfo{
		Name:         obj.Name,
		SelfLink:     obj.SelfLink,
		MediaLink:    mediaLink,
		ContentType:  obj.ContentType,
		Size:         obj.Size,
		ETag:         obj.Etag,
		LastModified: updated,
	}, nil
}

// ServiceAccountClient create http client with service account
func ServiceAccountClient(serviceAccountFilePath string) *http.Client {
	data, err := ioutil.ReadFile(serviceAccountFilePath)
	if err != nil {
		panic(err)
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/devstorage.full_control")
	if err != nil {
		panic(err)
	}
	client := conf.Client(oauth2.NoContext)
	client.Get("...")
	return client
}
