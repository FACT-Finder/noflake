package database

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrations embed.FS

// New creates a db instance.
func New(path string) *sqlx.DB {
	err := createDirectory(path)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't create database directory")
	}

	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't open database")
	}

	// Only use one connection otherwise we would have to handle "database is locked" errors.
	// See https://github.com/mattn/go-sqlite3/issues/274
	db.SetMaxOpenConns(1)
	db.MustExec("PRAGMA foreign_keys = ON")

	log.Debug().Msg("Initializing database")

	err = initDB(db)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't migrate database")
	}

	log.Debug().Msg("Database initialized")
	return db
}

func createDirectory(path string) error {
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
			return err
		}
	}
	return nil
}

func initDB(db *sqlx.DB) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite3", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}
