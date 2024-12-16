package model

import "fmt"

type InvalidTime struct {
	Err error
}

func (e *InvalidTime) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

type DateInfo struct {
	Value string
}

type HolidayClassifier = func(aDate *DateInfo) bool
type Statitic struct {
	Category string
	Hours    uint8
	Minutes  uint8
	Overtime bool
}

type DailyStatistic struct {
	Category string
	Dirty    Statitic
	Daily    Statitic
}

type WeeklyStatistic struct {
	Category string
	Dirty    Statitic
	Weekly   Statitic
}

type MonthlyStatistic struct {
	Category      string
	Monthly       Statitic
	RequiredHours uint8
}
