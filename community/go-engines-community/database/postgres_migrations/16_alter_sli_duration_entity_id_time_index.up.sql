BEGIN;

DROP INDEX IF EXISTS sli_duration_entity_id_time_type_idx;

ALTER TABLE sli_duration ADD COLUMN start TIMESTAMP DEFAULT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS sli_duration_entity_id_time_start_type_idx ON sli_duration (entity_id, time DESC, start, type);

END;
