BEGIN;
DROP MATERIALIZED VIEW IF EXISTS instruction_execution_hourly;
DROP TABLE IF EXISTS instruction_execution;
DROP TABLE IF EXISTS instruction_mod_stats;
DROP INDEX IF EXISTS idx_instruction_time;
DROP TABLE IF EXISTS instruction_execution_by_modified_on;
COMMIT;
