package db_test

import (
	"context"
	dbp "github.com/jborkows/timesheets/internal/db"
	queries "github.com/jborkows/timesheets/internal/db"
	"github.com/jborkows/timesheets/internal/model"
	"log"
	"os"
	"testing"
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
