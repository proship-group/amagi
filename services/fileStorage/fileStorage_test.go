package fileStorage

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// go test -run=TestFPutObject ./services/fileStorage -v
func TestFPutObject(t *testing.T) {
	os.Setenv("FILE_STORAGE_FLAG", "minio")
	os.Setenv("ENV", "local")

	file, _ := os.Open(fmt.Sprintf("%v/src/github.com/b-eee/amagi/services/mailSender/templates/confirm.html", os.Getenv("GOPATH")))

	fi := strings.Split(file.Name(), "/")
	f := File{
		File:       file,
		ObjectName: fi[len(fi)-1],
		BucketName: "test",
	}
	if err := f.PutObject(); err != nil {
		t.Error(err)
	}
}
