package accounts

import (
	"net/http"

	"github.com/ekiru/kanna/models"
	"github.com/ekiru/kanna/routes"
	"github.com/ekiru/kanna/sessions"
	"github.com/ekiru/kanna/views"
)

// AddRoutes registers the routes related to accounts on the Router.
func AddRoutes(router *routes.Router) {
	router.Route([]interface{}{routes.Method{"GET"}, "auth"}, http.HandlerFunc(authGet))
	router.Route([]interface{}{routes.Method{"POST"}, "auth"}, http.HandlerFunc(authPost))
	router.Route([]interface{}{"auth", "logout"}, http.HandlerFunc(authLogout))
}

func authGet(w http.ResponseWriter, r *http.Request) {
	if user := sessions.Get(r.Context()).User; user != nil {
		type data struct {
			User *models.Account
		}
		authLoggedInTemplate.Render(w, r, data{User: user})
	} else {
		authShowForm.ServeHTTP(w, r)
	}
}

var authLoggedInTemplate = views.HtmlTemplate("auth/logged_in.html")

var authShowForm = views.HtmlTemplate("auth/form.html")

func authPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	username, password := r.PostForm.Get("username"), r.PostForm.Get("password")
	user, err := models.Authenticate(r.Context(), username, password)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("login failed"))
		return
	}
	sess := sessions.Get(r.Context())
	sess.User = user
	sess.Save()
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

func authLogout(w http.ResponseWriter, r *http.Request) {
	sessions.Close(r.Context())
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}
