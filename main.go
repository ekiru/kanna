package main

import (
	"log"
	"net/http"

	"github.com/ekiru/kanna/accounts"
	"github.com/ekiru/kanna/actors"
	"github.com/ekiru/kanna/db"
	"github.com/ekiru/kanna/pages"
	"github.com/ekiru/kanna/posts"
	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/sessions"
)

func main() {
	routes := buildRoutes()
	log.Fatal(http.ListenAndServe(":9123", routes))
}

func buildRoutes() http.Handler {
	var router routes.Router

	router.Middleware(sessions.Middleware())

	if err := db.InitParams(&router); err != nil {
		log.Fatal(err)
	}

	router.Route([]interface{}{}, pages.Home)

	accounts.AddRoutes(&router)
	actors.AddRoutes(&router)
	posts.AddRoutes(&router)

	router.NotFound(pages.NotFound)
	router.Error(pages.Error)

	return &router
}
