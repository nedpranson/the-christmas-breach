package db

import (
	"database/sql"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	migrate_sqlite "github.com/golang-migrate/migrate/v4/database/sqlite3"
	migrations "github.com/nedpranson/the-christmas-breach/db"
)

const Path = "./data/app.db"

func Open() (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(Path), 0o755); err != nil {
		return nil, err
	}

	dsn := "file:" + Path + "?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=1"
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err 
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func Migrate(conn *sql.DB) error {
	src, err := iofs.New(migrations.FS, "migrations")
	if err != nil {
		return err
	}

	drv, err := migrate_sqlite .WithInstance(conn, &migrate_sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite3", drv)
	if err != nil {
		return err; 
	}

	// Note: no m.Close(). It would close both the source and the database
	// driver, and the latter is our shared *sql.DB.

	if err := m.Up(); err != nil && !noMigration(err) {
		return err
	}

	return nil
}

func noMigration(err error) bool {
	return errors.Is(err, migrate.ErrNoChange) || errors.Is(err, fs.ErrNotExist)
}
