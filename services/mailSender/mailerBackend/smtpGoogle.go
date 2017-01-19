package mailerBackend

import (
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"bytes"

	utils "github.com/b-eee/amagi"
	"github.com/b-eee/amagi/services/fileStorage"
)

type (
	// GTemplateData google email template data
	GTemplateData struct {
		ID   string
		Name string
		URL  string
	}
)

// SMTPGoogleSendEmail smpt google send email backend
func SMTPGoogleSendEmail(r Request) error {
	body, err := ParseTemplate(r.TemplateData, r.SGTemplateID)
	if err != nil {
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "testing" + "!\n"
	msg := []byte(subject + mime + "\n" + body)
	addr := "smtp.gmail.com:587"
	auth := smtp.PlainAuth("", "demo@b-eee.com", "demo@beee", "smtp.gmail.com")

	if err := smtp.SendMail(addr, auth, r.Sender, r.Receivers, msg); err != nil {
		return err
	}
	return nil
}

// ParseTemplate parse email template
func ParseTemplate(templateData map[string]interface{}, SGTemplateID string) (string, error) {
	storage := fileStorage.File{
		BucketName: "test",
		ObjectName: fmt.Sprintf("%v.html", SGTemplateID),
	}

	object, err := storage.GetObject()
	if err != nil {
		return "", err
	}

	filePath := templateFiles(storage.ObjectName)
	defer os.Remove(filePath)

	if err := fileStorage.MIOExtractAndStoreObject(object, filePath); err != nil {
		return "", err
	}

	t, err := template.ParseFiles(filePath)
	if err != nil {
		utils.Error(fmt.Sprintf("error ParseTemplate %v", err))
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, templateData); err != nil {
		utils.Error(fmt.Sprintf("error ParseTemplate %v", err))
		return "", err
	}

	return buf.String(), nil
}

func templateFiles(filename string) string {
	return fmt.Sprintf("%v/src/github.com/b-eee/amagi/services/mailSender/templates/%v", os.Getenv("GOPATH"), filename)
}
