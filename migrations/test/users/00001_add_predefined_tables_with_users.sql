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

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

INSERT INTO employee (id, username) VALUES
    ('a25e6411-fd87-4dd4-82c4-655d782259ec', 'in-org1-1'),
    ('94019f37-d6cd-4dda-9d59-056f07b4f53c', 'in-org1-2'),
    ('94975bdd-fdd4-4ed9-8cdf-e9bbd586b6e7', 'in-org2-1'),
    ('d971db11-c43d-4842-add7-10f4730526b7', 'not-in-org');

INSERT INTO organization (id, name) VALUES
    ('97a5bb3d-0265-4b1f-aad0-f3f93a8bbfca', 'org-1'),
    ('50d8765a-2081-4d09-bd8e-784163dfbb65', 'org-2');

INSERT INTO organization_responsible (organization_id, user_id) VALUES
    ('97a5bb3d-0265-4b1f-aad0-f3f93a8bbfca', 'a25e6411-fd87-4dd4-82c4-655d782259ec'),
    ('97a5bb3d-0265-4b1f-aad0-f3f93a8bbfca', '94019f37-d6cd-4dda-9d59-056f07b4f53c'),
    ('50d8765a-2081-4d09-bd8e-784163dfbb65', '94975bdd-fdd4-4ed9-8cdf-e9bbd586b6e7');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
