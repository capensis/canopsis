CREATE TABLE IF NOT EXISTS entity_categories
(
    id     UUID PRIMARY KEY,
    name   VARCHAR(255),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
