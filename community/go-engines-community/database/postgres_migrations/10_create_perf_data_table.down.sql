BEGIN;

DROP INDEX IF EXISTS perf_data_name_name_idx;
DROP MATERIALIZED VIEW IF EXISTS perf_data_name;
DROP INDEX IF EXISTS perf_data_name_entity_id_time_idx;
DROP TABLE IF EXISTS perf_data;

END;
