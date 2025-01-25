-- name: FindStatistics :many
select * from daily_report_data where date = :date;
