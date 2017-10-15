package routes

import "net/http"

// A Middleware processes the ResponseWriter and Request for a route.
// A Router can be configured to call a Middleware prior to executing
// each route.
type Middleware interface {
	// HandleMiddleware performs the processing for the middleware.
	HandleMiddleware(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request)
}
