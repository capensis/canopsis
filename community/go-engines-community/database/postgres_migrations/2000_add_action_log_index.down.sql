BEGIN;

DROP INDEX IF EXISTS action_log_type_value_type_value_id_idx;

DROP INDEX IF EXISTS entity_categories_updated_at_idx;

COMMIT;
