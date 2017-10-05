package views

import (
	"bytes"
	"html/template"
	"net/http"
)

type htmlView struct {
	content string
}

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

type HtmlTemplate struct {
	template *template.Template
}

func ParseHtmlTemplate(source string) *HtmlTemplate {
	return &HtmlTemplate{
		template: template.Must(template.New("").Parse(source)),
	}
}

func (template *HtmlTemplate) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var output bytes.Buffer
	if err := template.template.Execute(&output, data); err != nil {
		panic(err)
	}
	sendHtml(w, output.Bytes())
}
