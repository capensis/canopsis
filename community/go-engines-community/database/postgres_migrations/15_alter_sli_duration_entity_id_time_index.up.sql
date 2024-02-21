BEGIN;

DROP INDEX IF EXISTS sli_duration_entity_id_time_idx;
CREATE UNIQUE INDEX IF NOT EXISTS sli_duration_entity_id_time_type_idx ON sli_duration (entity_id, time DESC, type);

END;
