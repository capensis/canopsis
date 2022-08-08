BEGIN;

CREATE TABLE IF NOT EXISTS instructions
(
    id         SERIAL PRIMARY KEY,
    custom_id  VARCHAR(255),
    title      VARCHAR(255),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (custom_id)
);

CREATE INDEX IF NOT EXISTS instructions_name_idx ON instructions (title);

COMMIT;
