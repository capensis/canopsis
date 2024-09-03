package mongo

import (
	"fmt"
	"os"
	"os/exec"
)

type ScriptExecutor interface {
	Exec(file string) error
}

func NewScriptExecutor() ScriptExecutor {
	return &scriptExecutor{}
}

type scriptExecutor struct{}

func (scriptExecutor) Exec(file string) error {
	cmd := fmt.Sprintf("mongosh %q %s", os.Getenv(EnvURL), file)
	result := exec.Command("bash", "-c", cmd)
	output, err := result.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute command %q: %w: %s", cmd, err, string(output))
	}

	return nil
}
