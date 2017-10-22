package actors

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ekiru/kanna/activitystreams"
	"github.com/ekiru/kanna/models"
	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/views"
)

// AddRoutes registers the routes related to actors on the Router.
func AddRoutes(router *routes.Router) {
	router.Route([]interface{}{"actor", routes.Param("actor")}, http.HandlerFunc(showActor))
}

var showActorTemplate = views.HtmlTemplate("actors/show.html")

func showActor(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Actor *models.Actor
		Posts []*models.Post
	}
	actorKey := r.Context().Value(routes.Param("actor")).(string)
	if strings.HasSuffix(actorKey, ".json") {
		// overrides to expect activity streams.
		actorKey = strings.TrimSuffix(actorKey, ".json")
		r.Header.Set("Accept", activitystreams.ContentType)
	}
	actorId := fmt.Sprintf("http://kanna.example/actor/%s", actorKey)
	if actor, err := models.ActorById(r.Context(), actorId); err == nil {
		if posts, err := models.PostsByActor(r.Context(), actor); err == nil {
			switch r.Header.Get("Accept") {
			case activitystreams.ContentType:
				views.ActivityStream(actor).ServeHTTP(w, r)
			default:
				showActorTemplate.Render(w, r, data{Actor: actor, Posts: posts})
			}
			return
		} else {
			panic(routes.Error(err))
		}
	} else if err == sql.ErrNoRows {
		panic(routes.NotFound)
	} else {
		panic(routes.Error(err))
	}
}
