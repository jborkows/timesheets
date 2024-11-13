package model

type CategoryType = string

type WorkItem interface {
	IsHoliday() bool
}

type Holiday struct {
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

func (h *Holiday) IsHoliday() bool {
	return true
}

func (t *TimesheetEntry) IsHoliday() bool {
	return false
}

func NewTimesheet(date string) *Timesheet {
	return &Timesheet{Date: date}
}

func (t *Timesheet) Clear() {
	t.Entries = nil
}

func (t *Timesheet) AddHoliday() {
	t.Entries = append(t.Entries, &Holiday{})
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
	return 0
}

func (t *Timesheet) PotentialWorkingTime() uint8 {
	return 0
}

func (t *Timesheet) WorkingTime() float32 {
	return 0
}
