package db

import (
	"context"
	"fmt"
	"strconv"
	"time"

	. "github.com/jborkows/timesheets/internal/model"
)

type impl struct {
	queries  *Queries
	overtime func(CategoryType) bool
}

func Repository(queries *Queries, overtime func(CategoryType) bool) *impl {
	return &impl{
		queries:  queries,
		overtime: overtime,
	}
}

func (repository *impl) Save(timesheet *Timesheet) error {
	err := repository.queries.ClearTimesheetData(context.TODO(), dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to clear time sheet data: %w", err)
	}
	err = repository.saveTimeSheet(timesheet, "SAVED")
	if err != nil {
		return fmt.Errorf("failed to save time sheet: %w", err)
	}
	err = repository.PendingSave(timesheet)
	if err != nil {
		return fmt.Errorf("failed to save pending: %w", err)
	}
	return nil
}

func (self *impl) PendingSave(timesheet *Timesheet) error {
	err := self.queries.ClearPending(context.TODO(), dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to clear pending: %w", err)
	}
	return self.saveTimeSheet(timesheet, "PENDING")
}

func (self *impl) saveTimeSheet(timesheet *Timesheet, pendingInput string) error {
	pending := pendingInput == "PENDING"

	err := self.queries.CreateTimesheet(context.TODO(), dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to create time sheet: %w", err)
	}
	for _, entry := range timesheet.Entries {
		holiday, ok := entry.(*Holiday)
		if ok {
			err := self.queries.AddHoliday(context.TODO(), AddHolidayParams{
				Holiday:       true,
				Pending:       pending,
				TimesheetDate: dayAsInteger(&holiday.Date),
			})
			if err != nil {
				return fmt.Errorf("failed to insert holiday: %w", err)
			}
			continue
		}
		entry, ok := entry.(*TimesheetEntry)
		if !ok {
			continue
		}
		savingDate := AddEntryParams{
			Holiday:       false,
			Pending:       pending,
			TimesheetDate: dayAsInteger(&timesheet.Date),
			Hours:         int64(entry.Hours),
			Minutes:       int64(entry.Minutes),
			Comment:       entry.Comment,
			Task:          *entry.Task,
			Category:      entry.Category,
		}
		err := self.queries.AddEntry(context.TODO(), savingDate)
		fmt.Printf("saving entry: %v", savingDate)

		if err != nil {
			return fmt.Errorf("failed to insert pending: %w", err)
		}
	}
	return nil
}

func (self *impl) Daily(knowsAboutDate KnowsAboutDate) ([]DailyStatistic, error) {
	values, err := self.queries.FindStatistics(context.TODO(), dayAsInteger(knowsAboutDate.Day()))

	if err != nil {
		return nil, fmt.Errorf("failed to find statistics: %w", err)
	}
	bucket := make(map[CategoryType]*DailyStatistic, len(values)/2+1)
	for _, value := range values {
		_, ok := bucket[value.Category]
		overtime := self.overtime(CategoryType(value.Category))
		if !ok {
			pointer := DailyStatistic{
				Dirty: Statitic{
					Category: value.Category,
					Hours:    0,
					Minutes:  0,
					Overtime: overtime,
				},
				Daily: Statitic{
					Category: value.Category,
					Hours:    0,
					Minutes:  0,
					Overtime: overtime,
				},
			}

			bucket[value.Category] = &pointer
		}

	}
	for _, value := range values {
		values, ok := bucket[value.Category]
		if !ok {
			panic("should not happen")
		}
		if value.Pending {
			values.Dirty.Hours += uint8(value.Hours)
			values.Dirty.Minutes += uint8(value.Minutes)
		} else {
			values.Daily.Hours += uint8(value.Hours)
			values.Daily.Minutes += uint8(value.Minutes)
		}

	}
	bucketSlice := make([]DailyStatistic, 0, len(bucket))
	for _, value := range bucket {
		bucketSlice = append(bucketSlice, *value)
	}
	return bucketSlice, nil
}

func (repository *impl) Weekly(knowsAboutWeek KnowsAboutWeek) ([]WeeklyStatistic, error) {
	return nil, nil
}

func (repository *impl) Monthly(knowsAboutMonth KnowsAboutMonth) ([]MonthlyStatistic, error) {
	return nil, nil
}

func dayAsInteger(d *Day) int64 {
	value := time.Time(*d).Format("20060102")
	v, e := strconv.Atoi(value)
	if e != nil {
		panic(e)
	}
	return int64(v)
}
