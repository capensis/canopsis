BEGIN;

DROP INDEX IF EXISTS sli_duration_entity_id_time_start_type_idx;

ALTER TABLE sli_duration DROP COLUMN start;

CREATE UNIQUE INDEX IF NOT EXISTS sli_duration_entity_id_time_type_idx ON sli_duration (entity_id, time DESC, type);

END;
