BEGIN;

-- it's required to drop and create materialized views again because they are depended on column.
DROP MATERIALIZED VIEW IF EXISTS cps_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS cps_event_summary_daily;

-- to alter compressed hypertable it's required to manually decompress all chunks and remove all policies.
SELECT remove_compression_policy('cps_event');

CREATE TEMP TABLE cps_event_compressed_chunks AS
SELECT chunk_schema || '.' || chunk_name AS full_chunk_name
FROM chunk_compression_stats('cps_event')
WHERE compression_status = 'Compressed';

DO $$
    DECLARE
        chunk TEXT;
    BEGIN
        FOR chunk IN SELECT full_chunk_name FROM cps_event_compressed_chunks
            LOOP
                PERFORM decompress_chunk(chunk);
            END LOOP;
    END;
$$;

ALTER TABLE cps_event SET (timescaledb.compress = false);
ALTER TABLE cps_event ALTER COLUMN interval TYPE BIGINT;
ALTER TABLE cps_event SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');

DO $$
    DECLARE
        chunk TEXT;
    BEGIN
        FOR chunk IN SELECT full_chunk_name FROM cps_event_compressed_chunks
            LOOP
                PERFORM compress_chunk(chunk);
            END LOOP;
    END;
$$;

SELECT add_compression_policy('cps_event', INTERVAL '1 day');

CREATE MATERIALIZED VIEW IF NOT EXISTS cps_event_summary_hourly
            (time, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
FROM cps_event
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS cps_event_summary_daily
            (day, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, count(*), avg(interval)
FROM cps_event
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;

SELECT add_continuous_aggregate_policy('cps_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('cps_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('cps_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('cps_event_summary_daily', INTERVAL '30 days');

ALTER MATERIALIZED VIEW cps_event_summary_hourly SET (timescaledb.compress = true, timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW cps_event_summary_daily SET (timescaledb.compress = true, timescaledb.materialized_only = false);

SELECT add_compression_policy('cps_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('cps_event_summary_daily', compress_after=>'3 days'::interval);

COMMIT;
