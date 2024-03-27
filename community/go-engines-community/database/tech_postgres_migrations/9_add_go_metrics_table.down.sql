BEGIN;

SELECT remove_compression_policy('go_metrics');
ALTER TABLE go_metrics SET (timescaledb.compress = false);

SELECT remove_retention_policy('go_metrics');
DROP TABLE IF EXISTS go_metrics;

COMMIT;
