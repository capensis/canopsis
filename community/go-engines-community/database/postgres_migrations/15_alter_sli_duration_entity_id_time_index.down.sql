BEGIN;

DROP INDEX IF EXISTS sli_duration_entity_id_time_type_idx;
CREATE INDEX IF NOT EXISTS sli_duration_entity_id_time_idx ON sli_duration (entity_id, time DESC);

END;
