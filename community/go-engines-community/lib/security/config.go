package security

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"gopkg.in/yaml.v2"
)

// Config providers which auth methods must be used.
type Config struct {
	Security struct {
		AuthProviders []string `yaml:"auth_providers"`
		Basic         struct {
			InactivityInterval string `yaml:"inactivity_interval"`
			ExpirationInterval string `yaml:"expiration_interval"`
		} `yaml:"basic"`
		Ldap struct {
			InactivityInterval string `yaml:"inactivity_interval"`
			ExpirationInterval string `yaml:"expiration_interval"`
		} `yaml:"ldap"`
		Cas struct {
			InactivityInterval string `yaml:"inactivity_interval"`
			ExpirationInterval string `yaml:"expiration_interval"`
		} `yaml:"cas"`
		Saml struct {
			InactivityInterval      string            `yaml:"inactivity_interval"`
			ExpirationInterval      string            `yaml:"expiration_interval"`
			X509Cert                string            `yaml:"x509_cert"`
			X509Key                 string            `yaml:"x509_key"`
			IdpMetadataUrl          string            `yaml:"idp_metadata_url"`
			IdpMetadataXml          string            `yaml:"idp_metadata_xml"`
			IdpAttributesMap        map[string]string `yaml:"idp_attributes_map"`
			CanopsisSamlUrl         string            `yaml:"canopsis_saml_url"`
			DefaultRole             string            `yaml:"default_role"`
			InsecureSkipVerify      bool              `yaml:"skip_verify"`
			CanopsisSSOBinding      string            `yaml:"canopsis_sso_binding"`
			CanopsisACSBinding      string            `yaml:"canopsis_acs_binding"`
			SignAuthRequest         bool              `yaml:"sign_auth_request"`
			NameIdFormat            string            `yaml:"name_id_format"`
			SkipSignatureValidation bool              `yaml:"skip_signature_validation"`
			ACSIndex                *int              `yaml:"acs_index"`
			AutoUserRegistration    bool              `yaml:"auto_user_registration"`
		} `yaml:"saml"`
	} `yaml:"security"`
}

const (
	AuthMethodBasic  = "basic"
	AuthMethodApiKey = "apikey"
	AuthMethodCas    = "cas"
	AuthMethodSaml   = "saml"
	AuthMethodLdap   = "ldap"
)

const (
	LdapConfigID = "cservice.ldapconfig"
	CasConfigID  = "cservice.casconfig"
	SamlConfigID = "cservice.saml2config"
)

const DefaultInactivityInterval = 24 // hours

const configPath = "/api/security/config.yml"

// LoadConfig creates Config by config file.
func LoadConfig(configDir string) (Config, error) {
	buf, err := ioutil.ReadFile(filepath.Join(configDir, configPath))
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
		_, err := types.ParseDurationWithUnit(config.Security.Basic.ExpirationInterval)
		if err != nil {
			return fmt.Errorf("invalid basic.expiration_interval: %w", err)
		}
	}
	if config.Security.Basic.InactivityInterval != "" {
		_, err := types.ParseDurationWithUnit(config.Security.Basic.InactivityInterval)
		if err != nil {
			return fmt.Errorf("invalid basic.inactivity_interval: %w", err)
		}
	}
	if config.Security.Ldap.ExpirationInterval != "" {
		_, err := types.ParseDurationWithUnit(config.Security.Ldap.ExpirationInterval)
		if err != nil {
			return fmt.Errorf("invalid ldap.expiration_interval: %w", err)
		}
	}
	if config.Security.Ldap.InactivityInterval != "" {
		_, err := types.ParseDurationWithUnit(config.Security.Ldap.InactivityInterval)
		if err != nil {
			return fmt.Errorf("invalid ldap.inactivity_interval: %w", err)
		}
	}
	if config.Security.Cas.ExpirationInterval != "" {
		_, err := types.ParseDurationWithUnit(config.Security.Cas.ExpirationInterval)
		if err != nil {
			return fmt.Errorf("invalid cas.expiration_interval: %w", err)
		}
	}
	if config.Security.Cas.InactivityInterval != "" {
		_, err := types.ParseDurationWithUnit(config.Security.Cas.InactivityInterval)
		if err != nil {
			return fmt.Errorf("invalid cas.inactivity_interval: %w", err)
		}
	}
	if config.Security.Saml.ExpirationInterval != "" {
		_, err := types.ParseDurationWithUnit(config.Security.Saml.ExpirationInterval)
		if err != nil {
			return fmt.Errorf("invalid saml.expiration_interval: %w", err)
		}
	}
	if config.Security.Saml.InactivityInterval != "" {
		_, err := types.ParseDurationWithUnit(config.Security.Saml.InactivityInterval)
		if err != nil {
			return fmt.Errorf("invalid saml.inactivity_interval: %w", err)
		}
	}

	return nil
}
