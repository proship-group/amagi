package storage

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
)

// EncodeFileToBase64 encode file(io.Reader) to base64
func EncodeFileToBase64(file multipart.File) string {
	b := new(bytes.Buffer)
	b.ReadFrom(file)

	// TODO better error handler! -JP
	if len(b.Bytes()) == 0 {
		panic("empty file!")
	}

	return base64.StdEncoding.EncodeToString(b.Bytes())
}
