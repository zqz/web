package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/zqz/upl/filedb"
)

type Server struct {
	config   config
	database *sql.DB
	logger   *log.Logger
}

func (s Server) Log(x ...interface{}) {
	s.logger.Println(x...)
}

func Init(path string) (Server, error) {
	s := Server{}
	s.logger = log.New(os.Stdout, "", log.LstdFlags)

	cfg, err := parseConfig(path)
	if err != nil {
		return s, err
	}
	s.logger.Println("Parsed Config")

	db, err := cfg.DBConfig.loadDatabase()
	if err != nil {
		return s, err
	}
	s.logger.Println("Connected to DB")

	s.database = db
	s.config = cfg

	return s, nil
}

func (s Server) Close() {
	s.database.Close()
}

func (s Server) runInsecure(r http.Handler) error {
	listenPort := fmt.Sprintf(":%d", s.config.Port)

	s.logger.Println("[server] listening for HTTP traffic on port", listenPort)

	return http.ListenAndServe(listenPort, r)
}

func instrumentNewRelic(app *newrelic.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			txn := app.StartTransaction(r.URL.Path)
			defer txn.End()
			txn.SetWebRequestHTTP(r)
			next.ServeHTTP(txn.SetWebResponse(w), r)
		}
		return http.HandlerFunc(mw)
	}
}

func (s Server) Run() error {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("zqz.ca"),
		newrelic.ConfigLicense("eu01xx5c9482d2075bd8fb489adfa40d4b2aNRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if err != nil {
		s.logger.Println("failed to start newrelic", err.Error())
	} else {
		s.logger.Println("looks like new relic started")
	}

	fdb := filedb.NewServer(
		filedb.NewFileDB(
			filedb.NewDiskPersistence(),
			filedb.NewDBMetaStorage(s.database),
		),
	)

	r := chi.NewRouter()

	logger := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: s.logger})
	r.Use(logger)

	r.Mount("/api", fdb.Router())

	// ra.Get("/{slug}", fdb.GetDataWithSlug)

	s.logger.Println("Listening for web traffic")

	mux := http.NewServeMux()
	mux.HandleFunc(
		newrelic.WrapHandleFunc(
			app,
			"/",
			func(w http.ResponseWriter, rx *http.Request) {
				r.ServeHTTP(w, rx)
			},
		),
	)

	return s.run(mux)
}

func (s Server) run(r http.Handler) error {
	return s.runInsecure(r)
}
