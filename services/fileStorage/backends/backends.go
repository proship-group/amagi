package backends

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go"

	storage "google.golang.org/api/storage/v1"
)

type (
	// FileObject file object request interface
	FileObject struct {
		BucketName  string
		ObjectName  string
		FilePath    string
		ContentType string

		File io.Reader

		GCPObjectFile multipart.File
		GCPObject     *storage.Object
		GCPAcls       []*storage.ObjectAccessControl
	}

	// BackendConfigs backend configurations for multiple storage clients/servers
	BackendConfigs struct {
		BackendName string
		Minio       *minio.Client
	}
)

var (
	// FileStorageFlagENV file storage flag to use backend
	FileStorageFlagENV = "FILE_STORAGE_FLAG"

	// CurrentBackendStorage initalized backend storage setting
	CurrentBackendStorage *BackendConfigs
)

// CreateBackend create storage backend
func CreateBackend(name string) {
	switch name {
	case "minio":
		client, err := MIOCreateClient()
		if err != nil {
			panic("can't create minio backend..")
		}
		CurrentBackendStorage = &BackendConfigs{
			Minio: client,
		}
	case "gcp":
	}
}

// GetBackendFromENV get backend name from config.env or environment variable
func GetBackendFromENV() string {
	return os.Getenv(FileStorageFlagENV)
}

// PutObject put object or upload object to file storage with io.Reader
func PutObject(fo FileObject) (interface{}, error) {
	var err error
	var resp interface{}

	switch os.Getenv(FileStorageFlagENV) {
	case "minio":
		resp, err = MIOPutObject(fo)
	}

	return resp, err
}

// GetObject get object
func GetObject(fo FileObject) (interface{}, error) {
	var err error
	var obj interface{}

	switch os.Getenv(FileStorageFlagENV) {
	case "minio":
		obj, err = MIOGetObject(fo)
	}

	return obj, err
}
