BEGIN;

CREATE TABLE IF NOT EXISTS go_metrics
(
    time   TIMESTAMP    NOT NULL,
    metric VARCHAR(60)  NOT NULL,
    source VARCHAR(500) NOT NULL,
    value  REAL         NOT NULL
);
SELECT create_hypertable('go_metrics', 'time', if_not_exists => TRUE);

SELECT add_retention_policy('go_metrics', INTERVAL '1 day');

ALTER TABLE go_metrics
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'metric');
SELECT add_compression_policy('go_metrics', INTERVAL '1 day');

COMMIT;
