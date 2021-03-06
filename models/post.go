package models

import (
	"context"
	"database/sql"

	"github.com/ekiru/kanna/db"
)

func (post *Post) FromRow(rows *sql.Rows) error {
	post.Author = &Actor{}
	actor := post.Author.Scanners()
	return rows.Scan(
		db.URLScanner{&post.id},
		&post.typ,
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

func PostsByActor(ctx context.Context, actor *Actor) ([]*Post, error) {
	var posts []*Post
	rows, err := db.DB(ctx).QueryContext(ctx,
		"select post.id, post.type, post.audience, post.content, post.published, "+
			"post.authorId, act.type, act.name, act.inbox, act.outbox "+
			"from Posts post join Actors act on post.authorId = act.id "+
			"where act.id = ?",
		actor.ID().String())
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post Post
		if err = post.FromRow(rows); err != nil {
			return posts, err
		}
		posts = append(posts, &post)
	}
	return posts, rows.Err()
}
