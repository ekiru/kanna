package models

import (
	"net/url"
)

type Post struct {
	id        *url.URL
	typ       string
	Published string
	Audience  string
	Author    *Actor
	Content   string
}

func (model *Post) ID() *url.URL {
	return model.id
}

func (model *Post) Types() []string {
	return []string{model.typ}
}

func (model *Post) HasType(t string) bool {
	return t == model.typ
}

func (model *Post) Props() []string {
	return []string{"id", "type", "published", "audience", "author", "content"}
}

func (model *Post) GetProp(prop string) (interface{}, bool) {
	switch prop {
	case "id":
		return model.id, true
	case "type":
		return model.typ, true
	case "published":
		return model.Published, true
	case "audience":
		return model.Audience, true
	case "author":
		return model.Author, true
	case "content":
		return model.Content, true
	default:
		return nil, false
	}
}
