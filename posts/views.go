package posts

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ekiru/kanna/models"
	"github.com/ekiru/kanna/pages"
	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/views"
)

// AddRoutes registers the routes related to posts on the Router.
func AddRoutes(router *routes.Router) {
	router.Route([]interface{}{"post", routes.Param("post")}, http.HandlerFunc(showPost))
}

func showPost(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Post *models.Post
	}
	postKey := r.Context().Value(routes.Param("post")).(string)
	postId := fmt.Sprintf("http://kanna.example/post/%s", postKey)
	log.Println(postKey, postId)
	if post, err := models.PostById(r.Context(), postId); err == nil {
		views.HtmlTemplate("post_show.html").Render(w, r, data{Post: post})
	} else {
		// TODO expose a way to serve a NotFound via the Router at this point
		// TODO expose a way to serve an error
		// TODO distinguish not found from other errors
		log.Println(err)
		pages.NotFound.ServeHTTP(w, r)
	}
}
