package security

import (
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
)

const DefaultInactivityInterval = 24 // hours

const configPath = "/api/security/config.yml"

// Config providers which auth methods must be used.
type Config struct {
	Security struct {
		AuthProviders []string    `yaml:"auth_providers"`
		Basic         BasicConfig `yaml:"basic"`
		Ldap          LdapConfig  `yaml:"ldap"`
		Cas           CasConfig   `yaml:"cas"`
		Saml          SamlConfig  `yaml:"saml"`
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

func validateConfig(config Config) error {
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

	return nil
}
