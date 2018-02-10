package filedb

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"time"
)

func ParseMeta(r io.Reader) (*Meta, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	m := Meta{}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return &m, err
}

type Meta struct {
	Alias         string    `json:"alias"`
	Name          string    `json:"name"`
	Hash          string    `json:"hash"`
	Slug          string    `json:"slug"`
	ContentType   string    `json:"type"`
	Path          string    `json:"path"`
	Size          int       `json:"size"`
	Date          time.Time `json:"date"`
	BytesReceived int       `json:"bytes_received"`
}
