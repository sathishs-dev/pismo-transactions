package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

const (
	driver = "postgres"
)

type config struct {
	DSN         string `envconfig:"DB_DSN" required:"true"`
	ScriptsPath string `envconfig:"SCRIPTS_PATH" required:"true"`
	DBName      string `envconfig:"DB_NAME" required:"true"`
}

func main() {
	var cfg config
	failOnError(loadConfig(&cfg), "failed to load config")

	db, err := sql.Open(driver, cfg.DSN)
	failOnError(err, "failed to open database connection")
	defer func() {
		if err := db.Close(); err != nil {
			logOnError(err, "closing db connection failed")
		}
	}()

	db.SetMaxOpenConns(1)

	err = waitDB(5*time.Minute, db)
	failOnError(err, "database unreachable")

	fPath, err := os.ReadDir(cfg.ScriptsPath)
	failOnError(err, "unable to find directory")

	for _, file := range fPath {
		if file.IsDir() || !strings.Contains(file.Name(), ".sql") {
			log.Warn().Msgf("migration file not found, expected .sql file")
			continue
		}

		filePath := path.Join("/migrations", file.Name())

		in, err := os.ReadFile(filePath)
		failOnError(err, "reading file failed")

		out := os.ExpandEnv(string(in))
		failOnError(os.WriteFile(filePath, []byte(out), 0), "writing file failed")
	}

	d, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable:  fmt.Sprintf("%s_migrations", cfg.DBName),
		DatabaseName:     cfg.DBName,
		StatementTimeout: time.Minute * 10,
	})
	failOnError(err, "creating db instance failed")

	failOnError(Run(strings.Join([]string{"file://", "/migrations"}, ""), cfg.DBName, d), "running migrations failed")
}

func Run(sourceURL, dbName string, dbDriver database.Driver) error {
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		dbName,
		dbDriver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}

func waitDB(duration time.Duration, db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for {
		select {
		// validates the context timeout
		case <-ctx.Done():
			return errors.New("database wait timeout")
		default:
		}

		err := db.PingContext(ctx)
		if err == nil {
			break
		}

		logOnError(err, "database not reachable")

		// wait before retrying
		time.Sleep(5 * time.Second)
	}

	return nil
}

func loadConfig(cfg *config) error {
	return envconfig.Process("", cfg)
}

func logOnError(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}
func failOnError(err error, msg string) {
	if err == nil {
		return
	}

	logOnError(err, msg)
	os.Exit(1)
}
