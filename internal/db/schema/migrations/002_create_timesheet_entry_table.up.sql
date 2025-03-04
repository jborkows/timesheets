CREATE TABLE timesheet_entry_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    holiday BOOLEAN NOT NULL,
    pending BOOLEAN NOT NULL,
    timesheet_date INTEGER NOT NULL,
    Hours    integer not null,
    Minutes  integer not null default 0,
    Comment  text not null default '',
    Task     text not null default '',
    Category text not null default '',
    constraint fk_timesheet_entry_timesheet
        FOREIGN KEY (timesheet_date)
        REFERENCES timesheet_data (date)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
