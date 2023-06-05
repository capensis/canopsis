BEGIN;

CREATE TABLE IF NOT EXISTS cps_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('cps_event', 'time', if_not_exists => TRUE);

-- FIFO
CREATE TABLE IF NOT EXISTS fifo_queue
(
    time   TIMESTAMP NOT NULL,
    length INT       NOT NULL
);
SELECT create_hypertable('fifo_queue', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS fifo_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('fifo_event', 'time', if_not_exists => TRUE);
-- END FIFO

-- CHE
CREATE TABLE IF NOT EXISTS che_event
(
    time                TIMESTAMP   NOT NULL,
    type                VARCHAR(30) NOT NULL,
    entity_type         VARCHAR(30) NOT NULL,
    interval            INT         NOT NULL,
    is_new_entity       BOOLEAN     NOT NULL,
    is_infos_updated    BOOLEAN     NOT NULL,
    is_services_updated BOOLEAN     NOT NULL
);
SELECT create_hypertable('che_event', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS che_infos
(
    time TIMESTAMP   NOT NULL,
    name VARCHAR(30) NOT NULL
);
SELECT create_hypertable('che_infos', 'time', if_not_exists => TRUE);
-- END CHE

-- AXE
CREATE TABLE IF NOT EXISTS axe_event
(
    time              TIMESTAMP   NOT NULL,
    type              VARCHAR(30) NOT NULL,
    entity_type       VARCHAR(30) NOT NULL,
    interval          INT         NOT NULL,
    alarm_change_type VARCHAR(30) NOT NULL
);
SELECT create_hypertable('axe_event', 'time', if_not_exists => TRUE);

CREATE TABLE IF NOT EXISTS axe_periodical
(
    time     TIMESTAMP NOT NULL,
    interval INT       NOT NULL,
    events   INT       NOT NULL
);
SELECT create_hypertable('axe_periodical', 'time', if_not_exists => TRUE);
-- END AXE

-- PBEHAVIOR
CREATE TABLE IF NOT EXISTS pbehavior_periodical
(
    time       TIMESTAMP NOT NULL,
    interval   INT       NOT NULL,
    events     INT       NOT NULL,
    entities   INT       NOT NULL,
    pbehaviors INT       NOT NULL
);
SELECT create_hypertable('pbehavior_periodical', 'time', if_not_exists => TRUE);
-- END PBEHAVIOR

-- CORRELATION
CREATE TABLE IF NOT EXISTS correlation_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('correlation_event', 'time', if_not_exists => TRUE);
-- END CORRELATION

-- SERVICE
CREATE TABLE IF NOT EXISTS service_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('service_event', 'time', if_not_exists => TRUE);
-- END SERVICE

-- DYNAMIC INFOS
CREATE TABLE IF NOT EXISTS dynamic_infos_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('dynamic_infos_event', 'time', if_not_exists => TRUE);
-- END DYNAMIC INFOS

-- ACTION
CREATE TABLE IF NOT EXISTS action_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('action_event', 'time', if_not_exists => TRUE);
-- END ACTION

-- API
CREATE TABLE IF NOT EXISTS api_requests
(
    time     TIMESTAMP    NOT NULL,
    method   VARCHAR(7)   NOT NULL,
    url      VARCHAR(255) NOT NULL,
    interval INT          NOT NULL
);
SELECT create_hypertable('api_requests', 'time', if_not_exists => TRUE);
-- END API

COMMIT;
