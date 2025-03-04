package model

import "context"

type Saver interface {
	Save(ctx context.Context, timesheet *Timesheet) error
	PendingSave(ctx context.Context, timesheet *Timesheet) error
}

type KnowsAboutWeek interface {
	Week() *Week
}

type KnowsAboutMonth interface {
	Month() *Month
}

type KnowsAboutDate interface {
	Day() *Day
}

type DayEntry struct {
	Holiday       bool
	Pending       bool
	TimesheetDate int64
	Hours         int64
	Minutes       int64
	Comment       string
	Task          string
	Category      string
}

type Queryer interface {
	Daily(ctx context.Context, knowsAboutDate KnowsAboutDate) ([]DailyStatistic, error)
	Weekly(ctx context.Context, knowsAboutWeek KnowsAboutWeek) ([]WeeklyStatistic, error)
	Monthly(ctx context.Context, knowsAboutMonth KnowsAboutMonth) ([]MonthlyStatistic, error)
	DaySummary(ctx context.Context, knowsAboutDate KnowsAboutDate) ([]DayEntry, error)
}

type Repository interface {
	Transactional(ctx context.Context, operation func(context.Context, Saver, Queryer) error) error
}
