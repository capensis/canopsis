BEGIN;

DROP MATERIALIZED VIEW IF EXISTS cps_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS cps_event_summary_daily;

DROP MATERIALIZED VIEW IF EXISTS fifo_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS fifo_event_summary_daily;

DROP MATERIALIZED VIEW IF EXISTS che_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS che_event_summary_daily;
DROP MATERIALIZED VIEW IF EXISTS che_infos_summary_daily;

DROP MATERIALIZED VIEW IF EXISTS axe_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS axe_event_summary_daily;

DROP MATERIALIZED VIEW IF EXISTS correlation_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS correlation_event_summary_daily;

DROP MATERIALIZED VIEW IF EXISTS service_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS service_event_summary_daily;

DROP MATERIALIZED VIEW IF EXISTS dynamic_infos_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS dynamic_infos_event_summary_daily;

DROP MATERIALIZED VIEW IF EXISTS action_event_summary_hourly;
DROP MATERIALIZED VIEW IF EXISTS action_event_summary_daily;

COMMIT;
