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

type Session struct {
	sd   *sessionData
	User *models.Account
}

func (s *Session) Save() {
	sd := s.sd
	if s.User != nil {
		sd.username = s.User.Username
	} else {
		sd.username = ""
	}
}
