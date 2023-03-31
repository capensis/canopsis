package proxy

//go:generate mockgen -destination=../../../mocks/lib/security/proxy/proxy.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/proxy AccessConfig

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const configPath = "/api/security/routes_permissions_map.yml"

// AccessConfig checks if proxy uri must be protected by permission.
type AccessConfig interface {
	// Get returns obj and act if proxy uri is protected by permission.
	Get(uri, method string) (string, string)
}

// fileAccessConfig uses file which contains map : uri prefix, http method -> object, act.
type fileAccessConfig struct {
	Map map[string]map[string][]string `yaml:"map"`
}

func (m *fileAccessConfig) Get(uri, method string) (string, string) {
	for prefix, conf := range m.Map {
		if strings.HasPrefix(uri, prefix) {
			if perm, ok := conf[method]; ok {
				if len(perm) == 2 {
					return perm[0], perm[1]
				}
			}
		}
	}

	return "", ""
}

// LoadAccessConfig creates AccessConfig by config file.
func LoadAccessConfig(configDir string) (AccessConfig, error) {
	buf, err := os.ReadFile(filepath.Join(configDir, configPath))

	if err != nil {
		return nil, err
	}

	var config fileAccessConfig
	err = yaml.Unmarshal(buf, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
