package backends

import (
	"io"
	"os"
)

type (
	// FileObject file object request interface
	FileObject struct {
		BucketName  string
		ObjectName  string
		FilePath    string
		ContentType string

		File io.Reader
	}
)

var (
	// FileStorageFlagENV file storage flag to use backend
	FileStorageFlagENV = "FILE_STORAGE_FLAG"
)

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
