package models

import (
	"github.com/ekiru/kanna/db"
)

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
