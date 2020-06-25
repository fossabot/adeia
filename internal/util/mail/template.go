package mail

import (
	"html/template"
	"path/filepath"
)

const templateDir = "web/email_templates"

// Template represents an email template.
type Template struct {
	name string
	// Data is exposed so that the consumer can copy the struct, update its fields
	// and pass it to Mailer's `Send()`.
	Data interface{}
}

var (
	// EmailTemplateEmailVerify is the email verification template.
	EmailTemplateEmailVerify = Template{
		name: "email_verify.tmpl",
		Data: struct {
			Link string
		}{},
	}
)

// newTemplateCache creates a new in-memory cache (a map) of templates.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// iterate over all pages
	pages, err := filepath.Glob(filepath.Join(dir, "pages/*.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// parse pages
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// parse required layouts
		ts, err = ts.ParseGlob(filepath.Join(dir, "layouts/*.tmpl"))
		if err != nil {
			return nil, err
		}

		// parse required partials
		ts, err = ts.ParseGlob(filepath.Join(dir, "partials/*.tmpl"))
		if err != nil {
			return nil, err
		}

		// store template to cache
		cache[name] = ts
	}

	return cache, nil
}
