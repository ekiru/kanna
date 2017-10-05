package main

import (
	"log"
	"net/http"

	"github.com/ekiru/kanna/pages"
	"github.com/ekiru/kanna/routes"
)

func main() {
	routes := buildRoutes()
	log.Fatal(http.ListenAndServe(":9123", routes))
}

func buildRoutes() http.Handler {
	var router routes.Router

	router.Route([]interface{}{}, pages.Home)

	router.NotFound(pages.NotFound)
	router.Error(pages.Error)

	return &router
}
