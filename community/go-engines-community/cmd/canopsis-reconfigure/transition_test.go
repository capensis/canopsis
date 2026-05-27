package main

import (
	"os"
	"path/filepath"
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
	stateNew := func(v uint) migrationState {
		p := parseMigrationVersion(v)
		return migrationState{hasVersion: true, version: v, parsed: &p}
	}
	stateOld := func(v uint) migrationState {
		p := parseMigrationVersion(v)
		return migrationState{hasVersion: true, version: v, parsed: &p}
	}
	reg := func(checkpoint int64) *registryRow {
		return &registryRow{releaseLine: "test.x", lastOldFormatVersion: checkpoint}
	}

	cases := []struct {
		name        string
		state       migrationState
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
			state:     migrationState{hasVersion: false},
			dir:       dirInfo{},
			wantForce: false,
		},
		{
			name:      "no-version-mixed-dir",
			state:     migrationState{hasVersion: false},
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
			dir:       dirInfo{maxOldAvailable: 26},
			registry:  reg(26),
			wantForce: false,
		},
		{
			// pure new-format directory: no old sequential migrations at all
			name:      "new-format-pure-new-dir",
			state:     stateNew(260401001),
			dir:       dirInfo{maxOldAvailable: 0, minNewAvailable: 261000001},
			registry:  reg(30),
			wantForce: false,
		},
		{
			// 25.10.x -> 26.04.x: tech_metrics checkpoint 14, maxOld 14 -> no force
			name:      "new-format-tech-metrics-no-force",
			state:     stateNew(251099001),
			dir:       dirInfo{maxOldAvailable: 14},
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
			if plan.needsForce != tc.wantForce {
				t.Errorf("needsForce: got %v, want %v", plan.needsForce, tc.wantForce)
			}
			if !tc.wantForce {
				return
			}
			if plan.forceVersion != tc.wantForceTo {
				t.Errorf("forceVersion: got %d, want %d", plan.forceVersion, tc.wantForceTo)
			}
			if plan.replayWindow != tc.wantWindow {
				t.Errorf("replayWindow: got %v, want %v", plan.replayWindow, tc.wantWindow)
			}
			if plan.sourceReleaseLine == "" {
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
		name       string
		files      []string
		wantMaxOld uint
		wantMinNew uint
	}{
		{
			name:       "empty-directory",
			files:      nil,
			wantMaxOld: 0,
			wantMinNew: 0,
		},
		{
			name:       "old-sequential-only",
			files:      []string{"1_init.up.sql", "26_add_column.up.sql", "28_alter_table.up.sql"},
			wantMaxOld: 28,
			wantMinNew: 0,
		},
		{
			name:       "new-format-only",
			files:      []string{"261000001_transition_tables.up.sql", "261000002_another.up.sql"},
			wantMaxOld: 0,
			wantMinNew: 261000001,
		},
		{
			name: "mixed-old-and-new",
			files: []string{
				"28_alter.up.sql", "30_add_index.up.sql",
				"260401001_transition.up.sql", "260401002_other.up.sql",
			},
			wantMaxOld: 30,
			wantMinNew: 260401001,
		},
		{
			name: "new-format-min-is-lowest",
			files: []string{
				"261000002_second.up.sql",
				"261000001_first.up.sql",
				"261000003_third.up.sql",
			},
			wantMaxOld: 0,
			wantMinNew: 261000001,
		},
		{
			name:       "down-files-ignored",
			files:      []string{"26_add.up.sql", "26_add.down.sql", "261000001_t.down.sql"},
			wantMaxOld: 26,
			wantMinNew: 0,
		},
		{
			name:       "files-without-underscore-ignored",
			files:      []string{"README.up.sql", "26_valid.up.sql"},
			wantMaxOld: 26,
			wantMinNew: 0,
		},
		{
			name:       "non-numeric-prefix-ignored",
			files:      []string{"abc_foo.up.sql", "26_valid.up.sql"},
			wantMaxOld: 26,
			wantMinNew: 0,
		},
		{
			// actual migration file present in the repository
			name:       "real-transition-file",
			files:      []string{"261000001_create_migration_transition_tables.up.sql"},
			wantMaxOld: 0,
			wantMinNew: 261000001,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			touch(t, dir, tc.files...)

			got, err := scanMigrationDirectory(dir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.maxOldAvailable != tc.wantMaxOld {
				t.Errorf("maxOldAvailable: got %d, want %d", got.maxOldAvailable, tc.wantMaxOld)
			}
			if got.minNewAvailable != tc.wantMinNew {
				t.Errorf("minNewAvailable: got %d, want %d", got.minNewAvailable, tc.wantMinNew)
			}
		})
	}
}

func TestScanMigrationDirectoryNotFound(t *testing.T) {
	_, err := scanMigrationDirectory("/nonexistent/path/that/cannot/exist")
	if err == nil {
		t.Error("expected error for nonexistent directory, got nil")
	}
}
