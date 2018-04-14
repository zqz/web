package server

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zqz/upl/filedb"
)

type Server struct {
	config   config
	database *sql.DB
	logger   *log.Logger
}

func Init(path string, l *log.Logger) (Server, error) {
	s := Server{}

	cfg, err := parseConfig(path)
	if err != nil {
		return s, err
	}

	l.Println("Parsed Config")

	db, err := cfg.DBConfig.loadDatabase()
	if err != nil {
		return s, err
	}

	l.Println("Connected to DB")

	s.database = db
	s.config = cfg
	s.logger = l

	return s, nil
}

func (s Server) Close() {
	s.database.Close()
}

func (s Server) runInsecure(r chi.Router) error {
	listenPort := fmt.Sprintf(":%d", s.config.Port)

	s.logger.Println("[server] listening for HTTP traffic on port", listenPort)

	return http.ListenAndServe(listenPort, r)
}

func (s Server) runSecure(r chi.Router) error {
	c := autocert.DirCache("./")
	m := autocert.Manager{
		Cache:  c,
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(
			"zqz.ca",
		),
	}

	tlsPort := fmt.Sprintf(":%d", s.config.TLSPort)

	h := &http.Server{
		Addr:      tlsPort,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   r,
		ErrorLog:  s.logger,
	}

	listenPort := fmt.Sprintf(":%d", s.config.Port)
	go http.ListenAndServe(listenPort, m.HTTPHandler(s.secureRedirect()))

	s.logger.Println("[server] listening for TLS traffic on port", s.config.TLSPort)
	s.logger.Println("[server] redirecting HTTP traffic on port", s.config.Port, "to HTTPS")

	// start https server
	return h.ListenAndServeTLS("", "")
}

func (s Server) Run() error {
	db := s.database

	fdb := filedb.NewServer(
		filedb.NewFileDB(
			filedb.NewDiskPersistence(),
			filedb.NewDBMetaStorage(db),
			filedb.NewDBThumbnailStorage(db),
		),
	)
	// fdb.SetLogger(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/api", fdb.Router())
	r.Get("/*", serveIndex)
	serveAssets(r)

	s.logger.Println("Listening for web traffic")

	return s.run(r)
}

func (s Server) run(r chi.Router) error {
	if s.config.Secure {
		return s.runSecure(r)
	} else {
		return s.runInsecure(r)
	}
}

func (s Server) secureRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		redir := "https://" + req.Host + req.RequestURI

		s.logger.Println("[server] redirected request to", redir)
		http.Redirect(w, req, redir, http.StatusMovedPermanently)
	}
}
