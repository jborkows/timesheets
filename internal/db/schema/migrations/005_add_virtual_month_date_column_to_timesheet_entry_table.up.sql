alter table timesheet_entry_data add column month generated always as (Round(timesheet_date/100,0)) virtual;
