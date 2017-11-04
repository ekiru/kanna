
package models

import (
	
	"context"
	"database/sql"
	
	"net/url"
	
	"github.com/ekiru/kanna/db"
	
)

type Actor struct {
	id *url.URL
	typ string
	Inbox *url.URL
	Name string
	Outbox *url.URL
}

func (model *Actor) ID() *url.URL {
	return model.id
}

func (model *Actor) Types() []string {
	return []string{ model.typ }
}

func (model *Actor) HasType(t string) bool {
	return t == model.typ
}

func (model *Actor) Props() []string {
	return []string{ "id", "type", "inbox","name","outbox", }
}

func (model *Actor) GetProp(prop string) (interface{}, bool) {
	switch prop {
	case "id":
		return model.id, true
	case "type":
		return model.typ, true
	case "inbox":
		return model.Inbox, true
	case "name":
		return model.Name, true
	case "outbox":
		return model.Outbox, true
	default:
		return nil, false
	}
}

func ActorById(ctx context.Context, id string) (*Actor, error) {
	var model Actor
	rows, err := db.DB(ctx).QueryContext(ctx, "select id, type, inbox, name, outbox from Actors where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	err = rows.Scan(
		db.URLScanner{ &model.id },
		&model.typ,
		db.URLScanner{ &model.Inbox },
		&model.Name,
		db.URLScanner{ &model.Outbox },
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
