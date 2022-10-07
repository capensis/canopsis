package postgres

import (
	"fmt"
	"os/exec"
)

// Dump calls pg_dump binary.
func Dump(connStr, filepath string) error {
	cmd := fmt.Sprintf("pg_dump %s -Fc -f %s", connStr, filepath)
	result := exec.Command("bash", "-c", cmd)
	output, err := result.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute command %q: %w: %s", cmd, err, string(output))
	}

	return nil
}
