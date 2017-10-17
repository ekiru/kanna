package actors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ekiru/kanna/activitystreams"
	"github.com/ekiru/kanna/models"
	"github.com/ekiru/kanna/pages"
	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/views"
)

// AddRoutes registers the routes related to actors on the Router.
func AddRoutes(router *routes.Router) {
	router.Route([]interface{}{"actor", routes.Param("actor")}, http.HandlerFunc(showActor))
}

var showActorTemplate = views.HtmlTemplate("actor_show.html")

func showActor(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Actor *models.Actor
	}
	actorKey := r.Context().Value(routes.Param("actor")).(string)
	if strings.HasSuffix(actorKey, ".json") {
		// overrides to expect activity streams.
		actorKey = strings.TrimSuffix(actorKey, ".json")
		r.Header.Set("Accept", activitystreams.ContentType)
	}
	actorId := fmt.Sprintf("http://kanna.example/actor/%s", actorKey)
	if actor, err := models.ActorById(r.Context(), actorId); err == nil {
		switch r.Header.Get("Accept") {
		case activitystreams.ContentType:
			views.ActivityStream(actor).ServeHTTP(w, r)
		default:
			showActorTemplate.Render(w, r, data{Actor: actor})
		}
	} else {
		// TODO expose a way to serve a NotFound via the Router at this point
		// TODO expose a way to serve an error
		// TODO distinguish not found from other errors
		pages.NotFound.ServeHTTP(w, r)
	}
}
