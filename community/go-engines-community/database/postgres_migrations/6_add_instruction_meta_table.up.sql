BEGIN;

CREATE TABLE IF NOT EXISTS manual_instruction_assigned_alarms
(
    time           TIMESTAMP NOT NULL,
    entity_id      INT,
    value          INT
);
SELECT create_hypertable('manual_instruction_assigned_alarms', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS manual_instruction_executed_alarms
(
    time           TIMESTAMP NOT NULL,
    entity_id      INT,
    value          INT
);
SELECT create_hypertable('manual_instruction_executed_alarms', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS manual_instruction_assigned_instructions
(
    time           TIMESTAMP NOT NULL,
    instruction_id VARCHAR(255),
    value          INT
);
SELECT create_hypertable('manual_instruction_assigned_instructions', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS manual_instruction_executed_instructions
(
    time           TIMESTAMP NOT NULL,
    instruction_id VARCHAR(255),
    value          INT
);
SELECT create_hypertable('manual_instruction_executed_instructions', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS manual_instruction_assigned_alarms_entity_id_time_idx ON manual_instruction_assigned_alarms (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS manual_instruction_executed_alarms_entity_id_time_idx ON manual_instruction_executed_alarms (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS manual_instruction_assigned_instructions_instruction_id_time_idx ON manual_instruction_assigned_instructions (instruction_id, time DESC);
CREATE INDEX IF NOT EXISTS manual_instruction_executed_instructions_instruction_id_time_idx ON manual_instruction_executed_instructions (instruction_id, time DESC);

COMMIT;
