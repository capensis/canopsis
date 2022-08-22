BEGIN;

CREATE TABLE IF NOT EXISTS fifo_event
(
    time     TIMESTAMP   NOT NULL,
    type     VARCHAR(30) NOT NULL,
    interval INT         NOT NULL
);
SELECT create_hypertable('fifo_event', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('fifo_event', INTERVAL '1 day');

COMMIT;
