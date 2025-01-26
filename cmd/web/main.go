package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AyanDas-99/snippetbox/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	//Define cmd-line flag with name 'addr', with defualt value of ":4000"

	addr := flag.String("addr", ":4000", "Http network address")
	//Data source name for db connection
	dsn := flag.String("dns", "web:password@/snippetbox?parseTime=true", "Mysql data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)

	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new template cache
	templaceCache, err := newTemplateCache()

	if err != nil {
		errLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := application{
		errorLog:       errLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templaceCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	// config to hold non-default TLS settings
	tlsConfig := &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}

	srv := http.Server{
		Addr:         *addr,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	infoLog.Printf("Starting server on %s", *addr)
	// Starting HTTPS server, we pass the tls key and certificate
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
