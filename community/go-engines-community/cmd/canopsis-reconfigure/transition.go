package main

// transition.go implements the helpers that decide whether a cross-line
// Force() is required before running postgres migrations and record the
// outcome afterwards.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// newVersionThreshold is the smallest possible MMmmppsss value.
const newVersionThreshold uint = 10_00_00_000 // MMmmppsss with MM=10, mm=00, pp=00, sss=000

// versionFormat classifies a migration version integer.
type versionFormat int

const (
	versionFormatOld versionFormat = iota // plain sequential: 1, 2, 3 ...
	versionFormatNew                      // MMmmppsss, e.g. 260412001
)

const (
	versionFormatOldLabel = "old-sequential"
	versionFormatNewLabel = "new-MMmmppsss"
	versionFormatUnknown  = "unknown"
)

// parsedVersion holds the interpreted fields of a schema_migrations version.
type parsedVersion struct {
	raw         uint
	format      versionFormat
	releaseLine string // non-empty only for versionFormatNew, e.g. "26.04.x"
}

// migrationState is the current database migration state.
type migrationState struct {
	hasVersion bool
	version    uint
	dirty      bool
	parsed     *parsedVersion // non-nil when hasVersion == true
}

// dirInfo summarises the migration files present in one directory.
type dirInfo struct {
	maxOldAvailable uint // highest old-format sequence number; 0 = none found
	minNewAvailable uint // lowest new-format version number; 0 = none found
}

// registryRow is a single row from schema_migration_transition_registry.
type registryRow struct {
	releaseLine          string
	lastOldFormatVersion int64
}

// transitionPlan is the decision produced by decideTransitionPlan.
type transitionPlan struct {
	needsForce   bool
	forceVersion uint // only meaningful when needsForce == true

	// Diagnostic fields, populated regardless of needsForce.
	sourceReleaseLine string  // set when needsForce == true
	replayWindow      [2]uint // [first, last] inclusive; zero-value when !needsForce
}

// parseMigrationVersion classifies v and extracts its release line when it
// is new-format -- MMmmppsss:
//
//	MM  – release major  (2 digits, e.g. 26)
//	mm  – release minor  (2 digits, e.g. 04)
//	pp  – patch number   (2 digits, e.g. 12)
//	sss – within-patch sequence (3 digits, e.g. 001)
func parseMigrationVersion(v uint) parsedVersion {
	if v < newVersionThreshold {
		return parsedVersion{raw: v, format: versionFormatOld}
	}
	major := v / 10_000_000
	minor := (v / 100_000) % 100
	return parsedVersion{
		raw:         v,
		format:      versionFormatNew,
		releaseLine: fmt.Sprintf("%02d.%02d.x", major, minor),
	}
}

// loadMigrationState reads the current schema_migrations state via the
// already-opened migrate instance. If no version has ever been applied,
// the returned state has hasVersion == false and is not an error.
func loadMigrationState(m *migrate.Migrate) (migrationState, error) {
	v, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return migrationState{hasVersion: false}, nil
	}
	if err != nil {
		return migrationState{}, fmt.Errorf("reading migration version: %w", err)
	}
	parsed := parseMigrationVersion(v)
	return migrationState{
		hasVersion: true,
		version:    v,
		dirty:      dirty,
		parsed:     &parsed,
	}, nil
}

// scanMigrationDirectory lists all *.up.sql files in dirPath to return highest
// old-format version and lowest new-format version. The former is needed
// for deciding whether a cross-line transition is required.
func scanMigrationDirectory(dirPath string) (dirInfo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return dirInfo{}, fmt.Errorf("reading migration directory %q: %w", dirPath, err)
	}

	var info dirInfo
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		// Filename: "<version>_<description>.up.sql"
		prefix, _, ok := strings.Cut(name, "_")
		if !ok {
			continue
		}
		n, convErr := strconv.ParseUint(prefix, 10, 64)
		if convErr != nil {
			continue
		}
		v := uint(n)
		if v < newVersionThreshold {
			info.maxOldAvailable = max(info.maxOldAvailable, v)
		} else if info.minNewAvailable == 0 || v < info.minNewAvailable {
			info.minNewAvailable = v
		}
	}

	return info, nil
}

// loadTransitionRegistry queries schema_migration_transition_registry for the
// given release line.  Returns nil (and no error) when the row is not found,
// meaning this release line has no recorded checkpoint.
//
// The table is guaranteed to exist when this function is called: it is created
// by the first new-format migration file, which must already have run for the
// database to be in new-format state.
func loadTransitionRegistry(ctx context.Context, pool *pgxpool.Pool, releaseLine string) (*registryRow, error) {
	const q = `SELECT release_line, last_old_format_version
	           FROM schema_migration_transition_registry
	           WHERE release_line = $1`
	row := &registryRow{}
	err := pool.QueryRow(ctx, q, releaseLine).Scan(&row.releaseLine, &row.lastOldFormatVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("querying transition registry for %q: %w", releaseLine, err)
	}
	return row, nil
}

// hasTransitionHistory reports whether any cross-line Force() + Up() transition
// has ever been recorded for this database.
func hasTransitionHistory(ctx context.Context, pool *pgxpool.Pool) (exists bool, err error) {
	const q = `SELECT EXISTS (SELECT 1 FROM schema_migration_transition_history LIMIT 1)`
	if err = pool.QueryRow(ctx, q).Scan(&exists); err != nil {
		if postgres.IsTableMissingErr(err) {
			return false, nil
		}
		return false, fmt.Errorf("querying transition history: %w", err)
	}
	return exists, nil
}

// decideTransitionPlan resolves the transition plan.
//
// Inputs:
//
//	state    – current DB state; the caller must already have checked dirty
//	dir      – content of the target migration directory
//	registry – registry row for state.parsed.releaseLine; nil when state is
//	           old-format or has no version (no lookup needed in those cases)
//
// Returns an error only when the input combination is logically inconsistent
// (e.g. new-format source but no matching registry row).
func decideTransitionPlan(state migrationState, dir dirInfo, registry *registryRow) (transitionPlan, error) {
	// only Up(), no force needed.
	if !state.hasVersion || state.parsed.format == versionFormatOld {
		return transitionPlan{needsForce: false}, nil
	}

	// registry row must exist with valid lastOldFormatVersion for a new-format source database.
	if registry == nil {
		return transitionPlan{}, fmt.Errorf(
			"no transition registry row found for release line %q; cannot safely plan upgrade",
			state.parsed.releaseLine,
		)
	}

	if registry.lastOldFormatVersion < 0 {
		return transitionPlan{}, fmt.Errorf(
			"invalid last_old_format_version %d in transition registry for release line %q",
			registry.lastOldFormatVersion, registry.releaseLine,
		)
	}

	checkpoint := uint(registry.lastOldFormatVersion)
	if dir.maxOldAvailable <= checkpoint {
		// All old-format migrations that exist in the directory were already applied.
		// Up(), no force needed.
		return transitionPlan{needsForce: false}, nil
	}

	// The directory contains old-format migrations beyond the checkpoint.
	// Force() back to the checkpoint then Up() replays from checkpoint+1.
	return transitionPlan{
		needsForce:        true,
		forceVersion:      checkpoint,
		sourceReleaseLine: state.parsed.releaseLine,
		replayWindow:      [2]uint{checkpoint + 1, dir.maxOldAvailable},
	}, nil
}

// recordTransitionHistory inserts a permanent audit row into
// schema_migration_transition_history documenting a completed cross-line
// Force() + Up() transition.
//
// Must be called only after m.Up() returns without error. The history table
// is created by the first new-format migration file, so it is guaranteed to
// exist at this point.
func recordTransitionHistory(ctx context.Context, pool *pgxpool.Pool, state migrationState, plan transitionPlan, dir dirInfo) error {
	const q = `INSERT INTO schema_migration_transition_history
	           (source_version, source_release_line, forced_checkpoint_version, target_max_old_available)
	           VALUES ($1, $2, $3, $4)`
	_, err := pool.Exec(ctx, q,
		int64(state.version),
		plan.sourceReleaseLine,
		int64(plan.forceVersion),
		int64(dir.maxOldAvailable),
	)
	if err != nil {
		return fmt.Errorf("recording transition history: %w", err)
	}
	return nil
}

// printMigrationDiagnostics emits a structured INFO log summarising the
// current migration state, the directory scan, and the chosen plan.
func printMigrationDiagnostics(logger zerolog.Logger, state migrationState, dir dirInfo, plan transitionPlan) {
	e := logger.Info()
	strVersion := "none"
	if state.hasVersion {
		strVersion = strconv.FormatUint(uint64(state.version), 10)
		e = e.Bool("dirty", state.dirty).
			Str("version_format", versionFormatLabel(state.parsed.format))
		if state.parsed.format == versionFormatNew {
			e = e.Str("source_release_line", state.parsed.releaseLine)
		}
	}
	e = e.Str("current_version", strVersion)
	e = e.Uint("max_old_available", dir.maxOldAvailable)
	if dir.minNewAvailable != 0 {
		e = e.Uint("min_new_available", dir.minNewAvailable)
	}
	e = e.Bool("needs_force", plan.needsForce)
	if plan.needsForce {
		e = e.Uint("force_to_version", plan.forceVersion).
			Uint("replay_from", plan.replayWindow[0]).
			Uint("replay_to", plan.replayWindow[1])
	}
	e.Msg("migration diagnostics")
}

func upTransitionPreflight(ctx context.Context, pool *pgxpool.Pool, state migrationState, migrationDirectory string, unsafe, diagnose bool, logger zerolog.Logger) (transitionPlan, dirInfo, bool, error) {
	// Scan the directory to classify available migration versions.
	dir, err := scanMigrationDirectory(migrationDirectory)
	if err != nil {
		return transitionPlan{}, dirInfo{}, false, err
	}

	// Load the transition registry row when the source database is already
	// in new-format state, unless the unsafe override is active.
	var registry *registryRow
	if state.hasVersion && state.parsed.format == versionFormatNew && !unsafe {
		registry, err = loadTransitionRegistry(ctx, pool, state.parsed.releaseLine)
		if err != nil {
			return transitionPlan{}, dirInfo{}, false, err
		}
	}

	plan, planErr := decideTransitionPlan(state, dir, registry)
	if planErr != nil {
		if !unsafe {
			return transitionPlan{}, dirInfo{}, false, planErr
		}
		logger.Warn().Err(planErr).Msg("transition plan error bypassed due to -postgres-migration-unsafe flag")
	}

	if diagnose {
		logger.Info().Msg("diagnose mode: no migration will be executed, only diagnostics will be printed")
	}
	printMigrationDiagnostics(logger, state, dir, plan)
	if diagnose {
		return transitionPlan{}, dirInfo{}, true, nil
	}

	return plan, dir, false, nil
}

// -----------------------------------------------------------------------
// Internal helpers
// -----------------------------------------------------------------------

// versionFormatLabel returns a human-readable label for a versionFormat.
func versionFormatLabel(f versionFormat) string {
	switch f {
	case versionFormatOld:
		return versionFormatOldLabel
	case versionFormatNew:
		return versionFormatNewLabel
	default:
		return versionFormatUnknown
	}
}
