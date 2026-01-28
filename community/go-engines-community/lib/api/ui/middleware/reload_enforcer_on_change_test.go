package middleware

import (
	"net/http"
	"testing"

	mock_response "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/ui/response"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestReloadEnforcerPolicyOnChange_GivenOkResponse_ShouldLoadPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		LoadPolicy().
		Return(nil)
	mockResponder := mock_response.NewMockErrorResponder(ctrl)
	router := gin.New()
	router.GET(
		okURL,
		okHandler,
		ReloadEnforcerPolicyOnChange(mockEnforcer, mockResponder),
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
	mockResponder := mock_response.NewMockErrorResponder(ctrl)
	router := gin.New()
	router.GET(
		okURL,
		func(c *gin.Context) {
			c.Status(http.StatusBadRequest)
		},
		ReloadEnforcerPolicyOnChange(mockEnforcer, mockResponder),
	)

	_ = performRequest(router, "GET", okURL)
}
