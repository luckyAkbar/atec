-- +migrate Up notransaction

CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- to be able to access uuid_generate_v4 function

CREATE TYPE roles AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    username TEXT NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    roles roles NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS children (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_user_id UUID NOT NULL,
    name TEXT NOT NULL,
    date_of_birth TIMESTAMPTZ NOT NULL,
    gender BOOLEAN NOT NULL, -- male is true, female is false <- just for clarity
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT fk_parent_user_id FOREIGN KEY (parent_user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS packages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_by UUID NOT NULL,
    questionnaire JSONB NOT NULL,
    indication_categories JSONB NOT NULL,
    name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT fk_created_by FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    package_id UUID NOT NULL,
    child_id UUID DEFAULT NULL,
    created_by UUID DEFAULT NULL,
    answer JSONB NOT NULL,
    result JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,

    CONSTRAINT fk_package_id FOREIGN KEY (package_id) REFERENCES packages(id),
    CONSTRAINT fk_child_id FOREIGN KEY (child_id) REFERENCES children(id)
);


-- +migrate Down

DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS packages;
DROP TABLE IF EXISTS children;
DROP TABLE IF EXISTS users;