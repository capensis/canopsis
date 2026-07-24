-- Remove the 25.04.x registry row.
-- The registry and history tables are NOT dropped: they are permanent
-- infrastructure shared across all release lines. Dropping them would
-- erase audit records needed to block unintended Down() operations.
DELETE FROM schema_migration_transition_registry
WHERE release_line = '25.04.x';