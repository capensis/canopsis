BEGIN;

CREATE INDEX IF NOT EXISTS action_log_type_value_type_value_id_idx ON action_log (type, value_type, value_id);

CREATE INDEX IF NOT EXISTS entity_categories_updated_at_idx ON entity_categories (updated_at);

COMMIT;