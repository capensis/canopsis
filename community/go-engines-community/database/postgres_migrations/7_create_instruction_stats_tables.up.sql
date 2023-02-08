BEGIN;

CREATE TABLE IF NOT EXISTS instruction_execution
(
    time             TIMESTAMP    NOT NULL,
    instruction      VARCHAR(500) NOT NULL,
    successful       BOOLEAN      NOT NULL,
    complete_time    INT,
    init_alarm_state INT,
    res_alarm_state  INT
);
SELECT create_hypertable('instruction_execution', 'time', if_not_exists => TRUE);

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
            WITH (timescaledb.continuous)
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

SELECT add_retention_policy('instruction_execution', INTERVAL '48 hours');
SELECT add_continuous_aggregate_policy('instruction_execution_hourly', '12 hours', '1 hour', '1 hour');
ALTER TABLE instruction_execution
    SET (timescaledb.compress = true, timescaledb.compress_segmentby = 'instruction');
ALTER MATERIALIZED VIEW instruction_execution_hourly SET (timescaledb.compress = true);
SELECT add_compression_policy('instruction_execution', compress_after=>'2 days'::interval);
SELECT add_compression_policy('instruction_execution_hourly', compress_after=>'7 days'::interval);

CREATE TABLE IF NOT EXISTS instruction_execution_by_modified_on
(
    time              TIMESTAMP    NOT NULL,
    instruction       VARCHAR(500) NOT NULL,
    execution_count   INT          NOT NULL,
    successful        INT          NOT NULL,
    avg_complete_time INT          NOT NULL,
    init_critical     INT          NOT NULL DEFAULT 0,
    init_major        INT          NOT NULL DEFAULT 0,
    init_minor        INT          NOT NULL DEFAULT 0,
    res_critical      INT          NOT NULL DEFAULT 0,
    res_major         INT          NOT NULL DEFAULT 0,
    res_minor         INT          NOT NULL DEFAULT 0,
    res_ok            INT          NOT NULL DEFAULT 0
);
SELECT create_hypertable('instruction_execution_by_modified_on', 'time', if_not_exists => TRUE);

CREATE UNIQUE INDEX idx_instruction_time ON instruction_execution_by_modified_on (instruction, time);

COMMIT;
