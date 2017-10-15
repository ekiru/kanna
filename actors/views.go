package actors

import (
	"fmt"
	"net/http"

	"github.com/ekiru/kanna/models"
	"github.com/ekiru/kanna/pages"
	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/views"
)

// AddRoutes registers the routes related to actors on the Router.
func AddRoutes(router *routes.Router) {
	router.Route([]interface{}{"actor", routes.Param("actor")}, http.HandlerFunc(showActor))
}

var showActorTemplate = views.ParseHtmlTemplate(`<!doctype html>
<title>Kanna - {{.Actor.Name}} ({{.Actor.ID}}</title>
<h1>The {{.Actor.Type}} named <a href={{.Actor.ID}}>{{.Actor.Name}}</a></h1>

<nav>
	<ul>
		<li><a href={{.Actor.Inbox}}>Inbox</a>
		<li><a href={{.Actor.Outbox}}>Outbox</a>
	</ul>
</nav>`)

func showActor(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Actor *models.Actor
	}
	actorKey := r.Context().Value(routes.Param("actor")).(string)
	actorId := fmt.Sprintf("http://kanna.example/actor/%s", actorKey)
	if actor, err := models.ActorById(r.Context(), actorId); err == nil {
		showActorTemplate.Render(w, r, data{Actor: actor})
	} else {
		// TODO expose a way to serve a NotFound via the Router at this point
		// TODO expose a way to serve an error
		// TODO distinguish not found from other errors
		pages.NotFound.ServeHTTP(w, r)
	}
}
