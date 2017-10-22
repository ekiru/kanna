// The routes package defines the routing system and interface that
// the rest of Kanna uses to route requests to different endpoints
// and parse parameters out of request paths.
package routes

import (
	"context"
	"net/http"
	"strings"
)

// A Router multiplexes HTTP requests similar to the net/http package's
// ServeMux, but with more flexible patterns that allow conveniently
// parsing dynamic components out of request paths as well as
// automating some other details.
type Router struct {
	routes       []*route
	notFound     http.Handler
	errorHandler http.Handler
	baseParams   []baseParam
	middleware   []Middleware
}

type baseParam struct {
	key, value interface{}
}

// A Param matches any path component in a URL pattern and stores the
// matched path component in the request context. The matched portion
// of the path is stored within the context with the Param value for
// the parameter name as the key. The value stored in the context will
// always be a string.
//
// TODO define accessor functions for the context keys for these
type Param string

// A Rest parameter in a URL pattern matches the rest of a request path
// and stores that rest of the path in the context with the Rest value
// for the parameter name as the key. The value stored in the context
// will be the remainder of the path as a string.
type Rest string

// A Method component in a pattern matches if the request's HTTP method
// is equal to one of the values in the Method.
type Method []string

// ServeHTTP fulfills the http.Handler interface and handles HTTP
// requests by dispatching them to the appropriate route or the
// NotFound handler if no route matches. Routes are tested against in
// the order that they were defined. The first matching route will be
// served.
//
// If the route panics, then the panic will be recovered here and
// displayed using the Error handler, passing the recovered value in
// the Param("error") context key. The error value is not guaranteed to
// satisfy the error interface.
//
// Any BaseParams specified on the Router will be added to the request
// context before calling any handler.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			switch err := err.(type) {
			case Failure:
				err.HandleFailure(router, w, r)
			default:
				r = r.WithContext(
					context.WithValue(r.Context(), Param("error"), err),
				)
				router.errorHandler.ServeHTTP(w, r)
			}
		}
	}()
	ctx := r.Context()
	for _, param := range router.baseParams {
		ctx = context.WithValue(ctx, param.key, param.value)
	}
	r = r.WithContext(ctx)
	for _, mw := range router.middleware {
		w, r = mw.HandleMiddleware(w, r)
	}
	var ok bool
	for _, route := range router.routes {
		if ok, r = route.Match(r); ok {
			route.ServeHTTP(w, r)
			return
		}
	}
	router.notFound.ServeHTTP(w, r)
}

// BaseParam defines base parameters that the Router will define on the
// request's context prior to calling any handler.
func (router *Router) BaseParam(param interface{}, value interface{}) {
	router.baseParams = append(router.baseParams, baseParam{param, value})
}

// Middleware adds a middleware to execute on all matching requests.
func (router *Router) Middleware(mw Middleware) {
	// TODO Should we run middleware on notfound/error requests?
	router.middleware = append(router.middleware, mw)
}

// Route maps a pattern to a http.Handler. The Router's ServeHTTP
// method will dispatch requests to the earliest-defined route with a
// matching pattern. The pattern must consist only of strings, Params,
// and Rests. Strings match path component identical to the string.
// Param values match any single path component. Rest values match the
// entire rest of the request path. A Rest may only appear at the end
// of a pattern.
func (router *Router) Route(pattern []interface{}, handler http.Handler) {
	// Validate the pattern
	seenRest := false
	pathPattern := make([]interface{}, 0, len(pattern))
	var methods Method
	for _, component := range pattern {
		if seenRest {
			panic("Invalid route: contained additional components after a Rest")
		}
		switch component := component.(type) {
		case Method:
			if methods != nil {
				panic("Invalid route: contained multiple Method components.")
			}
			methods = component
		case Rest:
			seenRest = true
			pathPattern = append(pathPattern, component)
		case string, Param:
			pathPattern = append(pathPattern, component)
		default:
			panic("Invalid route type")
		}
	}
	router.routes = append(router.routes, &route{
		patternComponents: pathPattern,
		methods:           methods,
		handler:           handler,
	})
}

// NotFound specifies a handler that will be called when no defined
// Route matches a request. A NotFound handler must be defined.
func (router *Router) NotFound(handler http.Handler) {
	router.notFound = handler
}

// Error specifies a handler that will be called when a Route panics.
// The value passed to panic will be passed to the Error handler in the
// Param("error") context key. Router does not guarantee that this
// value will implement the error interface.
func (router *Router) Error(handler http.Handler) {
	router.errorHandler = handler
}

type route struct {
	patternComponents []interface{}
	methods           Method
	handler           http.Handler
}

func (route *route) Match(r *http.Request) (bool, *http.Request) {
	if route.methods != nil {
		var found bool
		for _, m := range route.methods {
			if r.Method == m {
				found = true
				break
			}
		}
		if !found {
			return false, r
		}
	}
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
