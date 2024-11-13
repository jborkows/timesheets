package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"time"
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
	return &Database{db}, nil
}

func (d *Database) Optimize() error {
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
