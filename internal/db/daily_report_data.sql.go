// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: daily_report_data.sql

package db

import (
	"context"
)

const findMonthlyStatistics = `-- name: FindMonthlyStatistics :many
select month, pending, category, holiday, hours, minutes from monthly_report_data where month = ?1
`

func (q *Queries) FindMonthlyStatistics(ctx context.Context, date interface{}) ([]MonthlyReportDatum, error) {
	rows, err := q.db.QueryContext(ctx, findMonthlyStatistics, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MonthlyReportDatum
	for rows.Next() {
		var i MonthlyReportDatum
		if err := rows.Scan(
			&i.Month,
			&i.Pending,
			&i.Category,
			&i.Holiday,
			&i.Hours,
			&i.Minutes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findStatistics = `-- name: FindStatistics :many
select date, pending, category, holiday, hours, minutes from daily_report_data where date = ?1
`

func (q *Queries) FindStatistics(ctx context.Context, date int64) ([]DailyReportDatum, error) {
	rows, err := q.db.QueryContext(ctx, findStatistics, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DailyReportDatum
	for rows.Next() {
		var i DailyReportDatum
		if err := rows.Scan(
			&i.Date,
			&i.Pending,
			&i.Category,
			&i.Holiday,
			&i.Hours,
			&i.Minutes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findWeeklyStatistics = `-- name: FindWeeklyStatistics :many
select week_begin_date, week_end_date, pending, category, holiday, hours, minutes from weekly_report_data where week_begin_date = ?1 and week_end_date = ?2
`

type FindWeeklyStatisticsParams struct {
	StartDate int64
	EndDate   int64
}

func (q *Queries) FindWeeklyStatistics(ctx context.Context, arg FindWeeklyStatisticsParams) ([]WeeklyReportDatum, error) {
	rows, err := q.db.QueryContext(ctx, findWeeklyStatistics, arg.StartDate, arg.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WeeklyReportDatum
	for rows.Next() {
		var i WeeklyReportDatum
		if err := rows.Scan(
			&i.WeekBeginDate,
			&i.WeekEndDate,
			&i.Pending,
			&i.Category,
			&i.Holiday,
			&i.Hours,
			&i.Minutes,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
