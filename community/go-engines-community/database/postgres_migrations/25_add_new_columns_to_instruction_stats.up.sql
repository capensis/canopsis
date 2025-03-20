BEGIN;

ALTER TABLE instruction_execution ADD COLUMN alarm_ok_timeout INT;
ALTER TABLE instruction_execution_by_modified_on ADD COLUMN sum_alarm_ok_timeout INT, ADD COLUMN count_alarm_ok_timeout INT;

DROP MATERIALIZED VIEW IF EXISTS instruction_execution_hourly;

CREATE MATERIALIZED VIEW IF NOT EXISTS instruction_execution_hourly
            (
             time,
             instruction,
             execution_count,
             successful,
             avg_complete_time,
             sum_alarm_ok_timeout,
             count_alarm_ok_timeout,
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
       sum(alarm_ok_timeout) FILTER ( WHERE alarm_ok_timeout != 0 ),
       count(alarm_ok_timeout) FILTER ( WHERE alarm_ok_timeout != 0 ),
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
