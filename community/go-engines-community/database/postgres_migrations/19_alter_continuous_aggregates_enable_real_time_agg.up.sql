-- Since Timescale 2.13 real time aggregates are DISABLED by default.
-- The migration manually enables real time aggregates as it's supposed to be.

BEGIN;

ALTER MATERIALIZED VIEW instruction_execution_hourly SET (timescaledb.materialized_only = false);
ALTER MATERIALIZED VIEW message_rate_hourly SET (timescaledb.materialized_only = false);

END;
