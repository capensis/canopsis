BEGIN;

CREATE TABLE IF NOT EXISTS entity_infos_update
(
    time      TIMESTAMP    NOT NULL,
    rule_id   VARCHAR(255) NOT NULL,
    entity_id BIGINT,
    name      VARCHAR(255) NOT NULL,
    value     JSONB
);
SELECT create_hypertable('entity_infos_update', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS entity_infos_update_entity_id_time_idx ON entity_infos_update (entity_id, time DESC);

SELECT add_retention_policy('entity_infos_update', INTERVAL '1 week');

ALTER TABLE entity_infos_update
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'name');
SELECT add_compression_policy('entity_infos_update', INTERVAL '1 day');

COMMIT;
