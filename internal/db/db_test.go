package db_test

import (
	"fmt"
	queries "github.com/jborkows/timesheets/internal/db"
	"log"
	"os"
	"testing"
)

func TestMigrations(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	defer func() {
		tempFile.Close()

		err = os.Remove(tempFile.Name())
		if err != nil {
			fmt.Println("Error removing temporary file:", err)
		} else {
			fmt.Println("Temporary file removed.")
		}
	}()

	_, err = queries.NewDatabase(tempFile.Name())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Add test logic here, e.g., checking if tables were created
}
