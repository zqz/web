package file

import (
	"time"

	"github.com/zqz/web/backend/internal/models"
)

type File struct {
	models.File

	BytesReceived int
	Thumbnail     string `json:"thumbnail,omitempty"`
}

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
	Private       bool      `json:"private"`
	Comment       string    `json:"comment"`
	UserID        int       `json:"-"`

	Thumbnail string `json:"thumbnail,omitempty"`
}

func (m File) Finished() bool {
	return m.Size == m.BytesReceived
}
