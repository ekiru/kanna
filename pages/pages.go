package pages

import "github.com/ekiru/kanna/views"

var Home = views.Html(
	`<!doctype html>
<title>Kanna - Hoooommmmeeeee</title>
<p>
	This is just stubbing in a home page to have something existing.
</p>`)

var NotFound = views.Html(
	`<!doctype html>
<title>Kanna - Page Not Found</title>
<p>
	Kanna can't find yr page. T_T Please give her headpats before she starts crying.
</p>`)

var Error = views.Html(
	`<!doctype html>
<title>Kanna - Page Not Found</title>
<p>
	Oh no, something went wrong. :( Kanna is _not_ happy.
</p>`)
