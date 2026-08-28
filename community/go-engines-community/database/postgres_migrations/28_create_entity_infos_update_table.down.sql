BEGIN;

SELECT remove_compression_policy(to_regclass('entity_infos_update'), if_exists => TRUE)
WHERE to_regclass('entity_infos_update') IS NOT NULL;
ALTER TABLE IF EXISTS entity_infos_update
    SET (timescaledb.compress = false);

SELECT remove_retention_policy(to_regclass('entity_infos_update'), if_exists => TRUE)
WHERE to_regclass('entity_infos_update') IS NOT NULL;

DROP INDEX IF EXISTS entity_infos_update_entity_id_time_idx;

DROP TABLE IF EXISTS entity_infos_update;

COMMIT;
