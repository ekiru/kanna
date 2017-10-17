package views

import (
	"bytes"
	"html/template"
	"net/http"
)

var templates *template.Template = template.Must(template.ParseGlob("templates/*.html"))

// HtmlTemplate views serve an HTML document from a template defined
// using the template/html package.
type HtmlTemplate string

// Render renders the template to the response writer, passing the
// supplied data to the template.
func (template HtmlTemplate) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var output bytes.Buffer
	if err := templates.ExecuteTemplate(&output, string(template), data); err != nil {
		panic(err)
	}
	sendHtml(w, output.Bytes())
}

// ServeHTTP on an HtmlTemplate calls Render with nil data.
func (template HtmlTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	template.Render(w, r, nil)
}
