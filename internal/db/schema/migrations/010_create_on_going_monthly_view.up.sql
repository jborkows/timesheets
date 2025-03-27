create view monthly_ongoing_report_data as
SELECT 
    te.month,
    te.pending,
    COUNT(te.timesheet_date) AS counted_days
FROM (
    SELECT DISTINCT timesheet_date, month, pending
    FROM timesheet_entry_data
    WHERE holiday = 0
) te
GROUP BY te.month, te.pending;

