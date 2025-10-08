package flag

import (
	"flag"
	"slices"

	"github.com/rs/zerolog"
)

// FindDeprecatedFlags should be used for deprecated flags at least for one release.
func FindDeprecatedFlags(flags ...string) []string {
	var deprecatedFlags []string
	if len(flags) == 0 {
		return deprecatedFlags
	}

	flag.Visit(func(f *flag.Flag) {
		if slices.Contains(flags, f.Name) {
			deprecatedFlags = append(deprecatedFlags, f.Name)
		}
	})

	return deprecatedFlags
}

func LogDeprecatedFlags(logger zerolog.Logger, deprecatedFlags []string) {
	if len(deprecatedFlags) != 0 {
		logger.Warn().Strs("flags", deprecatedFlags).Msg("Deprecated flags will be removed in a future versions.")
	}
}
