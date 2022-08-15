BEGIN;

CREATE TABLE IF NOT EXISTS instructions
(
    id         SERIAL PRIMARY KEY,
    custom_id  VARCHAR(255),
    name       VARCHAR(255),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE (custom_id)
);

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
    instruction_id INT,
    value          INT
);
SELECT create_hypertable('manual_instruction_assigned_instructions', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS manual_instruction_executed_instructions
(
    time           TIMESTAMP NOT NULL,
    instruction_id INT,
    value          INT
);
SELECT create_hypertable('manual_instruction_executed_instructions', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS instructions_name_idx ON instructions (name);

CREATE INDEX IF NOT EXISTS manual_instruction_assigned_alarms_entity_id_time_idx ON manual_instruction_assigned_alarms (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS manual_instruction_executed_alarms_entity_id_time_idx ON manual_instruction_executed_alarms (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS manual_instruction_assigned_instructions_instruction_id_time_idx ON manual_instruction_assigned_instructions (instruction_id, time DESC);
CREATE INDEX IF NOT EXISTS manual_instruction_executed_instructions_instruction_id_time_idx ON manual_instruction_executed_instructions (instruction_id, time DESC);

COMMIT;
