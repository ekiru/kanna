// The actors package defines request handlers and database helpers for
// retrieving and modifying information about ActivityPub Actors.
package actors

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/ekiru/kanna/db"
)

// actors.Model represents an Actor, either on this server or elsewhere.
type Model struct {
	// Inbox is the URL of the Actor's inbox, to which Activities
	// can be posted to deliver them to the Actor and from which the
	// Actor can read those Activities.
	Inbox *url.URL `json:"inbox"`
	// Outbox is the URL of the Actor's outbox, to which the Actor
	// can post Activities to deliver them to other Actors and from
	// which other Actors can read (some) of the Activities the
	// Actor has posted.
	Outbox *url.URL `json:"outbox"`
	// Name is a "simple, human-readable, plain-text name for the"
	// Actor.
	Name string `json:"name"`
	// Type represents the type of Actor; common values include
	// Person, Group, Organization, and Application.
	Type string `json:"type"`
	// ID is the URL which uniquely identifies the Actor.
	ID *url.URL `json:"id"`
}

var exampleInbox, exampleOutbox *url.URL

func init() {
	var err error
	exampleInbox, err = url.Parse("http://kanna.example/actor/srn/inbox")
	if err != nil {
		panic("failed parsing example inbox url")
	}
	exampleOutbox, err = url.Parse("http://kanna.example/actor/srn/outbox")
	if err != nil {
		panic("failed parsing example outbox url")
	}
}

// FromRow fills a Model with the data from a row returned by a
// database query from the Actors table.
func (m *Model) FromRow(rows *sql.Rows) error {
	err := db.FromRow(rows, map[string]interface{}{
		"inbox":  db.URLScanner{&m.Inbox},
		"outbox": db.URLScanner{&m.Outbox},
		"name":   &m.Name,
		"type":   &m.Type,
		"id":     db.URLScanner{&m.ID},
	})
	return err
}

// ById retrieves an actor.Model from the database with the specified
// ID if they exist. If no such Actor exists, database/sql.ErrNoRows
// will be returned as the error. Other errors may be returned.
func ById(ctx context.Context, id string) (*Model, error) {
	var model Model
	rows, err := db.DB(ctx).QueryContext(ctx, "select id, type, name, inbox, outbox from Actors where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	if err = model.FromRow(rows); err != nil {
		return nil, err
	}
	return &model, nil
}
