BEGIN;

CREATE TABLE IF NOT EXISTS fifo_queue
(
    time   TIMESTAMP NOT NULL,
    length INT       NOT NULL
);
SELECT create_hypertable('fifo_queue', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('fifo_queue', INTERVAL '30 days');

CREATE TABLE IF NOT EXISTS fifo_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('fifo_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('fifo_event', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS che_event
(
    time                TIMESTAMP   NOT NULL,
    interval            INT         NOT NULL,
    type                VARCHAR(30) NOT NULL,
    entity_type         VARCHAR(30) NOT NULL,
    is_new_entity       BOOLEAN     NOT NULL,
    is_infos_updated    BOOLEAN     NOT NULL,
    is_services_updated BOOLEAN     NOT NULL
);
SELECT create_hypertable('che_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('che_event', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS axe_event
(
    time              TIMESTAMP   NOT NULL,
    interval          INT         NOT NULL,
    type              VARCHAR(30) NOT NULL,
    entity_type       VARCHAR(30) NOT NULL,
    alarm_change_type VARCHAR(30) NOT NULL
);
SELECT create_hypertable('axe_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('axe_event', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS axe_periodical
(
    time     TIMESTAMP NOT NULL,
    interval INT       NOT NULL
);
SELECT create_hypertable('axe_periodical', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('axe_periodical', INTERVAL '7 days');

CREATE TABLE IF NOT EXISTS pbehavior_periodical
(
    time     TIMESTAMP NOT NULL,
    interval INT       NOT NULL
);
SELECT create_hypertable('pbehavior_periodical', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('pbehavior_periodical', INTERVAL '7 days');

CREATE TABLE IF NOT EXISTS che_infos
(
    time TIMESTAMP   NOT NULL,
    name VARCHAR(30) NOT NULL
);
SELECT create_hypertable('che_infos', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('che_infos', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS api_requests
(
    time     TIMESTAMP NOT NULL,
    interval INT       NOT NULL
);
SELECT create_hypertable('api_requests', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('api_requests', INTERVAL '7 days');

CREATE TABLE IF NOT EXISTS pbehavior_event
(
    time              TIMESTAMP   NOT NULL,
    interval          INT         NOT NULL,
    type              VARCHAR(30) NOT NULL
);
SELECT create_hypertable('pbehavior_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('pbehavior_event', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS correlation_event
(
    time              TIMESTAMP   NOT NULL,
    interval          INT         NOT NULL,
    type              VARCHAR(30) NOT NULL
);
SELECT create_hypertable('correlation_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('correlation_event', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS service_event
(
    time              TIMESTAMP   NOT NULL,
    interval          INT         NOT NULL,
    type              VARCHAR(30) NOT NULL
);
SELECT create_hypertable('service_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('service_event', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS dynamic_infos_event
(
    time              TIMESTAMP   NOT NULL,
    interval          INT         NOT NULL,
    type              VARCHAR(30) NOT NULL
);
SELECT create_hypertable('dynamic_infos_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('dynamic_infos_event', INTERVAL '1 day');

CREATE TABLE IF NOT EXISTS action_event
(
    time              TIMESTAMP   NOT NULL,
    interval          INT         NOT NULL,
    type              VARCHAR(30) NOT NULL
    );
SELECT create_hypertable('action_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('action_event', INTERVAL '1 day');

COMMIT;
