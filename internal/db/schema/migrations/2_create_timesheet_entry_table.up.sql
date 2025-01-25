CREATE TABLE timesheet_entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    holiday BOOLEAN NOT NULL,
    pending BOOLEAN NOT NULL,
    timesheet_id INTEGER NOT NULL,
    Hours    integer not null,
    Minutes  integer,
    Comment  text,
    Task     text,
    Category text,
    constraint fk_timesheet_entry_timesheet
        FOREIGN KEY (timesheet_id)
        REFERENCES timesheet (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
