package httpprovider

import (
	"git.canopsis.net/canopsis/go-engines/lib/security"
	mock_security "git.canopsis.net/canopsis/go-engines/mocks/lib/security"
	"github.com/golang/mock/gomock"
	"net/url"
	"testing"
)

func TestQueryProvider_Auth_GivenUsernameAndPasswordByQueryParam_ShouldAuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	username := "testuser"
	password := "testpassword"
	expectedUser := &security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Eq(username), gomock.Eq(password)).
		Return(expectedUser, nil)

	p := NewQueryProvider(mockProvider)
	r := newRequest()
	r.URL.RawQuery = url.Values{
		"username": []string{username},
		"password": []string{password},
	}.Encode()
	user, err, ok := p.Auth(r)

	if !ok {
		t.Errorf("expected true but got %v", ok)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != expectedUser {
		t.Errorf("expected user: %v but got %v", expectedUser, user)
	}
}

func TestQueryProvider_Auth_GivenNoQueryParam_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any(), gomock.Any()).
		Times(0)

	p := NewQueryProvider(mockProvider)
	r := newRequest()
	user, err, ok := p.Auth(r)

	if ok {
		t.Errorf("expected false but got %v", ok)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestQueryProvider_Auth_GivenInvalidCredentials_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	username := "testuser"
	password := "testpassword"
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Eq(username), gomock.Eq(password)).
		Return(nil, nil)

	p := NewQueryProvider(mockProvider)
	r := newRequest()
	r.URL.RawQuery = url.Values{
		"username": []string{username},
		"password": []string{password},
	}.Encode()
	user, err, ok := p.Auth(r)

	if !ok {
		t.Errorf("expected true but got %v", ok)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}
