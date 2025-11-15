CREATE TABLE users (
    user_name varchar(100),
    email varchar(100),
    password varchar(100)
);

CREATE TABLE clip_data (
    id serial primary key,
    user_name varchar(100),
    title varchar,
    text text,
    tag varchar,
    added_at timestamp default now()
);