package views

import (
	"net/http"

	"github.com/ekiru/kanna/activitystreams"
)

type activityStreamsView struct {
	obj activitystreams.AsObject
}

// ActivityStream creates a handler that serializes an object as an
// Activity Stream and serves it.
func ActivityStream(obj activitystreams.AsObject) http.Handler {
	return activityStreamsView{obj}
}

func (view activityStreamsView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf, err := activitystreams.Marshal(view.obj)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", activitystreams.ContentType)
	w.Write(buf)
}
