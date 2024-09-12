BEGIN;

CREATE TABLE IF NOT EXISTS alarm_state_change_number
(
    time      TIMESTAMP NOT NULL,
    entity_id BIGINT
);
SELECT create_hypertable('alarm_state_change_number', 'time', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS alarm_state_change_number_entity_id_time_idx ON alarm_state_change_number (entity_id, time DESC);

END;
