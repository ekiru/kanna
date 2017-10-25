package views

import (
	"context"
	"net/http"

	"github.com/ekiru/kanna/routes"
)

type paramMapper struct {
	handler http.Handler
	mapper  func(context.Context, interface{}) interface{}
	param   string
}

func MapParam(handler http.Handler, param string, mapper func(context.Context, interface{}) interface{}) http.Handler {
	return paramMapper{
		handler: handler,
		mapper:  mapper,
		param:   param,
	}
}

func (pm paramMapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	val := pm.mapper(ctx, ctx.Value(routes.Param(pm.param)))
	ctx = context.WithValue(ctx, routes.Param(pm.param), val)
	pm.handler.ServeHTTP(w, r.WithContext(ctx))
}
