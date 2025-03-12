package model

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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

func (self *Service) ParseLine(line string, date time.Time) WorkItem {
	dateInfo := DateInfoFrom(date)
	parseLine := self.parser.ParseLine(dateInfo)
	workItem, err := parseLine(line)
	if err != nil {
		return nil
	}
	switch e := workItem.(type) {
	case *TimesheetEntry:
		err := e.Validate()
		if err != nil {
			return nil
		}
		return e
	case *Holiday:
		return e
	}
	return nil

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

type FilePath string

func printDayStatistics(reportFile *os.File, dayEntries []DayEntry) {
	fmt.Fprintln(reportFile, "Daily statistics")
	var categoryFiles map[string][]DayEntry = make(map[string][]DayEntry)
	for _, entry := range dayEntries {
		if _, ok := categoryFiles[entry.Category]; !ok {
			categoryFiles[entry.Category] = make([]DayEntry, 0)
		}
		if !entry.Pending {
			continue
		}
		categoryFiles[entry.Category] = append(categoryFiles[entry.Category], entry)
	}
	categories := make([]string, 0, len(categoryFiles))
	for category := range categoryFiles {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	for _, category := range categories {
		entries := categoryFiles[category]
		sumHours := 0
		sumMinutes := 0
		for _, entry := range entries {
			sumHours += int(entry.Hours)
			sumMinutes += int(entry.Minutes)
		}
		fmt.Fprintf(reportFile, "%s %d:%02d\n", category, sumHours+sumMinutes/60, sumMinutes%60)
		for _, entry := range entries {
			var description strings.Builder
			if entry.Task != "" {
				description.WriteString(entry.Task)
				description.WriteString(" ")
			}
			description.WriteString(entry.Comment)
			minutes := entry.Minutes * 10 / 6
			if minutes%10 == 0 {
				fmt.Fprintf(reportFile, "%d.%d %s\n", entry.Hours, minutes/10, description.String())
			} else {
				fmt.Fprintf(reportFile, "%d.%d %s\n", entry.Hours, entry.Minutes*10/6, description.String())
			}
		}
	}
}

func printMinutesAsDecimal(minutes uint8) string {
	value := uint16(minutes) * 10 / 6
	if value%10 == 0 {
		return fmt.Sprintf("%d", value/10)
	} else {
		return fmt.Sprintf("%d", value)
	}
}

func printWeeklyStatistics(reportFile *os.File, weeklyStatistics []WeeklyStatistic) {
	fmt.Fprintln(reportFile, "Weekly statistics")
	sort.Slice(weeklyStatistics, func(i, j int) bool {
		return weeklyStatistics[i].Category < weeklyStatistics[j].Category
	})

	for _, entry := range weeklyStatistics {
		fmt.Fprintf(reportFile, "%s %d.%s\n", entry.Category, int(entry.Weekly.Hours), printMinutesAsDecimal(entry.Weekly.Minutes))
	}
}

func printMonthly(reportFile *os.File, statistics []MonthlyStatistic) {
	fmt.Fprintln(reportFile, "Monthly statistics")
	sort.Slice(statistics, func(i, j int) bool {
		return statistics[i].Category < statistics[j].Category
	})

	for _, entry := range statistics {
		fmt.Fprintf(reportFile, "%s %d.%s\n", entry.Category, int(entry.Monthly.Hours), printMinutesAsDecimal(entry.Monthly.Minutes))

	}
}

func (self *Service) ShowDailyStatistics(date time.Time) (FilePath, error) {

	fileName := fmt.Sprintf("report-timesheet-%s.txt", date.Format("2006-01-02"))
	fileName = filepath.Join(os.TempDir(), fileName)

	if path, err := os.Stat(fileName); err == nil {
		if err := os.Remove(path.Name()); err != nil {
			if !os.IsNotExist(err) {
				return "", fmt.Errorf("failed to remove report file: %w", err)
			}
		}
	}

	reportFile, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create report file: %w", err)
	}
	defer reportFile.Close()
	fmt.Fprintf(reportFile, "For %s\n", date.Format("2006-01-02"))

	dailyStatistics, err := self.DayStatistics(date)
	if err != nil {
		return "", fmt.Errorf("failed to get daily statistics: %w", err)
	}
	printDayStatistics(reportFile, dailyStatistics)
	fmt.Fprintln(reportFile)
	weeklyStatistics, err := self.WeeklyStatistics(date)
	if err != nil {
		return "", fmt.Errorf("failed to get weekly statistics: %w", err)
	}
	printWeeklyStatistics(reportFile, weeklyStatistics)
	fmt.Fprintln(reportFile)
	monthlyStatistics, err := self.MonthlyStatistics(date)
	if err != nil {
		return "", fmt.Errorf("failed to get monthly statistics: %w", err)
	}
	printMonthly(reportFile, monthlyStatistics)

	return FilePath(reportFile.Name()), nil
}

func (self *Service) ValidCategory(category string) bool {
	return self.config.IsCategory(category)
}

func (self *Service) SemanaticTokenFrom(content []Line, date time.Time) []TokenReady {
	if len(content) < 1 {
		return []TokenReady{}
	}

	var tokens []TokenReady
	for i, line := range content {
		parsed := self.ParseLine(line, date)
		if parsed == nil {
			continue
		}
		switch e := parsed.(type) {
		case *TimesheetEntry:
			err := e.Validate()
			if err != nil {
				return nil
			}
			tokens = append(tokens, TokenReady{
				Line:   i,
				Column: 0,
				Length: len(e.Category),
				Type:   ClassType,
			})
			startOfTime := 0

			for index, char := range line {
				log.Printf("Char: %s %d\n", string(char), index)
			}
			log.Printf("After category: '%s', %d", line[len(e.Category):], len(e.Category))
			j := len(e.Category)
			for ; j < len(line); j++ {
				if line[j] != ' ' && startOfTime == 0 {
					startOfTime = j
					continue
				}
				if line[j] == ' ' && startOfTime > 0 {
					log.Printf("Time: '%s'  -> %d %d %d", line[startOfTime:j], startOfTime, j, j-startOfTime)
					tokens = append(tokens, TokenReady{
						Line:   i,
						Column: startOfTime,
						Length: j - startOfTime,
						Type:   PropertyType,
					})
					break
				}
			}
			words := tokenizeFromIndex(line, j)
			for _, tocken := range words {
				tokens = append(tokens, TokenReady{
					Line:   i,
					Column: tocken.Index + 1,
					Length: len(tocken.Word),
					Type:   StringType,
				})
				log.Printf("Token %d: %s %d", i, tocken.Word, tocken.Index+1)
			}

		default:
			continue
		}

	}
	if len(tokens) == 0 {
		return []TokenReady{}
	} else {
		return tokens
	}
}
