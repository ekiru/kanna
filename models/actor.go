package models

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/ekiru/kanna/activitystreams"
	"github.com/ekiru/kanna/db"
)

// An Actor represents an ActivityPub Actor, either on this server or
// elsewhere.
type Actor struct {
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

// AsObject serializes the Actor to an Activity Streams Object.
func (a *Actor) AsObject() *activitystreams.Object {
	return &activitystreams.Object{
		ID:   a.ID,
		Type: a.Type,
		Props: map[string]interface{}{
			"inbox":  a.Inbox,
			"outbox": a.Outbox,
			"name":   a.Name,
		},
	}
}

// FromRow fills an Actor with the data from a row returned by a
// database query from the Actors table.
func (a *Actor) FromRow(rows *sql.Rows) error {
	err := db.FromRow(rows, a.Scanners())
	return err
}

// Scanners returns a map of scanners that will scan database columns
// into the fields of the Actor.
func (a *Actor) Scanners() map[string]interface{} {
	return map[string]interface{}{
		"inbox":  db.URLScanner{&a.Inbox},
		"outbox": db.URLScanner{&a.Outbox},
		"name":   &a.Name,
		"type":   &a.Type,
		"id":     db.URLScanner{&a.ID},
	}
}

// ActorById retrieves an Actor from the database with the specified ID
// if they exist. If no such Actor exists, database/sql.ErrNoRows will
// be returned as the error. Other errors may be returned.
func ActorById(ctx context.Context, id string) (*Actor, error) {
	var actor Actor
	rows, err := db.DB(ctx).QueryContext(ctx, "select id, type, name, inbox, outbox from Actors where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	if err = actor.FromRow(rows); err != nil {
		return nil, err
	}
	return &actor, nil
}
