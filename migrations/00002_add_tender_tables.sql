-- +goose Up
-- +goose StatementBegin
-- Table: bid_versions
CREATE TABLE tender (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),      -- Tender ID (string)
    status VARCHAR(100),                  -- Status (string)
    tender_version_id INT,                 -- Version number (int)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),  -- Creation timestamp (time.Time)
    organization_id VARCHAR(255),         -- Organization ID (string)
    creator_username VARCHAR(255)         -- Creator's username (string)
);

-- Table: tender_versions
CREATE TABLE tender_version (
    id SERIAL PRIMARY KEY,                -- Integer primary key
    tender_id uuid,      -- Tender ID (string)
    name VARCHAR(255) NOT NULL,           -- Name (string)
    description TEXT,                     -- Description (string)
    service_type VARCHAR(255),            -- Service type (string)
    version INT NOT NULL                 -- Version number (int)
);

-- Optional: Add foreign key constraints if needed
 ALTER TABLE tender_version ADD CONSTRAINT fk_tender_id FOREIGN KEY (tender_id) REFERENCES tender(id);
 ALTER TABLE tender ADD CONSTRAINT fk_tender_version_id FOREIGN KEY (tender_version_id) REFERENCES tender_version(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
