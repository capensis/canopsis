package flag

import (
	"flag"
	"slices"

	"github.com/rs/zerolog"
)

func FindDeprecatedFlags(list []string) []string {
	var deprecatedFlags []string

	flag.Visit(func(f *flag.Flag) {
		if slices.Contains(list, f.Name) {
			deprecatedFlags = append(deprecatedFlags, f.Name)
		}
	})

	return deprecatedFlags
}

func LogDeprecatedFlags(logger zerolog.Logger, deprecatedFlags []string) {
	for _, f := range deprecatedFlags {
		logger.Warn().Msgf("The flag %q is deprecated and will be removed in a future versions.", f)
	}
}
