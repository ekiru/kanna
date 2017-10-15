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

// Middleware returns a middleware that will
func Middleware() routes.Middleware {
	return sessionMiddleware{
		sessions: make(map[string]*sessionData),
	}
}

const (
	cookieName = "kannaSession"
	idLen      = 42
)

type contextKey struct{}

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
	r = r.WithContext(context.WithValue(r.Context(), contextKey{}, session.load(r.Context())))
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

// Get retrieves the Session from the request context.
func Get(ctx context.Context) *Session {
	// TODO maybe check this
	return ctx.Value(contextKey{}).(*Session)
}
