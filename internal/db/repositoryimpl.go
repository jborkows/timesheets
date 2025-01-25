package db

import (
	. "github.com/jborkows/timesheets/internal/model"
)

type impl struct {
}

func Repository() *impl {
	return &impl{}
}

func (repository *impl) Save(timesheet *Timesheet) error {
	return nil
}

func (repository *impl) PendingSave(timesheet *Timesheet) error {
	return nil
}

func (repository *impl) Daily(knowsAboutDate KnowsAboutDate) ([]DailyStatistic, error) {
	return nil, nil
}

func (repository *impl) Weekly(knowsAboutWeek KnowsAboutWeek) ([]WeeklyStatistic, error) {
	return nil, nil
}

func (repository *impl) Monthly(knowsAboutMonth KnowsAboutMonth) ([]MonthlyStatistic, error) {
	return nil, nil
}
