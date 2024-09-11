-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA public;

CREATE TABLE IF NOT EXISTS employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--INSERT INTO employee (id, username, first_name, last_name)
--       VALUES ('c81e51fe-3140-4dda-9be6-a4adf506d5be', 'aboba', 'xd', 'xd');
--
--INSERT INTO employee (id, username, first_name, last_name)
--       VALUES ('ed0cd1f2-5c98-4ba4-bf7a-0ab28078bfac', 'aboba2', 'xd', 'xd');
--
--INSERT INTO employee (id, username, first_name, last_name)
--       VALUES ('33ef51b7-4297-4d81-82ef-3c50e00f5db1', 'sex', 'xd', 'xd');

--CREATE TYPE organization_type AS ENUM (
--    'IE',
--    'LLC',
--    'JSC'
--);

CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--INSERT INTO organization (id, name, description, type)
--       VALUES ('78acd561-76e9-485c-a9c0-7e02c2305cbf', 'aboba', 'bg', 'IE');

CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

--INSERT INTO organization_responsible (organization_id, user_id)
--       VALUES ('78acd561-76e9-485c-a9c0-7e02c2305cbf', 'c81e51fe-3140-4dda-9be6-a4adf506d5be');
--
--INSERT INTO organization_responsible (organization_id, user_id)
--       VALUES ('78acd561-76e9-485c-a9c0-7e02c2305cbf', 'ed0cd1f2-5c98-4ba4-bf7a-0ab28078bfac');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
