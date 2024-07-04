BEGIN;

DROP INDEX IF EXISTS perf_data_name_name_idx;
DROP MATERIALIZED VIEW IF EXISTS perf_data_name;

CREATE TABLE IF NOT EXISTS perf_data_name
(
    name       VARCHAR(255) NOT NULL,
    unit       VARCHAR(5)   NOT NULL,
    updated_at TIMESTAMP    NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS perf_data_name_name_idx ON perf_data_name (name);

INSERT INTO perf_data_name (name, unit, updated_at)
SELECT name, unit, MAX(time)
FROM perf_data
GROUP BY name, unit;

END;
