-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA public;

CREATE TABLE IF NOT EXISTS tender (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    status VARCHAR(100),
    tender_version_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    organization_id uuid,
    creator_username VARCHAR(50)
);

CREATE TYPE tender_version_service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TABLE IF NOT EXISTS tender_version (
    id SERIAL PRIMARY KEY,
    tender_id uuid,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    service_type tender_version_service_type,
    version INT NOT NULL
);

ALTER TABLE tender ADD CONSTRAINT fk_tender_version_id FOREIGN KEY (tender_version_id) REFERENCES tender_version(id);
ALTER TABLE tender ADD CONSTRAINT fk_organization_id FOREIGN KEY (organization_id) REFERENCES organization(id);
ALTER TABLE tender ADD CONSTRAINT fk_creator_username FOREIGN KEY (creator_username) REFERENCES employee(username);

ALTER TABLE tender_version ADD CONSTRAINT fk_tender_id FOREIGN KEY (tender_id) REFERENCES tender(id) ON DELETE CASCADE;-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
