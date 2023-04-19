-- +migrate Up
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50),
                       token TEXT,
                       code TEXT
);

CREATE TABLE template (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER REFERENCES users(id)

);

CREATE TABLE images (
                      id SERIAL PRIMARY KEY,
                      user_id INTEGER REFERENCES users(id),
                      template_id INTEGER REFERENCES template(id),
                      img_bytes bytea
);

-- +migrate Down
DROP TABLE users;
DROP TABLE template;
DROP TABLE images;

