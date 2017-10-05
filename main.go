package main

import (
	"log"
	"net/http"

	"kanna/routes"
)

func main() {
	routes := routes.Routes()
	log.Fatal(http.ListenAndServe(":9123", routes))
}
