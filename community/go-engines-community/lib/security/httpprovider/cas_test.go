package httpprovider

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_http "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/http"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/golang/mock/gomock"
)

func TestCasProvider_Auth_GivenTicketByQueryParam_ShouldAuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ticket := "testticket"
	externalID := "testexternal"
	expectedUser := &security.User{
		ID:         "testid",
		Name:       externalID,
		ExternalID: externalID,
		Source:     security.SourceCas,
		IsEnabled:  true,
	}
	service := "http://test-service"
	config := security.CasConfig{
		ValidateUrl: "http://test-validate",
	}
	mockDoer := mock_http.NewMockDoer(ctrl)
	casRequest, _ := http.NewRequest(http.MethodGet, "http://test-validate", nil)
	casRequest.URL.RawQuery = url.Values{
		"service": []string{service + "?service=" + url.QueryEscape(service)},
		"ticket":  []string{ticket},
	}.Encode()
	casBody := `<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>
	<cas:authenticationSuccess><cas:user>` + externalID + `</cas:user></cas:authenticationSuccess>
</cas:serviceResponse>`
	casResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(casBody)),
	}
	mockDoer.
		EXPECT().
		Do(gomock.Eq(casRequest)).
		Return(casResponse, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Eq(externalID), gomock.Eq(security.SourceCas)).
		Return(expectedUser, nil)

	p := NewCasProvider(mockDoer, config, mockUserProvider)
	r := newRequest()
	r.URL.Host = "test-service"
	r.Host = "test-service"
	r.URL.RawQuery = url.Values{
		security.QueryParamCasTicket:  []string{ticket},
		security.QueryParamCasService: []string{service},
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

func TestCasProvider_Auth_GivenNoQueryParam_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDoer := mock_http.NewMockDoer(ctrl)
	mockDoer.
		EXPECT().
		Do(gomock.Any()).
		Times(0)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	p := NewCasProvider(mockDoer, security.CasConfig{}, mockUserProvider)
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

func TestCasProvider_Auth_GivenInvalidTicketInQueryParam_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ticket := "testticket"
	service := "http://test-service"
	config := security.CasConfig{
		ValidateUrl: "http://test-validate",
	}
	mockDoer := mock_http.NewMockDoer(ctrl)
	casRequest, _ := http.NewRequest(http.MethodGet, "http://test-validate", nil)
	casRequest.URL.RawQuery = url.Values{
		"service": []string{service + "?service=" + url.QueryEscape(service)},
		"ticket":  []string{ticket},
	}.Encode()
	casBody := `<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>
		<cas:authenticationFailure code="INVALID_TICKET"></cas:authenticationFailure>
	</cas:serviceResponse>`
	casResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(casBody)),
	}
	mockDoer.
		EXPECT().
		Do(gomock.Eq(casRequest)).
		Return(casResponse, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	p := NewCasProvider(mockDoer, config, mockUserProvider)
	r := newRequest()
	r.URL.Host = "test-service"
	r.Host = "test-service"
	r.URL.RawQuery = url.Values{
		security.QueryParamCasTicket:  []string{ticket},
		security.QueryParamCasService: []string{service},
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

func TestCasProvider_Auth_GivenTicketByQueryParamAndNoUserInStore_ShouldCreateNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ticket := "testticket"
	externalID := "testexternal"
	expectedUser := &security.User{
		Name:       externalID,
		ExternalID: externalID,
		Source:     security.SourceCas,
		IsEnabled:  true,
		Roles:      []string{"admin"},
	}
	service := "http://test-service"
	config := security.CasConfig{
		ValidateUrl: "http://test-validate",
		DefaultRole: "admin",
	}
	mockDoer := mock_http.NewMockDoer(ctrl)
	casRequest, _ := http.NewRequest(http.MethodGet, "http://test-validate", nil)
	casRequest.URL.RawQuery = url.Values{
		"service": []string{service + "?service=" + url.QueryEscape(service)},
		"ticket":  []string{ticket},
	}.Encode()
	casBody := `<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>
	<cas:authenticationSuccess><cas:user>` + externalID + `</cas:user></cas:authenticationSuccess>
</cas:serviceResponse>`
	casResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(casBody)),
	}
	mockDoer.
		EXPECT().
		Do(gomock.Eq(casRequest)).
		Return(casResponse, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Eq(externalID), gomock.Eq(security.SourceCas)).
		Return(nil, nil)
	mockUserProvider.
		EXPECT().
		Save(gomock.Any(), gomock.Eq(expectedUser)).
		Return(nil)

	p := NewCasProvider(mockDoer, config, mockUserProvider)
	r := newRequest()
	r.URL.RawQuery = url.Values{
		security.QueryParamCasTicket:  []string{ticket},
		security.QueryParamCasService: []string{service},
	}.Encode()
	_, _, _ = p.Auth(r)
}

func TestCasProvider_Auth_GivenTicketByQueryParamAndUserInStore_ShouldNotUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ticket := "testticket"
	externalID := "testexternal"
	expectedUser := &security.User{
		ID:         "testid",
		Name:       externalID,
		ExternalID: externalID,
		Source:     security.SourceCas,
		IsEnabled:  true,
	}
	service := "http://test-service"
	config := security.CasConfig{
		ValidateUrl: "http://test-validate",
	}
	mockDoer := mock_http.NewMockDoer(ctrl)
	casRequest, _ := http.NewRequest(http.MethodGet, "http://test-validate", nil)
	casRequest.URL.RawQuery = url.Values{
		"service": []string{service + "?service=" + url.QueryEscape(service)},
		"ticket":  []string{ticket},
	}.Encode()
	casBody := `<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>
	<cas:authenticationSuccess><cas:user>` + externalID + `</cas:user></cas:authenticationSuccess>
</cas:serviceResponse>`
	casResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(casBody)),
	}
	mockDoer.
		EXPECT().
		Do(gomock.Eq(casRequest)).
		Return(casResponse, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Eq(externalID), gomock.Eq(security.SourceCas)).
		Return(expectedUser, nil)
	mockUserProvider.
		EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Times(0)

	p := NewCasProvider(mockDoer, config, mockUserProvider)
	r := newRequest()
	r.URL.RawQuery = url.Values{
		security.QueryParamCasTicket:  []string{ticket},
		security.QueryParamCasService: []string{service},
	}.Encode()
	_, _, _ = p.Auth(r)
}
