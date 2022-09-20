BEGIN;

SELECT remove_continuous_aggregate_policy('cps_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('cps_event_summary_daily');
SELECT remove_retention_policy('cps_event');
SELECT remove_retention_policy('cps_event_summary_hourly');
SELECT remove_retention_policy('cps_event_summary_daily');

SELECT remove_retention_policy('fifo_queue');
SELECT remove_continuous_aggregate_policy('fifo_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('fifo_event_summary_daily');
SELECT remove_retention_policy('fifo_event');
SELECT remove_retention_policy('fifo_event_summary_hourly');
SELECT remove_retention_policy('fifo_event_summary_daily');

SELECT remove_continuous_aggregate_policy('che_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('che_event_summary_daily');
SELECT remove_retention_policy('che_event');
SELECT remove_retention_policy('che_event_summary_hourly');
SELECT remove_retention_policy('che_event_summary_daily');
SELECT remove_continuous_aggregate_policy('che_infos_summary_daily');
SELECT remove_retention_policy('che_infos');
SELECT remove_retention_policy('che_infos_summary_daily');

SELECT remove_continuous_aggregate_policy('axe_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('axe_event_summary_daily');
SELECT remove_retention_policy('axe_event');
SELECT remove_retention_policy('axe_event_summary_hourly');
SELECT remove_retention_policy('axe_event_summary_daily');
SELECT remove_retention_policy('axe_periodical');

SELECT remove_retention_policy('pbehavior_periodical');

SELECT remove_continuous_aggregate_policy('correlation_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('correlation_event_summary_daily');
SELECT remove_retention_policy('correlation_event');
SELECT remove_retention_policy('correlation_event_summary_hourly');
SELECT remove_retention_policy('correlation_event_summary_daily');

SELECT remove_continuous_aggregate_policy('service_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('service_event_summary_daily');
SELECT remove_retention_policy('service_event');
SELECT remove_retention_policy('service_event_summary_hourly');
SELECT remove_retention_policy('service_event_summary_daily');

SELECT remove_continuous_aggregate_policy('dynamic_infos_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('dynamic_infos_event_summary_daily');
SELECT remove_retention_policy('dynamic_infos_event');
SELECT remove_retention_policy('dynamic_infos_event_summary_hourly');
SELECT remove_retention_policy('dynamic_infos_event_summary_daily');

SELECT remove_continuous_aggregate_policy('action_event_summary_hourly');
SELECT remove_continuous_aggregate_policy('action_event_summary_daily');
SELECT remove_retention_policy('action_event');
SELECT remove_retention_policy('action_event_summary_hourly');
SELECT remove_retention_policy('action_event_summary_daily');

SELECT remove_retention_policy('api_requests');

COMMIT;
