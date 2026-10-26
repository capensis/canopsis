package webhook

import (
	"context"
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type stubUserRoleReader struct {
	roleNames []string
	err       error
}

func (r *stubUserRoleReader) GetRoleNames(_ context.Context, _ string) ([]string, error) {
	return r.roleNames, r.err
}

func TestPrepareAdditionalDataErrors(t *testing.T) {
	t.Run("non-empty user without reader fails", func(t *testing.T) {
		if _, err := PrepareAdditionalData(context.Background(), nil, types.AdditionalData{User: "user-1"}); err == nil {
			t.Fatal("expected an error when no role reader is configured")
		}
	})

	t.Run("reader errors are returned", func(t *testing.T) {
		wantErr := errors.New("lookup failed")
		_, err := PrepareAdditionalData(context.Background(), &stubUserRoleReader{err: wantErr}, types.AdditionalData{User: "user-1"})
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected wrapped reader error, got %v", err)
		}
	})
}
