package backends

import (
	"fmt"
	"os"

	"github.com/b-eee/amagi/services/configctl"

	utils "github.com/b-eee/amagi"
	minio "github.com/minio/minio-go"
)

func getMIOCredentials() configctl.Environment {
	env := configctl.GetDBCfgStngWEnvName("minio", os.Getenv("ENV"))

	return env
}

// MIOCreateClient create client for minio
func MIOCreateClient() (*minio.Client, error) {
	env := getMIOCredentials()
	fmt.Println(env, "client ============================")

	client, err := minio.New(env.Host, env.MIOAccessKeyID, env.MIOSecretAccessKey, false)
	if err != nil {
		utils.Error(fmt.Sprintf("error MIOCreateClient %v", err))
		return client, err
	}

	return client, nil
}

// MIOPutObject put object to minio with io.Reader
func MIOPutObject(fo FileObject) error {
	client, _ := MIOCreateClient()

	res, err := client.PutObject(fo.BucketName, fo.ObjectName, fo.File, "application/octet-stream")
	if err != nil {
		utils.Error(fmt.Sprintf("error MIOPutObject %v", err))
		return err
	}

	fmt.Println(res)

	return nil
}

// MIOGetObject minio get object
func MIOGetObject(fo FileObject) (interface{}, error) {
	client, err := MIOCreateClient()
	if err != nil {
		utils.Error(fmt.Sprintf("error MIOGetObject %v", err))
		return nil, err
	}
	fmt.Println(fo.BucketName, fo.ObjectName)
	obj, err := client.GetObject(fo.BucketName, fo.ObjectName)
	if err != nil {
		utils.Error(fmt.Sprintf("error MIOGetObject %v", err))
		return nil, err
	}

	return obj, nil
}
