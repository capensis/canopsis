package provider

//go:generate mockgen -destination=../../../mocks/lib/security/provider/provider.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/provider LdapDialer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/go-ldap/ldap/v3"
)

// LdapDialer interface is used to implement creation of LDAP connection.
type LdapDialer interface {
	DialURL(config *security.LdapConfig) (ldap.Client, error)
}

// baseDialer wraps ldap lib function.
type baseDialer struct{}

// NewLdapDialer creates new dialer.
func NewLdapDialer() LdapDialer {
	return &baseDialer{}
}

func (baseDialer) DialURL(config *security.LdapConfig) (ldap.Client, error) {
	return ldap.DialURL(config.Url, ldap.DialWithTLSConfig(&tls.Config{
		InsecureSkipVerify: config.InsecureSkipVerify,
	}))
}

// ldapProvider implements LDAP authentication.
type ldapProvider struct {
	configProvider security.ConfigProvider
	userProvider   security.UserProvider
	ldapDialer     LdapDialer
}

// NewLdapProvider creates new provider.
func NewLdapProvider(
	cp security.ConfigProvider,
	up security.UserProvider,
	d LdapDialer,
) security.Provider {
	return &ldapProvider{
		ldapDialer:     d,
		configProvider: cp,
		userProvider:   up,
	}
}

func (p *ldapProvider) Auth(username, password string) (*security.User, error) {
	config, err := p.configProvider.LoadLdapConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot find ldap config: %v", err)
	}

	if config == nil {
		return nil, errors.New("ldap config not found")
	}

	entry, err := p.validateUser(username, password, config)
	if err != nil {
		return nil, err
	}

	if entry != nil {
		return p.saveUser(username, config, entry)
	}

	return nil, nil
}

// validateUser calls LDAP server to authenticate user.
func (p *ldapProvider) validateUser(
	username, password string,
	config *security.LdapConfig,
) (*ldap.Entry, error) {
	l, err := p.ldapDialer.DialURL(config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to ldap server: %v", err)
	}
	defer l.Close()

	// First bind with a admin
	err = l.Bind(config.AdminUsername, config.AdminPassword)
	if err != nil {
		return nil, fmt.Errorf("cannot verify admin: %v", err)
	}

	// Search for the given username
	searchRequest := genSearchRequest(config, username)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot search user: %v", err)
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	if len(sr.Entries) > 1 {
		return nil, errors.New("too many entries returned")
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

		return nil, fmt.Errorf("cannot verify password: %v", err)
	}

	return entry, nil
}

// genSearchRequest creates a new search request by config.
func genSearchRequest(config *security.LdapConfig, username string) *ldap.SearchRequest {
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
	username string,
	config *security.LdapConfig,
	entry *ldap.Entry,
) (*security.User, error) {
	if config.UsernameAttr != "" {
		name := entry.GetAttributeValue(config.UsernameAttr)
		if name != "" {
			username = name
		}
	}

	ldapID := entry.GetAttributeValue("cn")
	user, err := p.userProvider.FindByExternalSource(ldapID, security.SourceLdap)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %v", err)
	}

	if user == nil {
		user = &security.User{
			Name:       username,
			Role:       config.DefaultRole,
			IsEnabled:  true,
			ExternalID: ldapID,
			Source:     security.SourceLdap,
		}
	} else if !user.IsEnabled {
		return nil, nil
	}

	for field, attr := range config.Attributes {
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

	err = p.userProvider.Save(user)
	if err != nil {
		return nil, fmt.Errorf("cannot save user: %v", err)
	}

	return user, nil
}