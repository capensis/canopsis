package api

import (
	"net/url"
	"os"
)

// GetLegacyURL returns old API url from env var or default url.
func GetLegacyURL() (parsed *url.URL) {
	var err error
	legacy := os.Getenv(EnvOldApiURL)
	if legacy != "" {
		parsed, err = url.Parse(legacy)
	}
	if parsed == nil || err != nil {
		return &url.URL{
			Scheme: "http",
			Host:   DefaultOldAPI,
		}
	}
	return parsed
}
