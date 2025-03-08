create view weekly_report_data as
select t.week_begin_date,
    t.week_end_date,
    te.pending,
    te.category,
    te.holiday,
    sum(te.hours) + Round(sum(te.minutes)/60,0) as hours,
    sum(te.Minutes)%60 minutes 
from timesheet_data t
join timesheet_entry_data te on t.date = te.timesheet_date
group by t.week_begin_date, t.week_end_date, te.pending, te.holiday, te.category;
