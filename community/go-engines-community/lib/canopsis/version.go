package canopsis

import (
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// Version is a version of service.
var Version string

// BuildDate is a Unix timestamp (as a string).
var BuildDate string

type BuildInfo struct {
	Name        string
	Version     string
	Edition     string
	Date        time.Time
	VcsRevision string
	GoVersion   string
	Os          string
}

func PrintVersionInfo() {
	bi := GetBuildInfo()

	fmt.Printf("%s version %s %s build %s %s %s %s\n", bi.Name, bi.Edition, bi.Version, bi.VcsRevision,
		bi.Date.Format(time.RFC3339), bi.GoVersion, bi.Os)
}

func GetBuildInfo() BuildInfo {
	bi := BuildInfo{}

	bi.Version = Version
	if bi.Version == "" {
		bi.Version = "development"
	}

	timestamp, err := strconv.ParseInt(BuildDate, 10, 64)
	if err == nil {
		bi.Date = time.Unix(timestamp, 0).UTC()
	}

	buildOs := ""
	buildArch := ""
	vcsModified := false
	const revLen = 12
	if runtimeBi, ok := debug.ReadBuildInfo(); ok {
		bi.Name = path.Base(runtimeBi.Path)
		editions := []string{"community", "pro"}
		for _, edition := range editions {
			if strings.Contains(runtimeBi.Path, edition) {
				bi.Edition = edition
				break
			}
		}

		bi.GoVersion = runtimeBi.GoVersion
		for _, setting := range runtimeBi.Settings {
			switch setting.Key {
			case "GOOS":
				buildOs = setting.Value
			case "GOARCH":
				buildArch = setting.Value
			case "vcs.revision":
				bi.VcsRevision = setting.Value
				if len(bi.VcsRevision) > revLen {
					bi.VcsRevision = bi.VcsRevision[:revLen]
				}
			case "vcs.modified":
				vcsModified, _ = strconv.ParseBool(setting.Value)
			}
		}
	}

	if bi.Name == "" {
		bi.Name = path.Base(os.Args[0])
	}

	if bi.VcsRevision == "" {
		bi.VcsRevision = "unknown"
	} else if vcsModified {
		bi.VcsRevision += " (dirty)"
	}

	bi.Os = fmt.Sprintf("%s/%s", buildOs, buildArch)

	return bi
}
