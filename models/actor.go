package models

import (
	"context"
	"database/sql"
	"net/url"

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
	id *url.URL `json:"id"`
}

func (a *Actor) ID() *url.URL {
	return a.id
}

func (a *Actor) Types() []string {
	return []string{a.Type}
}

func (a *Actor) HasType(t string) bool {
	return a.Type == t
}

func (a *Actor) GetProp(name string) (interface{}, bool) {
	switch name {
	case "id":
		return a.id, true
	case "type":
		return a.Type, true
	case "name":
		return a.Name, true
	case "inbox":
		return a.Inbox, true
	case "outbox":
		return a.Outbox, true
	default:
		return nil, false
	}
}

func (a *Actor) Props() []string {
	return []string{"id", "type", "name", "inbox", "outbox"}
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
		"id":     db.URLScanner{&a.id},
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
