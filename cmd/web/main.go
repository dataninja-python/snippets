package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.ajigherighe.net/internal/models"
	"time"
)

// create an application wide struct for logging
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// add better configuration management
	addr := flag.String("addr", ":4000", "Http network address")
	// add database command line flag
	dsn := flag.String("dsn", "web:z@rchN3rd2024@/snippets?parseTime=true", "MySQL data source name")
	flag.Parse()

	// add new logging functionality
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// database function used to clean up code
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// defer closing the database to ensure it closes when main exits
	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Use scs.New() to initialize a new session manager.  Configure MySQL to implement session use.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// And add the session manager to our application dependencies
	// initialize application instance of our struct with the dependencies
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize and use http.Server struct using same network address and routes as before.
	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes log entries at Error level,
		// and assign it to the ErrorLog field.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Print a log a message to say that the server is starting.
	logger.Info("starting server", "addr", srv.Addr)
	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address
	// simplify the original function using the new http.Server struct
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

// OpenDB() function that wraps sql.Open functionality
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
