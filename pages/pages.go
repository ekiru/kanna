// The pages package defines a few miscellaneous pages that don't
// belong to some more specific group of functionality.
package pages

import (
	"log"
	"net/http"

	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/views"
)

// Home handles requests to the root path and currently doesn't really
// do anything.
var Home = views.Html(
	`<!doctype html>
<title>Kanna - Hoooommmmeeeee</title>
<p>
	This is just stubbing in a home page to have something existing.
</p>`)

// NotFound is displayed when a request does not match any Route.
var NotFound = views.Html(
	`<!doctype html>
<title>Kanna - Page Not Found</title>
<p>
	Kanna can't find yr page. T_T Please give her headpats before she starts crying.
</p>`)

// Error is displayed when an error occurs while processing a request
// handler.
var Error = http.HandlerFunc(errorPage)

func errorPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Context().Value(routes.Param("error")))
	views.Html(
		`<!doctype html>
<title>Kanna - Page Not Found</title>
<p>
	Oh no, something went wrong. :( Kanna is _not_ happy.
</p>`).ServeHTTP(w, r)
}
