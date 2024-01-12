BEGIN;
ALTER SEQUENCE "entities_id_seq" AS bigint MAXVALUE 9223372036854775807;

ALTER TABLE entities ALTER COLUMN id type bigint USING id::bigint;

ALTER TABLE ack_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE ack_duration ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE cancel_ack_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE correlation_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE instruction_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE manual_instruction_assigned_alarms ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE manual_instruction_executed_alarms ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE non_displayed_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE not_acked_in_day_alarms ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE not_acked_in_four_hours_alarms ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE not_acked_in_hour_alarms ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE pbh_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE resolve_duration ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE sli_duration ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE ticket_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;
ALTER TABLE total_alarm_number ALTER COLUMN entity_id type bigint USING entity_id::bigint;

COMMIT;
