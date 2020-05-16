package filedb

import "time"

type Meta struct {
	ID            int       `json:"-"`
	Alias         string    `json:"alias,omitempty"`
	Name          string    `json:"name"`
	Hash          string    `json:"hash"`
	Slug          string    `json:"slug"`
	ContentType   string    `json:"type"`
	Path          string    `json:"path,omitempty"`
	Size          int       `json:"size"`
	Date          time.Time `json:"date"`
	BytesReceived int       `json:"bytes_received,omitempty"`

	Thumbnail string `json:"thumbnail,omitempty"`
}

func (m Meta) finished() bool {
	return m.Size == m.BytesReceived
}
