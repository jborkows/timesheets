package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jborkows/timesheets/internal/logs"
	"github.com/jborkows/timesheets/internal/lspserver"
	"github.com/jborkows/timesheets/internal/model"
	"github.com/jborkows/timesheets/internal/rpc"
)

func main() {
	logger, err := logs.Initialize(logs.FileLogger("timesheets.log"))
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	var versionFlag = flag.Bool("version", false, "Print version")
	var configFlag = flag.String("c", "", "Path to config file")
	var reloadFlag = flag.Bool("lsptesting", false, "For air and lsp testing")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("Version: %s", model.Version)
		return
	}
	if *reloadFlag {
		waitForTerminationSignal()
	}

	if *configFlag == "" {
		log.Fatal("Config file is required")
	}

	file, err := os.Open(*configFlag)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	config, err := model.ReadConfig(file)

	log.Printf("Config: %+v", config)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	writer := os.Stdout
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			log.Printf("Got an error: %s", err)
			continue
		}

		handleMessage(writer, method, contents)
	}
}

func waitForTerminationSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Waiting for SIGINT or SIGTERM...")
	sig := <-sigChan
	fmt.Printf("Received signal: %v\n", sig)
}

func handleMessage(writer io.Writer, method string, contents []byte) {
	response, err := lspserver.Route(method, contents)
	if err != nil {
		log.Printf("Got an error: %s", err)
		return
	}
	if response != nil {
		log.Printf("Sending response for %s", method)
		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))

}
