package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/zqz/upl/filedb"
)

var con *sql.DB
var db filedb.FileDB

func connect() (*sql.DB, error) {
	host := "localhost"
	port := 5432
	user := "dylan"
	dbname := "zqz2-dev"

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname,
	)

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// defer db.Close()

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// os.Mkdir(tmpPath, 0744)
	// os.Mkdir(finalPath, 0744)

	db, err := connect()
	if err != nil {
		fmt.Println("welp, db null, end game", err.Error())
		return
	}

	s := filedb.NewServer(
		filedb.NewFileDB(
			filedb.NewMemoryPersistence(),
			filedb.NewDBMetaStorage(db),
		),
	)

	s.EnableLogging = true

	http.ListenAndServe(":3001", s.Router())
}
