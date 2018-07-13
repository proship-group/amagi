package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	utils "github.com/b-eee/amagi"
)

type (
	// Service service interface
	Service interface {
		// Initialize execute needed initialization
		Initialize() error
		// CreateObject create object
		CreateObject(objectName string, file io.Reader, contentType string) (*ObjectInfo, error)
		// SaveObject saves the data with a unique random name
		SaveObject(file io.Reader) (*ObjectInfo, error)
		// CreatePublicObject creates the object in the public storage/folder
		CreatePublicObject(objectName string, file io.Reader, contentType string) (*ObjectInfo, error)
		// SavePublicObject save the object in public storage
		SavePublicObject(file io.Reader) (*ObjectInfo, error)
		// DownloadObjectDest save object to file system
		DownloadObjectDest(objectName, destFilename string) (localPath string, err error)
		// DownloadObject save object to file system with a uniques name
		DownloadObject(objectName string) (localPath string, err error)
		// DeleteObject base function for deleting objects from storage
		DeleteObject(objectName string) error
		// GetObject get object from storage
		GetObject(objectName string) (io.ReadCloser, error)
		// GetObjectInfo returns object information
		GetObjectInfo(objectName string) (*ObjectInfo, error)

		// getters for common config
		// GetBucketName returns the configured bucketname
		GetBucketName() string
		// GetPublicPath returns the configured PublicPath
		GetPublicPath() string
		// GetRegion returns the configured region
		GetRegion() string
		// GetEndpoint returns the configured endpoint
		GetEndpoint() string
	}

	// FileObject file with objectinfo
	FileObject struct {
		File       io.ReadCloser // the object reader
		ObjectInfo               // object info
	}

	// ObjectInfo object details
	ObjectInfo struct {
		Name         string    // the object name
		SelfLink     string    // link to object for client consumption
		MediaLink    string    // link to object relative to storage service
		ContentType  string    // mime-type
		Size         uint64    // bytes written
		ETag         string    // etag meta data
		LastModified time.Time // last datetime modified
	}
)

var (
	// ErrNotImplemented error for not implemented interface api
	ErrNotImplemented = fmt.Errorf("The API is not yet implemented")
	// RandObjNameLength random object name length
	RandObjNameLength = 256
	tmpPathName       = "hexalink"
)

// NewTempFile prepares a string for temp file
func NewTempFile(name string, prefixDir ...string) string {
	folder := fmt.Sprintf("%s%s", os.TempDir(), tmpPathName)
	for _, prefix := range prefixDir {
		folder = fmt.Sprintf("%s/%s", folder, prefix)
	}
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		utils.Error(fmt.Sprintf("cant mkdir to %v", err))
		return ""
	}
	return filepath.FromSlash(fmt.Sprintf("%s/%s", folder, name))
}
