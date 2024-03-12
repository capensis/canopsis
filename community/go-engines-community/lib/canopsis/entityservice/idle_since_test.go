package entityservice_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
)

func TestMark(t *testing.T) {
	m := entityservice.NewServicesIdleSinceMap()

	if !m.Mark("1", 2) {
		t.Fatalf("Should be able to mark")
	}

	if m.Mark("1", 3) {
		t.Fatalf("Shouldn't be able to mark")
	}

	if !m.Mark("1", 1) {
		t.Fatalf("Should be able to mark")
	}
}
