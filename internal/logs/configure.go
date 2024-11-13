package logs

import (
	"io"
	"log"
	"os"
	"strings"
)

type configuration struct {
	writers []io.Writer
	cleanup func() error
}

// Write writes data to each writer and returns the number of bytes written to the first writer and any error encountered
func (t *configuration) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
	}
	return
}

type Logger func(*configuration) error

func FileLogger(outputFile string) Logger {
	return func(multiWriter *configuration) error {
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		multiWriter.writers = append(multiWriter.writers, file)
		multiWriter.cleanup = func() error {

			error := file.Close()
			if error != nil {
				log.Println("Failed to close log file: ", error)
			}
			return nil
		}
		return nil
	}
}

func StdoutLogger() Logger {
	return func(multiWriter *configuration) error {
		multiWriter.writers = append(multiWriter.writers, os.Stdout)
		return nil
	}
}

func Initialize(options ...Logger) (io.Closer, error) {
	logger := &configuration{}
	for _, option := range options {
		err := option(logger)
		if err != nil {
			return nil, err
		}
	}
	if len(logger.writers) == 0 {
		logger.writers = append(logger.writers, os.Stdout)
	}

	log.SetOutput(logger)
	log.Println(strings.Repeat("#", 80))
	log.Println("Starting logging...")
	return logger, nil
}

func (t *configuration) Close() error {
	log.Println("Closing logging...")
	log.Println(strings.Repeat("#", 80))
	if t.cleanup != nil {
		return t.cleanup()
	}
	return nil
}
