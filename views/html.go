package views

import (
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
