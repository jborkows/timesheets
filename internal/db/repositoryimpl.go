package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jborkows/timesheets/internal/model"
)

type TimesheetStatus bool

const (
	Pending TimesheetStatus = true
	Saved   TimesheetStatus = false
)

type impl struct {
	queries  *Queries
	overtime func(model.CategoryType) bool
}

func Repository(queries *Queries, overtime func(model.CategoryType) bool) *impl {
	return &impl{
		queries:  queries,
		overtime: overtime,
	}
}

func (repository *impl) Save(timesheet *model.Timesheet) error {
	err := repository.queries.ClearTimesheetData(context.TODO(), dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to clear time sheet data: %w", err)
	}
	err = repository.saveTimeSheet(timesheet, Saved)
	if err != nil {
		return fmt.Errorf("failed to save time sheet: %w", err)
	}
	err = repository.PendingSave(timesheet)
	if err != nil {
		return fmt.Errorf("failed to save pending: %w", err)
	}
	return nil
}

func (self *impl) PendingSave(timesheet *model.Timesheet) error {
	err := self.queries.ClearPending(context.TODO(), dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to clear pending: %w", err)
	}
	return self.saveTimeSheet(timesheet, Pending)
}

func (self *impl) saveTimeSheet(timesheet *model.Timesheet, state TimesheetStatus) error {
	pending := state == Pending
	err := self.queries.CreateTimesheet(context.TODO(), dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to create time sheet: %w", err)
	}
	for _, entry := range timesheet.Entries {
		switch e := entry.(type) {
		case *model.Holiday:
			err := self.queries.AddHoliday(context.TODO(), AddHolidayParams{
				Holiday:       true,
				Pending:       pending,
				TimesheetDate: dayAsInteger(&e.Date),
			})
			if err != nil {
				return fmt.Errorf("failed to insert holiday: %w", err)
			}
		case *model.TimesheetEntry:
			savingDate := AddEntryParams{
				Holiday:       false,
				Pending:       pending,
				TimesheetDate: dayAsInteger(&timesheet.Date),
				Hours:         int64(e.Hours),
				Minutes:       int64(e.Minutes),
				Comment:       e.Comment,
				Task:          *e.Task,
				Category:      e.Category,
			}
			err := self.queries.AddEntry(context.TODO(), savingDate)
			fmt.Printf("saving entry: %v", savingDate)

			if err != nil {
				return fmt.Errorf("failed to insert pending: %w", err)
			}
		}
	}
	return nil
}

func (self *impl) Daily(knowsAboutDate model.KnowsAboutDate) ([]model.DailyStatistic, error) {
	values, err := self.queries.FindStatistics(context.TODO(), dayAsInteger(knowsAboutDate.Day()))

	if err != nil {
		return nil, fmt.Errorf("failed to find statistics: %w", err)
	}
	bucket := make(map[model.CategoryType]*model.DailyStatistic, len(values)/2+1)
	for _, value := range values {
		overtime := self.overtime(model.CategoryType(value.Category))
		if _, ok := bucket[value.Category]; !ok {
			pointer := model.DailyStatistic{
				Dirty: model.Statitic{
					Category: value.Category,
					Overtime: overtime,
				},
				Daily: model.Statitic{
					Category: value.Category,
					Overtime: overtime,
				},
			}

			bucket[value.Category] = &pointer
		}

		statistics := bucket[value.Category]
		if value.Pending {
			statistics.Dirty.Hours += uint8(value.Hours)
			statistics.Dirty.Minutes += uint8(value.Minutes)
		} else {
			statistics.Daily.Hours += uint8(value.Hours)
			statistics.Daily.Minutes += uint8(value.Minutes)
		}

	}
	bucketSlice := make([]model.DailyStatistic, 0, len(bucket))
	for _, value := range bucket {
		bucketSlice = append(bucketSlice, *value)
	}
	return bucketSlice, nil
}

func (repository *impl) Weekly(knowsAboutWeek model.KnowsAboutWeek) ([]model.WeeklyStatistic, error) {
	return nil, nil
}

func (repository *impl) Monthly(knowsAboutMonth model.KnowsAboutMonth) ([]model.MonthlyStatistic, error) {
	return nil, nil
}

func dayAsInteger(d *model.Day) int64 {
	value := time.Time(*d).Format("20060102")
	v, e := strconv.Atoi(value)
	if e != nil {
		panic(e)
	}
	return int64(v)
}
