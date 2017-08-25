package fileStorage

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// go test -run=TestFPutObject ./services/fileStorage -v
func TestFPutObject(t *testing.T) {
	os.Setenv("FILE_STORAGE_FLAG", "minio")
	os.Setenv("ENV", "local")

	file, _ := os.Open(fmt.Sprintf("%v/src/github.com/b-eee/amagi/services/fileStorage/test.txt", os.Getenv("GOPATH")))

	fi := strings.Split(file.Name(), "/")
	f := File{
		File:       file,
		ObjectName: fmt.Sprintf("%s_%v", time.Now(), fi[len(fi)-1]),
		BucketName: "test",
	}
	if _, err := f.PutObject(); err != nil {
		t.Error(err)
	}
}
