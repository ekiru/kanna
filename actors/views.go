package actors

import (
	"fmt"
	"net/http"

	"github.com/ekiru/kanna/pages"
	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/views"
)

func AddRoutes(router *routes.Router) {
	router.Route([]interface{}{"actor", routes.Param("actor")}, http.HandlerFunc(showActor))
}

var showActorTemplate = views.ParseHtmlTemplate(`<!doctype html>
<title>Kanna - {{.Actor.Name}}</title>
<h1>The {{.Actor.Type}} named {{.Actor.Name}}</h1>

<nav>
	<ul>
		<li><a href={{.Actor.Inbox}}>Inbox</a>
		<li><a href={{.Actor.Outbox}}>Outbox</a>
	</ul>
</nav>`)

func showActor(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Actor *Model
	}
	actorKey := r.Context().Value(routes.Param("actor")).(string)
	actorId := fmt.Sprintf("http://kanna.example/actor/%s", actorKey)
	if ok, actor := ById(actorId); ok {
		showActorTemplate.Render(w, r, data{Actor: actor})
	} else {
		// TODO expose a way to serve a NotFound via the Router at this point
		pages.NotFound.ServeHTTP(w, r)
	}
}
