package mailSender

import (
	"os"
	"testing"
)

func TestInitMailSender(t *testing.T) {
	os.Setenv("EMAIL_BACKEND_FLAG", "smtp.google")
	os.Setenv("FILE_STORAGE_FLAG", "minio")
	os.Setenv("ENV", "local")

	mailer := InitMailSender()

	if err := mailer.Send(); err != nil {
		t.Error(err)
	}
}
