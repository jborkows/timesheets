create view daily_report as
select date,pending,category,holiday,hours + floor(minutes/60) as hours, mod(Minutes,60) minutes from (
select t.date, te.pending, te.holiday, te.category, sum(te.Hours) hours, sum(te.Minutes) minutes
from timesheet t
join timesheet_entry te on t.id = te.timesheet_id
group by t.date, te.pending, te.holiday, te.category
)
