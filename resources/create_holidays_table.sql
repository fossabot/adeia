CREATE TABLE holidays
(
    id   SERIAL PRIMARY KEY,
    name varchar(255) NOT NULL,
    type varchar(255) NOT NULL,
    date date         NOT NULL
);
