package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sathishs-dev/pismo-transactions/internal/meta/signal"
	"github.com/sathishs-dev/pismo-transactions/pkg/handler"
	"github.com/sathishs-dev/pismo-transactions/pkg/repository"
)

type env struct {
	PismoDBDSN      string        `envconfig:"PISMO_DB_DSN" required:"true"`
	Loglevel        string        `envconfig:"LOG_LEVEL" default:"info"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
}

func main() {
	var conf env

	failOnError(loadConfig(&conf), "failed to load the env config")

	configureLogger(conf.Loglevel)

	ctx := signal.NewWithContext(context.Background())

	db, err := sql.Open("postgres", conf.PismoDBDSN)
	failOnError(err, "failed to open database connection")
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("closing db connection failed")
		}
	}()

	db.SetMaxOpenConns(3)

	dbx := sqlx.NewDb(db, "postgres")

	repo := repository.NewPismoRepo(dbx)

	h := handler.NewHandler(repo)

	webServer := initWebServer(log.Output(os.Stderr), h)
	go func() {
		if err := webServer.Start(); err != nil {
			logOnError(err, "failed to start webserver")
			signal.Shutdown()
		}
	}()

	signal.Add(func() {
		shutdownCtx, shutDownCanecl := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
		defer shutDownCanecl()
		logOnError(webServer.Stop(shutdownCtx), "failed to stop webserver")
	})

	<-ctx.Done()
}

// configureLogger sets the globalLogLevel for zerolog instance
func configureLogger(level string) {
	zerolog.SetGlobalLevel(getLoglevel(level))
}

// loadConfig wrapper method for envconfig.Process, it receives config var and loads the env to it
func loadConfig(cfg *env) error {
	return envconfig.Process("", cfg)
}

// logOnError receives error and message, if there is an error it logs
func logOnError(err error, msg string) {
	if err != nil {
		log.Error().Err(err).Msg(msg)
	}
}

// failOnError receives error and message, if there is an error it logs the error and exits the execution
func failOnError(err error, msg string) {
	if err == nil {
		return
	}

	logOnError(err, msg)
	os.Exit(1)
}
