package integrationtests

import (
	"log"
	"os"
	"testing"

	mydb "github.com/jborkows/timesheets/internal/db"
	"github.com/jborkows/timesheets/internal/model"
)

func TestStartWorkspace(t *testing.T) {
	config := model.NewConfig([]string{"aaa", "bbb", "ccc"})
	useWorkspace(config, func(service *model.Service) {
		// Add test logic here, e.g., creating a new project
	})
}

func useWorkspace(config *model.Config, usage func(service *model.Service)) {

	repository, cleanup := initDB(config)
	defer cleanup()

	root, err := os.CreateTemp("", "projectroot-*.txt")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer root.Close()
	defer os.Remove(root.Name())
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	service := model.NewService(root.Name(), config, repository)
	usage(service)
}

type cleanupFunction func()

func initDB(config *model.Config) (model.Repository, cleanupFunction) {
	dbPath, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	db, err := mydb.NewDatabase(dbPath.Name())
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	repository := mydb.CreateRepository(db, config)
	return repository, func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing db: %s", err)
		}
		mydb.RemoveDatabase(dbPath.Name())
	}
}
