BEGIN;

DROP INDEX IF EXISTS sli_duration_entity_id_time_idx;

CREATE TEMP TABLE sli_duration_tmp AS
SELECT time, entity_id, type, sum(value) as value
FROM sli_duration
GROUP BY time, entity_id, type;

TRUNCATE sli_duration;

INSERT INTO sli_duration (time, entity_id, type, value)
SELECT time, entity_id, type, value
FROM sli_duration_tmp;

DROP TABLE sli_duration_tmp;

CREATE UNIQUE INDEX IF NOT EXISTS sli_duration_entity_id_time_type_idx ON sli_duration (entity_id, time DESC, type);

END;
