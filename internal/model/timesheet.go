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
	Date        string
	Description string
}

type TimesheetEntry struct {
	Hours    uint8
	Minutes  uint8
	Comment  string
	Task     *string
	Category CategoryType
}

type Timesheet struct {
	Date    string
	Entries []WorkItem
}

func NewHoliday(date string, description string) (*Holiday, error) {
	return &Holiday{Date: date, Description: description}, nil
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

func validateDate(date string) error {
	_, error := time.Parse("2006-01-02", date)
	if error != nil {
		return &ValidateDateError{Date: date, Err: error}
	}
	return nil
}

func NewTimesheet(date string) (*Timesheet, error) {
	if err := validateDate(date); err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}
	return &Timesheet{Date: date}, nil
}

func (t *Timesheet) Clear() {
	t.Entries = nil
}

func validate(entry *TimesheetEntry) error {
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
func (t *Timesheet) AddEntry(entry TimesheetEntry) error {
	if err := validate(&entry); err != nil {
		return fmt.Errorf("invalid entry: %w", err)
	}
	t.Entries = append(t.Entries, &entry)
	return nil
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
