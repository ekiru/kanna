package models

//go:generate kanna-genmodel -output actor_gen.go actor.json

import (
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
