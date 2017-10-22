package routes

import (
	"context"
	"net/http"
)

type Failure interface {
	HandleFailure(*Router, http.ResponseWriter, *http.Request)
}

var NotFound Failure = notFoundFailure{}

type notFoundFailure struct{}

func (_ notFoundFailure) HandleFailure(router *Router, w http.ResponseWriter, r *http.Request) {
	router.notFound.ServeHTTP(w, r)
}

type errorFailure struct {
	err interface{}
}

func Error(err interface{}) Failure {
	return errorFailure{err}
}

func (err errorFailure) HandleFailure(router *Router, w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(
		context.WithValue(r.Context(), Param("error"), err.err),
	)
	router.errorHandler.ServeHTTP(w, r)
}
