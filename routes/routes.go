package routes

import (
	"context"
	"net/http"
	"strings"

	"kanna/pages"
)

func Routes() http.Handler {
	var router router
	router.Route([]interface{}{}, pages.Home)
	router.NotFound(pages.NotFound)
	router.Error(pages.Error)
	return &router
}

type router struct {
	routes       []*route
	notFound     http.Handler
	errorHandler http.Handler
}

type Param string
type Rest string

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			r = r.WithContext(
				context.WithValue(r.Context(), Param("error"), err),
			)
			router.errorHandler.ServeHTTP(w, r)
		}
	}()
	var ok bool
	for _, route := range router.routes {
		if ok, r = route.Match(r); ok {
			route.ServeHTTP(w, r)
			return
		}
	}
	router.notFound.ServeHTTP(w, r)
}

func (router *router) Route(pattern []interface{}, handler http.Handler) {
	// Validate the pattern
	seenRest := false
	for _, component := range pattern {
		if seenRest {
			panic("Invalid route: contained additional components after a Rest")
		}
		switch component.(type) {
		case Rest:
			seenRest = true
		case string, Param:
			continue
		default:
			panic("Invalid route type")
		}
	}
	router.routes = append(router.routes, &route{
		patternComponents: pattern,
		handler:           handler,
	})
}

func (router *router) NotFound(handler http.Handler) {
	router.notFound = handler
}

func (router *router) Error(handler http.Handler) {
	router.errorHandler = handler
}

type route struct {
	patternComponents []interface{}
	handler           http.Handler
}

func (route *route) Match(r *http.Request) (bool, *http.Request) {
	ctx := r.Context()
	urlPath := []string{}
	for _, part := range strings.Split(r.URL.Path, "/") {
		if part != "" {
			urlPath = append(urlPath, part)
		}
	}
	handledThrough := 0
	for i, component := range route.patternComponents {
		switch component := component.(type) {
		case string:
			if len(urlPath) <= i || component != urlPath[i] {
				return false, r
			}
			handledThrough = i + 1
		case Param:
			if len(urlPath) <= i {
				return false, r
			}
			ctx = context.WithValue(ctx, component, urlPath[i])
			handledThrough = i + 1
		case Rest:
			ctx = context.WithValue(ctx, component, urlPath[i:])
			handledThrough = len(urlPath)
		default:
			panic("unreachable")
		}
	}
	if handledThrough != len(urlPath) {
		return false, r
	}
	return true, r.WithContext(ctx)
}

func (route *route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route.handler.ServeHTTP(w, r)
}
