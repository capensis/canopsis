package api

import (
	"net/url"
	"os"

	"github.com/rs/zerolog"
)

// GetLegacyURL returns old API url from env var.
func GetLegacyURL(logger zerolog.Logger) *url.URL {
	var err error
	legacy := os.Getenv(EnvOldApiURL)
	if legacy == "" {
		logger.Warn().Msgf("%s is empty", EnvOldApiURL)
		return nil
	}

	parsed, err := url.Parse(legacy)
	if err != nil {
		logger.Err(err).Str("url", legacy).Msgf("cannot parse %s", EnvOldApiURL)
		return nil
	}

	return parsed
}
