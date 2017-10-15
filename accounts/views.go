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
	router.Route([]interface{}{"auth"}, http.HandlerFunc(authDispatcher))
}

func authDispatcher(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		authPost(w, r)
	} else if user := sessions.Get(r.Context()).User; user != nil {
		type data struct {
			User *models.Account
		}
		authLoggedInTemplate.Render(w, r, data{User: user})
	} else {
		authShowForm.ServeHTTP(w, r)
	}
}

var authLoggedInTemplate = views.ParseHtmlTemplate(`<!doctype html>
	<title>Kanna - Logged In!</title>
	<p>
		You're already logged in as {{.User.Username}}!
	</p>
`)

var authShowForm = views.Html(`<!doctype html>
	<title>Kanna - Login</title>
	<form method=post>
		<p>
			<label for=username>Username</label>
			<input type=text name=username />
		</p>
		<p>
			<label for=password>Password</label>
			<input type=password name=password />
		</p>
		<p>
			<input type=submit value="log in" />
		</p>
	</form>
`)

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
