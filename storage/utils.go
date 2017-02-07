package storage

import (
	"bytes"
	"mime/multipart"
)

// EncodeFileToBase64 encode file(io.Reader) to base64
func EncodeFileToBase64(file multipart.File) error {
	b := new(bytes.Buffer)

	b.ReadFrom(file)
	b.Bytes()

	return nil
}
