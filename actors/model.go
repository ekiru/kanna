package actors

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/ekiru/kanna/db"
)

type Model struct {
	Inbox  *url.URL `json:"inbox"`
	Outbox *url.URL `json:"outbox"`
	Name   string   `json:"name"`
	Type   string   `json:"type"`
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

func (m *Model) FromRow(rows *sql.Rows) error {
	var inbox, outbox string
	err := db.FromRow(rows, map[string]interface{}{
		"inbox":  &inbox,
		"outbox": &outbox,
		"name":   &m.Name,
		"type":   &m.Type,
	})
	if err != nil {
		return err
	}
	if m.Inbox, err = url.Parse(inbox); err != nil {
		return err
	}
	if m.Outbox, err = url.Parse(outbox); err != nil {
		return err
	}
	return nil
}

func ById(ctx context.Context, id string) (*Model, error) {
	var model Model
	rows, err := db.DB(ctx).QueryContext(ctx, "select type, name, inbox, outbox from Actors where id = ?", id)
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
