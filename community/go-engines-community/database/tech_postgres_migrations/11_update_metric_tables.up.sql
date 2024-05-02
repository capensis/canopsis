BEGIN;

ALTER TABLE cps_event ADD COLUMN is_ok_state BOOLEAN DEFAULT NULL;

ALTER TABLE fifo_event ADD COLUMN external_requests JSONB DEFAULT NULL;

ALTER TABLE che_event ADD COLUMN is_state_settings_updated BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE che_event ADD COLUMN executed_enrich_rules INT DEFAULT NULL;
ALTER TABLE che_event ADD COLUMN external_requests JSONB DEFAULT NULL;

ALTER TABLE axe_event ADD COLUMN is_counters_updated BOOLEAN DEFAULT NULL;
ALTER TABLE axe_event ADD COLUMN is_ok_state BOOLEAN DEFAULT NULL;

ALTER TABLE correlation_event ADD COLUMN matched_rules INT DEFAULT NULL;
ALTER TABLE correlation_event ADD COLUMN matched_rule_types VARCHAR(20)[] DEFAULT NULL;

ALTER TABLE dynamic_infos_event ADD COLUMN executed_rules INT DEFAULT NULL;

ALTER TABLE action_event ADD COLUMN executed_rules INT DEFAULT NULL;
ALTER TABLE action_event ADD COLUMN executed_webhooks INT DEFAULT NULL;

ALTER TABLE axe_periodical ADD COLUMN idle_events INT DEFAULT NULL;

COMMIT;
