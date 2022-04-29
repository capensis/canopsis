package alarm

//go:generate mockgen -destination=../../../mocks/lib/canopsis/alarm/alarm.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm Adapter,Service,EventProcessor

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

type Adapter interface {
	// Insert insert an alarm
	Insert(alarm types.Alarm) error

	// Update update an alarm
	Update(alarm types.Alarm) error

	PartialUpdateOpen(ctx context.Context, alarm *types.Alarm) error

	// RemoveAll remove all alarms
	RemoveAll() error

	// RemoveId remove an alarm from its id
	RemoveId(id string) error

	Get(filter map[string]interface{}, alarms *[]types.Alarm) error

	// GetAlarmsByID finds all alarms with an entity id.
	GetAlarmsByID(id string) ([]types.Alarm, error)

	// GetAlarmsWithCancelMark returns all alarms where v.cancel is not null
	GetAlarmsWithCancelMark() ([]types.Alarm, error)

	// GetAlarmsWithDoneMark returns all alarms where v.done is not null
	GetAlarmsWithDoneMark() ([]types.Alarm, error)

	// GetAlarmsWithSnoozeMark returns all alarms where v.snooze is not null
	GetAlarmsWithSnoozeMark() ([]types.Alarm, error)

	// GetAlarmsWithFlappingStatus returns all alarms whose status is flapping
	GetAlarmsWithFlappingStatus() ([]types.Alarm, error)

	// GetAllOpenedResourceAlarmsByComponent returns all ongoing alarms for component
	GetAllOpenedResourceAlarmsByComponent(component string) ([]types.AlarmWithEntity, error)

	// GetUnacknowledgedAlarmsByComponent returns all ongoing alarms which have
	// not been acknowledged, given a component's name.
	GetUnacknowledgedAlarmsByComponent(component string) ([]types.Alarm, error)

	// GetAlarmsWithoutTicketByComponent returns all ongoing alarms which do
	// not have a ticket, given a component's name.
	GetAlarmsWithoutTicketByComponent(component string) ([]types.Alarm, error)

	GetOpenedAlarmByAlarmId(id string) (types.Alarm, error)

	GetAlarmByAlarmId(id string) (types.Alarm, error)

	// GetOpenedAlarm find one opened alarm with his entity id.
	// Note : a control is added to prevent fetching future alarms.
	GetOpenedAlarm(connector, connectorName, id string) (types.Alarm, error)

	GetOpenedMetaAlarm(ruleId string, valuePath string) (types.Alarm, error)

	// GetLastAlarm find the last alarm with an id
	GetLastAlarm(connector, connectorName, id string) (types.Alarm, error)

	// GetUnresolved returns all alarms that have v.resolved to null or absent field.
	GetUnresolved() ([]types.Alarm, error)

	// GetOpenedAlarmsByIDs gets ongoing alarms related the provided entity ids
	GetOpenedAlarmsByIDs(ids []string, alarms *[]types.Alarm) error
	GetOpenedAlarmsWithEntityByIDs(ids []string, alarms *[]types.AlarmWithEntity) error
	GetCountOpenedAlarmsByIDs(ids []string) (int64, error)

	// GetOpenedAlarmsByAlarmIDs gets ongoing alarms related the provided alarm ids
	GetOpenedAlarmsByAlarmIDs(ids []string, alarms *[]types.Alarm) error
	GetOpenedAlarmsWithEntityByAlarmIDs(ids []string, alarms *[]types.AlarmWithEntity) error

	MassUpdate(alarms []types.Alarm, notUpdateResolved bool) error

	MassUpdateWithEntity(alarms []types.AlarmWithEntity) error

	// MassPartialUpdateOpen updates opened alarms matching by list of IDs, applying partial update from alarm
	MassPartialUpdateOpen(context.Context, *types.Alarm, []string) error

	GetOpenedAlarmsWithLastDatesBefore(ctx context.Context, time types.CpsTime) (mongo.Cursor, error)

	GetOpenedAlarmsByConnectorIdleRules(ctx context.Context) ([]types.Alarm, error)

	CountResolvedAlarm(alarmList []string) (int, error)

	GetLastAlarmByEntityID(ctx context.Context, entityID string) (*types.Alarm, error)
}

type EventProcessor interface {
	// Process processes an event and updates the corresponding
	// alarm. It enriches the event with this alarm, and returns an AlarmChange
	// representing the change that occured on this alarm and its previous
	// state.
	Process(ctx context.Context, event *types.Event) (types.AlarmChange, error)
}

type Service interface {
	// ResolveAlarms that have v.resolved to null
	ResolveAlarms(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)

	// ResolveCancels close canceled alarms when time has expired
	ResolveCancels(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)

	// ResolveDone close one alarms when time has expired
	ResolveDone(ctx context.Context) ([]types.Alarm, error)

	// ResolveSnoozes remove snooze state when snooze time has expired
	ResolveSnoozes(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)

	// UpdateFlappingAlarms updates the status of the flapping alarms, removing
	// the flapping status if needed.
	UpdateFlappingAlarms(ctx context.Context, alarmConfig config.AlarmConfig) ([]types.Alarm, error)
}
