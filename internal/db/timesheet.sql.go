// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: timesheet.sql

package db

import (
	"context"
)

const createTimesheet = `-- name: CreateTimesheet :exec
INSERT or IGNORE INTO timesheet_data (date) VALUES (?1)
`

func (q *Queries) CreateTimesheet(ctx context.Context, date int64) error {
	_, err := q.db.ExecContext(ctx, createTimesheet, date)
	return err
}

const findTimesheet = `-- name: FindTimesheet :one
SELECT date FROM timesheet_data WHERE date = (?1)
`

func (q *Queries) FindTimesheet(ctx context.Context, date int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, findTimesheet, date)
	err := row.Scan(&date)
	return date, err
}
