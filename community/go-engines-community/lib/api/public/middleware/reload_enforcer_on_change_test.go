package middleware

import (
	"net/http"
	"testing"

	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestReloadEnforcerPolicyOnChange_GivenOkResponse_ShouldLoadPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		LoadPolicy().
		Return(nil)
	r := chi.NewRouter()
	r.Use(ReloadEnforcerPolicyOnChange(mockEnforcer, zerolog.Nop()))
	r.Get(okURL, okHandler)
	_ = performRequest(r, "GET", okURL)
}

func TestReloadEnforcerPolicyOnChange_GivenNotOkResponse_ShouldNotLoadPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		LoadPolicy().
		Times(0)
	r := chi.NewRouter()
	r.Use(ReloadEnforcerPolicyOnChange(mockEnforcer, zerolog.Nop()))
	r.Get(okURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	_ = performRequest(r, "GET", okURL)
}
