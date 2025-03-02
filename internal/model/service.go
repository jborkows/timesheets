package model

import (
	"log"
	"strings"
	"time"
)

type Service struct {
	projectRoot string
	config      *Config
	parser      *Parser
}

func NewService(projectRoot string, config *Config) *Service {
	parser := Parser{
		HolidayClassifier: func(a *DateInfo) bool { return config.IsHoliday(a) },
		IsCategory:        func(text string) bool { return config.IsCategory(text) },
		IsTask:            func(text string) bool { return config.IsTask(text) },
	}
	return &Service{
		projectRoot: projectRoot,
		config:      config,
		parser:      &parser,
	}
}

type LineError struct {
	LineNumber int
	LineLength int
	Err        error
}

func (self *Service) ParseText(text string, date time.Time) ([]WorkItem, []LineError) {

	var workItems []WorkItem = nil
	var errors []LineError = nil
	dateInfo := DateInfoFrom(date)
	parseLine := self.parser.ParseLine(dateInfo)
	lines := strings.Split(text, "\n")
	for counter, line := range lines {
		if counter == len(lines)-1 && line == "" {
			continue
		}
		workItem, err := parseLine(line)
		if err != nil {
			errors = append(errors, LineError{LineNumber: counter, LineLength: len(line), Err: err})
		} else {
			workItems = append(workItems, workItem)
		}
	}
	log.Printf("Parsed %+v items", workItems)
	log.Printf("Parsed %+v errors", errors)
	return workItems, errors
}

func (self *Service) PossibleCategories() []string {
	return self.config.PossibleCategories()
}

func (self *Service) ParseDateFromName(uri string) (time.Time, error) {
	return DateFromFile(DateFromFileNameParams{URI: uri, ProjectRoot: self.projectRoot})
}
