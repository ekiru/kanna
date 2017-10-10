package db

import (
	"errors"
	"net/url"
)

// A URLScanner wraps the type **url.URL to allow scanning string and
// binary data from a SQL database as *url.URL values.
type URLScanner struct {
	URL **url.URL
}

func (u URLScanner) Scan(src interface{}) error {
	var err error
	if src == nil {
		*u.URL = nil
		return nil
	}
	var urlstr string
	switch src := src.(type) {
	case string:
		urlstr = src
	case []byte:
		urlstr = string(src)
	default:
		return errors.New("can only Scan into URLScanner from strings")
	}
	*u.URL, err = url.Parse(urlstr)
	return err
}
