package helper

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("a-secret-string"))
var flashStorageName = "flash-messages"

func AddFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	session, err := store.Get(r, flashStorageName)
	if err != nil {
		return err
	}

	session.AddFlash(msg)
	fmt.Println("adding flash", msg)
	return session.Save(r, w)
}

func Flash(ctx context.Context) ([]string, bool) {
	messages, ok := ctx.Value("flashes").([]any)

	if !ok {
		return []string{}, false
	}

	if len(messages) == 0 {
		return []string{}, false
	}

	notices := make([]string, 0)

	for _, m := range messages {
		notices = append(notices, m.(string))
	}

	return notices, true
}

func Flashes(w http.ResponseWriter, r *http.Request) []any {
	session, err := store.Get(r, flashStorageName)
	if err != nil {
		return nil
	}

	flashes := session.Flashes()
	session.Save(r, w)
	return flashes
}
