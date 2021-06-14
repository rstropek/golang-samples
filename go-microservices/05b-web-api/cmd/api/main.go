package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	cr "github.com/rstropek/golang-samples/go-microservices/05b-web-api/internal/customerrepository"
	"github.com/urfave/negroni"
)

// Struct holding configuration settings.
type config struct {
	port int
	env  string
}

// Struct holding dependencies for our HTTP handlers, helpers, and middleware.
type application struct {
	config     config
	logger     *log.Logger
	repository *cr.CustomerRepository
}

func main() {
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Create customer repository
	repo := cr.NewCustomerRepository()

	// Declare an instance of the application struct.
	app := &application{
		config:     cfg,
		logger:     logger,
		repository: &repo,
	}

	// Add negroni middleware for logging, static files, and recovery.
	n := negroni.Classic()
	n.UseHandler(app.routes())
	n.Use(cors.AllowAll())

	// Declare a HTTP server.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      n,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
