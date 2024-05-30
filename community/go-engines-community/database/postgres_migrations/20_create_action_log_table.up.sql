BEGIN;

CREATE TABLE IF NOT EXISTS action_log_object(
    id            BIGSERIAL PRIMARY KEY,
    value_type    VARCHAR(255) NOT NULL,
    value_id      VARCHAR(255) NOT NULL,
    initial_value JSONB,
    created       TIMESTAMP NOT NULL DEFAULT NOW(),
    created_by    VARCHAR(255),
    deleted       TIMESTAMP,
    deleted_by    VARCHAR(255),
    UNIQUE (value_type, value_id)
);

CREATE TABLE IF NOT EXISTS action_log_object_changes(
    id                 BIGSERIAL PRIMARY KEY,
    object_id          BIGINT NOT NULL,
    author             VARCHAR(255) NOT NULL,
    time               TIMESTAMP NOT NULL DEFAULT NOW(),
    update_description JSONB
);

COMMIT;
