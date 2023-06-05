CREATE TABLE IF NOT EXISTS export
(
    id        SERIAL PRIMARY KEY,
    filepath  VARCHAR(500) NOT NULL,
    status    INT          NOT NULL,
    created   TIMESTAMP    NOT NULL,
    started   TIMESTAMP,
    completed TIMESTAMP
);
