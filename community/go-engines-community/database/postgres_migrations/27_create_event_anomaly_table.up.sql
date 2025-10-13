BEGIN;

CREATE TABLE IF NOT EXISTS event_anomaly
(
    time TIMESTAMP NOT NULL,
    connector_name VARCHAR(255) DEFAULT NULL,
    count BIGINT NOT NULL,
    anomaly BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (time, connector_name)
);
SELECT create_hypertable('event_anomaly', 'time', if_not_exists => TRUE);

CREATE MATERIALIZED VIEW IF NOT EXISTS event_anomaly_hourly
            (time, connector_name, count, anomaly)
            WITH (timescaledb.continuous)
AS
SELECT
  time_bucket('1 hour', time) AS time,
  connector_name,
  sum(count) AS count,
  (COUNT(*) FILTER (WHERE anomaly)
   > COUNT(*) FILTER (WHERE NOT anomaly)) AS anomaly
FROM event_anomaly
GROUP BY time_bucket('1 hour', time), connector_name
    WITH NO DATA;

SELECT add_continuous_aggregate_policy('event_anomaly_hourly', '3 hours', '1 hour', '1 hour');
SELECT add_retention_policy('event_anomaly', INTERVAL '24 hours');

ALTER MATERIALIZED VIEW event_anomaly_hourly SET (timescaledb.compress = true);
SELECT add_compression_policy('event_anomaly_hourly', compress_after=>'1 day'::interval);

END;
