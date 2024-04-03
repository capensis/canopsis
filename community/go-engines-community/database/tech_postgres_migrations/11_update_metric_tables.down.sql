BEGIN;

ALTER TABLE cps_event DROP COLUMN is_ok_state;

ALTER TABLE fifo_event DROP COLUMN external_requests;

ALTER TABLE che_event DROP COLUMN is_state_settings_updated;
ALTER TABLE che_event DROP COLUMN executed_enrich_rules;
ALTER TABLE che_event DROP COLUMN external_requests;

ALTER TABLE axe_event DROP COLUMN is_counters_updated;
ALTER TABLE axe_event DROP COLUMN is_ok_state;

ALTER TABLE correlation_event DROP COLUMN matched_rules;
ALTER TABLE correlation_event DROP COLUMN matched_rule_types;

ALTER TABLE dynamic_infos_event DROP COLUMN executed_rules;

ALTER TABLE action_event DROP COLUMN executed_rules;
ALTER TABLE action_event DROP COLUMN executed_webhooks;

ALTER TABLE axe_periodical DROP COLUMN idle_events;

COMMIT;
