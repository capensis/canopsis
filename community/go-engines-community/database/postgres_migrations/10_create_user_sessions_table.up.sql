BEGIN;

CREATE TABLE IF NOT EXISTS user_sessions
(
    time    TIMESTAMP NOT NULL,
    user_id VARCHAR(255),
    value   INT
);
SELECT create_hypertable('user_sessions', 'time', if_not_exists => TRUE);
CREATE INDEX IF NOT EXISTS user_sessions_user_id_time_idx ON user_sessions (user_id, time DESC);

END;
