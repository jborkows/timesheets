package main

import (
	"log"

	"github.com/jborkows/timesheets/internal/example"
	"github.com/jborkows/timesheets/internal/logs"
)

func main() {
	logger, err := logs.Initialize(logs.FileLogger("timesheets.log"))
	if err != nil {
		panic(err)
	}
	log.Println("Starting timesheets")
	log.Println("1 + 2 = ", example.Example(1, 2))
	defer logger.Close()

}
