-- +migrate Up

ALTER TABLE children
    ADD COLUMN IF NOT EXISTS guardian_name TEXT DEFAULT NULL;

-- +migrate Down

ALTER TABLE children
    DROP COLUMN IF EXISTS guardian_name;


