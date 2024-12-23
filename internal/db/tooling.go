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

type Database struct {
	*sql.DB
}

func NewDatabase(filePath string) (*Database, error) {

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
	runMigrations(db)
	database := &Database{db}
	if err := database.optimize(); err != nil {
		return nil, fmt.Errorf("While optimizing %w", err)
	}
	return database, nil
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

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	return nil
}

func (d *Database) optimize() error {
	return Execute(context.Background(), d, func(tx *sql.Tx) error {
		_, err := tx.Exec("PRAGMA optimize")
		return err
	})
}

func (d *Database) Close() error {
	return d.DB.Close()
}
func Execute(ctx context.Context, db *Database, exec func(*sql.Tx) error) error {
	_, err := RunTransaction(ctx, db, func(tx *sql.Tx) (*interface{}, error) {
		return nil, exec(tx)
	})
	return err
}

func RunTransaction[T any](ctx context.Context, db *Database, exec func(*sql.Tx) (*T, error)) (*T, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				panic(fmt.Errorf("failed to rollback transaction: %w", err))
			}
			panic(r) // Re-panic after rollback
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				panic(fmt.Errorf("failed to rollback transaction: %w", err))
			}
		}
	}()

	value, err := exec(tx)
	if err != nil {
		return nil, fmt.Errorf("Executing exec %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return value, nil
}
