package views

import (
	"bytes"
	"html/template"
	"net/http"
)

type htmlView struct {
	content string
}

// Html views serve the passed string as an HTML document.
func Html(doc string) http.Handler {
	return &htmlView{
		content: doc,
	}
}

func sendHtml(w http.ResponseWriter, buf []byte) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(buf)
}

func (view *htmlView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sendHtml(w, []byte(view.content))
}

// HtmlTemplate views serve an HTML document from a template defined
// using the template/html package.
type HtmlTemplate struct {
	template *template.Template
}

// ParseHtmlTemplate parses the source string as a template and
// constructs a HtmlTemplate view for the template. If an error occurs
// parsing the template, ParseHtmlTemplate will panic.
func ParseHtmlTemplate(source string) *HtmlTemplate {
	return &HtmlTemplate{
		template: template.Must(template.New("").Parse(source)),
	}
}

// Render renders the template to the response writer, passing the
// supplied data to the template.
func (template *HtmlTemplate) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var output bytes.Buffer
	if err := template.template.Execute(&output, data); err != nil {
		panic(err)
	}
	sendHtml(w, output.Bytes())
}
