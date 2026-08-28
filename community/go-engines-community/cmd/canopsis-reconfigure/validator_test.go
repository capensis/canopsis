package main

import (
	"testing"
)

func TestValidateRecoveryRequest(t *testing.T) {
	valid := recoveryRequest{forceVersion: 26, forceClearDirty: true}

	cases := []struct {
		name          string
		migrate       bool
		request       recoveryRequest
		modeExplicit  bool
		stepsExplicit bool
		diagnose      bool
		unsafe        bool
		wantErr       bool
	}{
		{
			name:    "not-requested",
			request: recoveryRequest{forceVersion: -1},
			wantErr: false,
		},
		{
			name:    "version-without-clear-dirty",
			migrate: true,
			request: recoveryRequest{forceVersion: 26},
			wantErr: true,
		},
		{
			name:    "clear-dirty-without-version",
			migrate: true,
			request: recoveryRequest{forceVersion: -1, forceClearDirty: true},
			wantErr: true,
		},
		{
			name:    "negative-version",
			migrate: true,
			request: recoveryRequest{forceVersion: -2, forceClearDirty: true},
			wantErr: true,
		},
		{
			name:    "requires-migrate-flag",
			request: valid,
			wantErr: true,
		},
		{
			name:         "reject-mode",
			migrate:      true,
			request:      valid,
			modeExplicit: true,
			wantErr:      true,
		},
		{
			name:          "reject-steps",
			migrate:       true,
			request:       valid,
			stepsExplicit: true,
			wantErr:       true,
		},
		{
			name:     "reject-diagnose",
			migrate:  true,
			request:  valid,
			diagnose: true,
			wantErr:  true,
		},
		{
			name:    "reject-unsafe",
			migrate: true,
			request: valid,
			unsafe:  true,
			wantErr: true,
		},
		{
			name:    "valid-request",
			migrate: true,
			request: valid,
			wantErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateRecoveryRequest(
				dbPostgresName+" recovery",
				tc.migrate,
				tc.request,
				tc.modeExplicit,
				tc.stepsExplicit,
				tc.diagnose,
				tc.unsafe,
				flagMigratePostgres,
				flagPostgresMigrationForceVersion,
				flagPostgresMigrationForceClearDirty,
				flagPostgresMigrationMode,
				flagPostgresMigrationSteps,
				flagPostgresMigrationUnsafe,
			)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
