package alarm

//go:generate mockgen -destination=../../../mocks/lib/canopsis/alarm/alarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm Adapter,Service,EventProcessor

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

type Adapter interface {
	// Insert insert an alarm
	Insert(ctx context.Context, alarm types.Alarm) error

	// Update update an alarm
	Update(ctx context.Context, alarm types.Alarm) error

	PartialUpdateOpen(ctx context.Context, alarm *types.Alarm) error

	PartialMassUpdateOpen(ctx context.Context, alarms []types.Alarm) error

	// GetAlarmsByID finds all alarms with an entity id.
	GetAlarmsByID(ctx context.Context, id string) ([]types.Alarm, error)

	// GetAlarmsWithCancelMark returns all alarms where v.cancel is not null
	GetAlarmsWithCancelMark(ctx context.Context) ([]types.Alarm, error)

	// GetAlarmsWithDoneMark returns all alarms where v.done is not null
	GetAlarmsWithDoneMark(ctx context.Context) ([]types.Alarm, error)

	// GetAlarmsWithSnoozeMark returns all alarms where v.snooze is not null
	GetAlarmsWithSnoozeMark(ctx context.Context) ([]types.Alarm, error)

	// GetAlarmsWithFlappingStatus returns all alarms whose status is flapping
	GetAlarmsWithFlappingStatus(ctx context.Context) ([]types.AlarmWithEntity, error)

	// GetAllOpenedResourceAlarmsByComponent returns all ongoing alarms for component
	GetAllOpenedResourceAlarmsByComponent(ctx context.Context, component string) ([]types.AlarmWithEntity, error)

	// GetUnacknowledgedAlarmsByComponent returns all ongoing alarms which have
	// not been acknowledged, given a component's name.
	GetUnacknowledgedAlarmsByComponent(ctx context.Context, component string) ([]types.AlarmWithEntity, error)

	// GetAlarmsWithoutTicketByComponent returns all ongoing alarms which do
	// not have a ticket, given a component's name.
	GetAlarmsWithoutTicketByComponent(ctx context.Context, component string) ([]types.AlarmWithEntity, error)

	GetOpenedAlarmByAlarmId(ctx context.Context, id string) (types.Alarm, error)
	GetAlarmByAlarmId(ctx context.Context, id string) (types.Alarm, error)

	// GetOpenedAlarm find one opened alarm with his entity id.
	// Note : a control is added to prevent fetching future alarms.
	GetOpenedAlarm(ctx context.Context, connector, connectorName, id string) (types.Alarm, error)

	GetOpenedMetaAlarm(ctx context.Context, ruleId string, valuePath string) (types.Alarm, error)
	GetOpenedMetaAlarmWithEntity(ctx context.Context, ruleId string, valuePath string) (types.AlarmWithEntity, error)

	// GetLastAlarm find the last alarm with an id
	GetLastAlarm(ctx context.Context, connector, connectorName, id string) (types.Alarm, error)
	GetLastAlarmWithEntity(ctx context.Context, connector, connectorName, id string) (types.AlarmWithEntity, error)

	// GetOpenedAlarmsByIDs gets ongoing alarms related the provided entity ids
	GetOpenedAlarmsByIDs(ctx context.Context, ids []string, alarms *[]types.Alarm) error
	GetOpenedAlarmsWithEntityByIDs(ctx context.Context, ids []string, alarms *[]types.AlarmWithEntity) error
	GetCountOpenedAlarmsByIDs(ctx context.Context, ids []string) (int64, error)
	GetOpenedAlarmsWithEntity(ctx context.Context) (mongo.Cursor, error)

	// GetOpenedAlarmsByAlarmIDs gets ongoing alarms related the provided alarm ids
	GetOpenedAlarmsByAlarmIDs(ctx context.Context, ids []string, alarms *[]types.Alarm) error
	GetOpenedAlarmsWithEntityByAlarmIDs(ctx context.Context, ids []string, alarms *[]types.AlarmWithEntity) error

	MassUpdate(ctx context.Context, alarms []types.Alarm, notUpdateResolved bool) error

	// MassPartialUpdateOpen updates opened alarms matching by list of IDs, applying partial update from alarm
	MassPartialUpdateOpen(context.Context, *types.Alarm, []string) error

	GetOpenedAlarmsWithLastDatesBefore(ctx context.Context, time types.CpsTime) (mongo.Cursor, error)

	GetOpenedAlarmsByConnectorIdleRules(ctx context.Context) ([]types.Alarm, error)

	GetOpenedAlarmsWithEntityAfter(ctx context.Context, createdAfter types.CpsTime) (mongo.Cursor, error)

	CountResolvedAlarm(ctx context.Context, alarmList []string) (int, error)

	GetLastAlarmByEntityID(ctx context.Context, entityID string) (*types.Alarm, error)

	// DeleteResolvedAlarms deletes resolved alarms from resolved collection after some duration
	DeleteResolvedAlarms(ctx context.Context, duration time.Duration) error

	// DeleteResolvedAlarms deletes resolved alarms from archived collection after some duration
	DeleteArchivedResolvedAlarms(ctx context.Context, duration time.Duration) (int64, error)

	// CopyAlarmToResolvedCollection copies alarm to resolved alarm collection
	CopyAlarmToResolvedCollection(ctx context.Context, alarm types.Alarm) error

	// ArchiveResolvedAlarms archives alarm to archived alarm collection
	ArchiveResolvedAlarms(ctx context.Context, duration time.Duration) (int64, error)

	FindToCheckPbehaviorInfo(ctx context.Context, createdAfter types.CpsTime, idsWithPbehaviors []string) (mongo.Cursor, error)
}

type EventProcessor interface {
	// Process processes an event and updates the corresponding
	// alarm. It enriches the event with this alarm, and returns an AlarmChange
	// representing the change that occured on this alarm and its previous
	// state.
	Process(ctx context.Context, event *types.Event) (types.AlarmChange, error)
}

type Service interface {
	// ResolveClosed close ok alarms.
	ResolveClosed(ctx context.Context) ([]types.Alarm, error)

	// ResolveCancels close canceled alarms when time has expired
	ResolveCancels(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)

	// ResolveDone close one alarms when time has expired
	ResolveDone(ctx context.Context) ([]types.Alarm, error)

	// ResolveSnoozes remove snooze state when snooze time has expired
	ResolveSnoozes(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)

	// UpdateFlappingAlarms updates the status of the flapping alarms, removing
	// the flapping status if needed.
	UpdateFlappingAlarms(ctx context.Context) ([]types.Alarm, error)
}
