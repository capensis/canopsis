BEGIN;

CREATE TABLE IF NOT EXISTS ticket_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    user_id   VARCHAR(255),
    value     INT
);
SELECT create_hypertable('ticket_number', 'time', if_not_exists => TRUE);

CREATE INDEX IF NOT EXISTS ticket_number_entity_id_time_idx ON ticket_number (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS ticket_number_user_id_time_idx ON ticket_number (user_id, time DESC);

END;
