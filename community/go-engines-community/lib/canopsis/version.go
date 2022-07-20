package canopsis

import (
	"fmt"
	"os"
	"path"
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
	GoVersion   string
	Os          string
}

func PrintVersionInfo() {
	bi := GetBuildInfo()

	fmt.Printf("%s version %s build %s %s %s %s\n", path.Base(os.Args[0]), bi.Version, bi.VcsRevision,
		bi.Date.Format(time.RFC3339), bi.GoVersion, bi.Os)
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

	goVersion := ""
	buildOs := ""
	buildArch := ""
	vcsRevision := ""
	vcsModified := false
	const revLen = 12
	if bi, ok := debug.ReadBuildInfo(); ok {
		goVersion = bi.GoVersion
		for _, setting := range bi.Settings {
			switch setting.Key {
			case "GOOS":
				buildOs = setting.Value
			case "GOARCH":
				buildArch = setting.Value
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
		GoVersion:   goVersion,
		Os:          fmt.Sprintf("%s/%s", buildOs, buildArch),
	}
}
