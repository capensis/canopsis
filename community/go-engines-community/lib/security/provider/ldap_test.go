package provider

import (
	"context"
	"testing"

	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_ldap "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/github.com/go-ldap/ldap"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	mock_provider "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security/provider"
	"github.com/go-ldap/ldap/v3"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const validRole = "valid"

func TestLdapProvider_Auth_GivenUsernameAndPassword_ShouldReturnUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	expectedUser := &security.User{
		ID:        "testid",
		IsEnabled: true,
	}
	config := security.LdapConfig{
		Url:           "ldaps://test",
		AdminUsername: "testadminname",
		AdminPassword: "testadminpass",
		DefaultRole:   validRole,
	}
	entry := &ldap.Entry{
		DN:         username,
		Attributes: []*ldap.EntryAttribute{ldap.NewEntryAttribute("cn", []string{username})},
	}
	mockClient := mock_ldap.NewMockClient(ctrl)
	mockClient.
		EXPECT().
		Bind(gomock.Eq(config.AdminUsername), gomock.Eq(config.AdminPassword)).
		Return(nil)
	mockClient.
		EXPECT().
		Search(gomock.Any()).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{entry}}, nil)
	mockClient.EXPECT().Close()
	mockClient.
		EXPECT().
		Bind(gomock.Eq(username), gomock.Eq(password)).
		Return(nil)
	mockLdapDialer := mock_provider.NewMockLdapDialer(ctrl)
	mockLdapDialer.EXPECT().
		DialURL(gomock.Eq(config)).
		Return(mockClient, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Eq(username), gomock.Eq(security.SourceLdap)).
		Return(expectedUser, nil)
	mockUserProvider.
		EXPECT().
		Save(gomock.Any(), gomock.Eq(expectedUser)).
		Return(nil)

	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.EXPECT().LoadPolicy().Return(nil)

	p := NewLdapProvider(getDbClientMock(ctrl), config, mockUserProvider, mockLdapDialer, mockEnforcer)
	user, err := p.Auth(ctx, username, password)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != expectedUser {
		t.Errorf("expected user: %v but got %v", expectedUser, user)
	}
}

func TestLdapProvider_Auth_GivenInvalidAdminCredentials_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	config := security.LdapConfig{
		Url:           "ldaps://test",
		AdminUsername: "testadminname",
		AdminPassword: "testadminpass",
		DefaultRole:   validRole,
	}
	mockClient := mock_ldap.NewMockClient(ctrl)
	mockClient.
		EXPECT().
		Bind(gomock.Eq(config.AdminUsername), gomock.Eq(config.AdminPassword)).
		Return(&ldap.Error{ResultCode: ldap.LDAPResultInvalidCredentials})
	mockClient.
		EXPECT().
		Search(gomock.Any()).
		Times(0)
	mockClient.EXPECT().Close()
	mockLdapDialer := mock_provider.NewMockLdapDialer(ctrl)
	mockLdapDialer.EXPECT().
		DialURL(gomock.Eq(config)).
		Return(mockClient, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	mockEnforcer := mock_security.NewMockEnforcer(ctrl)

	p := NewLdapProvider(getDbClientMock(ctrl), config, mockUserProvider, mockLdapDialer, mockEnforcer)
	user, err := p.Auth(ctx, username, password)

	if err == nil {
		t.Error("expected error but got nil")
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestLdapProvider_Auth_GivenInvalidUsername_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	config := security.LdapConfig{
		Url:           "ldaps://test",
		AdminUsername: "testadminname",
		AdminPassword: "testadminpass",
		DefaultRole:   validRole,
	}
	mockClient := mock_ldap.NewMockClient(ctrl)
	mockClient.
		EXPECT().
		Bind(gomock.Eq(config.AdminUsername), gomock.Eq(config.AdminPassword)).
		Return(nil)
	mockClient.
		EXPECT().
		Search(gomock.Any()).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	mockClient.EXPECT().Close()
	mockLdapDialer := mock_provider.NewMockLdapDialer(ctrl)
	mockLdapDialer.EXPECT().
		DialURL(gomock.Eq(config)).
		Return(mockClient, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	mockEnforcer := mock_security.NewMockEnforcer(ctrl)

	p := NewLdapProvider(getDbClientMock(ctrl), config, mockUserProvider, mockLdapDialer, mockEnforcer)
	user, err := p.Auth(ctx, username, password)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestLdapProvider_Auth_GivenInvalidPassword_ShouldReturnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	config := security.LdapConfig{
		Url:           "ldaps://test",
		AdminUsername: "testadminname",
		AdminPassword: "testadminpass",
		DefaultRole:   validRole,
	}
	entry := &ldap.Entry{
		DN: username,
	}
	mockClient := mock_ldap.NewMockClient(ctrl)
	mockClient.
		EXPECT().
		Bind(gomock.Eq(config.AdminUsername), gomock.Eq(config.AdminPassword)).
		Return(nil)
	mockClient.
		EXPECT().
		Search(gomock.Any()).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{entry}}, nil)
	mockClient.EXPECT().Close()
	mockClient.
		EXPECT().
		Bind(gomock.Eq(username), gomock.Eq(password)).
		Return(&ldap.Error{ResultCode: ldap.LDAPResultInvalidCredentials})
	mockLdapDialer := mock_provider.NewMockLdapDialer(ctrl)
	mockLdapDialer.EXPECT().
		DialURL(gomock.Eq(config)).
		Return(mockClient, nil)
	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	mockEnforcer := mock_security.NewMockEnforcer(ctrl)

	p := NewLdapProvider(getDbClientMock(ctrl), config, mockUserProvider, mockLdapDialer, mockEnforcer)
	user, err := p.Auth(ctx, username, password)

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func TestLdapProvider_Auth_GivenUsernameAndPasswordAndNoUserInStore_ShouldCreateNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	config := security.LdapConfig{
		Url:          "ldaps://test",
		UsernameAttr: "testnameattr",
		Attributes: map[string]string{
			"mail":      "testemailattr",
			"firstname": "testfirstnameattr",
			"lastname":  "testlastnameattr",
		},
		DefaultRole: validRole,
	}
	externalID := "testid"
	entry := &ldap.Entry{
		DN: username,
		Attributes: []*ldap.EntryAttribute{
			ldap.NewEntryAttribute("testnameattr", []string{"testnewname"}),
			ldap.NewEntryAttribute("testemailattr", []string{"johndoe@canopsis.net"}),
			ldap.NewEntryAttribute("testfirstnameattr", []string{"testfirstname"}),
			ldap.NewEntryAttribute("testlastnameattr", []string{"testlastname"}),
			ldap.NewEntryAttribute("cn", []string{externalID}),
		},
	}
	expectedUser := &security.User{
		Name:       "testnewname",
		Email:      "johndoe@canopsis.net",
		Firstname:  "testfirstname",
		Lastname:   "testlastname",
		Roles:      []string{validRole},
		ExternalID: externalID,
		Source:     security.SourceLdap,
		IsEnabled:  true,
	}
	mockClient := mock_ldap.NewMockClient(ctrl)
	mockClient.
		EXPECT().
		Bind(gomock.Any(), gomock.Any()).
		Times(2).
		Return(nil)
	mockClient.
		EXPECT().
		Search(gomock.Any()).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{entry}}, nil)
	mockClient.EXPECT().Close()
	mockLdapDialer := mock_provider.NewMockLdapDialer(ctrl)
	mockLdapDialer.EXPECT().
		DialURL(gomock.Any()).
		Return(mockClient, nil)

	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Eq(externalID), gomock.Any()).
		Return(nil, nil)
	mockUserProvider.
		EXPECT().
		Save(gomock.Any(), gomock.Eq(expectedUser)).
		Return(nil)

	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.EXPECT().LoadPolicy().Return(nil)

	p := NewLdapProvider(getDbClientMock(ctrl), config, mockUserProvider, mockLdapDialer, mockEnforcer)
	_, _ = p.Auth(ctx, username, password)
}

func TestLdapProvider_Auth_GivenUsernameAndPasswordAndUserInStore_ShouldUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	config := security.LdapConfig{
		Url:          "ldaps://test",
		UsernameAttr: "testnameattr",
		Attributes: map[string]string{
			"mail":      "testemailattr",
			"firstname": "testfirstnameattr",
			"lastname":  "testlastnameattr",
		},
		DefaultRole: validRole,
	}
	externalID := "testid"
	entry := &ldap.Entry{
		DN: username,
		Attributes: []*ldap.EntryAttribute{
			ldap.NewEntryAttribute("testnameattr", []string{"testnewname"}),
			ldap.NewEntryAttribute("testemailattr", []string{"newjohndoe@canopsis.net"}),
			ldap.NewEntryAttribute("testfirstnameattr", []string{"testnewfirstname"}),
			ldap.NewEntryAttribute("testlastnameattr", []string{"testnewlastname"}),
			ldap.NewEntryAttribute("cn", []string{externalID}),
		},
	}
	expectedUser := &security.User{
		ID:         "testname",
		Name:       "testname",
		Email:      "newjohndoe@canopsis.net",
		Firstname:  "testnewfirstname",
		Lastname:   "testnewlastname",
		ExternalID: externalID,
		Source:     security.SourceLdap,
		IsEnabled:  true,
	}
	mockClient := mock_ldap.NewMockClient(ctrl)
	mockClient.
		EXPECT().
		Bind(gomock.Any(), gomock.Any()).
		Times(2).
		Return(nil)
	mockClient.
		EXPECT().
		Search(gomock.Any()).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{entry}}, nil)
	mockClient.EXPECT().Close()
	mockLdapDialer := mock_provider.NewMockLdapDialer(ctrl)
	mockLdapDialer.EXPECT().
		DialURL(gomock.Any()).
		Return(mockClient, nil)

	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Eq(externalID), gomock.Any()).
		Return(&security.User{
			ID:         "testname",
			Name:       "testname",
			Email:      "johndoe@canopsis.net",
			Firstname:  "testfirstname",
			Lastname:   "testlastname",
			ExternalID: externalID,
			Source:     security.SourceLdap,
			IsEnabled:  true,
		}, nil)
	mockUserProvider.
		EXPECT().
		Save(gomock.Any(), gomock.Eq(expectedUser)).
		Return(nil)

	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.EXPECT().LoadPolicy().Return(nil)

	p := NewLdapProvider(getDbClientMock(ctrl), config, mockUserProvider, mockLdapDialer, mockEnforcer)
	_, _ = p.Auth(ctx, username, password)
}

func TestLdapProvider_Auth_GivenUsernameAndPasswordAndNotFoundRole_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	username := "testname"
	password := "testpass"
	config := security.LdapConfig{
		Url:          "ldaps://test",
		UsernameAttr: "testnameattr",
		Attributes: map[string]string{
			"mail":      "testemailattr",
			"firstname": "testfirstnameattr",
			"lastname":  "testlastnameattr",
		},
		DefaultRole: "not-found-role",
	}
	externalID := "testid"
	entry := &ldap.Entry{
		DN: username,
		Attributes: []*ldap.EntryAttribute{
			ldap.NewEntryAttribute("testnameattr", []string{"testnewname"}),
			ldap.NewEntryAttribute("testemailattr", []string{"johndoe@canopsis.net"}),
			ldap.NewEntryAttribute("testfirstnameattr", []string{"testfirstname"}),
			ldap.NewEntryAttribute("testlastnameattr", []string{"testlastname"}),
			ldap.NewEntryAttribute("cn", []string{externalID}),
		},
	}
	mockClient := mock_ldap.NewMockClient(ctrl)
	mockClient.
		EXPECT().
		Bind(gomock.Any(), gomock.Any()).
		Times(2).
		Return(nil)
	mockClient.
		EXPECT().
		Search(gomock.Any()).
		Return(&ldap.SearchResult{Entries: []*ldap.Entry{entry}}, nil)
	mockClient.EXPECT().Close()
	mockLdapDialer := mock_provider.NewMockLdapDialer(ctrl)
	mockLdapDialer.EXPECT().
		DialURL(gomock.Any()).
		Return(mockClient, nil)

	mockUserProvider := mock_security.NewMockUserProvider(ctrl)
	mockUserProvider.
		EXPECT().
		FindByExternalSource(gomock.Any(), gomock.Eq(externalID), gomock.Any()).
		Return(nil, nil)

	p := NewLdapProvider(getDbClientMock(ctrl), config, mockUserProvider, mockLdapDialer, mock_security.NewMockEnforcer(ctrl))
	user, err := p.Auth(ctx, username, password)

	if err == nil {
		t.Error("expected error but got none")
	}

	if user != nil {
		t.Errorf("expected no user but got %v", user)
	}
}

func getDbClientMock(ctrl *gomock.Controller) *mock_mongo.MockDbClient {
	singleResultHelperSuccess := mock_mongo.NewMockSingleResultHelper(ctrl)
	singleResultHelperSuccess.EXPECT().Err().Return(nil).AnyTimes()

	singleResultHelperNotFound := mock_mongo.NewMockSingleResultHelper(ctrl)
	singleResultHelperNotFound.EXPECT().Err().Return(mongo.ErrNoDocuments).AnyTimes()

	unexpectedResultHelper := mock_mongo.NewMockSingleResultHelper(ctrl)
	unexpectedResultHelper.EXPECT().Err().Return(mongo.ErrNoDocuments).AnyTimes()

	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, query bson.M, _ ...*options.FindOneOptions) libmongo.SingleResultHelper {
			role, ok := query["name"]
			if !ok {
				return unexpectedResultHelper
			}

			if role == validRole {
				return singleResultHelperSuccess
			}

			return singleResultHelperNotFound
		},
	).AnyTimes()

	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Collection(libmongo.RoleCollection).Return(mockDbCollection)

	return mockDbClient
}
