BEGIN;

CREATE TABLE IF NOT EXISTS perf_data
(
    time      TIMESTAMP    NOT NULL,
--     TODO BIGINT
    entity_id INT          NOT NULL,
    name      VARCHAR(255) NOT NULL,
    value     NUMERIC      NOT NULL,
    unit      VARCHAR(5)   NOT NULL
);
SELECT create_hypertable('perf_data', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS perf_data_entity_id_time_idx ON perf_data (entity_id, time DESC);

END;
