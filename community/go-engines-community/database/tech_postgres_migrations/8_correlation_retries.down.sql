BEGIN;

DROP MATERIALIZED VIEW IF EXISTS correlation_retries_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS correlation_retries_summary_daily;
DROP TABLE IF EXISTS correlation_retries;

COMMIT;
