CREATE TABLE IF NOT EXISTS user_logins
(
    time    TIMESTAMP NOT NULL,
    user_id VARCHAR(255),
    value   INT
);
SELECT create_hypertable('user_logins', 'time', if_not_exists => TRUE);
