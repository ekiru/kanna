package actors

import "net/url"

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

func ById(id string) (bool, *Model) {
	if id == "http://kanna.example/actor/srn" {
		return true, &Model{
			Inbox:  exampleInbox,
			Outbox: exampleOutbox,
			Name:   "Surinna",
			Type:   "Person",
		}
	} else {
		return false, nil
	}
}
