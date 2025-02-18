package model

// TODO: Issue #3
type Saver interface {
	Save(timesheet *Timesheet) error
	PendingSave(timesheet *Timesheet) error
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
	Daily(knowsAboutDate KnowsAboutDate) ([]DailyStatistic, error)
	Weekly(knowsAboutWeek KnowsAboutWeek) ([]WeeklyStatistic, error)
	Monthly(knowsAboutMonth KnowsAboutMonth) ([]MonthlyStatistic, error)
}
