package httpprovider

import (
	"encoding/base64"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/golang/mock/gomock"
)

func TestBasicProvider_Auth_GivenAuthorizationHeader_ShouldAuthUser(t *testing.T) {
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

	p := NewBasicProvider(mockProvider)
	r := newRequest()
	base := []byte(fmt.Sprintf("%s:%v", username, password))
	r.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(base))
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

func TestBasicProvider_Auth_GivenNoHeader_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	p := NewBasicProvider(mockProvider)
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

func TestBasicProvider_Auth_GivenInvalidCredentials_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	username := "testuser"
	password := "testpassword"
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any(), gomock.Eq(username), gomock.Eq(password)).
		Return(nil, nil)

	p := NewBasicProvider(mockProvider)
	r := newRequest()
	base := []byte(fmt.Sprintf("%s:%v", username, password))
	r.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(base))
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

func TestBasicProvider_Auth_GivenInvalidHeader_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProvider := mock_security.NewMockProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	p := NewBasicProvider(mockProvider)
	r := newRequest()
	r.Header.Set("Authorization", "Basic test")
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
