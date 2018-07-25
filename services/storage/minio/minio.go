package minio

import (
	"fmt"
	"io"
	"time"

	utils "github.com/b-eee/amagi"
	"github.com/b-eee/amagi/helpers"
	"github.com/b-eee/amagi/services/storage"

	"github.com/minio/minio-go"
)

type (
	// Service storage configuration
	Service struct {
		Endpoint   string
		AccessID   string
		SecretKey  string
		UseSSL     bool
		Region     string
		BucketName string
		PublicPath string

		// helper variables / cache variables
		cacheVars struct {
			client *minio.Client
		}
	}
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
	if ok := (svc.Endpoint != "" &&
		svc.AccessID != "" &&
		svc.SecretKey != "" &&
		svc.Region != "" &&
		svc.BucketName != "" &&
		svc.PublicPath != ""); !ok {
		return fmt.Errorf("Invalid service struct: %v", svc)
	}
	client, err := minio.NewWithRegion(
		svc.Endpoint,
		svc.AccessID,
		svc.SecretKey,
		svc.UseSSL,
		svc.Region,
	)
	if err != nil {
		return fmt.Errorf("Unable to create storage service: %v", err)
	}
	var found bool
	found, err = createBucket(client, svc.BucketName, svc.Region)
	// set policy to public read
	if !found && err == nil {
		err = client.SetBucketPolicy(
			svc.BucketName,
			svc.publicPolicy(),
		)
	}
	if err != nil {
		return fmt.Errorf("Failed to ensure public folder: %v", err)
	}
	svc.cacheVars.client = client
	return nil
}

// NewObject base method in creating objects
func (svc *Service) NewObject(bucket, objectName string, file io.Reader, contentType string) (*storage.ObjectInfo, error) {
	n, err := svc.cacheVars.client.PutObject(svc.BucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		utils.Error(fmt.Sprintf("minio.PutObject failed: %v", err))
		return nil, err
	}
	return &storage.ObjectInfo{
		Name:         objectName,
		MediaLink:    fmt.Sprintf("/storage/%s/%s", svc.BucketName, objectName),
		SelfLink:     fmt.Sprintf("%s/%s", svc.BucketName, objectName),
		ContentType:  contentType,
		Size:         uint64(n),
		ETag:         "",
		LastModified: time.Now(),
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
		file, contentType)
}

// SavePublicObject save the object in public storage
func (svc *Service) SavePublicObject(file io.Reader) (*storage.ObjectInfo, error) {
	return svc.CreatePublicObject(
		helpers.RandString6(storage.RandObjNameLength),
		file,
		"application/octet-stream",
	)
}

// DownloadObjectDest save the object in public storage with a unique random name
func (svc *Service) DownloadObjectDest(objectName, destFilename string) (string, error) {
	destFilename = storage.NewTempFile(destFilename)
	err := svc.cacheVars.client.FGetObject(
		svc.BucketName,
		objectName,
		destFilename,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return "", err
	}
	utils.Info(fmt.Sprintf("file downloaded: %s", destFilename))
	return destFilename, nil
}

// DownloadObject save object to file system with a uniques name
func (svc *Service) DownloadObject(objectName string) (localPath string, err error) {
	destFilename := helpers.RandString6(storage.RandObjNameLength)
	return svc.DownloadObjectDest(objectName, destFilename)
}

// GetObject get object from storage
func (svc *Service) GetObject(objectName string) (io.ReadCloser, error) {
	object, err := svc.cacheVars.client.GetObject(svc.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		utils.Error(fmt.Sprintf("Error GetObject '%s/%s': %v", svc.BucketName, objectName, err))
		return nil, err
	}
	return object, nil
}

// DeleteObject base function for deleting objects from storage
func (svc *Service) DeleteObject(objectName string) error {
	if err := svc.cacheVars.client.RemoveObject(svc.BucketName, objectName); err != nil {
		utils.Warn(fmt.Sprintf("error DeleteObject: %v[%v/%v]", err, svc.BucketName, objectName))
		return err
	}
	return nil
}

// GetObjectInfo get object info
func (svc *Service) GetObjectInfo(objectName string) (*storage.ObjectInfo, error) {
	obj, err := svc.cacheVars.client.StatObject(svc.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &storage.ObjectInfo{
		Name:         objectName,
		MediaLink:    fmt.Sprintf("/storage/%s/%s", svc.BucketName, objectName),
		SelfLink:     fmt.Sprintf("%s/%s", svc.BucketName, objectName),
		ContentType:  obj.ContentType,
		Size:         uint64(obj.Size),
		ETag:         obj.ETag,
		LastModified: obj.LastModified,
	}, nil
}
