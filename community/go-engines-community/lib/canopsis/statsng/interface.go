package statsng

//go:generate mockgen -destination=../../../mocks/lib/canopsis/statsng/statsng.go git.canopsis.net/canopsis/go-engines/lib/canopsis/statsng Service

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type Service interface {
	ProcessAlarmChange(
		ctx context.Context,
		change types.AlarmChange,
		timestamp types.CpsTime,
		alarm types.Alarm,
		entity types.Entity,
		author, eventType string,
		location *time.Location,
	) error

	ProcessResolvedAlarm(alarm types.Alarm, entity types.Entity) error
}
