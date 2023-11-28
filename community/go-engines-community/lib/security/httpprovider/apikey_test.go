package httpprovider

import (
	"net/http"
	"net/url"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/golang/mock/gomock"
)

func newRequest() *http.Request {
	return &http.Request{
		Method: http.MethodGet,
		Host:   "www.google.com",
		URL: &url.URL{
			Scheme: "http",
			Host:   "www.google.com",
			Path:   "/search",
		},
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
	}
}

func TestApikeyProvider_Auth_GivenApiKeyByQueryParam_ShouldAuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiKey := "testkey"
	expectedUser := &security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
		IsEnabled:      true,
	}
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByAuthApiKey(gomock.Any(), gomock.Eq(apiKey)).
		Return(expectedUser, nil)

	p := NewApikeyProvider(mockUserProvider)
	r := newRequest()
	r.URL.RawQuery = url.Values{security.QueryParamApiKey: []string{apiKey}}.Encode()
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

func TestApikeyProvider_Auth_GivenApiKeyByHeader_ShouldAuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiKey := "testkey"
	expectedUser := &security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
		IsEnabled:      true,
	}
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByAuthApiKey(gomock.Any(), gomock.Eq(apiKey)).
		Return(expectedUser, nil)

	p := NewApikeyProvider(mockUserProvider)
	r := newRequest()
	r.Header.Set(security.HeaderApiKey, apiKey)
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

func TestApikeyProvider_Auth_GivenNoQueryParamAndNoHeader_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByAuthApiKey(gomock.Any(), gomock.Any()).
		Times(0)

	p := NewApikeyProvider(mockUserProvider)
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

func TestApikeyProvider_Auth_GivenInvalidApiKeyInQueryParam_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiKey := "testkey"
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByAuthApiKey(gomock.Any(), gomock.Eq(apiKey)).
		Return(nil, nil)

	p := NewApikeyProvider(mockUserProvider)
	r := newRequest()
	r.URL.RawQuery = url.Values{security.QueryParamApiKey: []string{apiKey}}.Encode()
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

func TestApikeyProvider_Auth_GivenInvalidApiKeyInHeader_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiKey := "testkey"
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByAuthApiKey(gomock.Any(), gomock.Eq(apiKey)).
		Return(nil, nil)

	p := NewApikeyProvider(mockUserProvider)
	r := newRequest()
	r.Header.Set(security.HeaderApiKey, apiKey)
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
