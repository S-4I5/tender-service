-- +goose Up
-- +goose StatementBegin
CREATE TABLE bid (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    bid_version_id uuid,
    status VARCHAR(100),
    tender_id VARCHAR(255),               -- Tender ID (string)
    author_type VARCHAR(100),             -- Author type (string)
    author_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()  -- Creation timestamp (time.Time)
);

CREATE TABLE bid_version (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),                -- Integer primary key
    bid_id uuid DEFAULT public.uuid_generate_v4(),         -- Bid ID (string)
    name VARCHAR(255) NOT NULL,           -- Name (string)
    description TEXT,                     -- Description (string)
    version INT NOT NULL                 -- Version number (int)
);

-- Table: decisions
CREATE TABLE decisions (
    id uuid PRIMARY KEY  DEFAULT public.uuid_generate_v4(),                  -- UUID primary key
    verdict VARCHAR(50) NOT NULL,         -- Verdict (enum or varchar to represent Verdict)
    username VARCHAR(255) NOT NULL,       -- Username (string)
    bid_id uuid NOT NULL      -- Bid ID (string)
);

-- Table: feedback
CREATE TABLE feedback (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),                  -- UUID primary key
    bid_id uuid NOT NULL,         -- Bid ID (string)
    description TEXT,                     -- Description (string)
    username VARCHAR(255) NOT NULL,       -- Username (string)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()   -- Creation timestamp (time.Time)
);

 ALTER TABLE bid_version ADD CONSTRAINT fk_bidid FOREIGN KEY (bid_id) REFERENCES bid(id);
 ALTER TABLE decisions ADD CONSTRAINT fk_bidid FOREIGN KEY (bid_id) REFERENCES bid(id);
 ALTER TABLE feedback ADD CONSTRAINT fk_bidid FOREIGN KEY (bid_id) REFERENCES bid(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
