package client

import (
	"fmt"
	"os"
	"time"

	utils "amagi"
	"amagi/services/configctl"
	"amagi/services/storage"
	"amagi/services/storage/gcs"
	"amagi/services/storage/minio"
)

const (
	storageCredentialKey = "storage"

	// IMPORTANT: Required env vars

	// BucketNameEnv env name for the bucket name config
	BucketNameEnv = "STORAGE_BUCKET"
	// PublicNameEnv env name for the public path config
	PublicNameEnv = "STORAGE_PUBLIC_PATH"
	// ProjectIDEnv env name for project id config
	ProjectIDEnv = "STORAGE_PROJECT_ID"
	// ServiceAccountEnv env name for service account config
	ServiceAccountEnv = "STORAGE_SERVICE_ACCOUNT"
)

var (
	// Client the storage service client
	Client storage.Service
)

// InitStorageClient initialize the storage client config
func InitStorageClient() {
	s := time.Now()
	if Client != nil {
		return
	}
	cfg, err := configctl.APIrequestGetter(storageCredentialKey, "")
	if err != nil {
		panic(fmt.Errorf("Error initializing storage service: %v", err))
	}
	switch cfg["storageService"] {
	case "gcs":
		Client = &gcs.Service{
			Endpoint:      cfg["endpoint"],
			Project:       os.Getenv(ProjectIDEnv),
			ServieAccount: os.Getenv(ServiceAccountEnv),
			Region:        cfg["region"],
			BucketName:    os.Getenv(BucketNameEnv),
			PublicPath:    os.Getenv(PublicNameEnv),
		}
	case "s3", "minio":
		Client = &minio.Service{
			Endpoint:   cfg["endpoint"],
			AccessID:   cfg["accessID"],
			SecretKey:  cfg["secretKey"],
			UseSSL:     cfg["useSSL"] == "true",
			Region:     cfg["region"],
			BucketName: os.Getenv(BucketNameEnv),
			PublicPath: os.Getenv(PublicNameEnv),
		}
	default:
		panic(fmt.Errorf("Invalid storage service: %s", cfg["storageService"]))
	}
	utils.Info(fmt.Sprintf("InitStorageClient %v", Client))
	if err := Client.Initialize(); err != nil {
		panic(err)
	}
	utils.Info(fmt.Sprintf("InitStorageClient took: %v", time.Since(s)))
}
