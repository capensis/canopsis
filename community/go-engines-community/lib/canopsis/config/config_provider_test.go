package config

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestUserInterfaceConfigProviderDateFormat(t *testing.T) {
	var logs bytes.Buffer
	provider := NewUserInterfaceConfigProvider(UserInterfaceConf{DateFormat: "initial"}, zerolog.New(&logs))
	if got := provider.Get().DateFormat; got != "initial" {
		t.Fatalf("expected initial date format to be preserved, got %q", got)
	}
	if got := strings.Count(logs.String(), "DateFormat of user interface config"); got != 1 {
		t.Fatalf("expected constructor to log the initial date format once, got %d date format logs", got)
	}

	provider.Update(UserInterfaceConf{DateFormat: "initial"})
	if got := strings.Count(logs.String(), "DateFormat of user interface config"); got != 1 {
		t.Fatalf("expected unchanged date format not to be logged, got %d date format logs", got)
	}

	provider.Update(UserInterfaceConf{DateFormat: "updated"})
	if got := provider.Get().DateFormat; got != "updated" {
		t.Fatalf("expected updated date format, got %q", got)
	}
	if got := strings.Count(logs.String(), "DateFormat of user interface config"); got != 2 {
		t.Fatalf("expected changed date format to be logged once, got %d date format logs", got)
	}

	provider.Update(UserInterfaceConf{})
	if got := provider.Get().DateFormat; got != "" {
		t.Fatalf("expected empty date format after reset, got %q", got)
	}
	if got := strings.Count(logs.String(), "DateFormat of user interface config"); got != 3 {
		t.Fatalf("expected reset date format to be logged once, got %d date format logs", got)
	}
}
