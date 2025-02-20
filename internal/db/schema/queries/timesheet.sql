-- name: CreateTimesheet :exec
INSERT or IGNORE INTO timesheet_data (date,week_begin_date,week_end_date) VALUES (:date,:week_start,:week_end);

-- name: FindTimesheet :one
SELECT * FROM timesheet_data WHERE date = (:date);


