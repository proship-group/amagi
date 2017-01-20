package mailSender

import (
	"github.com/b-eee/amagi/services/mailSender/mailerBackend"
)

type (
	// Mail mail interface
	Mail struct {
		Request mailerBackend.Request
	}
)

// InitMailSender init mail sender interface
func InitMailSender() *Mail {
	m := new(Mail)

	return m
}

// Send send email
func (m *Mail) Send() error {
	if err := mailerBackend.Send(&m.Request); err != nil {
		return err
	}

	return nil
}
