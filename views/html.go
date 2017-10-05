package views

import (
	"io"
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

func (view *htmlView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, view.content)
}
