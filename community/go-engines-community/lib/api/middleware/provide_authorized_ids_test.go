package middleware

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestProvideAuthorizedIds_GivenAuthorizedUser_ShouldReturnIds(t *testing.T) {
	ctrl := gomock.NewController(t)
	subj := "testsubj"
	obj1 := "testobj1"
	obj2 := "testobj2"
	act := "testact"
	role := "testrole"
	expectedCode := http.StatusOK
	expectedIds := []string{obj1, obj2}
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		GetRolesForUser(gomock.Eq(subj)).
		Return([]string{role}, nil)
	mockEnforcer.
		EXPECT().
		GetPermissionsForUser(gomock.Eq(role)).
		Return([][]string{
			{role, obj1, act},
			{role, obj1, "testanotheract"},
			{role, obj2, act},
		}, nil)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(auth.UserKey, subj)
	})
	router.GET(
		okURL,
		ProvideAuthorizedIds(act, mockEnforcer, nil),
		func(c *gin.Context) {
			ids, _ := c.Get(AuthorizedIds)
			c.JSON(http.StatusOK, ids)
		},
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}

	ids := make([]string, 0)
	err := json.Unmarshal(w.Body.Bytes(), &ids)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if !reflect.DeepEqual(ids, expectedIds) {
		t.Errorf("expected ids: %v but got %v", expectedIds, ids)
	}
}

func TestProvideAuthorizedIds_GivenNoUser_ShouldReturnUnauthorizedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	act := "testact"
	expectedCode := http.StatusUnauthorized
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		GetRolesForUser(gomock.Any()).
		Times(0)
	router := gin.New()
	router.GET(
		okURL,
		ProvideAuthorizedIds(act, mockEnforcer, nil),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestProvideAuthorizedIds_GivenNotAuthorizedUser_ShouldReturnEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	subj := "testsubj"
	act := "testact"
	expectedCode := http.StatusOK
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		GetRolesForUser(gomock.Eq(subj)).
		Return(nil, nil)
	mockEnforcer.
		EXPECT().
		GetPermissionsForUser(gomock.Any()).
		Times(0)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(auth.UserKey, subj)
	})
	router.GET(
		okURL,
		ProvideAuthorizedIds(act, mockEnforcer, nil),
		func(c *gin.Context) {
			ids, _ := c.Get(AuthorizedIds)
			c.JSON(http.StatusOK, ids)
		},
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}

	ids := make([]string, 0)
	err := json.Unmarshal(w.Body.Bytes(), &ids)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if len(ids) != 0 {
		t.Errorf("expected no ids but got %v", ids)
	}
}
