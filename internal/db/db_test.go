package db_test

import (
	"bytes"
	"context"
	"io"
	dbp "github.com/jborkows/timesheets/internal/db"
	queries "github.com/jborkows/timesheets/internal/db"
	"github.com/jborkows/timesheets/internal/model"
	"log"
	"os"
	"testing"
	"time"
)

func TestMigrations(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer dbp.RemoveDatabase(tempFile.Name())
	_, err = queries.NewDatabase(tempFile.Name())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Add test logic here, e.g., checking if tables were created
}

func useDb(t *testing.T, test func(saver model.Saver, querier model.Queryer)) {
	t.Parallel()
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	log.Printf("Using temporary file: %s", tempFile.Name())
	defer dbp.RemoveDatabase(tempFile.Name())
	db, err := dbp.NewDatabase(tempFile.Name())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	support := dbp.NewTransactionSupport(db)
	err = support.WithTransaction(context.Background(), func(ctx context.Context, q *dbp.Queries) error {
		repository := dbp.Repository(q, func(model.CategoryType) bool { return false })
		test(repository, repository)
		return nil
	})
	if err != nil {
		t.Errorf("Error in transaction: %v", err)
	}

}

// TestSaveDoesNotCorruptStdout ensures that saving a timesheet entry
// does not write anything to stdout, which would corrupt LSP protocol communication.
func TestSaveDoesNotCorruptStdout(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "testdb-stdout-*.db")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer dbp.RemoveDatabase(tempFile.Name())

	db, err := dbp.NewDatabase(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	support := dbp.NewTransactionSupport(db)
	err = support.WithTransaction(context.Background(), func(ctx context.Context, q *dbp.Queries) error {
		repository := dbp.Repository(q, func(model.CategoryType) bool { return false })

		// Create and save a timesheet with entries
		timesheet := model.TimesheetForDate(time.Now())
		task := "work"
		entry := model.TimesheetEntry{Hours: 4, Minutes: 0, Category: "work", Comment: "work", Task: &task}
		if err := timesheet.AddEntry(entry); err != nil {
			t.Errorf("Error adding entry: %v", err)
			return err
		}

		// Save - this is where the stdout corruption was happening
		if err := repository.Save(ctx, timesheet); err != nil {
			t.Errorf("Error saving timesheet: %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		t.Errorf("Error in transaction: %v", err)
	}

	// Restore stdout and capture output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify no stdout output that would corrupt LSP protocol
	if output != "" {
		t.Errorf("Save operation wrote to stdout, which corrupts LSP protocol. Output: %q", output)
	}
}
