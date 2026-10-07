package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	mock_response "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/ui/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

type topActiveStoreStub struct {
	Store
	user *User
	err  error
}

func (s topActiveStoreStub) GetTopByActiveConnects(context.Context, string, []string) (*User, error) {
	return s.user, s.err
}

func TestGetTopByActiveConnectsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("store errors use the error responder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockErrResponder := mock_response.NewMockErrorResponder(ctrl)
		storeErr := errors.New("store failure")
		recorder, ctx := topConnectContext(t.Context())
		authctx.SetUserKey(ctx, "root")
		authctx.SetRoles(ctx, []string{"admin"})
		mockErrResponder.EXPECT().Respond(ctx, storeErr).DoAndReturn(func(c *gin.Context, _ error) {
			c.Status(http.StatusInternalServerError)
		})
		api := &api{store: topActiveStoreStub{err: storeErr}, errorResponder: mockErrResponder}

		api.GetTopByActiveConnects(ctx)
		ctx.Writer.WriteHeaderNow()

		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
		}
	})
}

func topConnectContext(ctx context.Context) (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(recorder)
	ginCtx.Request = httptest.NewRequestWithContext(ctx, http.MethodGet, "/internal/users/top/active-connects", nil)
	return recorder, ginCtx
}
