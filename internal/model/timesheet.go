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
	Date string
}

type TimesheetEntry struct {
	Hours    uint8
	Minutes  uint8
	Comment  string
	Task     string
	Category CategoryType
}

type Timesheet struct {
	Date    string
	Entries []WorkItem
}

func NewHoliday(date string) (*Holiday, error) {
	if err := validateDate(date); err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}
	return &Holiday{Date: date}, nil
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

func (t *Timesheet) AddEntry(hours, minutes uint8, comment, task string, category CategoryType) {
	t.Entries = append(t.Entries, &TimesheetEntry{
		Hours:    hours,
		Minutes:  minutes,
		Comment:  comment,
		Task:     task,
		Category: category,
	})
}

func (t *Timesheet) PotentialTotalTime() uint8 {
	var total uint8 = uint8(0)
	for _, entry := range t.Entries {
		if !entry.IsHoliday() {
			total += entry.(*TimesheetEntry).Hours
		} else {
			return 8
		}
	}
	return total

}

func (t *Timesheet) PotentialWorkingTime() uint8 {
	return 0
}

func (t *Timesheet) WorkingTime() float32 {
	return 0
}
