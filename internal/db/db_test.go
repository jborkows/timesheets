package db_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	queries "github.com/jborkows/timesheets/internal/db"
)

func TestMigrations(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer cleanupFunc(tempFile)
	_, err = queries.NewDatabase(tempFile.Name())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Add test logic here, e.g., checking if tables were created
}

func cleanupFunc(tempFile *os.File) {
	cleanup := func() {
		tempFile.Close()
		err := os.Remove(tempFile.Name())
		if err != nil {
			fmt.Println("Error removing temporary file:", err)
		} else {
			fmt.Println("Temporary file removed.")
		}
		removeAdditionalDbFiles(tempFile.Name())
	}
	if r := recover(); r != nil {
		cleanup()
		panic(r)
	} else {
		cleanup()
	}
	cleanup()
}

func removeAdditionalDbFiles(fileName string) {
	removeAdditionalDbFile(fileName, "wal")
	removeAdditionalDbFile(fileName, "shm")
}

func removeAdditionalDbFile(fileName string, suffix string) {

	auxielieryFile := fmt.Sprintf("%s-%s", fileName, suffix)
	if _, err := os.Stat(auxielieryFile); err != nil {
		log.Printf("Error checking %s file: %s", suffix, err)
		return
	}
	if err := os.Remove(auxielieryFile); err != nil {
		log.Printf("Error removing %s file: %s", suffix, err)
	} else {
		log.Printf("%s file %s removed.", suffix, auxielieryFile)
	}
}
