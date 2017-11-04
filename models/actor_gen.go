
package models

import (
	"net/url"
)

type Actor struct {
	id *url.URL
	typ string
	Inbox *url.URL
	Outbox *url.URL
	Name string
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
	return []string{ "id", "type", "inbox","outbox","name", }
}

func (model *Actor) GetProp(prop string) (interface{}, bool) {
	switch prop {
	case "id":
		return model.id, true
	case "type":
		return model.typ, true
	case "inbox":
		return model.Inbox, true
	case "outbox":
		return model.Outbox, true
	case "name":
		return model.Name, true
	default:
		return nil, false
	}
}
