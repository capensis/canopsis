package author

import (
	"strings"
	"testing"
)

func TestGetTempField(t *testing.T) {
	s := "a"
	got := getTempField(s)
	lenDiff := len(got) - len(s)
	if lenDiff < 1 || got[lenDiff:] != s {
		t.Errorf("Expected XXX%s, got %s", s, got)
	}
	s = "a.bc.def"
	fieldPrefix := "a.bc."
	got = getTempField(s)
	lenDiff = len(got) - len(s)
	if lenDiff < 1 || !strings.HasPrefix(got, fieldPrefix) || got[lenDiff+len(fieldPrefix):] != "def" {
		t.Errorf("Expected XXX%s, got %s", s, got)
	}
}
