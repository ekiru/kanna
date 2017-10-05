package routes

import (
	"net/http"

	"kanna/pages"
)

func Routes() http.Handler {
	var router router
	router.Route([]interface{}{}, pages.Home)
	router.NotFound(pages.NotFound)
	router.Error(pages.Error)
	return &router
}
