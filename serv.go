package main

import (
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/zqz/upl/filedb"
)

var con *sqlx.DB
var db filedb.FileDB

func connect(str string) (*sqlx.DB, error) {
	if len(str) == 0 {
		return nil, errors.New("Empty DB string")
	}

	var err error
	if parsedURL, err := pq.ParseURL(str); err == nil && parsedURL != "" {
		str = parsedURL
	}

	var con *sqlx.DB
	if con, err = sqlx.Connect("postgres", str); err != nil {
		return nil, err
	}

	if err = con.Ping(); err != nil {
		return nil, err
	}

	return con, nil
}

func main() {
	// os.Mkdir(tmpPath, 0744)
	// os.Mkdir(finalPath, 0744)

	s := filedb.NewServer(
		filedb.NewFileDB(
			filedb.NewMemoryPersistence(),
			filedb.NewMemoryMetaStorage(),
		),
	)

	s.EnableLogging = true

	http.ListenAndServe(":3001", s.Router())
}
