package export

import (
	"context"
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type userInterfaceAdapterStub struct {
	conf config.UserInterfaceConf
	err  error
}

func (s userInterfaceAdapterStub) GetConfig(context.Context) (config.UserInterfaceConf, error) {
	return s.conf, s.err
}

func TestResolveTimeFormat(t *testing.T) {
	tests := []struct {
		name       string
		requested  string
		stored     string
		adapterErr error
		want       string
		wantErr    error
	}{
		{name: "request override", requested: "02/01/2006 (Mon) 15:04:05", stored: "Mon 02/01/2006 15:04:05", want: "02/01/2006 (Mon) 15:04:05"},
		{name: "stored allowed format", stored: "Mon 02/01/2006 15:04:05", want: "Mon 02/01/2006 15:04:05"},
		{name: "unknown stored value", stored: "=HYPERLINK(\"x\")", want: "2006-01-02T15:04:05"},
		{name: "missing config", adapterErr: mongodriver.ErrNoDocuments, want: "2006-01-02T15:04:05"},
		{name: "other config error", adapterErr: errRead, wantErr: errRead},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := userInterfaceAdapterStub{
				conf: config.UserInterfaceConf{DateFormat: tt.stored},
				err:  tt.adapterErr,
			}
			got, err := ResolveTimeFormat(t.Context(), tt.requested, adapter)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				if got != "" {
					t.Fatalf("expected empty format on error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected format %q, got %q", tt.want, got)
			}
		})
	}
}

var errRead = errors.New("read failed")
