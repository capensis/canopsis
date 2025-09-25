package postgres

import (
	"context"
	"fmt"
	"os/exec"
)

// Dump calls pg_dump binary.
func Dump(ctx context.Context, connStr, filepath string) error {
	result := exec.CommandContext(ctx, "pg_dump", connStr, "-Fc", "-f", filepath)
	output, err := result.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute command \"pg_dump\": %w: %s", err, string(output))
	}

	return nil
}
