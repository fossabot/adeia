package mailer

import (
	"io"

	"github.com/jordan-wright/email"
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
