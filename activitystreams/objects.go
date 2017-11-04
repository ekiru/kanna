package activitystreams

import "net/url"

// An Object represents an Object as defined in the Activity Streams
// specification (https://www.w3.org/TR/activitystreams-core/#object)
type Object interface {
	ID() *url.URL
	Types() []string
	HasType(string) bool
	GetProp(string) (interface{}, bool)
	Props() []string
}
