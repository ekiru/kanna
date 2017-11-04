package activitystreams

import (
	"encoding/json"
	"net/url"
)

// ContentType is the defined MIME content-type for Activity Streams
// and ActivityPub.
const ContentType = `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`

// A Link represent a Link as defined in the Activity Streams
// specification (https://www.w3.org/TR/activitystreams-core/#link)
type Link struct {
	Href  *url.URL
	Type  string
	Props map[string]interface{}
}

// Marshal serializes an object as an Activity Stream.
func Marshal(obj Object) ([]byte, error) {
	ser, err := serializeValue(obj)
	if err != nil {
		return nil, err
	}
	serMap := ser.(map[string]interface{})
	serMap["@context"] = "https://www.w3.org/ns/activitystreams"
	return json.Marshal(serMap)
}

func serializeValue(val interface{}) (interface{}, error) {
	switch val := val.(type) {
	case string:
		return val, nil
	case *url.URL:
		return val.String(), nil
	case []interface{}:
		vals := make([]interface{}, 0, len(val))
		for elem := range val {
			ser, err := serializeValue(elem)
			if err != nil {
				return nil, err
			}
			vals = append(vals, ser)
		}
		return vals, nil
	case Object:
		ser := map[string]interface{}{
			"id": val.ID().String(),
		}
		if types := val.Types(); len(types) == 1 {
			ser["type"] = types[0]
		} else {
			ser["type"] = types
		}
		for _, name := range val.Props() {
			propVal, _ := val.GetProp(name)
			serVal, err := serializeValue(propVal)
			if err != nil {
				return nil, err
			}
			ser[name] = serVal
		}
		return ser, nil
	case *Link:
		ser := map[string]interface{}{
			"type": val.Type,
			"href": val.Href.String(),
		}
		for name, val := range val.Props {
			serVal, err := serializeValue(val)
			if err != nil {
				return nil, err
			}
			ser[name] = serVal
		}
		return ser, nil
	default:
		panic("unrecognized value type")
	}
}
