package posts

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ekiru/kanna/activitystreams"
	"github.com/ekiru/kanna/models"
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
	postId := fmt.Sprintf("https://faew.ink/post/%s", postKey)
	if post, err := models.PostById(r.Context(), postId); err == nil {
		switch r.Header.Get("Accept") {
		case activitystreams.ContentType:
			views.ActivityStream(post).ServeHTTP(w, r)
		default:
			views.HtmlTemplate("posts/show.html").Render(w, r, data{Post: post})
		}
	} else if err == sql.ErrNoRows {
		panic(routes.NotFound)
	} else {
		panic(routes.Error(err))
	}
}
