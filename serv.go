package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
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

func fs(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		return
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	var css, js []string
	var err error

	assets := assetFS()

	css, err = assets.AssetDir("build/static/css")
	if err != nil {
		panic(err)
	}
	js, err = assets.AssetDir("build/static/js")
	if err != nil {
		panic(err)
	}

	tmplData := map[string]interface{}{
		"WSPath":  template.JSStr('/'),
		"ApiRoot": template.JSStr(fmt.Sprintf("http://%s/api", r.Host)),

		"Assets": map[string]interface{}{
			"Js":  js,
			"Css": css,
		},
	}

	tmplContent := `
<!DOCTYPE HTML>
<html>
  <head>
    <meta http-equiv='content-type' content='text/html; charset=utf-8'>
    <title>zqz.ca</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">

    <link rel='shortcut icon' href='/assets/favicon.ico'/>
    {{- range .Assets.Css }}
    <link rel='stylesheet' media='screen' href='/assets/static/css/{{ . }}'/>
    {{- end }}
  </head>
  <body>
    <script type='text/javascript'>
      window.apiRoot = {{.ApiRoot}};
    </script>
  	<noscript>You need to enable JavaScript to run this app.</noscript><div id="root"></div>

    {{- range .Assets.Js }}
    <script type='text/javascript' src='/assets/static/js/{{.}}'></script>
    {{- end }}
  </body>
</html>`

	t := template.New("App Index Template")
	t, err = t.Parse(tmplContent)
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	err = t.Execute(&output, tmplData)
	if err != nil {
		panic(err)
	}

	render.HTML(w, r, output.String())
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
			filedb.NewDiskPersistence(),
			filedb.NewDBMetaStorage(db),
			filedb.NewDBThumbnailStorage(db),
		),
	)

	s.EnableLogging = true

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/api", s.Router())
	fs(r, "/assets", assetFS())
	r.Get("/*", serveIndex)

	http.ListenAndServe(":3001", r)
}
