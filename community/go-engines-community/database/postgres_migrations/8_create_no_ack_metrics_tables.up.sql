BEGIN;

CREATE TABLE IF NOT EXISTS not_acked_in_hour_alarms
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('not_acked_in_hour_alarms', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS not_acked_in_four_hours_alarms
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('not_acked_in_four_hours_alarms', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS not_acked_in_day_alarms
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('not_acked_in_day_alarms', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS not_acked_in_hour_alarms_entity_id_time_idx ON not_acked_in_hour_alarms (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS not_acked_in_four_hours_alarms_entity_id_time_idx ON not_acked_in_four_hours_alarms (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS not_acked_in_day_alarms_entity_id_time_idx ON not_acked_in_day_alarms (entity_id, time DESC);

END;
