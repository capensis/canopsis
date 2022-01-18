BEGIN;

CREATE TABLE IF NOT EXISTS total_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('total_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS non_displayed_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('non_displayed_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS pbh_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('pbh_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS instruction_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('instruction_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS correlation_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('correlation_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS ticket_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    user_id   VARCHAR(255),
    value     INT
);
SELECT create_hypertable('ticket_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS ack_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    user_id   VARCHAR(255),
    value     INT
);
SELECT create_hypertable('ack_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS cancel_ack_alarm_number
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    user_id   VARCHAR(255),
    value     INT
);
SELECT create_hypertable('cancel_ack_alarm_number', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS ack_duration
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    user_id   VARCHAR(255),
    value     INT
);
SELECT create_hypertable('ack_duration', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS resolve_duration
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    value     INT
);
SELECT create_hypertable('resolve_duration', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS user_logins
(
    time    TIMESTAMP NOT NULL,
    user_id VARCHAR(255),
    value   INT
);
SELECT create_hypertable('user_logins', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS user_activity
(
    time    TIMESTAMP NOT NULL,
    user_id VARCHAR(255),
    value   INT
);
SELECT create_hypertable('user_activity', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS sli_duration
(
    time      TIMESTAMP NOT NULL,
    entity_id INT,
    type      SMALLINT,
    value     INT
);
SELECT create_hypertable('sli_duration', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS entities
(
    id              SERIAL PRIMARY KEY,
    custom_id       VARCHAR(500),
    name            VARCHAR(500),
    category        VARCHAR(255),
    impact_level    INT,
    type            VARCHAR(255),
    enabled         BOOLEAN,
    infos           JSONB,
    component_infos JSONB,
    component       VARCHAR(500),
    UNIQUE (custom_id)
);

CREATE TABLE IF NOT EXISTS users
(
    id       VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255),
    role     VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS metrics_criteria
(
    id      SERIAL PRIMARY KEY,
    type    INT,
    name    VARCHAR(255) UNIQUE,
    label   VARCHAR(255),
    enabled BOOLEAN
);

INSERT INTO metrics_criteria (type, name, label, enabled)
VALUES (1, 'username', 'username', true),
       (1, 'role', 'role', true),
       (0, 'category', 'category', true),
       (0, 'impact_level', 'impact_level', true);

COMMIT;
