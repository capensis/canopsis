package postgres

import (
	"fmt"
	"os/exec"
)

// Dump calls pg_dump binary.
func Dump(connStr, filepath string) error {
	result := exec.Command("pg_dump", connStr, "-Fc", "-f", filepath)
	output, err := result.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute command \"pg_dump\": %w: %s", err, string(output))
	}

	return nil
}
