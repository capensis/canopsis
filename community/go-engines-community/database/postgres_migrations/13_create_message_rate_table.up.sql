BEGIN;

CREATE TABLE IF NOT EXISTS message_rate
(
    time TIMESTAMP NOT NULL
);
SELECT create_hypertable('message_rate', 'time', if_not_exists => TRUE);

CREATE MATERIALIZED VIEW IF NOT EXISTS message_rate_hourly
            (time, count)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), count(*)
FROM message_rate
GROUP BY time_bucket('1 hour', time)
    WITH NO DATA;

SELECT add_continuous_aggregate_policy('message_rate_hourly', '3 hours', '1 hour', '1 hour');
SELECT add_retention_policy('message_rate', INTERVAL '24 hours');

ALTER MATERIALIZED VIEW message_rate_hourly SET (timescaledb.compress = true);
SELECT add_compression_policy('message_rate_hourly', compress_after=>'1 day'::interval);

END;
