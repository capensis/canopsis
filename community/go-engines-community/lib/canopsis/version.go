package canopsis

import (
	"fmt"
	"os"
)

// BuildDate ...
var BuildDate string

// BuildGitBranch ...
var BuildGitBranch string

// BuildGitCommit is the short version of git commit
var BuildGitCommit string

// PrintVersion outputs versions informations
func PrintVersion() {
	fmt.Printf("build date  : %s\n", BuildDate)
	fmt.Printf("build commit: %s\n", BuildGitCommit)
	fmt.Printf("build branch: %s\n", BuildGitBranch)
}

// PrintVersionExit calls PrintVersion then exit(0)
func PrintVersionExit() {
	PrintVersion()
	os.Exit(0)
}
