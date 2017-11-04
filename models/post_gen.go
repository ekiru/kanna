package models

import (
	
	"context"
	"database/sql"
	
	"net/url"
	
	"github.com/ekiru/kanna/db"
	
)

type Post struct {
	id *url.URL
	typ string
	Audience string
	Author *Actor
	Content string
	Published string
}

func (model *Post) ID() *url.URL {
	return model.id
}

func (model *Post) Types() []string {
	return []string{ model.typ }
}

func (model *Post) HasType(t string) bool {
	return t == model.typ
}

func (model *Post) Props() []string {
	return []string{ "id", "type", "audience","author","content","published", }
}

func (model *Post) GetProp(prop string) (interface{}, bool) {
	switch prop {
	case "id":
		return model.id, true
	case "type":
		return model.typ, true
	case "audience":
		return model.Audience, true
	case "author":
		return model.Author, true
	case "content":
		return model.Content, true
	case "published":
		return model.Published, true
	default:
		return nil, false
	}
}

func PostById(ctx context.Context, id string) (*Post, error) {
	var model Post
	rows, err := db.DB(ctx).QueryContext(ctx, "select Posts.id, Posts.type, Posts.audience, Posts.authorId, Posts.content, Posts.published, Actors.type, Actors.inbox, Actors.name, Actors.outbox from Posts join Actors on Posts.authorId = Actors.id where Posts.id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	model.Author = new(Actor)

	err = rows.Scan(
		db.URLScanner{ &model.id },
		&model.typ,
		&model.Audience,
		db.URLScanner{ &model.Author.id },
		&model.Content,
		&model.Published,
		&model.Author.typ,
		db.URLScanner{ &model.Author.Inbox },
		&model.Author.Name,
		db.URLScanner{ &model.Author.Outbox },
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
