create table tasks (
    day   integer not null,
    link  text not null,
    unique(day)
);