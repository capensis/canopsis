BEGIN;
ALTER SEQUENCE "entities_id_seq" AS int MAXVALUE 2147483647 RESTART;

DELETE FROM entities WHERE id > 2147483647;
ALTER TABLE entities ALTER COLUMN id type int USING id::int;

DELETE FROM ack_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE ack_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM ack_duration WHERE entity_id > 2147483647;
ALTER TABLE ack_duration ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM cancel_ack_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE cancel_ack_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM correlation_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE correlation_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM instruction_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE instruction_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM manual_instruction_assigned_alarms WHERE entity_id > 2147483647;
ALTER TABLE manual_instruction_assigned_alarms ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM manual_instruction_executed_alarms WHERE entity_id > 2147483647;
ALTER TABLE manual_instruction_executed_alarms ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM non_displayed_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE non_displayed_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM not_acked_in_day_alarms WHERE entity_id > 2147483647;
ALTER TABLE not_acked_in_day_alarms ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM not_acked_in_four_hours_alarms WHERE entity_id > 2147483647;
ALTER TABLE not_acked_in_four_hours_alarms ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM not_acked_in_hour_alarms WHERE entity_id > 2147483647;
ALTER TABLE not_acked_in_hour_alarms ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM pbh_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE pbh_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM perf_data WHERE entity_id > 2147483647;
ALTER TABLE perf_data ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM resolve_duration WHERE entity_id > 2147483647;
ALTER TABLE resolve_duration ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM sli_duration WHERE entity_id > 2147483647;
ALTER TABLE sli_duration ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM ticket_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE ticket_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM ticket_number WHERE entity_id > 2147483647;
ALTER TABLE ticket_number ALTER COLUMN entity_id type int USING entity_id::int;

DELETE FROM total_alarm_number WHERE entity_id > 2147483647;
ALTER TABLE total_alarm_number ALTER COLUMN entity_id type int USING entity_id::int;

COMMIT;
