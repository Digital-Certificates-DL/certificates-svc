-- +migrate Up
ALTER TABLE template DROP COLUMN is_default_template;

ALTER TABLE template ADD is_default_template BOOLEAN  DEFAULT False;


-- +migrate Down

ALTER TABLE template DROP COLUMN is_default_template;
