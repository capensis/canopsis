package canopsis

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// BuildDate is a Unix timestamp (as a string)
var BuildDate string

// BuildGitCommit is the short version of Git commit
var BuildGitCommit string

// PrintVersion outputs version information
func PrintVersion() {
	timestamp, err := strconv.ParseInt(BuildDate, 10, 64)
	if err != nil {
		timestamp = 0
	}

	fmt.Printf("build date  : %s\n", time.Unix(timestamp, 0).UTC().Format(time.RFC3339))
	fmt.Printf("build commit: %s\n", BuildGitCommit)
}

// PrintVersionExit calls PrintVersion then exit(0)
func PrintVersionExit() {
	PrintVersion()
	os.Exit(0)
}
