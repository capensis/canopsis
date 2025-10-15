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
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
			expectedNewConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
		},
		{
			name: "given configs in different order expect reordering",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
			},
			expectedNewConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
		},
		{
			name: "given configs completely reversed expect complete reordering",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
			},
			expectedNewConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
		},
		{
			name: "given single config expect no changes",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
			},
			expectedNewConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
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
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeBoolean}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeDateTime}},
			},
			expectedNewConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeDateTime}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeBoolean}},
			},
		},
		{
			name: "given different length configs expect error",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
			},
			expectedErrorMessage: "column config count mismatch",
		},
		{
			name: "given new configs longer than old configs expect error",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
			},
			expectedErrorMessage: "column config count mismatch",
		},
		{
			name: "given missing column in new configs expect error",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
			expectedErrorMessage: "no such column \"col3\"",
		},
		{
			name: "given duplicate column names in new configs expect error",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeNumber}},
			},
			expectedErrorMessage: "duplicate column name \"col1\"",
		},
		{
			name: "given large number of configs expect reordering",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col4", Type: externaldata.ColumnTypeDateTime}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col5", Type: externaldata.ColumnTypeTimestamp}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col6", Type: externaldata.ColumnTypeStringArray}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col6", Type: externaldata.ColumnTypeStringArray}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col4", Type: externaldata.ColumnTypeDateTime}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col5", Type: externaldata.ColumnTypeTimestamp}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
			},
			expectedNewConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "col1", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col2", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col3", Type: externaldata.ColumnTypeBoolean}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col4", Type: externaldata.ColumnTypeDateTime}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col5", Type: externaldata.ColumnTypeTimestamp}},
				{BaseColumnConfig: BaseColumnConfig{Name: "col6", Type: externaldata.ColumnTypeStringArray}},
			},
		},
		{
			name: "given configs where priority column already regexp but new column contain regexp type, expect error",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "priority", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "regexp", Type: externaldata.ColumnTypeNumber}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "priority", Type: externaldata.ColumnTypeString}},
				{BaseColumnConfig: BaseColumnConfig{Name: "regexp", Type: externaldata.ColumnTypeRegexp}},
			},
			expectedErrorMessage: "column \"regexp\" is regexp, but priority column already exists",
		},
		{
			name: "given configs where priority column already regexp but new column contain regexp type, expect error, alternate order",
			oldConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "regexp", Type: externaldata.ColumnTypeNumber}},
				{BaseColumnConfig: BaseColumnConfig{Name: "priority", Type: externaldata.ColumnTypeString}},
			},
			newConfigs: []ColumnConfig{
				{BaseColumnConfig: BaseColumnConfig{Name: "regexp", Type: externaldata.ColumnTypeRegexp}},
				{BaseColumnConfig: BaseColumnConfig{Name: "priority", Type: externaldata.ColumnTypeString}},
			},
			expectedErrorMessage: "column \"regexp\" is regexp, but priority column already exists",
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
