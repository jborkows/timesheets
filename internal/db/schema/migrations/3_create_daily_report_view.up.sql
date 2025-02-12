create view daily_report_data as
select t.date,
    te.pending,
    te.category,
    te.holiday,
    sum(te.hours) + Round(sum(te.minutes)/60,0) as hours,
    te.Minutes%60 minutes 
from timesheet_data t
join timesheet_entry_data te on t.date = te.timesheet_date
group by t.date, te.pending, te.holiday, te.category;
