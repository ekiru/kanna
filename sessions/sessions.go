// The sessions package provides a Middleware that manages user
// sessions, using a session identifier in the URL.
package sessions

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/ekiru/kanna/routes"
)

type sessionMiddleware struct {
	sessions map[string]*sessionData
}

// Middleware returns a middleware that uses a cookie to store a random
// session ID and stores a Session in the request context.
func Middleware() routes.Middleware {
	return sessionMiddleware{
		sessions: make(map[string]*sessionData),
	}
}

const (
	cookieName = "kannaSession"
	idLen      = 42
)

type mwContextKey struct{}
type sessionContextKey struct{}
type sessionIdContextKey struct{}

func (mw sessionMiddleware) HandleMiddleware(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request) {
	var sessionId string
	var session *sessionData
	if cookie, err := r.Cookie(cookieName); err == nil {
		sessionId = cookie.Value
		var ok bool
		if session, ok = mw.sessions[sessionId]; !ok {
			sessionId, session = mw.createSession(w)
		}
	} else {
		sessionId, session = mw.createSession(w)
	}
	ctx := r.Context()
	ctx = context.WithValue(ctx, mwContextKey{}, mw)
	ctx = context.WithValue(ctx, sessionContextKey{}, session.load(ctx))
	ctx = context.WithValue(ctx, sessionIdContextKey{}, sessionId)
	r = r.WithContext(ctx)
	return w, r
}

func (mw sessionMiddleware) createSession(w http.ResponseWriter) (string, *sessionData) {
	sessionId := newSessionId()
	session := newSession()
	mw.sessions[sessionId] = session
	w.Header().Add("Set-Cookie", (&http.Cookie{
		Name:  cookieName,
		Value: sessionId,
		Path:  "/",
		// Secure:   true,
		HttpOnly: true,
	}).String())
	return sessionId, session
}

func newSessionId() string {
	id := make([]byte, idLen)
	if _, err := rand.Read(id); err != nil {
		panic(err)
	}
	return hex.EncodeToString(id)
}

// Close invalidates the current Session.
func Close(ctx context.Context) {
	mw := ctx.Value(mwContextKey{}).(sessionMiddleware)
	id := ctx.Value(sessionIdContextKey{}).(string)
	delete(mw.sessions, id)
}

// Get retrieves the Session from the request context.
func Get(ctx context.Context) *Session {
	// TODO maybe check this
	return ctx.Value(sessionContextKey{}).(*Session)
}
