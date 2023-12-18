package middleware

import (
	"net/http"
	"testing"

	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestReloadEnforcerPolicyOnChange_GivenOkResponse_ShouldLoadPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		LoadPolicy().
		Return(nil)
	router := gin.New()
	router.GET(
		okURL,
		okHandler,
		ReloadEnforcerPolicyOnChange(mockEnforcer),
	)

	_ = performRequest(router, "GET", okURL)
}

func TestReloadEnforcerPolicyOnChange_GivenNotOkResponse_ShouldNotLoadPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		LoadPolicy().
		Times(0)
	router := gin.New()
	router.GET(
		okURL,
		func(c *gin.Context) {
			c.Status(http.StatusBadRequest)
		},
		ReloadEnforcerPolicyOnChange(mockEnforcer),
	)

	_ = performRequest(router, "GET", okURL)
}
