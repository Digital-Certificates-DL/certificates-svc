-- +migrate Up
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50),
                       token TEXT,
                       code TEXT
);

CREATE TABLE template (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER REFERENCES users(id),
                       template TEXT,
                       img_bytes bytea,
                       name TEXT,
);

-- +migrate Down
DROP TABLE template;
DROP TABLE images;

