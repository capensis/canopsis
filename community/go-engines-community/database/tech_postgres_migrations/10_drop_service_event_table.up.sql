BEGIN;

SELECT remove_compression_policy('service_event');
SELECT remove_compression_policy('service_event_summary_hourly');
SELECT remove_compression_policy('service_event_summary_daily');
ALTER TABLE service_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW service_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW service_event_summary_daily SET (timescaledb.compress = false);

DROP MATERIALIZED VIEW IF EXISTS service_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS service_event_summary_daily;
DROP TABLE IF EXISTS service_event;

COMMIT;
