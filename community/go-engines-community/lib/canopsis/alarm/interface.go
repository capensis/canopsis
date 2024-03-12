package alarm

//go:generate mockgen -destination=../../../mocks/lib/canopsis/alarm/alarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm Adapter,Service,EventProcessor,ActivationService,MetaAlarmEventProcessor

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

type Adapter interface {
	// GetAlarmsWithCancelMark returns all alarms where v.cancel is not null
	GetAlarmsWithCancelMark(ctx context.Context) ([]types.Alarm, error)

	// GetAlarmsWithSnoozeMark returns all alarms where v.snooze is not null
	GetAlarmsWithSnoozeMark(ctx context.Context) ([]types.Alarm, error)

	// GetAlarmsWithFlappingStatus returns all alarms whose status is flapping
	GetAlarmsWithFlappingStatus(ctx context.Context) ([]types.AlarmWithEntity, error)

	// GetAlarmsWithoutTicketByComponent returns all ongoing alarms which do
	// not have a ticket, given a component's name.
	GetAlarmsWithoutTicketByComponent(ctx context.Context, component string) ([]types.AlarmWithEntity, error)

	GetOpenedAlarmByAlarmId(ctx context.Context, id string) (types.Alarm, error)
	GetAlarmByAlarmId(ctx context.Context, id string) (types.Alarm, error)

	// GetOpenedAlarmsByIDs gets ongoing alarms related the provided entity ids
	GetOpenedAlarmsByIDs(ctx context.Context, ids []string, alarms *[]types.Alarm) error
	GetOpenedAlarmsWithEntityByIDs(ctx context.Context, ids []string, alarms *[]types.AlarmWithEntity) error
	GetCountOpenedAlarmsByIDs(ctx context.Context, ids []string) (int64, error)
	GetOpenedAlarmsWithEntity(ctx context.Context) (mongo.Cursor, error)

	// GetOpenedAlarmsByAlarmIDs gets ongoing alarms related the provided alarm ids
	GetOpenedAlarmsByAlarmIDs(ctx context.Context, ids []string, alarms *[]types.Alarm) error

	GetOpenedAlarmsWithLastDatesBefore(ctx context.Context, time datetime.CpsTime) (mongo.Cursor, error)

	GetOpenedAlarmsByConnectorIdleRules(ctx context.Context) ([]types.Alarm, error)

	CountResolvedAlarm(ctx context.Context, alarmList []string) (int, error)

	GetLastAlarmByEntityID(ctx context.Context, entityID string) (*types.Alarm, error)

	// DeleteResolvedAlarms deletes resolved alarms from resolved collection after some duration
	DeleteResolvedAlarms(ctx context.Context, duration time.Duration) error

	// CopyAlarmToResolvedCollection copies alarm to resolved alarm collection
	CopyAlarmToResolvedCollection(ctx context.Context, alarm types.Alarm) error

	FindToCheckPbehaviorInfo(ctx context.Context, createdBefore datetime.CpsTime, idsWithPbehaviors []string) (mongo.Cursor, error)

	GetWorstAlarmStateAndMaxLastEventDate(ctx context.Context, entityIds []string) (int64, int64, error)

	UpdateLastEventDate(ctx context.Context, entityIds []string, t datetime.CpsTime) error
}

type EventProcessor interface {
	// Process processes an event and updates the corresponding
	// alarm. It enriches the event with this alarm, and returns an AlarmChange
	// representing the change that occurred on this alarm and its previous
	// state.
	Process(ctx context.Context, event *types.Event) (types.AlarmChange, error)
}

type MetaAlarmEventProcessor interface {
	// ProcessAxeRpc handles related meta alarm parents and children after alarm change.
	ProcessAxeRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error
	// CreateMetaAlarm creates meta alarm by event.
	CreateMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, []types.Alarm, error)
	// AttachChildrenToMetaAlarm attaches children to meta alarm by event.
	AttachChildrenToMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, []types.Alarm, []types.Event, error)
	// DetachChildrenFromMetaAlarm detaches children from meta alarm by event.
	DetachChildrenFromMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, error)
}

type Service interface {
	// ResolveClosed close ok alarms.
	ResolveClosed(ctx context.Context) ([]types.Alarm, error)

	// ResolveCancels close canceled alarms when time has expired
	ResolveCancels(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)

	// ResolveSnoozes remove snooze state when snooze time has expired
	ResolveSnoozes(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)

	// UpdateFlappingAlarms updates the status of the flapping alarms, removing
	// the flapping status if needed.
	UpdateFlappingAlarms(ctx context.Context) ([]types.Alarm, error)
}
