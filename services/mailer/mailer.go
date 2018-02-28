package mailer

import (
	"fmt"
	"os"
	"time"

	utils "github.com/b-eee/amagi"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	// ConfirmEmailTplID email confirmation template ID
	ConfirmEmailTplID = "CONFIRM_EMAIL_TPL_ID"

	// ResetPasswordTplID reset password email template ID
	ResetPasswordTplID = "RESET_PASSWORD_TPL_ID"
)

type (
	// Email email struct for sending email
	Email struct {
		From      string    `bson:"from" json:"from"`
		To        string    `bson:"to" json:"to"`
		Subject   string    `bson:"subject" json:"subject"`
		CreatedAt time.Time `bson:"created_at" json:"created_at"`
		EmailSent bool      `bson:"email_sent" json:"email_sent"`

		PlainTextBody string `bson:"plain_text_body,omitempty" json:"plain_text_body,omitempty"`
		HTMLContent   string `bson:"html_content,omitempty" json:"html_content,omitempty"`

		Substitutions map[string]string `bson:"substitutions" json:"substitutions"`

		TemplateEnvName string
	}
)

// SendEmail send an email to user with sendgrid API
func (e *Email) SendEmail() error {
	var content *mail.Content
	s := time.Now()
	m := mail.NewV3Mail()
	from := mail.NewEmail(e.From, e.From)
	to := mail.NewEmail(e.To, e.To)

	if len(e.PlainTextBody) != 0 {
		content = mail.NewContent("text/plain", e.PlainTextBody)
	}

	m.SetFrom(from)
	m.AddContent(content)
	m.SetTemplateID(os.Getenv(e.TemplateEnvName))

	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	for k, v := range e.Substitutions {
		personalization.SetSubstitution(k, v)
	}

	personalization.Subject = e.Subject

	m.AddPersonalizations(personalization)
	sendgridAPIKEY := os.Getenv("SENDGRID_API_KEY")
	request := sendgrid.GetRequest(sendgridAPIKEY, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		return utils.Error(fmt.Sprintf("error on send Email error %v", err))
	}

	return utils.Info(fmt.Sprintf("SendEmail succes! took: %v statuscode: %v", time.Since(s), response.StatusCode))
}

// GetEmailTemplate get email template ID from env name
func GetEmailTemplate(envName string) string {
	return os.Getenv(envName)
}
