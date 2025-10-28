-- +migrate Up

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS nik VARCHAR DEFAULT NULL;

-- +migrate Down

ALTER TABLE users
    DROP COLUMN IF EXISTS nik;


