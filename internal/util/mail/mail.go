package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/smtp"

	"github.com/jordan-wright/email"
	config "github.com/spf13/viper"
)

// EmailBuilder is an interface for email builder.
type EmailBuilder interface {
	To([]string) EmailBuilder
	ReplyTo([]string) EmailBuilder
	Bcc([]string) EmailBuilder
	Cc([]string) EmailBuilder
	Subject(string) EmailBuilder
	Sender(string) EmailBuilder
	AttachFile(string) (EmailBuilder, error)
	Attach(r io.Reader, filename string, c string) (EmailBuilder, error)
	Build() *email.Email
}

// emailBuilderImpl implements the EmailBuilder interface.
type emailBuilderImpl struct {
	*email.Email
}

// NewEmailBuilder creates a new EmailBuilder.
func NewEmailBuilder() EmailBuilder {
	return &emailBuilderImpl{email.NewEmail()}
}

// To sets the slice of email addresses as the `To:` field of the email.
func (e *emailBuilderImpl) To(to []string) EmailBuilder {
	e.Email.To = to
	return e
}

// ReplyTo sets the slice of email addresses as the `ReplyTo:` field of the email.
func (e *emailBuilderImpl) ReplyTo(r []string) EmailBuilder {
	e.Email.ReplyTo = r
	return e
}

// Bcc sets the slice of email addresses as the `Bcc:` field of the email.
func (e *emailBuilderImpl) Bcc(bcc []string) EmailBuilder {
	e.Email.Bcc = bcc
	return e
}

// Cc sets the slice of email addresses as the `Cc:` field of the email.
func (e *emailBuilderImpl) Cc(cc []string) EmailBuilder {
	e.Email.Cc = cc
	return e
}

// Subject sets the `Subject:` field of the email.
func (e *emailBuilderImpl) Subject(s string) EmailBuilder {
	e.Email.Subject = s
	return e
}

// Sender sets the `Sender:` field of the email.
func (e *emailBuilderImpl) Sender(s string) EmailBuilder {
	e.Email.Sender = s
	return e
}

// AttachFile attaches the specified filename to the email.
func (e *emailBuilderImpl) AttachFile(f string) (EmailBuilder, error) {
	if _, err := e.Email.AttachFile(f); err != nil {
		return nil, err
	}
	return e, nil
}

// Attach creates an attachment from the io.Reader, adds the specified filename
// and content-type, and attaches it to the email.
func (e *emailBuilderImpl) Attach(r io.Reader, filename string, c string) (EmailBuilder, error) {
	if _, err := e.Email.Attach(r, filename, c); err != nil {
		return nil, err
	}
	return e, nil
}

// Build builds the email.
func (e *emailBuilderImpl) Build() *email.Email {
	return e.Email
}

// Mailer represents a mailer service.
type Mailer interface {
	Send(e *email.Email, template Template, data interface{}) error
}

// mailerImpl implements Mailer.
type mailerImpl struct {
	from          string
	smtpAddr      string
	templateCache map[string]*template.Template
	auth          smtp.Auth
}

// NewMailer creates a new Mailer.
func NewMailer() (Mailer, error) {
	// build new template cache
	tc, err := newTemplateCache(templateDir)
	if err != nil {
		return nil, err
	}

	u := config.GetString("mailer.username")
	p := config.GetString("mailer.password")
	host := config.GetString("mailer.smtp_host")
	port := config.GetString("mailer.smtp_port")
	addr := host + ":" + port

	// create auth
	auth := smtp.PlainAuth("", u, p, host)
	m := &mailerImpl{
		from:          u,
		templateCache: tc,
		auth:          auth,
		smtpAddr:      addr,
	}

	return m, nil
}

// Send sends the specified email with the template and data.
func (m *mailerImpl) Send(e *email.Email, template Template, data interface{}) error {
	e.From = m.from

	// read template from cache
	ts, ok := m.templateCache[template.name]
	if !ok {
		return fmt.Errorf("email template with name %v does not exist", template.name)
	}

	// execute template
	// we copy to buffer first to catch any runtime errors
	buf := new(bytes.Buffer)
	if err := ts.Execute(buf, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	e.HTML = buf.Bytes()

	if err := e.Send(m.smtpAddr, m.auth); err != nil {
		return fmt.Errorf("cannot send email: %v", err)
	}
	return nil
}
