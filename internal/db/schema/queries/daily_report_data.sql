-- name: FindStatistics :many
select * from daily_report_data where date = :date;


-- name: FindMonthlyStatistics :many
select * from monthly_report_data where month = :date;
