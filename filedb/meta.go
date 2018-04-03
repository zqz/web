package filedb

import "time"

type Meta struct {
	ID            int       `json:"-"`
	Alias         string    `json:"alias"`
	Name          string    `json:"name"`
	Hash          string    `json:"hash"`
	Slug          string    `json:"slug"`
	ContentType   string    `json:"type"`
	Path          string    `json:"path"`
	Size          int       `json:"size"`
	Date          time.Time `json:"date"`
	BytesReceived int       `json:"bytes_received"`

	Thumbnail string `json:"thumbnail,omitempty"`
}

func (m Meta) finished() bool {
	return m.Size == m.BytesReceived
}
