BEGIN;

DROP INDEX IF EXISTS perf_data_name_name_idx;
DROP TABLE IF EXISTS perf_data_name;

CREATE MATERIALIZED VIEW IF NOT EXISTS perf_data_name AS
SELECT name, unit
FROM perf_data
GROUP BY name, unit;
CREATE UNIQUE INDEX IF NOT EXISTS perf_data_name_name_idx ON perf_data_name (name);

END;
