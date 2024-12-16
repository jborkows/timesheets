package model

type Saver interface {
	Save(timesheet *Timesheet) error
	PendingSave(timesheet *Timesheet) error
}

type Queryer interface {
	Daily(timesheet *Timesheet) ([]DailyStatistic, error)
	Weekly(timesheet *Timesheet) ([]WeeklyStatistic, error)
	Monthly(timesheet *Timesheet) ([]MonthlyStatistic, error)
}
