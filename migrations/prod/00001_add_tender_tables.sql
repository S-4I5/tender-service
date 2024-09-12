-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" SCHEMA public;

CREATE TABLE IF NOT EXISTS tender (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    status VARCHAR(100),
    tender_version_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    organization_id VARCHAR(255),
    creator_username VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS tender_version (
    id SERIAL PRIMARY KEY,
    tender_id uuid,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    service_type VARCHAR(255),
    version INT NOT NULL
);

-- Optional: Add foreign key constraints if needed
 ALTER TABLE tender_version ADD CONSTRAINT fk_tender_id FOREIGN KEY (tender_id) REFERENCES tender(id);
 ALTER TABLE tender ADD CONSTRAINT fk_tender_version_id FOREIGN KEY (tender_version_id) REFERENCES tender_version(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
