BEGIN;

CREATE TABLE IF NOT EXISTS correlation_retries
(
    time    TIMESTAMP   NOT NULL,
    type    VARCHAR(30) NOT NULL,
    retries INT         NOT NULL
);
SELECT create_hypertable('correlation_retries', 'time', if_not_exists => TRUE);

CREATE MATERIALIZED VIEW IF NOT EXISTS correlation_retries_summary_hourly
            (time, type, retries)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, avg(retries)
FROM correlation_retries
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS correlation_retries_summary_daily
            (day, type, retries)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, avg(retries)
FROM correlation_retries
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;

SELECT add_continuous_aggregate_policy('correlation_retries_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('correlation_retries_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('correlation_retries', INTERVAL '49 hours');
SELECT add_retention_policy('correlation_retries_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('correlation_retries_summary_daily', INTERVAL '30 days');

ALTER TABLE correlation_retries
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW correlation_retries_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW correlation_retries_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('correlation_retries', INTERVAL '1 day');
SELECT add_compression_policy('correlation_retries_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('correlation_retries_summary_daily', compress_after=>'3 days'::interval);

COMMIT;
