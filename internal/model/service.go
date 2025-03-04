package model

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

type Service struct {
	projectRoot string
	config      *Config
	parser      *Parser
	repository  Repository
}

type CleanupFunction func()

func NewService(projectRoot string, config *Config, repository Repository) *Service {
	parser := Parser{
		HolidayClassifier: func(a *DateInfo) bool { return config.IsHoliday(a) },
		IsCategory:        func(text string) bool { return config.IsCategory(text) },
		IsTask:            func(text string) bool { return config.IsTask(text) },
	}

	return &Service{
		projectRoot: projectRoot,
		config:      config,
		parser:      &parser,
		repository:  repository,
	}
}

type LineError struct {
	LineNumber int
	LineLength int
	Err        error
}

type WriteMode int

const (
	DRAFT = iota
	SAVE
)

func (w WriteMode) String() string {
	switch w {
	case DRAFT:
		return "draft"
	case SAVE:
		return "save"
	default:
		log.Fatalf("Unknown write mode %d", w)
		return ""
	}
}

func (self *Service) ProcessForSave(text string, date time.Time) ([]WorkItem, []LineError) {
	return self.process(text, date, SAVE)
}
func (self *Service) ProcessForDraft(text string, date time.Time) ([]WorkItem, []LineError) {
	return self.process(text, date, DRAFT)
}

func (self *Service) process(text string, date time.Time, mode WriteMode) ([]WorkItem, []LineError) {
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
			continue
		}
		switch e := workItem.(type) {
		case *TimesheetEntry:
			err := e.Validate()
			if err != nil {
				errors = append(errors, LineError{LineNumber: counter, LineLength: len(line), Err: err})
				continue
			}
		}
		workItems = append(workItems, workItem)
	}
	log.Printf("Parsed %+v items", workItems)
	log.Printf("Parsed %+v errors", errors)
	if len(workItems) == 0 {
		return workItems, errors
	}
	timesheet := TimesheetForDate(date)
	for _, workItem := range workItems {
		switch e := workItem.(type) {
		case *Holiday:
			err := timesheet.AddHoliday(e)
			if err != nil {
				errors = append(errors, LineError{LineNumber: 0, LineLength: 0, Err: err})
			}
		case *TimesheetEntry:
			err := timesheet.Add(e)
			if err != nil {
				errors = append(errors, LineError{LineNumber: 0, LineLength: 0, Err: err})
			}
		}
	}
	err := self.saveData(timesheet, mode)
	if err != nil {
		errors = append(errors, LineError{LineNumber: 0, LineLength: 0, Err: err})
	}

	return workItems, errors
}

func (self *Service) saveData(timesheet *Timesheet, mode WriteMode) error {

	switch mode {
	case DRAFT:
		return self.repository.Transactional(context.TODO(), func(ctx context.Context, repository Saver, queryer Queryer) error {
			err := repository.PendingSave(ctx, timesheet)
			if err != nil {
				return fmt.Errorf("failed to save pending: %w", err)
			}
			return nil
		})
	case SAVE:
		return self.repository.Transactional(context.TODO(), func(ctx context.Context, repository Saver, queryer Queryer) error {
			err := repository.Save(ctx, timesheet)
			if err != nil {
				return fmt.Errorf("failed to save : %w", err)
			}
			return nil
		})

	}
	return nil
}

func (self *Service) PossibleCategories() []string {
	return self.config.PossibleCategories()
}

func (self *Service) ParseDateFromName(uri string) (time.Time, error) {
	return DateFromFile(DateFromFileNameParams{URI: uri, ProjectRoot: self.projectRoot})
}

type DataQuery[T any] = func(context.Context, Queryer, time.Time) (T, error)

func statistics[T any](self *Service, date time.Time, fn DataQuery[T]) (T, error) {
	var result T
	err := self.repository.Transactional(context.TODO(), func(ctx context.Context, repository Saver, queryer Queryer) error {
		r, error := fn(ctx, queryer, date)
		if error != nil {
			return fmt.Errorf("failed to get statistics: %w", error)
		}
		result = r
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (self *Service) DailyStatistics(date time.Time) ([]DailyStatistic, error) {
	return statistics(self, date, func(ctx context.Context, queryer Queryer, date time.Time) ([]DailyStatistic, error) {
		return queryer.Daily(ctx, TimesheetForDate(date))
	})
}
func (self *Service) DayStatistics(date time.Time) ([]DayEntry, error) {
	return statistics(self, date, func(ctx context.Context, queryer Queryer, date time.Time) ([]DayEntry, error) {
		return queryer.DaySummary(ctx, TimesheetForDate(date))
	})
}

func (self *Service) WeeklyStatistics(date time.Time) ([]WeeklyStatistic, error) {
	return statistics(self, date, func(ctx context.Context, queryer Queryer, date time.Time) ([]WeeklyStatistic, error) {
		return queryer.Weekly(ctx, TimesheetForDate(date))
	})
}

func (self *Service) MonthlyStatistics(date time.Time) ([]MonthlyStatistic, error) {
	return statistics(self, date, func(ctx context.Context, queryer Queryer, date time.Time) ([]MonthlyStatistic, error) {
		return queryer.Monthly(ctx, TimesheetForDate(date))
	})
}
