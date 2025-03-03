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

func (repository *impl) Save(ctx context.Context, timesheet *model.Timesheet) error {
	err := repository.queries.ClearTimesheetData(ctx, dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to clear time sheet data: %w", err)
	}
	err = repository.saveTimeSheet(ctx, timesheet, Saved)
	if err != nil {
		return fmt.Errorf("failed to save time sheet: %w", err)
	}
	err = repository.PendingSave(ctx, timesheet)
	if err != nil {
		return fmt.Errorf("failed to save pending: %w", err)
	}
	return nil
}

func (self *impl) PendingSave(ctx context.Context, timesheet *model.Timesheet) error {
	err := self.queries.ClearPending(ctx, dayAsInteger(&timesheet.Date))
	if err != nil {
		return fmt.Errorf("failed to clear pending: %w", err)
	}
	return self.saveTimeSheet(ctx, timesheet, Pending)
}

func (self *impl) saveTimeSheet(ctx context.Context, timesheet *model.Timesheet, state TimesheetStatus) error {
	pending := state == Pending
	err := self.queries.CreateTimesheet(ctx, CreateTimesheetParams{
		Date:      dayAsInteger(&timesheet.Date),
		WeekStart: dayAsInteger(&timesheet.Week().BeginDate),
		WeekEnd:   dayAsInteger(&timesheet.Week().EndDate),
	})
	if err != nil {
		return fmt.Errorf("failed to create time sheet: %w", err)
	}
	for _, entry := range timesheet.Entries {
		switch e := entry.(type) {
		case *model.Holiday:
			err := self.queries.AddHoliday(ctx, AddHolidayParams{
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
				Task:          e.TaskName(),
				Category:      e.Category,
			}
			err := self.queries.AddEntry(ctx, savingDate)
			fmt.Printf("saving entry: %v", savingDate)

			if err != nil {
				return fmt.Errorf("failed to insert pending: %w", err)
			}
		}
	}
	return nil
}

type selectOutput struct {
	Category string
	Hours    int64
	Minutes  int64
	Pending  bool
}

type dirtyWithData struct {
	Dirty model.Statitic
	Data  model.Statitic
}

type structDataParams[T any, V any] struct {
	self     *impl
	data     []T
	toSelect func(entry T) selectOutput
	toModel  func(entry *dirtyWithData) V
}

func groupData[T any, V any](params structDataParams[T, V]) []V {
	self := params.self
	data := params.data
	bucket := make(map[model.CategoryType]*dirtyWithData, len(data)/2+1)
	for _, input := range data {
		value := params.toSelect(input)
		overtime := self.overtime(model.CategoryType(value.Category))
		if _, ok := bucket[value.Category]; !ok {
			pointer := dirtyWithData{
				Dirty: model.Statitic{
					Category: value.Category,
					Overtime: overtime,
				},
				Data: model.Statitic{
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
			statistics.Data.Hours += uint8(value.Hours)
			statistics.Data.Minutes += uint8(value.Minutes)
		}
	}

	bucketSlice := make([]V, 0, len(bucket))
	for _, input := range bucket {
		value := params.toModel(input)
		bucketSlice = append(bucketSlice, value)
	}
	return bucketSlice
}

func (self *impl) Daily(ctx context.Context, knowsAboutDate model.KnowsAboutDate) ([]model.DailyStatistic, error) {
	values, err := self.queries.FindStatistics(context.TODO(), dayAsInteger(knowsAboutDate.Day()))

	if err != nil {
		return nil, fmt.Errorf("failed to find statistics: %w for %v", err, knowsAboutDate)
	}
	result := groupData(structDataParams[DailyReportDatum, model.DailyStatistic]{
		self: self,
		data: values,
		toSelect: func(entry DailyReportDatum) selectOutput {
			return selectOutput{
				Category: entry.Category,
				Hours:    entry.Hours,
				Minutes:  entry.Minutes,
				Pending:  entry.Pending,
			}
		},
		toModel: func(entry *dirtyWithData) model.DailyStatistic {
			return model.DailyStatistic{
				Category: entry.Data.Category,
				Dirty:    entry.Dirty,
				Daily:    entry.Data,
			}
		},
	})
	return result, nil
}

func (self *impl) Weekly(ctx context.Context, knowsAboutWeek model.KnowsAboutWeek) ([]model.WeeklyStatistic, error) {
	week := knowsAboutWeek.Week()
	values, err := self.queries.FindWeeklyStatistics(context.TODO(), FindWeeklyStatisticsParams{
		StartDate: dayAsInteger(&week.BeginDate),
		EndDate:   dayAsInteger(&week.EndDate),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find statistics: %w for %v", err, knowsAboutWeek)
	}

	result := groupData[WeeklyReportDatum, model.WeeklyStatistic](structDataParams[WeeklyReportDatum, model.WeeklyStatistic]{
		self: self,
		data: values,
		toSelect: func(entry WeeklyReportDatum) selectOutput {
			return selectOutput{
				Category: entry.Category,
				Hours:    entry.Hours,
				Minutes:  entry.Minutes,
				Pending:  entry.Pending,
			}
		},
		toModel: func(entry *dirtyWithData) model.WeeklyStatistic {
			return model.WeeklyStatistic{
				Dirty:  entry.Dirty,
				Weekly: entry.Data,
			}
		},
	})
	return result, nil
}

func (self *impl) Monthly(ctx context.Context, knowsAboutMonth model.KnowsAboutMonth) ([]model.MonthlyStatistic, error) {
	month := knowsAboutMonth.Month()
	values, err := self.queries.FindMonthlyStatistics(context.TODO(), dayAsInteger(&month.BeginDate)/100)

	if err != nil {
		return nil, fmt.Errorf("failed to find statistics: %w", err)
	}
	result := groupData[MonthlyReportDatum, model.MonthlyStatistic](structDataParams[MonthlyReportDatum, model.MonthlyStatistic]{
		self: self,
		data: values,
		toSelect: func(entry MonthlyReportDatum) selectOutput {
			return selectOutput{
				Category: entry.Category,
				Hours:    entry.Hours,
				Minutes:  entry.Minutes,
				Pending:  entry.Pending,
			}
		},
		toModel: func(entry *dirtyWithData) model.MonthlyStatistic {
			return model.MonthlyStatistic{
				Dirty:   entry.Dirty,
				Monthly: entry.Data,
			}
		},
	})
	return result, nil
}

func dayAsInteger(d *model.Day) int64 {
	value := time.Time(*d).Format("20060102")
	v, e := strconv.Atoi(value)
	if e != nil {
		panic(e)
	}
	return int64(v)
}
