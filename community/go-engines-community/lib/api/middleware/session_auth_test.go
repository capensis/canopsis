package middleware

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	mock_sessions "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/github.com/gorilla/sessions"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/sessions"
	"net/http"
	"testing"
)

func TestSessionAuth_GivenAuthUser_ShouldReturnResponseAndSetUserDataToContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusOK
	user := &model.Rbac{
		ID:         "testid",
		AuthApiKey: "testkey",
	}
	mockStore := mock_sessions.NewMockStore(ctrl)
	session := sessions.NewSession(mockStore, security.SessionKey)
	session.Values["user"] = user.ID
	mockStore.
		EXPECT().
		Get(gomock.Any(), security.SessionKey).
		Return(session, nil)
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Return(mockUserCursor(ctrl, user), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Eq(libmongo.RightsMongoCollection)).
		Return(mockDbCollection)
	router := gin.New()
	router.GET(
		okURL,
		SessionAuth(mockDbClient, mockStore),
		func(c *gin.Context) {
			c.String(
				expectedCode,
				"test %v %v",
				c.MustGet(auth.UserKey).(string),
				c.MustGet(auth.ApiKey).(string),
			)
		},
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}

	expectedResponse := fmt.Sprintf("test %v %v", user.ID, user.AuthApiKey)

	if w.Body.String() != expectedResponse {
		t.Errorf("expected response: \"%v\" but got \"%v\"", expectedResponse, w.Body.String())
	}
}

func TestSessionAuth_GivenNoSession_ShouldReturnResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusOK
	mockStore := mock_sessions.NewMockStore(ctrl)
	session := sessions.NewSession(mockStore, security.SessionKey)
	mockStore.
		EXPECT().
		Get(gomock.Any(), security.SessionKey).
		Return(session, nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Any()).
		Times(0)
	router := gin.New()
	router.GET(
		okURL,
		SessionAuth(mockDbClient, mockStore),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestSessionAuth_GivenInvalidUserSession_ShouldReturnUnauthorizedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusUnauthorized
	mockStore := mock_sessions.NewMockStore(ctrl)
	session := sessions.NewSession(mockStore, security.SessionKey)
	session.Values["user"] = "testid"
	mockStore.
		EXPECT().
		Get(gomock.Any(), security.SessionKey).
		Return(session, nil)
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Return(mockUserCursor(ctrl, nil), nil)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(gomock.Any()).
		Return(mockDbCollection)
	router := gin.New()
	router.GET(
		okURL,
		SessionAuth(mockDbClient, mockStore),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func mockUserCursor(ctrl *gomock.Controller, user *model.Rbac) libmongo.Cursor {
	mockCursor := mock_mongo.NewMockCursor(ctrl)

	if user != nil {
		mockCursor.EXPECT().Next(gomock.Any()).Return(true)
		mockCursor.
			EXPECT().
			Decode(gomock.Any()).
			Do(func(val interface{}) {
				if u, ok := val.(*model.Rbac); ok {
					*u = *user
				}
			}).
			Return(nil)
	} else {
		mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	}

	mockCursor.EXPECT().Close(gomock.Any())

	return mockCursor
}
