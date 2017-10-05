package pages

import (
	"io"
	"net/http"
)

var Home = http.HandlerFunc(home)

func home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w,
		`<!doctype html>
<title>Kanna - Hoooommmmeeeee</title>
<p>
	This is just stubbing in a home page to have something existing.
</p>`)
}

var NotFound = http.HandlerFunc(notFound)

func notFound(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w,
		`<!doctype html>
<title>Kanna - Page Not Found</title>
<p>
	Kanna can't find yr page. T_T Please give her headpats before she starts crying.
</p>`)
}

var Error = http.HandlerFunc(errorPage)

func errorPage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w,
		`<!doctype html>
<title>Kanna - Page Not Found</title>
<p>
	Oh no, something went wrong. :( Kanna is _not_ happy.
</p>`)
}
