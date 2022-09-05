BEGIN;

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

-- FIFO
CREATE MATERIALIZED VIEW IF NOT EXISTS fifo_event_summary_hourly
            (time, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
FROM fifo_event
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS fifo_event_summary_daily
            (day, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, count(*), avg(interval)
FROM fifo_event
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;
-- END FIFO

-- CHE
CREATE MATERIALIZED VIEW IF NOT EXISTS che_event_summary_hourly
            (time, type, entity_type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, entity_type, count(*), avg(interval)
FROM che_event
GROUP BY time_bucket('1 hour', time), type, entity_type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS che_event_summary_daily
            (day, type, entity_type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, entity_type, count(*), avg(interval)
FROM che_event
GROUP BY time_bucket('1 day', time), type, entity_type
WITH NO DATA;

CREATE MATERIALIZED VIEW che_infos_summary_daily
            (day, name, count)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), name, count(*)
FROM che_infos
GROUP BY time_bucket('1 day', time), name
WITH NO DATA;
-- END CHE

-- AXE
CREATE MATERIALIZED VIEW IF NOT EXISTS axe_event_summary_hourly
            (time, type, entity_type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, entity_type, count(*), avg(interval)
FROM axe_event
GROUP BY time_bucket('1 hour', time), type, entity_type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS axe_event_summary_daily
            (day, type, entity_type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, entity_type, count(*), avg(interval)
FROM axe_event
GROUP BY time_bucket('1 day', time), type, entity_type
WITH NO DATA;
-- END AXE

-- CORRELATION
CREATE MATERIALIZED VIEW IF NOT EXISTS correlation_event_summary_hourly
            (time, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
FROM correlation_event
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS correlation_event_summary_daily
            (day, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, count(*), avg(interval)
FROM correlation_event
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;
-- END CORRELATION

-- SERVICE
CREATE MATERIALIZED VIEW IF NOT EXISTS service_event_summary_hourly
            (time, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
FROM service_event
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS service_event_summary_daily
            (day, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, count(*), avg(interval)
FROM service_event
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;
-- END SERVICE

-- DYNAMIC INFOS
CREATE MATERIALIZED VIEW IF NOT EXISTS dynamic_infos_event_summary_hourly
            (time, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
FROM dynamic_infos_event
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS dynamic_infos_event_summary_daily
            (day, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, count(*), avg(interval)
FROM dynamic_infos_event
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;
-- END DYNAMIC INFOS

-- ACTION
CREATE MATERIALIZED VIEW IF NOT EXISTS action_event_summary_hourly
            (time, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 hour', time), type, count(*), avg(interval)
FROM action_event
GROUP BY time_bucket('1 hour', time), type
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS action_event_summary_daily
            (day, type, count, interval)
            WITH (timescaledb.continuous)
AS
SELECT time_bucket('1 day', time), type, count(*), avg(interval)
FROM action_event
GROUP BY time_bucket('1 day', time), type
WITH NO DATA;
-- END ACTION

COMMIT;
