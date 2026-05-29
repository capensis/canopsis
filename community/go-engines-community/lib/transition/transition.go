package transition

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

// versionFormat classifies a migration version integer.
type versionFormat int

const (
	versionFormatOld versionFormat = iota // plain sequential: 1, 2, 3 ...
	versionFormatNew                      // MMmmppsss, e.g. 260412001
)

const (
	// newVersionThreshold is the smallest possible MMmmppsss value with MM=10, mm=00, pp=00, sss=000.
	newVersionThreshold uint = 10_00_00_000

	versionFormatOldLabel = "old-sequential"
	versionFormatNewLabel = "new-MMmmppsss"
	versionFormatUnknown  = "unknown"

	downBlockReasonDirty   = "dirty migration state"
	downBlockReasonHistory = "cross-line transition history recorded"

	downMigrationBlockedPrefix = "down migration blocked: a cross-line Force() transition has been recorded for this database"
)

// parsedVersion holds the interpreted fields of a schema_migrations version.
type parsedVersion struct {
	raw         uint
	format      versionFormat
	releaseLine string // non-empty only for versionFormatNew, e.g. "26.04.x"
}

// MigrationState is the current database migration state.
type MigrationState struct {
	parsed     *parsedVersion // non-nil when hasVersion == true
	Version    uint
	HasVersion bool
	Dirty      bool
}

// dirInfo summarises the migration files present in one directory.
type dirInfo struct {
	maxOldAvailable uint // highest old-format sequence number; 0 = none found
	minNewAvailable uint // lowest new-format version number; 0 = none found
	maxAvailable    uint // highest migration version in the directory; 0 = none found
	versionFound    bool // current version migration file found in the directory; only meaningful when HasVersion == true
}

// registryRow is a single row from schema_migration_transition_registry.
type registryRow struct {
	releaseLine          string
	lastOldFormatVersion int64
}

// TransitionHistoryRow is the latest row from schema_migration_transition_history.
type TransitionHistoryRow struct {
	sourceReleaseLine       string
	sourceVersion           int64
	forcedCheckpointVersion int64
	targetMaxOldAvailable   int64
}

func (t *TransitionHistoryRow) NewDownMigrationBlockedError(unsafeFlagRef string) error {
	return fmt.Errorf(
		"%s (source_version=%d, forced_checkpoint_version=%d); restore from a pre-upgrade backup or pass %s to override",
		downMigrationBlockedPrefix,
		t.sourceVersion,
		t.forcedCheckpointVersion,
		unsafeFlagRef,
	)
}

type DirectUpgradeGuard struct {
	minSupportedSourceVersion    uint
	blockedTargetMaxOldAvailable uint
	stagedPath                   string
}

var (
	stagedPathDescr            = "23.10.x -> latest 24.04.x -> latest 24.10.x -> latest 25.04.x -> target 25.10.x or later"
	PostgresDirectUpgradeGuard = DirectUpgradeGuard{
		minSupportedSourceVersion:    26,
		blockedTargetMaxOldAvailable: 28,
		stagedPath:                   stagedPathDescr,
	}
	TechPostgresDirectUpgradeGuard = DirectUpgradeGuard{
		minSupportedSourceVersion:    13,
		blockedTargetMaxOldAvailable: 14,
		stagedPath:                   stagedPathDescr,
	}
)

// TransitionPlan is the decision produced by decideTransitionPlan.
type TransitionPlan struct {
	NeedsForce   bool
	ForceVersion uint // only meaningful when NeedsForce == true

	// Diagnostic fields, populated regardless of NeedsForce.
	SourceReleaseLine string  // set when NeedsForce == true
	replayWindow      [2]uint // [first, last] inclusive; zero-value when !NeedsForce
}

type validateDirectUpgradePathError error

func newValidateDirectUpgradePathError(version, maxOldAvailable uint, stagedPath string) validateDirectUpgradePathError {
	return fmt.Errorf(
		"unsupported direct upgrade attempt: old-sequential version %d cannot be upgraded directly to a target shipping old sequential migrations through %d; follow staged path %s",
		version,
		maxOldAvailable,
		stagedPath,
	)
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

// LoadMigrationState reads the current schema_migrations state via the
// already-opened migrate instance. If no version has ever been applied,
// the returned state has hasVersion == false and is not an error.
func LoadMigrationState(m *migrate.Migrate) (MigrationState, error) {
	v, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return MigrationState{HasVersion: false}, nil
	}
	if err != nil {
		return MigrationState{}, fmt.Errorf("reading migration version: %w", err)
	}
	parsed := parseMigrationVersion(v)
	return MigrationState{
		HasVersion: true,
		Version:    v,
		Dirty:      dirty,
		parsed:     &parsed,
	}, nil
}

// scanMigrationDirectory lists all *.up.sql files in dirPath to return highest
// old-format version and lowest new-format version. The former is needed
// for deciding whether a cross-line transition is required.
func scanMigrationDirectory(dirPath string, findVersion uint) (dirInfo, error) {
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
		info.maxAvailable = max(info.maxAvailable, v)
		if v < newVersionThreshold {
			info.maxOldAvailable = max(info.maxOldAvailable, v)
		} else if info.minNewAvailable == 0 || v < info.minNewAvailable {
			info.minNewAvailable = v
		}
		if v == findVersion {
			info.versionFound = true
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

// LoadLatestTransitionHistory returns the latest completed cross-line Force() +
// Up() transition recorded for this database. Nil means no row was found.
func LoadLatestTransitionHistory(ctx context.Context, pool *pgxpool.Pool) (*TransitionHistoryRow, error) {
	const q = `SELECT source_version, source_release_line, forced_checkpoint_version, target_max_old_available
	           FROM schema_migration_transition_history
	           ORDER BY id DESC
	           LIMIT 1`
	row := &TransitionHistoryRow{}
	err := pool.QueryRow(ctx, q).Scan(
		&row.sourceVersion,
		&row.sourceReleaseLine,
		&row.forcedCheckpointVersion,
		&row.targetMaxOldAvailable,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || postgres.IsTableMissingErr(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("querying transition history: %w", err)
	}

	return row, nil
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
func decideTransitionPlan(state MigrationState, dir dirInfo, registry *registryRow) (TransitionPlan, error) {
	// special case for version 2000: it was the first new-format version but the registry was introduced later.
	// It's effectively equivalent to an old-format version for transition planning purposes, so treat it as such.
	if state.HasVersion && state.Version == 2000 {
		return TransitionPlan{NeedsForce: true, ForceVersion: 26, SourceReleaseLine: "25.04.x"}, nil
	}
	// only Up(), no force needed.
	if !state.HasVersion || state.parsed.format == versionFormatOld {
		return TransitionPlan{NeedsForce: false}, nil
	}

	// registry row must exist with valid lastOldFormatVersion for a new-format source database.
	if registry == nil {
		return TransitionPlan{}, fmt.Errorf(
			"no transition registry row found for release line %q; cannot safely plan upgrade",
			state.parsed.releaseLine,
		)
	}

	if registry.lastOldFormatVersion < 0 {
		return TransitionPlan{}, fmt.Errorf(
			"invalid last_old_format_version %d in transition registry for release line %q",
			registry.lastOldFormatVersion, registry.releaseLine,
		)
	}

	checkpoint := uint(registry.lastOldFormatVersion)
	if dir.maxOldAvailable <= checkpoint && dir.versionFound {
		// All old-format migrations that exist in the directory were already applied.
		// Up(), no force needed.
		return TransitionPlan{NeedsForce: false}, nil
	}

	// The directory contains old-format migrations beyond the checkpoint.
	// Force() back to the checkpoint then Up() replays from checkpoint+1.
	return TransitionPlan{
		NeedsForce:        true,
		ForceVersion:      checkpoint,
		SourceReleaseLine: state.parsed.releaseLine,
		replayWindow:      [2]uint{checkpoint + 1, dir.maxOldAvailable},
	}, nil
}

// RecordTransitionHistory inserts a permanent audit row into
// schema_migration_transition_history documenting a completed cross-line
// Force() + Up() transition.
//
// Must be called only after m.Up() returns without error. The history table
// is created by the first new-format migration file, so it is guaranteed to
// exist at this point.
func RecordTransitionHistory(ctx context.Context, pool *pgxpool.Pool, state MigrationState, plan TransitionPlan, dir dirInfo) error {
	const q = `INSERT INTO schema_migration_transition_history
	           (source_version, source_release_line, forced_checkpoint_version, target_max_old_available)
	           VALUES ($1, $2, $3, $4)`
	_, err := pool.Exec(ctx, q,
		int64(state.Version),
		plan.SourceReleaseLine,
		int64(plan.ForceVersion),
		int64(dir.maxOldAvailable),
	)
	if err != nil {
		return fmt.Errorf("recording transition history: %w", err)
	}
	return nil
}

func evaluateDownPermission(state MigrationState, history *TransitionHistoryRow, unsafe bool) (bool, string) {
	if state.HasVersion && state.Dirty {
		return false, downBlockReasonDirty
	}
	if history != nil && !unsafe {
		return false, downBlockReasonHistory
	}

	return true, ""
}

func validateDirectUpgradePath(state MigrationState, dir dirInfo, guard DirectUpgradeGuard) error {
	if guard.minSupportedSourceVersion == 0 || guard.blockedTargetMaxOldAvailable == 0 {
		return nil
	}
	if !state.HasVersion || state.parsed.format != versionFormatOld {
		return nil
	}
	if state.Version >= guard.minSupportedSourceVersion {
		return nil
	}
	if dir.maxOldAvailable < guard.blockedTargetMaxOldAvailable {
		return nil
	}

	return newValidateDirectUpgradePathError(state.Version, dir.maxOldAvailable, guard.stagedPath)
}

// printMigrationDiagnostics emits a structured INFO log summarising the
// current migration state, the directory scan, the chosen plan, and the
// current Down() guard status.
func printMigrationDiagnostics(state MigrationState, dir dirInfo, plan TransitionPlan, history *TransitionHistoryRow, downPermitted bool, downBlockReason string, logger zerolog.Logger) {
	e := logger.Info()
	strVersion := "none"
	if state.HasVersion {
		strVersion = strconv.FormatUint(uint64(state.Version), 10)
		e = e.Bool("dirty", state.Dirty).
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
	if dir.maxAvailable != 0 {
		e = e.Uint("target_final_version_available", dir.maxAvailable)
	}
	e = e.Bool("transition_history_exists", history != nil)
	if history != nil {
		e = e.Int64("history_source_version", history.sourceVersion).
			Str("history_source_release_line", history.sourceReleaseLine).
			Int64("history_forced_checkpoint_version", history.forcedCheckpointVersion).
			Int64("history_target_max_old_available", history.targetMaxOldAvailable)
	}
	e = e.Bool("down_permitted", downPermitted)
	if downBlockReason != "" {
		e = e.Str("down_block_reason", downBlockReason)
	}
	e = e.Bool("needs_force", plan.NeedsForce)
	if plan.NeedsForce {
		e = e.Uint("force_to_version", plan.ForceVersion).
			Uint("replay_from", plan.replayWindow[0]).
			Uint("replay_to", plan.replayWindow[1])
	}
	e.Msg("migration diagnostics")
}

func UpTransitionPreflight(ctx context.Context, pool *pgxpool.Pool, state MigrationState, migrationDirectory string, unsafe bool, guard DirectUpgradeGuard, logger zerolog.Logger) (TransitionPlan, dirInfo, error) {
	// Scan the directory to classify available migration versions.
	dir, err := scanMigrationDirectory(migrationDirectory, state.Version)
	if err != nil {
		return TransitionPlan{}, dirInfo{}, err
	}

	// Load the transition registry row when the source database is already
	// in new-format state, unless the unsafe override is active.
	var registry *registryRow
	if state.HasVersion && state.parsed.format == versionFormatNew && !unsafe {
		registry, err = loadTransitionRegistry(ctx, pool, state.parsed.releaseLine)
		if err != nil {
			return TransitionPlan{}, dirInfo{}, err
		}
	}

	plan, planErr := decideTransitionPlan(state, dir, registry)
	if planErr != nil {
		if !unsafe {
			return TransitionPlan{}, dirInfo{}, planErr
		}
		logger.Warn().Err(planErr).Msg("transition plan error bypassed due to unsafe migration flag")
	}

	return plan, dir, validateDirectUpgradePath(state, dir, guard)
}

func DiagnoseMigrationState(ctx context.Context, pool *pgxpool.Pool, state MigrationState, migrationDirectory string, guard DirectUpgradeGuard, unsafe bool, logger zerolog.Logger) error {
	plan, dir, err := UpTransitionPreflight(ctx, pool, state, migrationDirectory, unsafe, guard, logger)
	if err != nil {
		if _, ok := errors.AsType[validateDirectUpgradePathError](err); !ok {
			return err
		}
		logger.Warn().Err(err).Msg("diagnose detected blocked unsupported direct upgrade path")
	}

	history, err := LoadLatestTransitionHistory(ctx, pool)
	if err != nil {
		return err
	}

	downPermitted, downBlockReason := evaluateDownPermission(state, history, unsafe)

	logger.Info().Msg("diagnose mode: no migration will be executed")
	printMigrationDiagnostics(state, dir, plan, history, downPermitted, downBlockReason, logger)

	return nil
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
