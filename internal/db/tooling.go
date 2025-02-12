package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for loading migrations from files
	_ "github.com/mattn/go-sqlite3"                      // SQLite driver
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"embed"

	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func NewDatabase(filePath string) (*sql.DB, error) {

	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=ON&_cache_size=2000&_busy_timeout=5000", filePath)
	db, err := otelsql.Open("sqlite3", dsn,
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName("mydb"))
	if err != nil {
		log.Printf("Error opening database: %v\n", err)
		return nil, fmt.Errorf("While opening %w", err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(10)                  // Max open connections
	db.SetMaxIdleConns(5)                   // Max idle connections
	db.SetConnMaxLifetime(30 * time.Minute) // Max connection lifetime

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("While pinging %w", err)
	}
	log.Println("Database connected successfully.")
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("While running migrations %w", err)
	}
	if err := optimize(db); err != nil {
		return nil, fmt.Errorf("While optimizing %w", err)
	}
	row := db.QueryRow("SELECT sqlite_version() as version")
	var version string
	err = row.Scan(&version)
	if err != nil {
		log.Printf("Failed to get SQLite version: %v", err)
	}
	log.Printf("SQLite version: %s", version)

	return db, nil
}

//go:embed schema/migrations/*.sql
var migrationFiles embed.FS

func runMigrations(db *sql.DB) error {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return err
	}
	fmt.Println("Working directory:", wd)
	d, err := iofs.New(migrationFiles, "schema/migrations")
	if err != nil {
		log.Fatalf("Failed to initialize migration source: %v", err)
	}
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	version, dirty, error := m.Version()
	if error != nil {
		log.Printf("Failed to get migration version: %v", error)
	}
	log.Printf("Current migration version: %d, dirty: %t", version, dirty)
	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	version, dirty, error = m.Version()
	if error != nil {
		log.Fatalf("Failed to get migration version: %v", error)
	}
	log.Printf("After migration version: %d, dirty: %t", version, dirty)

	return nil
}

func WithinTransaction(db *sql.DB, code func(tx *Queries) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	executor := New(db)
	txExecutor := executor.WithTx(tx)
	err = code(txExecutor)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return fmt.Errorf("While rolling back %w", err)
		} else {
			return fmt.Errorf("While running transaction %w", err)
		}
	}
	err = tx.Commit()
	return err
}

func optimize(database *sql.DB) error {
	return WithinTransaction(database, func(tx *Queries) error {
		_, err := tx.db.ExecContext(context.Background(), "PRAGMA optimize")
		return err
	})

}
