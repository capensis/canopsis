package utils_test

import (
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/kylelemons/godebug/pretty"
)

func TestUnique(t *testing.T) {
	dataSets := []struct {
		data     []string
		expected []string
	}{
		{
			data:     []string{"1", "2", "3"},
			expected: []string{"1", "2", "3"},
		},
		{
			data:     []string{"3", "1", "2", "2", "3", "2"},
			expected: []string{"3", "1", "2"},
		},
	}

	for i, set := range dataSets {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			result := utils.Unique(set.data)
			if diff := pretty.Compare(set.expected, result); diff != "" {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}
