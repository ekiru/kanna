package models

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/ekiru/kanna/db"
)

type Post struct {
	id        *url.URL
	Type      string
	Audience  string
	Author    *Actor
	Content   string
	Published string
}

func (p *Post) ID() *url.URL {
	return p.id
}

func (p *Post) Types() []string {
	return []string{p.Type}
}

func (p *Post) HasType(t string) bool {
	return p.Type == t
}

func (p *Post) GetProp(name string) (interface{}, bool) {
	switch name {
	case "id":
		return p.id, true
	case "type":
		return p.Type, true
	case "audience":
		return p.Audience, true
	case "author":
		return p.Author, true
	case "content":
		return p.Content, true
	case "published":
		return p.Published, true
	default:
		return nil, false
	}
}

func (p *Post) Props() []string {
	return []string{"id", "type", "audience", "author", "content", "published"}
}

func (post *Post) FromRow(rows *sql.Rows) error {
	post.Author = &Actor{}
	actor := post.Author.Scanners()
	return rows.Scan(
		db.URLScanner{&post.id},
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
