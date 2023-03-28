BEGIN;
DROP MATERIALIZED VIEW IF EXISTS instruction_execution_hourly;
DROP TABLE IF EXISTS instruction_execution;
DROP TABLE IF EXISTS instruction_mod_stats;
COMMIT;
