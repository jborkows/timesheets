// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

type DailyReportDatum struct {
	Date     int64
	Pending  bool
	Category string
	Holiday  bool
	Hours    int64
	Minutes  int64
}

type MonthlyReportDatum struct {
	Month    interface{}
	Pending  bool
	Category string
	Holiday  bool
	Hours    int64
	Minutes  int64
}

type TimesheetDatum struct {
	Date int64
}

type TimesheetEntryDatum struct {
	ID            int64
	Holiday       bool
	Pending       bool
	TimesheetDate int64
	Hours         int64
	Minutes       int64
	Comment       string
	Task          string
	Category      string
	Month         interface{}
}
