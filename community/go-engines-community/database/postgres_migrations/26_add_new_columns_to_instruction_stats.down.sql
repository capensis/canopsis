BEGIN;

DROP MATERIALIZED VIEW IF EXISTS instruction_execution_hourly;

ALTER TABLE instruction_execution DROP COLUMN alarm_ok_timeout;
ALTER TABLE instruction_execution_by_modified_on DROP COLUMN sum_alarm_ok_timeout, DROP COLUMN count_alarm_ok_timeout;

CREATE MATERIALIZED VIEW IF NOT EXISTS instruction_execution_hourly
            (
             time,
             instruction,
             execution_count,
             successful,
             avg_complete_time,
             init_critical,
             init_major,
             init_minor,
             res_critical,
             res_major,
             res_minor,
             res_ok
                )
            WITH (timescaledb.continuous, timescaledb.materialized_only = false)
AS
SELECT time_bucket('1 hour', time),
       instruction,
       count(*),
       sum(CASE WHEN successful = true THEN 1 ELSE 0 END),
       avg(complete_time) FILTER ( WHERE successful = true ),
       sum(CASE WHEN successful = true AND init_alarm_state = 3 THEN 1 ELSE 0 END),
       sum(CASE WHEN successful = true AND init_alarm_state = 2 THEN 1 ELSE 0 END),
       sum(CASE WHEN successful = true AND init_alarm_state = 1 THEN 1 ELSE 0 END),
       sum(CASE WHEN successful = true AND res_alarm_state = 3 THEN 1 ELSE 0 END),
       sum(CASE WHEN successful = true AND res_alarm_state = 2 THEN 1 ELSE 0 END),
       sum(CASE WHEN successful = true AND res_alarm_state = 1 THEN 1 ELSE 0 END),
       sum(CASE WHEN successful = true AND res_alarm_state = 0 THEN 1 ELSE 0 END)
FROM instruction_execution
GROUP BY time_bucket('1 hour', time), instruction
WITH NO DATA;

SELECT add_continuous_aggregate_policy('instruction_execution_hourly', '12 hours', '1 hour', '1 hour');
ALTER MATERIALIZED VIEW instruction_execution_hourly SET (timescaledb.compress = true);
SELECT add_compression_policy('instruction_execution_hourly', compress_after=>'7 days'::interval);

COMMIT;


