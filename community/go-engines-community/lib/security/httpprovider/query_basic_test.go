package httpprovider

import (
	"net/url"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/golang/mock/gomock"
)

func TestQueryBasicProvider_Auth_GivenUsernameAndPasswordByQueryParam_ShouldAuthUser(t *testing.T) {
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
		Auth(gomock.Any(), gomock.Eq(username), gomock.Eq(password)).
		Return(expectedUser, nil)

	p := NewQueryBasicProvider(mockProvider)
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

func TestQueryBasicProvider_Auth_GivenNoQueryParam_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	p := NewQueryBasicProvider(mockProvider)
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

func TestQueryBasicProvider_Auth_GivenInvalidCredentials_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	username := "testuser"
	password := "testpassword"
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any(), gomock.Eq(username), gomock.Eq(password)).
		Return(nil, nil)

	p := NewQueryBasicProvider(mockProvider)
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
