-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS bid_author_type;
CREATE TYPE bid_author_type AS ENUM (
    'Organization',
    'User'
);

DROP TYPE IF EXISTS bid_status;
CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
);

DROP TYPE IF EXISTS bid_decision_type;
CREATE TYPE bid_decision_type AS ENUM (
    'Approved',
    'None',
    'Rejected'
);

CREATE TABLE IF NOT EXISTS bid (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    bid_version_id uuid,
    status bid_status,
    tender_id uuid,
    author_type bid_author_type,
    author_id uuid,
    decision bid_decision_type DEFAULT 'None',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bid_version (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    bid_id uuid DEFAULT public.uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    version INT NOT NULL
);

DROP TYPE IF EXISTS decision_verdict_type;
CREATE TYPE decision_verdict_type AS ENUM (
    'Approved',
    'Rejected'
);

CREATE TABLE IF NOT EXISTS decision (
    id uuid PRIMARY KEY  DEFAULT public.uuid_generate_v4(),
    verdict decision_verdict_type NOT NULL,
    username VARCHAR(50) NOT NULL,
    bid_id uuid NOT NULL
);

CREATE TABLE IF NOT EXISTS feedback (
    id uuid PRIMARY KEY DEFAULT public.uuid_generate_v4(),
    bid_id uuid NOT NULL,
    description TEXT,
    username VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE bid ADD CONSTRAINT fk_tender_id FOREIGN KEY (tender_id) REFERENCES tender(id) ON DELETE CASCADE;
ALTER TABLE bid ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES employee(id);

ALTER TABLE bid_version ADD CONSTRAINT fk_bid_id FOREIGN KEY (bid_id) REFERENCES bid(id) ON DELETE CASCADE;

ALTER TABLE decision ADD CONSTRAINT fk_bid_id FOREIGN KEY (bid_id) REFERENCES bid(id) ON DELETE CASCADE;
ALTER TABLE decision ADD CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES employee(username);

ALTER TABLE feedback ADD CONSTRAINT fk_bid_id FOREIGN KEY (bid_id) REFERENCES bid(id) ON DELETE CASCADE;
ALTER TABLE feedback ADD CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES employee(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
