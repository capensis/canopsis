BEGIN;

CREATE TABLE IF NOT EXISTS service_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('service_event', 'time', if_not_exists => TRUE);

CREATE MATERIALIZED VIEW IF NOT EXISTS service_event_summary_hourly
            (time, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
FROM service_event
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS service_event_summary_daily
            (day, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, count(*), avg(interval)
FROM service_event
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;

ALTER TABLE service_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW service_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW service_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('service_event', INTERVAL '1 day');
SELECT add_compression_policy('service_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('service_event_summary_daily', compress_after=>'3 days'::interval);

COMMIT;
