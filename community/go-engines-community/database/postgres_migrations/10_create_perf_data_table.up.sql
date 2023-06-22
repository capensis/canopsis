BEGIN;

CREATE TABLE IF NOT EXISTS perf_data
(
    time      TIMESTAMP    NOT NULL,
    entity_id INT          NOT NULL,
    name      VARCHAR(255) NOT NULL,
    value     NUMERIC      NOT NULL,
    unit      VARCHAR(5)   NOT NULL
);
SELECT create_hypertable('perf_data', 'time', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS perf_data_name_entity_id_time_idx ON perf_data (name, entity_id, time DESC);

CREATE MATERIALIZED VIEW IF NOT EXISTS perf_data_name AS
SELECT name, unit FROM perf_data
GROUP BY name, unit;
CREATE UNIQUE INDEX IF NOT EXISTS perf_data_name_name_idx ON perf_data_name (name);

END;
