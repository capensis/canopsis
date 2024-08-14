package provider

//go:generate mockgen -destination=../../../mocks/lib/security/provider/provider.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/provider LdapDialer
//go:generate mockgen -destination=../../../mocks/github.com/go-ldap/ldap/client.go github.com/go-ldap/ldap/v3 Client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/go-ldap/ldap/v3"
)

// LdapDialer interface is used to implement creation of LDAP connection.
type LdapDialer interface {
	DialURL(config security.LdapConfig) (ldap.Client, error)
}

// baseDialer wraps ldap lib function.
type baseDialer struct{}

// NewLdapDialer creates new dialer.
func NewLdapDialer() LdapDialer {
	return &baseDialer{}
}

func (baseDialer) DialURL(config security.LdapConfig) (ldap.Client, error) {
	tc := &tls.Config{
		InsecureSkipVerify: config.InsecureSkipVerify, //nolint:gosec
	}
	tc.MinVersion = strToTlsVersion(config.MinTLSVersion)
	tc.MaxVersion = strToTlsVersion(config.MaxTLSVersion)

	return ldap.DialURL(config.Url, ldap.DialWithTLSConfig(tc))
}

// ldapProvider implements LDAP authentication.
type ldapProvider struct {
	config       security.LdapConfig
	userProvider security.UserProvider
	roleProvider security.RoleProvider
	ldapDialer   LdapDialer
	enforcer     security.Enforcer
}

// NewLdapProvider creates new provider.
func NewLdapProvider(
	config security.LdapConfig,
	userProvider security.UserProvider,
	roleProvider security.RoleProvider,
	ldapDialer LdapDialer,
	enforcer security.Enforcer,
) security.Provider {
	return &ldapProvider{
		config:       config,
		userProvider: userProvider,
		roleProvider: roleProvider,
		ldapDialer:   ldapDialer,
		enforcer:     enforcer,
	}
}

func (p *ldapProvider) GetName() string {
	return security.AuthMethodLdap
}

func (p *ldapProvider) Auth(ctx context.Context, username, password string) (*security.User, error) {
	entry, err := p.validateUser(ctx, username, password)
	if err != nil {
		return nil, err
	}

	if entry != nil {
		return p.saveUser(ctx, username, entry)
	}

	return nil, nil
}

// validateUser calls LDAP server to authenticate user.
func (p *ldapProvider) validateUser(
	ctx context.Context,
	username, password string,
) (*ldap.Entry, error) {
	l, err := p.ldapDialer.DialURL(p.config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to ldap server: %w", err)
	}
	defer l.Close()

	select {
	case <-ctx.Done():
		return nil, nil
	default:
	}

	// First bind with a admin
	err = l.Bind(p.config.AdminUsername, p.config.AdminPassword)
	if err != nil {
		return nil, fmt.Errorf("cannot verify admin: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, nil
	default:
	}

	// Search for the given username
	searchRequest := genSearchRequest(p.config, username)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot search user: %w", err)
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	if len(sr.Entries) > 1 {
		return nil, errors.New("too many entries returned")
	}

	select {
	case <-ctx.Done():
		return nil, nil
	default:
	}

	entry := sr.Entries[0]
	// Bind as the user to verify their password
	err = l.Bind(entry.DN, password)
	if err != nil {
		// Return nil if password verification failed.
		var ldapErr *ldap.Error
		if errors.As(err, &ldapErr) {
			if ldapErr.ResultCode == ldap.LDAPResultInvalidCredentials {
				return nil, nil
			}
		}

		return nil, fmt.Errorf("cannot verify password: %w", err)
	}

	return entry, nil
}

// genSearchRequest creates a new search request by config.
func genSearchRequest(config security.LdapConfig, username string) *ldap.SearchRequest {
	// Create attributes.
	attributes := make([]string, len(config.Attributes))
	i := 0
	for _, v := range config.Attributes {
		attributes[i] = v
		i++
	}
	if config.UsernameAttr != "" {
		attributes = append(attributes, config.UsernameAttr)
	}
	attributes = append(attributes, "cn")

	// Create filter.
	var filter string
	if len(config.Filter) > 0 {
		filter = fmt.Sprintf(config.Filter, username)
		if filter[0] != '(' {
			filter = fmt.Sprintf("(%s)", filter)
		}
	}

	return ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil,
	)
}

// saveUser creates or updates user data in storage.
func (p *ldapProvider) saveUser(
	ctx context.Context,
	username string,
	entry *ldap.Entry,
) (*security.User, error) {
	if p.config.UsernameAttr != "" {
		name := entry.GetAttributeValue(p.config.UsernameAttr)
		if name != "" {
			username = name
		}
	}

	ldapID := entry.GetAttributeValue("cn")
	user, err := p.userProvider.FindByExternalSource(ctx, ldapID, security.SourceLdap)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %w", err)
	}

	if user == nil {
		roleID, err := p.roleProvider.GetRoleID(ctx, p.config.DefaultRole)
		if err != nil {
			return nil, err
		}

		user = &security.User{
			Name:       username,
			Roles:      []string{roleID},
			IsEnabled:  true,
			ExternalID: ldapID,
			Source:     security.SourceLdap,
		}
	} else if !user.IsEnabled {
		return nil, nil
	}

	for field, attr := range p.config.Attributes {
		v := entry.GetAttributeValue(attr)
		switch field {
		case "mail":
			user.Email = v
		case "firstname":
			user.Firstname = v
		case "lastname":
			user.Lastname = v
		}
	}

	err = p.userProvider.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("cannot save user: %w", err)
	}

	err = p.enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("LdapProvider: reload enforcer error: %w", err)
	}

	return user, nil
}

func strToTlsVersion(str string) uint16 {
	switch str {
	case "tls10":
		return tls.VersionTLS10
	case "tls11":
		return tls.VersionTLS11
	case "tls12":
		return tls.VersionTLS12
	case "tls13":
		return tls.VersionTLS13
	}

	return 0
}
