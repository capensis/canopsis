package externaldatatable

import (
	"slices"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
)

func TestSyncColumnConfigsOrder(t *testing.T) {
	testCases := []struct {
		name                 string
		oldConfigs           []ColumnConfig
		newConfigs           []ColumnConfig
		expectedNewConfigs   []ColumnConfig
		expectedErrorMessage string
	}{
		{
			name: "given configs in same order expect no changes",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
			newConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
			expectedNewConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
		},
		{
			name: "given configs in different order expect reordering",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
			newConfigs: []ColumnConfig{
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
			},
			expectedNewConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
		},
		{
			name: "given configs completely reversed expect complete reordering",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
			newConfigs: []ColumnConfig{
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col1", Type: externaldata.ColumnTypeString},
			},
			expectedNewConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
		},
		{
			name: "given single config expect no changes",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
			},
			newConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
			},
			expectedNewConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
			},
		},
		{
			name:               "given empty configs expect no changes",
			oldConfigs:         []ColumnConfig{},
			newConfigs:         []ColumnConfig{},
			expectedNewConfigs: []ColumnConfig{},
		},
		{
			name: "given configs with different types but same names expect reordering with type preservation",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
			},
			newConfigs: []ColumnConfig{
				{Name: "col2", Type: externaldata.ColumnTypeBoolean},
				{Name: "col1", Type: externaldata.ColumnTypeDateTime},
			},
			expectedNewConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeDateTime},
				{Name: "col2", Type: externaldata.ColumnTypeBoolean},
			},
		},
		{
			name: "given different length configs expect error",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
			},
			newConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
			},
			expectedErrorMessage: "column config count mismatch",
		},
		{
			name: "given new configs longer than old configs expect error",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
			},
			newConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
			},
			expectedErrorMessage: "column config count mismatch",
		},
		{
			name: "given missing column in new configs expect error",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
			},
			newConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
			expectedErrorMessage: "no such column \"col3\"",
		},
		{
			name: "given duplicate column names in new configs expect error",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
			},
			newConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col1", Type: externaldata.ColumnTypeNumber},
			},
			expectedErrorMessage: "duplicate column name \"col1\"",
		},
		{
			name: "given large number of configs expect reordering",
			oldConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
				{Name: "col4", Type: externaldata.ColumnTypeDateTime},
				{Name: "col5", Type: externaldata.ColumnTypeTimestamp},
				{Name: "col6", Type: externaldata.ColumnTypeStringArray},
			},
			newConfigs: []ColumnConfig{
				{Name: "col6", Type: externaldata.ColumnTypeStringArray},
				{Name: "col4", Type: externaldata.ColumnTypeDateTime},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col5", Type: externaldata.ColumnTypeTimestamp},
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
			},
			expectedNewConfigs: []ColumnConfig{
				{Name: "col1", Type: externaldata.ColumnTypeString},
				{Name: "col2", Type: externaldata.ColumnTypeNumber},
				{Name: "col3", Type: externaldata.ColumnTypeBoolean},
				{Name: "col4", Type: externaldata.ColumnTypeDateTime},
				{Name: "col5", Type: externaldata.ColumnTypeTimestamp},
				{Name: "col6", Type: externaldata.ColumnTypeStringArray},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newConfigsCopy := make([]ColumnConfig, len(tc.newConfigs))
			copy(newConfigsCopy, tc.newConfigs)

			err := syncColumnConfigsOrder(tc.oldConfigs, newConfigsCopy)
			if err != nil {
				if tc.expectedErrorMessage == "" {
					t.Errorf("error is not expected: %v", err)
				} else if err.Error() != tc.expectedErrorMessage {
					t.Errorf("expected error containing %q but got %q", tc.expectedErrorMessage, err.Error())
				}
			} else if tc.expectedErrorMessage != "" {
				t.Error("expected error but got none")
			} else {
				if !slices.Equal(newConfigsCopy, tc.expectedNewConfigs) {
					t.Errorf("expected %v, got %v", tc.expectedNewConfigs, newConfigsCopy)
				}
			}
		})
	}
}
