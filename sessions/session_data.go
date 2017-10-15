package sessions

import (
	"context"

	"github.com/ekiru/kanna/models"
)

type sessionData struct {
	username string
}

func newSession() *sessionData {
	return &sessionData{
		username: "",
	}
}

func (sd *sessionData) load(ctx context.Context) *Session {
	var user *models.Account
	// If we can't find the user or we don't have a username in the
	// session, the User will be nil.
	if sd.username != "" {
		user, _ = models.AccountByUsername(ctx, sd.username)
	}
	return &Session{
		sd:   sd,
		User: user,
	}
}

// A Session stores session data for each client.
type Session struct {
	sd *sessionData
	// User is the Account model for the logged-in user, or nil if
	// the client is not logged-in. The User is retrieved at the
	// beginning of processing a request and is not automatically
	// updated if other copies of the Account are modified.
	User *models.Account
}

// Save saves changes to the session: particularly which user, if any,
// is currently logged-in.
func (s *Session) Save() {
	sd := s.sd
	if s.User != nil {
		sd.username = s.User.Username
	} else {
		sd.username = ""
	}
}
