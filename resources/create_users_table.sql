CREATE TABLE users
(
    id           SERIAL PRIMARY KEY,
    employee_id  CITEXT UNIQUE       NOT NULL,
    name         varchar(255)        NOT NULL,
    email        varchar(120) UNIQUE NOT NULL,
    password     varchar(128)        NOT NULL,
    designation  varchar(255)        NOT NULL,
    deleted_at   timestamp,
    is_activated boolean DEFAULT FALSE
);
