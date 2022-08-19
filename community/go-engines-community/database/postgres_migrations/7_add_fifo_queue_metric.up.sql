BEGIN;

CREATE TABLE IF NOT EXISTS fifo_queue
(
    time   TIMESTAMP NOT NULL,
    length INT       NOT NULL
);
SELECT create_hypertable('fifo_queue', 'time', if_not_exists => TRUE);
SELECT add_retention_policy('fifo_queue', INTERVAL '30 days');

COMMIT;
