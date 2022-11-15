package main

import (
    "context"
    "database/sql"
    "flag"
    "fmt"
    "net/http"
    "os"
    "sync"
    "time"

	"fitness.zioncastillo.net/internal/data"
    "fitness.zioncastillo.net/internal/jsonlog"
    "fitness.zioncastillo.net/internal/mailer"
    _ "github.com/lib/pq"
)

// Declare a string containing the application version number. Later in the book we'll 
// generate this automatically at build time, but for now we'll just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the 
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these  
// configuration settings from command-line flags when the application starts.

type config struct {
    port int
    env  string	
    db struct {
        dsn string
		maxOpenConns int
        maxIdleConns int
        maxIdleTime string
    }
    smtp struct {
        host     string
        port     int
        username string
        password string
        sender   string
    }
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a 
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
    config config
    logger *jsonlog.Logger
	models data.Models
    mailer mailer.Mailer
    wg sync.WaitGroup
}

func main() {
    // Declare an instance of the config struct.
    var cfg config
    // Read the value of the port and env command-line flags into the config struct. We
    // default to using the port number 4000 and the environment "development" if no
    // corresponding flags are provided.
    flag.IntVar(&cfg.port, "port", 4000, "API server port")
    flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
    flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("FITNESS_DB_DSN"), "PostgreSQL DSN")
    flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connection")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

    flag.StringVar(&cfg.smtp.host, "smpt-host", "smtp.mailtrap.io", "SMPT host")
    flag.IntVar(&cfg.smtp.port, "smpt-port", 25, "SMPT port")
    flag.StringVar(&cfg.smtp.username, "smpt-username", "2003542b049ac5", "SMPT username")
    flag.StringVar(&cfg.smtp.password, "smpt-password", "641e59bb93068e", "SMPT Password")
    flag.StringVar(&cfg.smtp.sender, "smpt-sender", "BIO <no-reply@fitness.zioncastillo.net>", "SMPT Sender")

	flag.Parse()
    // Initialize a new logger which writes messages to the standard out stream, 
    // prefixed with the current date and time.
    logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

    // Create a connection pool
    db, err := openDB(cfg)
    if err != nil {
        logger.PrintFatal(err, nil)
    }
    defer db.Close()
	logger.PrintInfo("datatbase connection pool established", nil)
    // Declare an instance of the application struct, containing the config struct and 
    // the logger.
    app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
        mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

    // Declare a HTTP server with some sensible timeout settings, which listens on the
    // port provided in the config struct and uses the servemux we created above as the 
    // handler.
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.port),
        Handler:      app.routes(),
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    // Start the HTTP server.
    logger.PrintInfo("starting server on port 4000", nil)
    err = srv.ListenAndServe()
    logger.PrintFatal(err, nil)
}

// Open DB function to return a *sql.DB connection pool
func openDB(cfg config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	// Create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil

} 