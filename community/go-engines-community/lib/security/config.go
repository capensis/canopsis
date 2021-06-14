package security

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"time"
)

// Config providers which auth methods must be used.
type Config struct {
	Security struct {
		AuthProviders []string `yaml:"auth_providers"`
	} `yaml:"security"`
	Session struct {
		StatsFrame time.Duration `yaml:"stats_frame"`
	} `yaml:"session"`
}

const AuthMethodBasic = "basic"
const AuthMethodApiKey = "apikey"
const AuthMethodCas = "cas"
const AuthMethodLdap = "ldap"

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
