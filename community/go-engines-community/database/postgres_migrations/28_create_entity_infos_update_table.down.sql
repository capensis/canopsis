BEGIN;

SELECT remove_compression_policy('entity_infos_update');
ALTER TABLE entity_infos_update
    SET (timescaledb.compress = false);

SELECT remove_retention_policy('entity_infos_update');

DROP INDEX IF EXISTS entity_infos_update_entity_id_time_idx;

DROP TABLE IF EXISTS entity_infos_update;

COMMIT;
