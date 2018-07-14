package gcs

import (
	gcs "google.golang.org/api/storage/v1"
)

var (
	timeLayout = "2006-01-02T15:04:05.000Z"
)

func publicReadACLTempl() *gcs.ObjectAccessControl {
	return &gcs.ObjectAccessControl{
		Role:   "READER",
		Entity: "allUsers",
	}
}

func (svc *Service) publicACL(objectName string) *gcs.ObjectAccessControl {
	acl := publicReadACLTempl()
	acl.Bucket = svc.BucketName
	acl.Object = objectName
	return acl
}
