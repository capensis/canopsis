package provider

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	libpassword "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	mock_password "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security/password"
	"github.com/golang/mock/gomock"
)

func TestBaseProvider_Auth_GivenUsernameAndPassword_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	expectedUser := &security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
		IsEnabled:      true,
	}
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByUsername(gomock.Any(), gomock.Eq(username)).
		Return(expectedUser, nil)
	mockEncoder := mock_password.NewMockEncoder(ctrl)
	mockEncoder.
		EXPECT().
		IsValidPassword(gomock.Eq([]byte(expectedUser.HashedPassword)), gomock.Eq([]byte(password))).
		Return(true, nil)

	p := NewBaseProvider(mockUserProvider, []libpassword.Encoder{mockEncoder})
	user, err := p.Auth(ctx, username, password)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != expectedUser {
		t.Errorf("expected user: %v but got %v", expectedUser, user)
	}
}

func TestBaseProvider_Auth_GivenInvalidUsername_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByUsername(gomock.Any(), gomock.Eq(username)).
		Return(nil, nil)
	mockEncoder := mock_password.NewMockEncoder(ctrl)
	mockEncoder.
		EXPECT().
		IsValidPassword(gomock.Any(), gomock.Any()).
		Times(0)

	p := NewBaseProvider(mockUserProvider, []libpassword.Encoder{mockEncoder})
	user, err := p.Auth(ctx, username, password)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestBaseProvider_Auth_GivenInvalidPassword_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	expectedUser := &security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
		IsEnabled:      true,
	}
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByUsername(gomock.Any(), gomock.Eq(username)).
		Return(expectedUser, nil)
	mockEncoder := mock_password.NewMockEncoder(ctrl)
	mockEncoder.
		EXPECT().
		IsValidPassword(gomock.Eq([]byte(expectedUser.HashedPassword)), gomock.Eq([]byte(password))).
		Return(false, nil)

	p := NewBaseProvider(mockUserProvider, []libpassword.Encoder{mockEncoder})
	user, err := p.Auth(ctx, username, password)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}
