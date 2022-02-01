package security

import (
	"context"
	"net/http"
)

// Provider interface is used to implement user authentication by username and password.
type Provider interface {
	GetName() string
	Auth(ctx context.Context, username, password string) (*User, error)
}

// HttpProvider interface is used to implement user authentication
// by credentials which are retrieved from http request.
type HttpProvider interface {
	Auth(*http.Request) (*User, error, bool)
}

// TokenProvider interface is used to implement user authentication by token.
type TokenProvider interface {
	Auth(ctx context.Context, token string) (*User, error)
}

// UserProvider is decorator for requests to user storage.
type UserProvider interface {
	// FindByUsername returns user with username or nil.
	FindByUsername(ctx context.Context, username string) (*User, error)
	// FindByAuthApiKey returns user with api key or nil.
	FindByAuthApiKey(ctx context.Context, apiKey string) (*User, error)
	// FindByID returns user with ID or nil.
	FindByID(ctx context.Context, id string) (*User, error)
	// FindByExternalSource returns user with ID from source or nil.
	FindByExternalSource(ctx context.Context, externalID string, source Source) (*User, error)
	// Save updates user or inserts user if not exist.
	Save(ctx context.Context, user *User) error
}

// ConfigProvider provides config from storage.
type ConfigProvider interface {
	LoadLdapConfig(ctx context.Context) (*LdapConfig, error)
	LoadCasConfig(ctx context.Context) (*CasConfig, error)
}

type LdapConfig struct {
	Url                string            `bson:"ldap_uri"`
	Host               string            `bson:"host"`
	Port               int64             `bson:"port"`
	AdminUsername      string            `bson:"admin_dn"`
	AdminPassword      string            `bson:"admin_passwd"`
	BaseDN             string            `bson:"user_dn"`
	Attributes         map[string]string `bson:"attrs"`
	UsernameAttr       string            `bson:"username_attr"`
	Filter             string            `bson:"ufilter"`
	DefaultRole        string            `bson:"default_role"`
	InsecureSkipVerify bool              `bson:"skip_verify"`
	MaxTLSVersion      string            `bson:"max_tls_ver"`
}

type CasConfig struct {
	LoginUrl    string `bson:"login_url"`
	ValidateUrl string `bson:"validate_url"`
	DefaultRole string `bson:"default_role"`
}
