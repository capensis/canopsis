package mongo

import "context"

const MigrationHelperComment = `// Available global functions:
// genID returns a new string UUID
// isInt checks if a value is integer
// toInt transforms value to integer`

// ScriptExecutor is used to execute JavaScript migration scripts.
// Each implementation has to support global functions described above.
type ScriptExecutor interface {
	Exec(ctx context.Context, file string) error
}
