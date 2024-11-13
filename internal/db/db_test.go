package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for loading migrations from files
	queries "github.com/jborkows/timesheets/internal/db"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"

	"github.com/jborkows/timesheets/internal/logs"
)

func runMigrations(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../db/migrations", // Path to your migration files
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func TestConnectionString(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	tempFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	db, err := queries.NewDatabase(tempFile.Name())

	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer db.Close()

	tracerClean := logs.Tracer("main")
	defer tracerClean()
	tracer := otel.Tracer("quering")

	ctx := context.Background()
	// Start a new span for this function
	_, span := tracer.Start(ctx, "doWork-span")
	defer span.End()
	err = queries.Execute(ctx, db, func(tx *sql.Tx) error {
		var journalMode, foreignKeys string
		var cacheSize int
		if err = tx.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
			return err
		}
		if err = tx.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
			return err
		}
		if err = tx.QueryRow("PRAGMA cache_size").Scan(&cacheSize); err != nil {
			return err
		}

		fmt.Printf("Journal mode: %s\n", journalMode)
		fmt.Printf("Foreign keys enabled: %s\n", foreignKeys)
		fmt.Printf("Cache size: %d pages\n", cacheSize)
		assert.Equal(t, "wal", journalMode)
		assert.Equal(t, "1", foreignKeys)
		assert.Equal(t, 2000, cacheSize)
		return nil
	})

	if err != nil {
		log.Fatalf("Quering database: %v\n", err)
	}
	err = db.Optimize()

	if err != nil {
		log.Fatalf("Optimizing database: %v\n", err)
	}

}

func TestMigrations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:") // Use in-memory DB for testing
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	// Add test logic here, e.g., checking if tables were created
}

func TestMigrations2(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	queries := queries.New(db)
	apps, error := queries.ListApps(context.Background())
	// Example test: check if table exists
	if error != nil {
		t.Fatalf("table 'app' does not exist: %v", err)
	}
	if len(apps) == 0 {
		t.Fatal("not found any ")
	}
}
