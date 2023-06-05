BEGIN;

SELECT remove_compression_policy('cps_event');
SELECT remove_compression_policy('cps_event_summary_hourly');
SELECT remove_compression_policy('cps_event_summary_daily');
ALTER TABLE cps_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW cps_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW cps_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('fifo_event');
SELECT remove_compression_policy('fifo_event_summary_hourly');
SELECT remove_compression_policy('fifo_event_summary_daily');
ALTER TABLE fifo_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW fifo_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW fifo_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('che_event');
SELECT remove_compression_policy('che_event_summary_hourly');
SELECT remove_compression_policy('che_event_summary_daily');
ALTER TABLE che_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW che_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW che_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('che_infos');
SELECT remove_compression_policy('che_infos_summary_daily');
ALTER TABLE che_infos
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW che_infos_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('axe_event');
SELECT remove_compression_policy('axe_event_summary_hourly');
SELECT remove_compression_policy('axe_event_summary_daily');
ALTER TABLE axe_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW axe_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW axe_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('correlation_event');
SELECT remove_compression_policy('correlation_event_summary_hourly');
SELECT remove_compression_policy('correlation_event_summary_daily');
ALTER TABLE correlation_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW correlation_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW correlation_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('service_event');
SELECT remove_compression_policy('service_event_summary_hourly');
SELECT remove_compression_policy('service_event_summary_daily');
ALTER TABLE service_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW service_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW service_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('dynamic_infos_event');
SELECT remove_compression_policy('dynamic_infos_event_summary_hourly');
SELECT remove_compression_policy('dynamic_infos_event_summary_daily');
ALTER TABLE dynamic_infos_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW dynamic_infos_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW dynamic_infos_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('action_event');
SELECT remove_compression_policy('action_event_summary_hourly');
SELECT remove_compression_policy('action_event_summary_daily');
ALTER TABLE action_event
    SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW action_event_summary_hourly SET (timescaledb.compress = false);
ALTER MATERIALIZED VIEW action_event_summary_daily SET (timescaledb.compress = false);

SELECT remove_compression_policy('api_requests');
ALTER TABLE api_requests
    SET (timescaledb.compress = false);

COMMIT;
