BEGIN;

ALTER TABLE cps_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW cps_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW cps_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('cps_event', INTERVAL '1 day');
SELECT add_compression_policy('cps_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('cps_event_summary_daily', compress_after=>'3 days'::interval);

-- -- FIFO
ALTER TABLE fifo_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW fifo_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW fifo_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('fifo_event', INTERVAL '1 day');
SELECT add_compression_policy('fifo_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('fifo_event_summary_daily', compress_after=>'3 days'::interval);
-- END FIFO

-- CHE
ALTER TABLE che_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW che_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW che_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('che_event', INTERVAL '1 day');
SELECT add_compression_policy('che_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('che_event_summary_daily', compress_after=>'3 days'::interval);

ALTER TABLE che_infos
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'name');
ALTER MATERIALIZED VIEW che_infos_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('che_infos', INTERVAL '1 day');
SELECT add_compression_policy('che_infos_summary_daily', compress_after=>'3 days'::interval);
-- END CHE

-- AXE
ALTER TABLE axe_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW axe_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW axe_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('axe_event', INTERVAL '1 day');
SELECT add_compression_policy('axe_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('axe_event_summary_daily', compress_after=>'3 days'::interval);
-- END AXE

-- CORRELATION
ALTER TABLE correlation_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW correlation_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW correlation_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('correlation_event', INTERVAL '1 day');
SELECT add_compression_policy('correlation_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('correlation_event_summary_daily', compress_after=>'3 days'::interval);
-- END CORRELATION

-- SERVICE
ALTER TABLE service_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW service_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW service_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('service_event', INTERVAL '1 day');
SELECT add_compression_policy('service_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('service_event_summary_daily', compress_after=>'3 days'::interval);
-- END SERVICE

-- DYNAMIC INFOS
ALTER TABLE dynamic_infos_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW dynamic_infos_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW dynamic_infos_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('dynamic_infos_event', INTERVAL '1 day');
SELECT add_compression_policy('dynamic_infos_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('dynamic_infos_event_summary_daily', compress_after=>'3 days'::interval);
-- END DYNAMIC INFOS

-- ACTION
ALTER TABLE action_event
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'type');
ALTER MATERIALIZED VIEW action_event_summary_hourly SET (timescaledb.compress = true);
ALTER MATERIALIZED VIEW action_event_summary_daily SET (timescaledb.compress = true);
SELECT add_compression_policy('action_event', INTERVAL '1 day');
SELECT add_compression_policy('action_event_summary_hourly', compress_after=>'1 day'::interval);
SELECT add_compression_policy('action_event_summary_daily', compress_after=>'3 days'::interval);
-- END ACTION

-- API
ALTER TABLE api_requests
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'url');
SELECT add_compression_policy('api_requests', INTERVAL '1 day');
-- END API

COMMIT;
