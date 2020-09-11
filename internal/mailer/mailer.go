package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"

	"adeia-api/internal/config"

	"github.com/jordan-wright/email"
)

type Mailer struct {
	from          string
	smtpAddr      string
	templateCache map[string]*template.Template
	auth          smtp.Auth
}

// New creates a new Mailer.
func New(conf *config.MailerConfig) (*Mailer, error) {
	// build new template cache
	tc, err := newTemplateCache(templateDir)
	if err != nil {
		return nil, err
	}

	u := conf.Username
	p := conf.Password
	host := conf.SMTPHost
	addr := host + ":" + strconv.Itoa(conf.SMTPPort)

	// create auth
	auth := smtp.PlainAuth("", u, p, host)
	m := &Mailer{
		from:          u,
		templateCache: tc,
		auth:          auth,
		smtpAddr:      addr,
	}

	return m, nil
}

// Send sends the specified email with the template and data.
func (m *Mailer) Send(e *email.Email, template string, data interface{}) error {
	e.From = m.from

	// read template from cache
	ts, ok := m.templateCache[template]
	if !ok {
		return fmt.Errorf("email template with name %v does not exist", template)
	}

	// execute template
	// we copy to buffer first to catch any runtime errors
	buf := new(bytes.Buffer)
	if err := ts.Execute(buf, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	e.HTML = buf.Bytes()

	// TODO: properly handle the error
	go e.Send(m.smtpAddr, m.auth)
	return nil
}
