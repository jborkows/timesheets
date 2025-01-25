-- name: ClearTimesheetData :exec
delete from timesheet_entry_data where timesheet_date = :timesheet_date;


-- name: ClearPending :exec
delete from timesheet_entry_data where timesheet_date = :timesheet_date and pending = 1;

-- name: AddHoliday :exec
insert into timesheet_entry_data (holiday, pending, hours, timesheet_date) values (:holiday, :pending, 8, :timesheet_date);

-- name: AddEntry :exec
insert into timesheet_entry_data (holiday, pending, timesheet_date, hours, minutes,comment,task, category) values (:holiday, :pending, :timesheet_date, :hours, :minutes, :comment, :task, :category);



