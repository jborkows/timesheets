package model

import (
	"fmt"
	"time"
)

type CategoryType = string

type WorkItem interface {
	IsHoliday() bool
}

type Holiday struct {
	Date        Day
	Description string
}

type TimesheetEntry struct {
	Hours    uint8
	Minutes  uint8
	Comment  string
	Task     *string
	Category CategoryType
}

func (self *TimesheetEntry) TaskName() string {
	if self.Task == nil {
		return ""
	}
	return *self.Task
}

type Timesheet struct {
	Date    Day
	Entries []WorkItem
}

func parseDate(date string) (*Day, error) {
	parsed, error := time.Parse("2006-01-02", date)

	if error != nil {
		return nil, &ValidateDateError{Date: date, Err: error}
	}
	asDay := Day(parsed)
	return &(asDay), nil
}

func NewHoliday(date string, description string) (*Holiday, error) {
	parsed, err := parseDate(date)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}
	return &Holiday{Date: *parsed, Description: description}, nil
}

func (h *Holiday) IsHoliday() bool {
	return true
}

func (t *TimesheetEntry) IsHoliday() bool {
	return false
}

func (t *Timesheet) AddHoliday(holiday *Holiday) error {
	if holiday.Date != t.Date {
		return fmt.Errorf("holiday date %v does not match timesheet date %v", holiday.Date, t.Date)
	}
	t.Entries = append(t.Entries, holiday)
	return nil
}

type ValidateDateError struct {
	Date string
	Err  error
}

func (e *ValidateDateError) Error() string {
	return fmt.Sprintf("invalid date: %v", e.Err)
}

type InvalidTimesheetEntry struct {
	Entry *TimesheetEntry
	Err   error
}

func (e *InvalidTimesheetEntry) Error() string {
	return fmt.Sprintf("invalid entry: %v", e.Err)
}

func TimesheetForDate(aTime time.Time) *Timesheet {
	aDay := Day(aTime)
	return &Timesheet{Date: aDay}
}

func NewTimesheet(date string) (*Timesheet, error) {
	parsed, error := parseDate(date)
	if error != nil {
		return nil, fmt.Errorf("invalid date: %w", error)
	}
	return &Timesheet{Date: *parsed}, nil
}

func (t *Timesheet) Clear() {
	t.Entries = nil
}

func (entry *TimesheetEntry) Validate() error {
	if entry.Hours >= 24 {
		return &InvalidTimesheetEntry{Entry: entry, Err: fmt.Errorf("hours cannot be more than 24")}
	}
	if entry.Minutes >= 60 {
		return &InvalidTimesheetEntry{Entry: entry, Err: fmt.Errorf("minutes cannot be more than 60")}
	}
	if entry.Hours == 0 && entry.Minutes == 0 {
		return &InvalidTimesheetEntry{Entry: entry, Err: fmt.Errorf("hours and minutes cannot be both 0")}
	}
	return nil
}
func validate(entry *TimesheetEntry) error {
	return entry.Validate()
}

func (t *Timesheet) Add(entry *TimesheetEntry) error {
	if err := validate(entry); err != nil {
		return fmt.Errorf("invalid entry: %w", err)
	}
	t.Entries = append(t.Entries, entry)
	return nil
}

func (t *Timesheet) AddEntry(entry TimesheetEntry) error {
	return t.Add(&entry)
}

func (t *Timesheet) PotentialTotalTime() uint8 {
	return 8
}

func (t *Timesheet) PotentialWorkingTime() uint8 {
	for _, entry := range t.Entries {
		if entry.IsHoliday() {
			return 0
		}
	}
	return 8
}

func (t *Timesheet) WorkingTime() float32 {

	var total float32 = 0
	for _, entry := range t.Entries {
		if entry.IsHoliday() {
			continue
		}
		entry, ok := entry.(*TimesheetEntry)
		if !ok {
			panic("invalid entry type")
		}
		total += float32(entry.Hours) + float32(entry.Minutes)/60
	}
	return total
}

func (t *Timesheet) Day() *Day {
	return &t.Date
}

func (t *Timesheet) startOfWeek() *Day {
	dayOfWeek := t.Date.DayOfWeek()
	if dayOfWeek == time.Monday {
		return &t.Date
	}
	if t.Date.DayOfMonth() == 1 {
		return &t.Date
	}
	backToTime := time.Time(t.Date)
	for i := int(dayOfWeek); i != int(time.Monday); i-- {
		backToTime = backToTime.AddDate(0, 0, -1)
		if backToTime.Weekday() == time.Monday {
			startOfWeek := Day(backToTime)
			return &startOfWeek
		}
		if backToTime.Day() == 1 {
			startOfWeek := Day(backToTime)
			return &startOfWeek
		}
	}
	panic("unreachable")
}

func (t *Timesheet) endOfWeek() *Day {
	dayOfWeek := t.Date.DayOfWeek()
	if dayOfWeek == time.Sunday {
		return &t.Date
	}
	backToTime := time.Time(t.Date)
	for i := int(dayOfWeek); i != int(time.Sunday); i++ {
		backToTime = backToTime.AddDate(0, 0, 1)
		if backToTime.Weekday() == time.Sunday {
			endOfWeek := Day(backToTime)
			return &endOfWeek
		}
		if backToTime.Day() == 1 {
			endOfWeek := Day(backToTime.AddDate(0, 0, -1))
			return &endOfWeek
		}
	}
	panic("unreachable")
}

func (t *Timesheet) Week() *Week {
	startOfWeek := t.startOfWeek()
	endOfWeek := t.endOfWeek()
	return &Week{
		BeginDate: *startOfWeek,
		EndDate:   *endOfWeek,
	}
}

func (t *Timesheet) Month() *Month {
	asTime := time.Time(t.Date)
	firstOfMonth := time.Date(asTime.Year(), asTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return &Month{
		BeginDate: Day(firstOfMonth),
		EndDate:   Day(lastOfMonth),
	}
}
