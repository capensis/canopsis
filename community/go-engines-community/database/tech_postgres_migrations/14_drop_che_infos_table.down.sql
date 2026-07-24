BEGIN;

CREATE TABLE IF NOT EXISTS che_infos
(
    time TIMESTAMP   NOT NULL,
    name VARCHAR(30) NOT NULL
);
SELECT create_hypertable('che_infos', 'time', if_not_exists => TRUE);

CREATE MATERIALIZED VIEW IF NOT EXISTS che_infos_summary_daily
            (day, name, count)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), name, count(*)
FROM che_infos
GROUP BY time_bucket('1 day', time), name
WITH NO DATA;

SELECT add_continuous_aggregate_policy(
           'che_infos_summary_daily',
           start_offset => INTERVAL '49 hours',
           end_offset => INTERVAL '1 hour',
           schedule_interval => INTERVAL '1 hour',
           if_not_exists => TRUE
       );
SELECT add_retention_policy('che_infos', drop_after => INTERVAL '49 hours', if_not_exists => TRUE);
SELECT add_retention_policy('che_infos_summary_daily', drop_after => INTERVAL '30 days', if_not_exists => TRUE);

ALTER TABLE che_infos
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'name');
ALTER MATERIALIZED VIEW che_infos_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('che_infos', compress_after => INTERVAL '1 day', if_not_exists => TRUE);
SELECT add_compression_policy('che_infos_summary_daily', compress_after => '3 days'::interval, if_not_exists => TRUE);

ALTER MATERIALIZED VIEW che_infos_summary_daily SET (timescaledb.materialized_only = false);

COMMIT;
