package models

//go:generate kanna-genmodel -output actor_gen.go actor.json

import (
	"context"
	"database/sql"

	"github.com/ekiru/kanna/db"
)

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
		"type":   &a.typ,
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
