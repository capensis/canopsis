package main

import (
	"bytes"
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Config struct {
	Api       ApiConfig `yaml:"api"`
	Pbehavior struct {
		Type   string `yaml:"type"`
		Reason string `yaml:"reason"`
	}
}

type ApiConfig struct {
	Host               string `yaml:"host"`
	InsecureSkipverify bool   `yaml:"insecure_skip_verify"`

	Username string `yaml:"-"`
	Password string `yaml:"-"`
}

// LoadConfig reads a file in configPath path and parses its YAML content.
func LoadConfig(configPath string) (Config, error) {
	config := Config{}
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(buf, &config)
	return config, err
}

func (c ApiConfig) CreateRequest(ctx context.Context, method, path string, b []byte, q url.Values) (*http.Request, error) {
	u, err := url.Parse(c.Host)
	if err != nil {
		return nil, fmt.Errorf("host %q is invalid url: %w", c.Host, err)
	}
	u.Path = path
	if q != nil {
		u.RawQuery = q.Encode()
	}

	var body io.Reader
	if len(b) > 0 {
		body = bytes.NewReader(b)
	}
	request, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", canopsis.JsonContentType)
	request.Header.Set("Accept", canopsis.JsonContentType)

	request.SetBasicAuth(c.Username, c.Password)

	return request, nil
}
