package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/AyanDas-99/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
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

	app := application{
		errorLog:      errLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{db},
		templateCache: templaceCache,
	}

	srv := http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
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
