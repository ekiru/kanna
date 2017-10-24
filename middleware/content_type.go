package middleware

import (
	"net/http"
	"strings"

	"github.com/ekiru/kanna/activitystreams"
	"github.com/ekiru/kanna/routes"
)

type contentTypeOverride struct{}

func ContentTypeOverride() routes.Middleware {
	return contentTypeOverride{}
}

func (_ contentTypeOverride) HandleMiddleware(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".json") {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, ".json")
		r.Header.Set("Accept", activitystreams.ContentType)
	}
	return w, r
}
