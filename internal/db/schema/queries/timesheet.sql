-- name: CreateTimesheet :exec
INSERT or IGNORE INTO timesheet_data (date) VALUES (:date);

-- name: FindTimesheet :one
SELECT * FROM timesheet_data WHERE date = (:date);


