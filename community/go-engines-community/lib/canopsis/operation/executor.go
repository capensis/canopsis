/*
Package operation implements alarm modification operations.
*/
package operation

//go:generate mockgen -destination=../../../mocks/lib/canopsis/operation/executor.go git.canopsis.net/canopsis/go-engines/lib/canopsis/operation Executor

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

// Executor interface is used to implement an alarm modification operation.
type Executor interface {
	// Exec modifies alarm.
	Exec(
		operation types.Operation,
		alarm *types.Alarm,
		timestamp types.CpsTime,
		role, initiator string,
	) (types.AlarmChangeType, error)
}
