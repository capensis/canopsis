BEGIN;

DROP INDEX IF EXISTS ticket_number_entity_id_time_idx;
DROP INDEX IF EXISTS ticket_number_user_id_time_idx;

DROP TABLE IF EXISTS ticket_number;

END;
