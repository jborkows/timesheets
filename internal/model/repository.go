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

type Queryer interface {
	Daily(ctx context.Context, knowsAboutDate KnowsAboutDate) ([]DailyStatistic, error)
	Weekly(ctx context.Context, knowsAboutWeek KnowsAboutWeek) ([]WeeklyStatistic, error)
	Monthly(ctx context.Context, knowsAboutMonth KnowsAboutMonth) ([]MonthlyStatistic, error)
}

type Repository interface {
	Transactional(ctx context.Context, operation func(context.Context, Saver, Queryer) error) error
}
