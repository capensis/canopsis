/*
Package operation implements alarm modification operations.
*/
package operation

//go:generate mockgen -destination=../../../mocks/lib/canopsis/operation/executor.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation Executor

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// Executor interface is used to implement an alarm modification operation.
type Executor interface {
	// Exec modifies alarm.
	Exec(
		ctx context.Context,
		operation types.Operation,
		alarm *types.Alarm,
		entity *types.Entity,
		timestamp types.CpsTime,
		userID, role, initiator string,
	) (types.AlarmChangeType, error)
}
