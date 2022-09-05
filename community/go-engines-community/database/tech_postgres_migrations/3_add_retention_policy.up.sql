BEGIN;

SELECT add_continuous_aggregate_policy('cps_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('cps_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('cps_event', INTERVAL '49 hours');
SELECT add_retention_policy('cps_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('cps_event_summary_daily', INTERVAL '30 days');

-- -- FIFO
SELECT add_retention_policy('fifo_queue', INTERVAL '30 days');

SELECT add_continuous_aggregate_policy('fifo_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('fifo_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('fifo_event', INTERVAL '49 hours');
SELECT add_retention_policy('fifo_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('fifo_event_summary_daily', INTERVAL '30 days');
-- END FIFO

-- CHE
SELECT add_continuous_aggregate_policy('che_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('che_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('che_event', INTERVAL '49 hours');
SELECT add_retention_policy('che_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('che_event_summary_daily', INTERVAL '30 days');

SELECT add_continuous_aggregate_policy('che_infos_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('che_infos', INTERVAL '49 hours');
SELECT add_retention_policy('che_infos_summary_daily', INTERVAL '30 days');
-- END CHE

-- AXE
SELECT add_continuous_aggregate_policy('axe_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('axe_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('axe_event', INTERVAL '49 hours');
SELECT add_retention_policy('axe_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('axe_event_summary_daily', INTERVAL '30 days');

SELECT add_retention_policy('axe_periodical', INTERVAL '7 days');
-- END AXE

-- PBEHAVIOR
SELECT add_retention_policy('pbehavior_periodical', INTERVAL '7 days');
-- END PBEHAVIOR

-- CORRELATION
SELECT add_continuous_aggregate_policy('correlation_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('correlation_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('correlation_event', INTERVAL '49 hours');
SELECT add_retention_policy('correlation_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('correlation_event_summary_daily', INTERVAL '30 days');
-- END CORRELATION

-- SERVICE
SELECT add_continuous_aggregate_policy('service_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('service_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('service_event', INTERVAL '49 hours');
SELECT add_retention_policy('service_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('service_event_summary_daily', INTERVAL '30 days');
-- END SERVICE

-- DYNAMIC INFOS
SELECT add_continuous_aggregate_policy('dynamic_infos_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('dynamic_infos_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('dynamic_infos_event', INTERVAL '49 hours');
SELECT add_retention_policy('dynamic_infos_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('dynamic_infos_event_summary_daily', INTERVAL '30 days');
-- END DYNAMIC INFOS

-- ACTION
SELECT add_continuous_aggregate_policy('action_event_summary_hourly', '12 hours', '1 hour', '1 hour');
SELECT add_continuous_aggregate_policy('action_event_summary_daily', '49 hours', '1 hour', '1 hour');
SELECT add_retention_policy('action_event', INTERVAL '49 hours');
SELECT add_retention_policy('action_event_summary_hourly', INTERVAL '7 days');
SELECT add_retention_policy('action_event_summary_daily', INTERVAL '30 days');
-- END ACTION

-- API
SELECT add_retention_policy('api_requests', INTERVAL '30 days');
-- END API

COMMIT;
