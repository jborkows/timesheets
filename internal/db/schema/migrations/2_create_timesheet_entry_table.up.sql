CREATE TABLE timesheet_entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    holiday BOOLEAN NOT NULL,
    pending BOOLEAN NOT NULL,
    timesheet_date INTEGER NOT NULL,
    Hours    integer not null,
    Minutes  integer,
    Comment  text,
    Task     text,
    Category text,
    constraint fk_timesheet_entry_timesheet
        FOREIGN KEY (timesheet_date)
        REFERENCES timesheet (date)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
