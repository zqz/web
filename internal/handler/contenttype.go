package handler

import "net/http"

const (
	ContentTypeHTML = "text/html; charset=utf-8"
	ContentTypeJSON = "application/json"
)

// SetContentType sets the response Content-Type header.
func SetContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}
