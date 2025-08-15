-- +migrate Up

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS phone_number VARCHAR DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS address TEXT DEFAULT NULL;

-- +migrate Down

ALTER TABLE users
    DROP COLUMN IF EXISTS phone_number,
    DROP COLUMN IF EXISTS address;


