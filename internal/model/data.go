package model

import (
	"fmt"
	"time"
)

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

type Day time.Time

func (d *Day) String() string {
	return time.Time(*d).Format("2006-01-02")
}

func (d *Day) DayOfWeek() time.Weekday {
	return time.Time(*d).Weekday()
}

func (d *Day) DayOfMonth() int8 {
	return int8(time.Time(*d).Day())
}

type Week struct {
	BeginDate Day
	EndDate   Day
}

func (w *Week) DaysInWeek() int8 {
	diff := int8(time.Time(w.EndDate).Day() - time.Time(w.BeginDate).Day())
	if diff == 0 {
		return 1
	} else {
		return diff
	}
}

func (w *Week) String() string {
	return fmt.Sprintf("%s - %s", w.BeginDate.String(), w.EndDate.String())
}

type Month struct {
	BeginDate Day
	EndDate   Day
}

func (w *Month) String() string {
	return fmt.Sprintf("%s - %s", w.BeginDate.String(), w.EndDate.String())
}

func (w *Month) DaysInMonth() int8 {
	return int8(time.Time(w.EndDate).Day() - time.Time(w.BeginDate).Day())
}
