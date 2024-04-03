package security

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"gopkg.in/yaml.v3"
)

const (
	AuthMethodBasic  = "basic"
	AuthMethodApiKey = "apikey"
	AuthMethodCas    = "cas"
	AuthMethodSaml   = "saml"
	AuthMethodLdap   = "ldap"
	AuthMethodOAuth2 = "oauth2"
)

const DefaultInactivityInterval = 24 // hours

const configPath = "/api/security/config.yml"

// Config providers which auth methods must be used.
type Config struct {
	Security struct {
		AuthProviders []string     `yaml:"auth_providers"`
		Basic         BasicConfig  `yaml:"basic"`
		Ldap          LdapConfig   `yaml:"ldap"`
		Cas           CasConfig    `yaml:"cas"`
		Saml          SamlConfig   `yaml:"saml"`
		OAuth2        OAuth2Config `yaml:"oauth2"`
	} `yaml:"security"`
}

type BasicConfig struct {
	InactivityInterval string `yaml:"inactivity_interval"`
	ExpirationInterval string `yaml:"expiration_interval"`
}

type LdapConfig struct {
	InactivityInterval string            `yaml:"inactivity_interval"`
	ExpirationInterval string            `yaml:"expiration_interval"`
	Url                string            `yaml:"url"`
	AdminUsername      string            `yaml:"admin_dn"`
	AdminPassword      string            `yaml:"admin_passwd"`
	BaseDN             string            `yaml:"user_dn"`
	Attributes         map[string]string `yaml:"attrs"`
	UsernameAttr       string            `yaml:"username_attr"`
	Filter             string            `yaml:"ufilter"`
	DefaultRole        string            `yaml:"default_role"`
	InsecureSkipVerify bool              `yaml:"insecure_skip_verify"`
	MinTLSVersion      string            `yaml:"min_tls_ver"`
	MaxTLSVersion      string            `yaml:"max_tls_ver"`
}

type CasConfig struct {
	InactivityInterval string `yaml:"inactivity_interval"`
	ExpirationInterval string `yaml:"expiration_interval"`
	Title              string `yaml:"title"`
	LoginUrl           string `yaml:"login_url"`
	ValidateUrl        string `yaml:"validate_url"`
	DefaultRole        string `yaml:"default_role"`
}

type SamlConfig struct {
	InactivityInterval      string            `yaml:"inactivity_interval"`
	ExpirationInterval      string            `yaml:"expiration_interval"`
	Title                   string            `yaml:"title"`
	X509Cert                string            `yaml:"x509_cert"`
	X509Key                 string            `yaml:"x509_key"`
	IdpMetadataUrl          string            `yaml:"idp_metadata_url"`
	IdpMetadataXml          string            `yaml:"idp_metadata_xml"`
	IdpAttributesMap        map[string]string `yaml:"idp_attributes_map"`
	CanopsisSamlUrl         string            `yaml:"canopsis_saml_url"`
	DefaultRole             string            `yaml:"default_role"`
	InsecureSkipVerify      bool              `yaml:"insecure_skip_verify"`
	CanopsisSSOBinding      string            `yaml:"canopsis_sso_binding"`
	CanopsisACSBinding      string            `yaml:"canopsis_acs_binding"`
	SignAuthRequest         bool              `yaml:"sign_auth_request"`
	NameIdFormat            string            `yaml:"name_id_format"`
	SkipSignatureValidation bool              `yaml:"skip_signature_validation"`
	ACSIndex                *int              `yaml:"acs_index"`
	AutoUserRegistration    bool              `yaml:"auto_user_registration"`
}

type OAuth2Config struct {
	Providers map[string]OAuth2ProviderConfig `yaml:"providers"`
}

type OAuth2ProviderConfig struct {
	InactivityInterval string            `yaml:"inactivity_interval"`
	ExpirationInterval string            `yaml:"expiration_interval"`
	Issuer             string            `yaml:"issuer"`
	ClientID           string            `yaml:"client_id"`
	ClientSecret       string            `yaml:"client_secret"`
	RedirectURL        string            `yaml:"redirect_url"`
	DefaultRole        string            `yaml:"default_role"`
	AuthURL            string            `yaml:"auth_url"`
	TokenURL           string            `yaml:"token_url"`
	UserURL            string            `yaml:"user_url"`
	UserID             string            `yaml:"user_id"`
	Scopes             []string          `yaml:"scopes"`
	AttributesMap      map[string]string `yaml:"attributes_map"`
	OpenID             bool              `yaml:"open_id"`
	PKCE               bool              `yaml:"pkce"`
}

// LoadConfig creates Config by config file.
func LoadConfig(configDir string) (Config, error) {
	buf, err := os.ReadFile(filepath.Join(configDir, configPath))
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return Config{}, err
	}

	err = validateConfig(config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func validateOAuth2Config(config OAuth2ProviderConfig) error {
	if config.ExpirationInterval != "" {
		_, err := datetime.ParseDurationWithUnit(config.ExpirationInterval)
		if err != nil {
			return fmt.Errorf("invalid expiration_interval: %w", err)
		}
	}

	if config.InactivityInterval != "" {
		_, err := datetime.ParseDurationWithUnit(config.InactivityInterval)
		if err != nil {
			return fmt.Errorf("invalid inactivity_interval: %w", err)
		}
	}

	if config.OpenID {
		if config.Issuer == "" {
			return errors.New("issuer shouldn't be empty")
		}
	} else {
		if config.UserID == "" {
			return errors.New("user_id shouldn't be empty")
		}

		if config.TokenURL == "" {
			return errors.New("token_url shouldn't be empty")
		}

		if config.AuthURL == "" {
			return errors.New("auth_url shouldn't be empty")
		}

		if config.UserURL == "" {
			return errors.New("user_url shouldn't be empty")
		}
	}

	if config.RedirectURL == "" {
		return errors.New("redirect_url shouldn't be empty")
	}

	if config.ClientID == "" {
		return errors.New("client_id shouldn't be empty")
	}

	if config.ClientSecret == "" {
		return errors.New("client_secret shouldn't be empty")
	}

	if config.DefaultRole == "" {
		return errors.New("default_role is required")
	}

	return nil
}

func validateConfig(config Config) error {
	for _, name := range config.Security.AuthProviders {
		switch name {
		case AuthMethodBasic:
			if config.Security.Basic.ExpirationInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Basic.ExpirationInterval)
				if err != nil {
					return fmt.Errorf("invalid basic.expiration_interval: %w", err)
				}
			}

			if config.Security.Basic.InactivityInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Basic.InactivityInterval)
				if err != nil {
					return fmt.Errorf("invalid basic.inactivity_interval: %w", err)
				}
			}
		case AuthMethodLdap:
			if config.Security.Ldap.ExpirationInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Ldap.ExpirationInterval)
				if err != nil {
					return fmt.Errorf("invalid ldap.expiration_interval: %w", err)
				}
			}

			if config.Security.Ldap.InactivityInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Ldap.InactivityInterval)
				if err != nil {
					return fmt.Errorf("invalid ldap.inactivity_interval: %w", err)
				}
			}

			if config.Security.Ldap.DefaultRole == "" {
				return errors.New("ldap.default_role is required")
			}
		case AuthMethodCas:
			if config.Security.Cas.ExpirationInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Cas.ExpirationInterval)
				if err != nil {
					return fmt.Errorf("invalid cas.expiration_interval: %w", err)
				}
			}

			if config.Security.Cas.InactivityInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Cas.InactivityInterval)
				if err != nil {
					return fmt.Errorf("invalid cas.inactivity_interval: %w", err)
				}
			}

			if config.Security.Cas.DefaultRole == "" {
				return errors.New("cas.default_role is required")
			}
		case AuthMethodSaml:
			if config.Security.Saml.ExpirationInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Saml.ExpirationInterval)
				if err != nil {
					return fmt.Errorf("invalid saml.expiration_interval: %w", err)
				}
			}

			if config.Security.Saml.InactivityInterval != "" {
				_, err := datetime.ParseDurationWithUnit(config.Security.Saml.InactivityInterval)
				if err != nil {
					return fmt.Errorf("invalid saml.inactivity_interval: %w", err)
				}
			}

			if config.Security.Saml.DefaultRole == "" {
				return errors.New("saml.default_role is required")
			}
		case AuthMethodOAuth2:
			for provider, cfg := range config.Security.OAuth2.Providers {
				if err := validateOAuth2Config(cfg); err != nil {
					return fmt.Errorf("invalid %s provider config: %w", provider, err)
				}
			}
		}
	}

	return nil
}
