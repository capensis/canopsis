package mongosh

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

//go:embed helpers.js
var helpersFuncs []byte

func NewScriptExecutor() mongo.ScriptExecutor {
	return &scriptExecutor{}
}

type scriptExecutor struct{}

func (scriptExecutor) Exec(ctx context.Context, file string) (resErr error) {
	h, err := os.CreateTemp("", "canopsis_migration_*.js")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = h.Write(helpersFuncs)
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	err = h.Close()
	if err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	defer func() {
		err = os.Remove(h.Name())
		if err != nil && resErr == nil {
			resErr = fmt.Errorf("failed to remove temp file: %w", err)
		}
	}()

	cmd := fmt.Sprintf("mongosh %s %s %s", os.Getenv(mongo.EnvURL), h.Name(), file)
	result := exec.CommandContext(ctx, "bash", "-c", cmd)
	output, err := result.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute command %w: %s", err, string(output))
	}

	return nil
}
