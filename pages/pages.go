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
var Home = views.HtmlTemplate("pages/home.html")

// NotFound is displayed when a request does not match any Route.
var NotFound = views.HtmlTemplate("pages/not_found.html")

// Error is displayed when an error occurs while processing a request
// handler.
var Error = http.HandlerFunc(errorPage)

func errorPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Context().Value(routes.Param("error")))
	views.HtmlTemplate("pages/error.html").ServeHTTP(w, r)
}
