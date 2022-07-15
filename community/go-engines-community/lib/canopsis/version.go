package canopsis

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

// Version is a version of service.
var Version string

// BuildDate is a Unix timestamp (as a string).
var BuildDate string

type BuildInfo struct {
	Version     string
	Date        time.Time
	VcsRevision string
}

// PrintVersion outputs version information
func PrintVersion() {
	bi := GetBuildInfo()

	fmt.Printf("build date  : %s\n", bi.Date.Format(time.RFC3339))
	fmt.Printf("build commit: %s\n", bi.VcsRevision)
}

// PrintVersionExit calls PrintVersion then exit(0)
func PrintVersionExit() {
	PrintVersion()
	os.Exit(0)
}

func PrintVersionInfo(name string) {
	bi := GetBuildInfo()

	fmt.Printf("%s version %s, build %s %s\n", name, bi.Version, bi.VcsRevision, bi.Date.Format(time.RFC3339))
}

func GetBuildInfo() BuildInfo {
	v := Version
	if v == "" {
		v = "development"
	}

	var buildDate time.Time
	timestamp, err := strconv.ParseInt(BuildDate, 10, 64)
	if err == nil {
		buildDate = time.Unix(timestamp, 0).UTC()
	}

	vcsRevision := ""
	vcsModified := false
	const revLen = 12
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range bi.Settings {
			switch setting.Key {
			case "vcs.revision":
				vcsRevision = setting.Value
				if len(vcsRevision) > revLen {
					vcsRevision = vcsRevision[:revLen]
				}
			case "vcs.modified":
				vcsModified, _ = strconv.ParseBool(setting.Value)
			}
		}
	}
	if vcsRevision == "" {
		vcsRevision = "unknown"
	} else if vcsModified {
		vcsRevision += " (dirty)"
	}

	return BuildInfo{
		Version:     v,
		Date:        buildDate,
		VcsRevision: vcsRevision,
	}
}
