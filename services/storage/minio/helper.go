package minio

import (
	"fmt"

	utils "github.com/b-eee/amagi"

	"github.com/minio/minio-go"
)

var (
	publicPolicyTpl = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Sid": "",
				"Effect": "Allow",
				"Principal": "*",
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/%s/*"]
			}
		]
	}`
)

func (svc *Service) publicPolicy() string {
	return fmt.Sprintf(publicPolicyTpl, svc.BucketName, svc.PublicPath)
}

// CreateBucket create bucket with name, returns right away if exists
func createBucket(client *minio.Client, bucketName, region string) (found bool, err error) {
	found, err = client.BucketExists(bucketName)
	if err != nil {
		utils.Info(fmt.Sprintf("Checking exist bucket '%s': %v\nCreating Bucket now", bucketName, err))
	}
	if !found {
		// create the bucket if not found
		err = client.MakeBucket(bucketName, region)
		if err != nil {
			err = fmt.Errorf("Error creating bucket '%s': %v", bucketName, err)
		}
	}
	return
}
