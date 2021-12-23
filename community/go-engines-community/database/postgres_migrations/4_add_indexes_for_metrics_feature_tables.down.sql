BEGIN;

DROP INDEX IF EXISTS ack_alarm_number_entity_id_time_idx;
DROP INDEX IF EXISTS ack_alarm_number_user_id_time_idx;

DROP INDEX IF EXISTS ack_duration_entity_id_time_idx;
DROP INDEX IF EXISTS ack_duration_user_id_time_idx;

DROP INDEX IF EXISTS cancel_ack_alarm_number_entity_id_time_idx;
DROP INDEX IF EXISTS cancel_ack_alarm_number_user_id_time_idx;

DROP INDEX IF EXISTS correlation_alarm_number_entity_id_time_idx;

DROP INDEX IF EXISTS instruction_alarm_number_entity_id_time_idx;

DROP INDEX IF EXISTS non_displayed_alarm_number_entity_id_time_idx;

DROP INDEX IF EXISTS pbh_alarm_number_entity_id_time_idx;

DROP INDEX IF EXISTS resolve_duration_entity_id_time_idx;

DROP INDEX IF EXISTS sli_duration_entity_id_time_idx;

DROP INDEX IF EXISTS ticket_alarm_number_entity_id_time_idx;
DROP INDEX IF EXISTS ticket_alarm_number_user_id_time_idx;

DROP INDEX IF EXISTS total_alarm_number_entity_id_time_idx;

DROP INDEX IF EXISTS user_activity_user_id_time_idx;

DROP INDEX IF EXISTS entities_category_idx;
DROP INDEX IF EXISTS entities_name_idx;
DROP INDEX IF EXISTS entities_impact_level_idx;
DROP INDEX IF EXISTS entities_infos_idx;

DROP INDEX IF EXISTS users_username_idx;
DROP INDEX IF EXISTS users_role_idx;

COMMIT;
