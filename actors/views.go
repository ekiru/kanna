package actors

import (
	"database/sql"
	"fmt"
	"net/http"

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
	actorId := fmt.Sprintf("http://kanna.example/actor/%s", actorKey)
	if actor, err := models.ActorById(r.Context(), actorId); err == nil {
		switch r.Header.Get("Accept") {
		case activitystreams.ContentType:
			views.ActivityStream(actor).ServeHTTP(w, r)
		default:
			if posts, err := models.PostsByActor(r.Context(), actor); err == nil {
				showActorTemplate.Render(w, r, data{Actor: actor, Posts: posts})
			} else {
				panic(routes.Error(err))
			}
		}
	} else if err == sql.ErrNoRows {
		panic(routes.NotFound)
	} else {
		panic(routes.Error(err))
	}
}
