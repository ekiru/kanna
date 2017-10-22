package routes

import "net/http"

type Failure interface {
	HandleFailure(*Router, http.ResponseWriter, *http.Request)
}

var NotFound Failure = notFoundFailure{}

type notFoundFailure struct{}

func (_ notFoundFailure) HandleFailure(router *Router, w http.ResponseWriter, r *http.Request) {
	router.notFound.ServeHTTP(w, r)
}
