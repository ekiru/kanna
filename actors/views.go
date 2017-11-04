package actors

import (
	"context"
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
	router.Route([]interface{}{"actor", routes.Param("actor")}, actorParam(http.HandlerFunc(showActor)))
}

func actorParam(handler http.Handler) http.Handler {
	return views.MapParam(handler, "actor", func(ctx context.Context, actorKey interface{}) interface{} {
		actorId := fmt.Sprintf("https://faew.ink/actor/%s", actorKey.(string))
		if actor, err := models.ActorById(ctx, actorId); err == nil {
			return actor
		} else if err == sql.ErrNoRows {
			panic(routes.NotFound)
		} else {
			panic(routes.Error(err))
		}
	})
}

var showActorTemplate = views.HtmlTemplate("actors/show.html")

func showActor(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Actor *models.Actor
		Posts []*models.Post
	}
	actor := r.Context().Value(routes.Param("actor")).(*models.Actor)
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
}
