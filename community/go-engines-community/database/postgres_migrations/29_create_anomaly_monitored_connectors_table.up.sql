BEGIN;

CREATE TABLE IF NOT EXISTS anomaly_monitored_connectors
(
    id            BIGSERIAL PRIMARY KEY,
    connector_name   VARCHAR(255) UNIQUE NOT NULL
);

COMMIT;
