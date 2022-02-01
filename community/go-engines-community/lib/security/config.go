package security

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config providers which auth methods must be used.
type Config struct {
	Security struct {
		AuthProviders []string `yaml:"auth_providers"`
		Saml          struct {
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

const configPath = "/api/security/config.yml"

// LoadConfig creates Config by config file.
func LoadConfig(configDir string) (*Config, error) {
	buf, err := ioutil.ReadFile(filepath.Join(configDir, configPath))

	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
