package transition

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------
// parseMigrationVersion
// -----------------------------------------------------------------------

func TestParseMigrationVersion(t *testing.T) {
	cases := []struct {
		name        string
		input       uint
		wantFormat  versionFormat
		wantRaw     uint
		wantRelease string
	}{
		{"zero", 0, versionFormatOld, 0, ""},
		{"sequential-1", 1, versionFormatOld, 1, ""},
		{"sequential-26", 26, versionFormatOld, 26, ""},
		{"sequential-30", 30, versionFormatOld, 30, ""},
		{"just-below-threshold", newVersionThreshold - 1, versionFormatOld, newVersionThreshold - 1, ""},
		{"at-threshold", newVersionThreshold, versionFormatNew, newVersionThreshold, "10.00.x"},
		{"25.04.x", 250499001, versionFormatNew, 250499001, "25.04.x"},
		{"25.04.6-first", 250406001, versionFormatNew, 250406001, "25.04.x"},
		{"25.10.x", 251003001, versionFormatNew, 251003001, "25.10.x"},
		{"26.04.x", 260401001, versionFormatNew, 260401001, "26.04.x"},
		{"26.10.x", 261000001, versionFormatNew, 261000001, "26.10.x"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseMigrationVersion(tc.input)
			if got.raw != tc.wantRaw {
				t.Errorf("raw: got %d, want %d", got.raw, tc.wantRaw)
			}
			if got.format != tc.wantFormat {
				t.Errorf("format: got %d, want %d", got.format, tc.wantFormat)
			}
			if got.releaseLine != tc.wantRelease {
				t.Errorf("releaseLine: got %q, want %q", got.releaseLine, tc.wantRelease)
			}
		})
	}
}

// -----------------------------------------------------------------------
// decideTransitionPlan
// -----------------------------------------------------------------------

func TestDecideTransitionPlan(t *testing.T) {
	// helpers to keep test table concise
	stateNew := func(v uint) MigrationState {
		p := parseMigrationVersion(v)
		return MigrationState{HasVersion: true, Version: v, parsed: &p}
	}
	stateOld := func(v uint) MigrationState {
		p := parseMigrationVersion(v)
		return MigrationState{HasVersion: true, Version: v, parsed: &p}
	}
	reg := func(checkpoint int64) *registryRow {
		return &registryRow{releaseLine: "test.x", lastOldFormatVersion: checkpoint}
	}

	cases := []struct {
		name        string
		state       MigrationState
		dir         dirInfo
		registry    *registryRow
		wantForce   bool
		wantForceTo uint
		wantWindow  [2]uint
		wantErr     bool
	}{
		// --- no version: always plain Up() ---
		{
			name:      "no-version-empty-dir",
			state:     MigrationState{HasVersion: false},
			dir:       dirInfo{},
			wantForce: false,
		},
		{
			name:      "no-version-mixed-dir",
			state:     MigrationState{HasVersion: false},
			dir:       dirInfo{maxOldAvailable: 26, minNewAvailable: 250406001},
			wantForce: false,
		},
		// --- old-sequential source: always plain Up() ---
		{
			name:      "old-seq-no-new-files",
			state:     stateOld(24),
			dir:       dirInfo{maxOldAvailable: 26},
			wantForce: false,
		},
		{
			name:        "25.04.6-state-2000-dir-mixed",
			state:       stateOld(2000),
			dir:         dirInfo{maxOldAvailable: 26, minNewAvailable: 261001001},
			wantForce:   true,
			wantForceTo: 26,
		},
		{
			name:      "old-seq-mixed-dir",
			state:     stateOld(24),
			dir:       dirInfo{maxOldAvailable: 26, minNewAvailable: 250406001},
			wantForce: false,
		},
		// --- new-format source, registry missing: error ---
		{
			name:     "new-format-registry-nil",
			state:    stateNew(250499001),
			dir:      dirInfo{maxOldAvailable: 28},
			registry: nil,
			wantErr:  true,
		},
		// --- new-format source, invalid registry: error ---
		{
			name:     "new-format-registry-negative-checkpoint",
			state:    stateNew(250499001),
			dir:      dirInfo{maxOldAvailable: 28},
			registry: &registryRow{releaseLine: "25.04.x", lastOldFormatVersion: -1},
			wantErr:  true,
		},
		// --- new-format source, maxOldAvailable <= checkpoint: plain Up() ---
		{
			// maxOld equals checkpoint: no old migration above checkpoint exists
			name:      "new-format-maxold-equals-checkpoint",
			state:     stateNew(250499001),
			dir:       dirInfo{maxOldAvailable: 26, versionFound: true},
			registry:  reg(26),
			wantForce: false,
		},
		{
			// 25.04.99 -> 25.10.x: checkpoint 26, maxOld 26, no migration file found '250499001_*' -> force 26, then Up() to 251000001_*
			name:        "new-format-not-found",
			state:       stateNew(250499001),
			dir:         dirInfo{maxOldAvailable: 30, versionFound: false},
			registry:    reg(26),
			wantForce:   true,
			wantForceTo: 26,
			wantWindow:  [2]uint{27, 30},
		},
		{
			// pure new-format directory: no old sequential migrations at all
			name:      "new-format-pure-new-dir",
			state:     stateNew(260401001),
			dir:       dirInfo{maxOldAvailable: 0, minNewAvailable: 261000001, versionFound: true},
			registry:  reg(30),
			wantForce: false,
		},
		{
			// 25.10.x -> 26.04.x: tech_metrics checkpoint 14, maxOld 14 -> no force
			name:      "new-format-tech-metrics-no-force",
			state:     stateNew(251099001),
			dir:       dirInfo{maxOldAvailable: 14, versionFound: true},
			registry:  reg(14),
			wantForce: false,
		},
		// --- new-format source, maxOldAvailable > checkpoint: Force() + Up() ---
		{
			// 25.04.x -> 25.10.x: checkpoint 26, maxOld 28 -> replay 27..28
			name:        "new-format-25.04-to-25.10",
			state:       stateNew(250499001),
			dir:         dirInfo{maxOldAvailable: 28, minNewAvailable: 251003001},
			registry:    reg(26),
			wantForce:   true,
			wantForceTo: 26,
			wantWindow:  [2]uint{27, 28},
		},
		{
			// 25.10.x -> 26.04.x: checkpoint 28, maxOld 30 -> replay 29..30
			name:        "new-format-25.10-to-26.04",
			state:       stateNew(251099001),
			dir:         dirInfo{maxOldAvailable: 30, minNewAvailable: 260401001},
			registry:    reg(28),
			wantForce:   true,
			wantForceTo: 28,
			wantWindow:  [2]uint{29, 30},
		},
		{
			// 25.04.x -> 26.04.x (skipped 25.10): checkpoint 26, maxOld 30 -> replay 27..30
			name:        "new-format-25.04-to-26.04-wide-window",
			state:       stateNew(250499001),
			dir:         dirInfo{maxOldAvailable: 30, minNewAvailable: 260401001},
			registry:    reg(26),
			wantForce:   true,
			wantForceTo: 26,
			wantWindow:  [2]uint{27, 30},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := decideTransitionPlan(tc.state, tc.dir, tc.registry)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if plan.NeedsForce != tc.wantForce {
				t.Errorf("needsForce: got %v, want %v", plan.NeedsForce, tc.wantForce)
			}
			if !tc.wantForce {
				return
			}
			if plan.ForceVersion != tc.wantForceTo {
				t.Errorf("forceVersion: got %d, want %d", plan.ForceVersion, tc.wantForceTo)
			}
			if plan.replayWindow != tc.wantWindow {
				t.Errorf("replayWindow: got %v, want %v", plan.replayWindow, tc.wantWindow)
			}
			if plan.SourceReleaseLine == "" {
				t.Error("sourceReleaseLine should be set when needsForce is true")
			}
		})
	}
}

// -----------------------------------------------------------------------
// scanMigrationDirectory
// -----------------------------------------------------------------------

func TestScanMigrationDirectory(t *testing.T) {
	touch := func(t *testing.T, dir string, names ...string) {
		t.Helper()
		for _, name := range names {
			if err := os.WriteFile(filepath.Join(dir, name), nil, 0o600); err != nil {
				t.Fatalf("creating file %s: %v", name, err)
			}
		}
	}

	cases := []struct {
		name             string
		files            []string
		currentVersion   uint
		wantMaxOld       uint
		wantMinNew       uint
		wantMaxAll       uint
		wantVersionFound bool
	}{
		{
			name:           "empty-directory",
			files:          nil,
			currentVersion: 0,
			wantMaxOld:     0,
			wantMinNew:     0,
			wantMaxAll:     0,
		},
		{
			name:             "old-sequential-only",
			files:            []string{"1_init.up.sql", "26_add_column.up.sql", "28_alter_table.up.sql"},
			currentVersion:   28,
			wantMaxOld:       28,
			wantMinNew:       0,
			wantMaxAll:       28,
			wantVersionFound: true,
		},
		{
			name:             "new-format-only",
			files:            []string{"261000001_transition_tables.up.sql", "261000002_another.up.sql"},
			currentVersion:   261000001,
			wantMaxOld:       0,
			wantMinNew:       261000001,
			wantMaxAll:       261000002,
			wantVersionFound: true,
		},
		{
			name: "mixed-old-and-new",
			files: []string{
				"28_alter.up.sql", "30_add_index.up.sql",
				"260401001_transition.up.sql", "260401002_other.up.sql",
			},
			currentVersion:   260401001,
			wantMaxOld:       30,
			wantMinNew:       260401001,
			wantMaxAll:       260401002,
			wantVersionFound: true,
		},
		{
			name: "new-format-min-is-lowest",
			files: []string{
				"261000002_second.up.sql",
				"261000001_first.up.sql",
				"261000003_third.up.sql",
			},
			currentVersion:   261000001,
			wantMaxOld:       0,
			wantMinNew:       261000001,
			wantMaxAll:       261000003,
			wantVersionFound: true,
		},
		{
			name:             "down-files-ignored",
			files:            []string{"26_add.up.sql", "26_add.down.sql", "261000001_t.down.sql"},
			currentVersion:   26,
			wantMaxOld:       26,
			wantMinNew:       0,
			wantMaxAll:       26,
			wantVersionFound: true,
		},
		{
			name: "transition-artifact-excluded-from-buckets",
			files: []string{
				"26_last_old.up.sql",
				"2000_transition_artifact.up.sql",
				"250406001_first_canonical_new.up.sql",
			},
			currentVersion:   transitionArtifactVersion,
			wantMaxOld:       26,
			wantMinNew:       250406001,
			wantMaxAll:       250406001,
			wantVersionFound: true,
		},
		{
			name:             "files-without-underscore-ignored",
			files:            []string{"README.up.sql", "26_valid.up.sql"},
			currentVersion:   28,
			wantMaxOld:       26,
			wantMinNew:       0,
			wantMaxAll:       26,
			wantVersionFound: false,
		},
		{
			name:             "non-numeric-prefix-ignored",
			files:            []string{"abc_foo.up.sql", "26_valid.up.sql"},
			currentVersion:   30,
			wantMaxOld:       26,
			wantMinNew:       0,
			wantMaxAll:       26,
			wantVersionFound: false,
		},
		{
			// actual migration file present in the repository
			name:             "real-transition-file",
			files:            []string{"261000001_create_migration_transition_tables.up.sql"},
			currentVersion:   261000001,
			wantMaxOld:       0,
			wantMinNew:       261000001,
			wantMaxAll:       261000001,
			wantVersionFound: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			touch(t, dir, tc.files...)

			got, err := scanMigrationDirectory(dir, tc.currentVersion)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.maxOldAvailable != tc.wantMaxOld {
				t.Errorf("maxOldAvailable: got %d, want %d", got.maxOldAvailable, tc.wantMaxOld)
			}
			if got.minNewAvailable != tc.wantMinNew {
				t.Errorf("minNewAvailable: got %d, want %d", got.minNewAvailable, tc.wantMinNew)
			}
			if got.maxAvailable != tc.wantMaxAll {
				t.Errorf("maxAvailable: got %d, want %d", got.maxAvailable, tc.wantMaxAll)
			}
			if got.versionFound != tc.wantVersionFound {
				t.Errorf("versionFound: got %v, want %v", got.versionFound, tc.wantVersionFound)
			}
		})
	}
}

func TestScanMigrationDirectoryNotFound(t *testing.T) {
	_, err := scanMigrationDirectory("/nonexistent/path/that/cannot/exist", 0)
	if err == nil {
		t.Error("expected error for nonexistent directory, got nil")
	}
}

func TestEvaluateDownPermission(t *testing.T) {
	history := &TransitionHistoryRow{sourceVersion: 250499001, forcedCheckpointVersion: 26}

	cases := []struct {
		name       string
		state      MigrationState
		history    *TransitionHistoryRow
		unsafe     bool
		wantOK     bool
		wantReason string
	}{
		{
			name:       "clean-without-history",
			state:      MigrationState{HasVersion: true, Version: 26},
			wantOK:     true,
			wantReason: "",
		},
		{
			name:       "dirty-blocks-down",
			state:      MigrationState{HasVersion: true, Version: 26, Dirty: true},
			wantOK:     false,
			wantReason: downBlockReasonDirty,
		},
		{
			name:       "history-blocks-down",
			state:      MigrationState{HasVersion: true, Version: 250499001},
			history:    history,
			wantOK:     false,
			wantReason: downBlockReasonHistory,
		},
		{
			name:       "unsafe-bypasses-history",
			state:      MigrationState{HasVersion: true, Version: 250499001},
			history:    history,
			unsafe:     true,
			wantOK:     true,
			wantReason: "",
		},
		{
			name:       "dirty-still-blocks-with-unsafe",
			state:      MigrationState{HasVersion: true, Version: 250499001, Dirty: true},
			history:    history,
			unsafe:     true,
			wantOK:     false,
			wantReason: downBlockReasonDirty,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotOK, gotReason := evaluateDownPermission(tc.state, tc.history, tc.unsafe)
			if gotOK != tc.wantOK {
				t.Errorf("permitted: got %v, want %v", gotOK, tc.wantOK)
			}
			if gotReason != tc.wantReason {
				t.Errorf("reason: got %q, want %q", gotReason, tc.wantReason)
			}
		})
	}
}

func TestValidateDirectUpgradePath(t *testing.T) {
	stateOld := func(v uint) MigrationState {
		p := parseMigrationVersion(v)
		return MigrationState{HasVersion: true, Version: v, parsed: &p}
	}
	stateNew := func(v uint) MigrationState {
		p := parseMigrationVersion(v)
		return MigrationState{HasVersion: true, Version: v, parsed: &p}
	}
	guard := DirectUpgradeGuard{
		minSupportedSourceVersion:    26,
		blockedTargetMaxOldAvailable: 28,
		stagedPath:                   stagedPathDescr,
	}

	cases := []struct {
		name    string
		state   MigrationState
		dir     dirInfo
		wantErr bool
	}{
		{
			name:    "block-direct-23-to-25.10",
			state:   stateOld(23),
			dir:     dirInfo{maxOldAvailable: 28, minNewAvailable: 251003001},
			wantErr: true,
		},
		{
			name:    "block-direct-24-to-26.04",
			state:   stateOld(24),
			dir:     dirInfo{maxOldAvailable: 30, minNewAvailable: 260401001},
			wantErr: true,
		},
		{
			name:    "allow-25.04-latest-to-25.10",
			state:   stateOld(26),
			dir:     dirInfo{maxOldAvailable: 28, minNewAvailable: 251003001},
			wantErr: false,
		},
		{
			name:    "allow-24.10-to-25.04",
			state:   stateOld(24),
			dir:     dirInfo{maxOldAvailable: 26, minNewAvailable: 250406001},
			wantErr: false,
		},
		{
			name:    "allow-new-format-source",
			state:   stateNew(250499001),
			dir:     dirInfo{maxOldAvailable: 28, minNewAvailable: 251003001},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDirectUpgradePath(tc.state, tc.dir, guard)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), guard.stagedPath) {
					t.Fatalf("expected staged path in error, got %q", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewDownMigrationBlockedError(t *testing.T) {
	const flagPostgresMigrationUnsafe = "-postgres-migration-unsafe"

	th := &TransitionHistoryRow{sourceVersion: 250499001, forcedCheckpointVersion: 26}
	err := th.NewDownMigrationBlockedError(flagPostgresMigrationUnsafe)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "source_version=250499001") {
		t.Fatalf("expected source_version in error, got %q", err)
	}
	if !strings.Contains(err.Error(), "forced_checkpoint_version=26") {
		t.Fatalf("expected forced_checkpoint_version in error, got %q", err)
	}
	if !strings.Contains(err.Error(), flagPostgresMigrationUnsafe) {
		t.Fatalf("expected unsafe flag in error, got %q", err)
	}
}
