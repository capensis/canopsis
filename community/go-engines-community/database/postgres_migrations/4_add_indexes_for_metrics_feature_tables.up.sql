BEGIN;

CREATE INDEX IF NOT EXISTS ack_alarm_number_entity_id_time_idx ON ack_alarm_number (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS ack_alarm_number_user_id_time_idx ON ack_alarm_number (user_id, time DESC);

CREATE INDEX IF NOT EXISTS ack_duration_entity_id_time_idx ON ack_duration (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS ack_duration_user_id_time_idx ON ack_duration (user_id, time DESC);

CREATE INDEX IF NOT EXISTS cancel_ack_alarm_number_entity_id_time_idx ON cancel_ack_alarm_number (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS cancel_ack_alarm_number_user_id_time_idx ON cancel_ack_alarm_number (user_id, time DESC);

CREATE INDEX IF NOT EXISTS correlation_alarm_number_entity_id_time_idx ON correlation_alarm_number (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS instruction_alarm_number_entity_id_time_idx ON instruction_alarm_number (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS non_displayed_alarm_number_entity_id_time_idx ON non_displayed_alarm_number (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS pbh_alarm_number_entity_id_time_idx ON pbh_alarm_number (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS resolve_duration_entity_id_time_idx ON resolve_duration (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS sli_duration_entity_id_time_idx ON sli_duration (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS ticket_alarm_number_entity_id_time_idx ON ticket_alarm_number (entity_id, time DESC);
CREATE INDEX IF NOT EXISTS ticket_alarm_number_user_id_time_idx ON ticket_alarm_number (user_id, time DESC);

CREATE INDEX IF NOT EXISTS total_alarm_number_entity_id_time_idx ON total_alarm_number (entity_id, time DESC);

CREATE INDEX IF NOT EXISTS user_activity_user_id_time_idx ON user_activity (user_id, time DESC);

CREATE INDEX IF NOT EXISTS entities_category_idx ON entities (category);
CREATE INDEX IF NOT EXISTS entities_name_idx ON entities (name);
CREATE INDEX IF NOT EXISTS entities_impact_level_idx ON entities (impact_level);
CREATE INDEX IF NOT EXISTS entities_infos_idx ON entities USING GIN (infos);

CREATE INDEX IF NOT EXISTS users_username_idx ON users (username);
CREATE INDEX IF NOT EXISTS users_role_idx ON users (role);

COMMIT;
