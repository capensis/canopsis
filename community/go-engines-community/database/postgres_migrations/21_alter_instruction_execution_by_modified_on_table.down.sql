BEGIN;

DROP INDEX IF EXISTS instruction_execution_by_modified_on_instruction_updated_time_idx;
ALTER TABLE instruction_execution_by_modified_on RENAME TO instruction_execution_by_modified_on_to_remove;

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

INSERT INTO instruction_execution_by_modified_on
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
SELECT
    instruction_updated as time,
    instruction,
    SUM(execution_count),
    SUM(successful),
    AVG(avg_complete_time),
    SUM(init_critical),
    SUM(init_major),
    SUM(init_minor),
    SUM(res_critical),
    SUM(res_major),
    SUM(res_minor),
    SUM(res_ok)
FROM instruction_execution_by_modified_on_to_remove
GROUP BY instruction, instruction_updated;

CREATE UNIQUE INDEX idx_instruction_time ON instruction_execution_by_modified_on (instruction, time);

DROP TABLE instruction_execution_by_modified_on_to_remove;

COMMIT;
