-- Create the transition registry.
-- One row per release line that has switched to the MMmmppsss format.
-- last_old_format_version is the highest sequential migration guaranteed
-- to exist on databases that were created or last fully migrated on that line.
CREATE TABLE IF NOT EXISTS schema_migration_transition_registry (
  release_line            TEXT    PRIMARY KEY,
  last_old_format_version BIGINT  NOT NULL,
  created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CHECK (last_old_format_version >= 0)
);

-- Create the transition history.
-- Written by canopsis-reconfigure after a successful cross-line Force()+Up().
-- The presence of any row permanently blocks Down() for this database.
CREATE TABLE IF NOT EXISTS schema_migration_transition_history (
  id                        BIGSERIAL   PRIMARY KEY,
  source_version            BIGINT      NOT NULL,
  source_release_line       TEXT        NOT NULL,
  forced_checkpoint_version BIGINT      NOT NULL,
  target_max_old_available  BIGINT      NOT NULL,
  created_at                TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Register 25.10.x: sequential migrations 1-28 were all applied before
-- the first new-format migration shipped.
INSERT INTO schema_migration_transition_registry (release_line, last_old_format_version)
VALUES ('25.10.x', 28)
ON CONFLICT (release_line) DO UPDATE
  SET last_old_format_version = EXCLUDED.last_old_format_version;