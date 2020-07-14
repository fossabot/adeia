CREATE TABLE sessions
(
    id                    SERIAL PRIMARY KEY,
    user_id               integer REFERENCES users (id),
    refresh_token         bytea UNIQUE NOT NULL,
    refresh_token_expires timestamp    NOT NULL
);
