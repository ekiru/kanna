package models

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/ekiru/kanna/db"
)

type Post struct {
	ID        *url.URL
	Type      string
	Audience  string
	Author    *Actor
	Content   string
	Published string
}

func (post *Post) FromRow(rows *sql.Rows) error {
	post.Author = &Actor{}
	actor := post.Author.Scanners()
	return rows.Scan(
		db.URLScanner{&post.ID},
		&post.Type,
		&post.Audience,
		&post.Content,
		&post.Published,
		actor["id"],
		actor["type"],
		actor["name"],
		actor["inbox"],
		actor["outbox"],
	)
}

func PostById(ctx context.Context, id string) (*Post, error) {
	var post Post
	rows, err := db.DB(ctx).QueryContext(ctx,
		"select post.id, post.type, post.audience, post.content, post.published, "+
			"post.authorId, act.type, act.name, act.inbox, act.outbox "+
			"from Posts post join Actors act on post.authorId = act.id "+
			"where post.id = ?",
		id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	if err = post.FromRow(rows); err != nil {
		return nil, err
	}
	return &post, nil
}
