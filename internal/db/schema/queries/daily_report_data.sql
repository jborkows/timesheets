-- name: FindStatistics :many
select * from daily_report_data where date = :date;


-- name: FindMonthlyStatistics :many
select * from monthly_report_data where month = :date;

-- name: FindWeeklyStatistics :many
select * from weekly_report_data where week_begin_date = :start_date and week_end_date = :end_date;
