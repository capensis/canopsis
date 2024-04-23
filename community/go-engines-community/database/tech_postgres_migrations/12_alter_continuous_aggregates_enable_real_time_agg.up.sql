-- Since Timescale 2.13 real time aggregates are DISABLED by default.
-- The migration manually enables real time aggregates as it's supposed to be.

BEGIN;

ALTER MATERIALIZED VIEW cps_event_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW cps_event_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW fifo_event_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW fifo_event_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW che_event_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW che_event_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW che_infos_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW axe_event_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW axe_event_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW correlation_event_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW correlation_event_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW dynamic_infos_event_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW dynamic_infos_event_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW action_event_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW action_event_summary_daily SET (timescaledb.materialized_only = false);

ALTER MATERIALIZED VIEW correlation_retries_summary_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW correlation_retries_summary_daily SET (timescaledb.materialized_only = false);

END;
