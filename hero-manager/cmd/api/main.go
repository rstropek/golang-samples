package main

import (
	"context"
	"database/sql"

	"flag"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"heroes.rainerstropek.com/internal/data"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	azure struct {
		tenantId string
	}
}

type application struct {
	config config
	logger *zerolog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("HEROES_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.StringVar(&cfg.azure.tenantId, "azure-tenant", os.Getenv("AZURE_TENANT"), "AAD Tenant")
	flag.Parse()

	//logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal().Err(err)
	}

	defer db.Close()
	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: &logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	logger.Fatal().Err(err)
}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
