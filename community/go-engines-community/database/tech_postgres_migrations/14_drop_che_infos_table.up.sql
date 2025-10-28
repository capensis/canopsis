BEGIN;

SELECT remove_compression_policy('che_infos');
SELECT remove_compression_policy('che_infos_summary_daily');
ALTER TABLE che_infos
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW che_infos_summary_daily SET (timescaledb.compress = false);

SELECT remove_continuous_aggregate_policy('che_infos_summary_daily');
SELECT remove_retention_policy('che_infos');
SELECT remove_retention_policy('che_infos_summary_daily');

DROP MATERIALIZED VIEW IF EXISTS che_infos_summary_daily;
DROP TABLE IF EXISTS che_infos;

COMMIT;
