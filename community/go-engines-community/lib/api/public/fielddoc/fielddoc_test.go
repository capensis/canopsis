package fielddoc

import (
	"strings"
	"testing"
)

func TestFields_AvailableValues(t *testing.T) {
	tests := map[string]struct {
		fields Fields
		want   string
	}{
		"labels set": {
			fields: Fields{Static: []string{"_id", "name"}, Labels: []string{"a", "infos.*.*"}},
			want:   "<i>Available values</i> : `a`, `infos.*.*`",
		},
		"labels empty falls back to static": {
			fields: Fields{Static: []string{"_id", "name"}},
			want:   "<i>Available values</i> : `_id`, `name`",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.fields.AvailableValues()
			if got != tc.want {
				t.Errorf("expected %q, got %q", tc.want, got)
			}

			for i, segment := range strings.Split(got, "`") {
				if i%2 == 0 && strings.Contains(segment, "*") {
					t.Errorf("unescaped * outside code span in %q", got)
				}
			}
		})
	}
}
