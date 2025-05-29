package api

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, o any) {
	b, err := json.Marshal(o)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, "failed to created json")
		return
	}

	w.Write(b)
}

func Error(w http.ResponseWriter, m string) {
	e := errorMessage{
		Message: m,
	}

	JSON(w, e)
}
