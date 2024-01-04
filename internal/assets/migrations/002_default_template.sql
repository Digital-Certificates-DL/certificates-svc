-- +migrate Up

ALTER TABLE template ADD is_default_template INTEGER;

ALTER TABLE template ADD CONSTRAINT ck_testbool_ischk CHECK (is_default_template IN (1,0));