package alarm

//go:generate mockgen -destination=../../../mocks/lib/canopsis/alarm/alarm.go git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm Adapter,Service,EventProcessor

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type Adapter interface {
	// Insert insert an alarm
	Insert(alarm types.Alarm) error

	// Update update an alarm
	Update(alarm types.Alarm) error

	PartialUpdateOpen(alarm *types.Alarm) error

	// RemoveAll remove all alarms
	RemoveAll() error

	// RemoveId remove an alarm from its id
	RemoveId(id string) error

	//changed bson.M to map[string]interface{} to supportability between mgo and mongo-driver bson.M without creating a new interface or functions
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
	GetCountOpenedAlarmsByIDs(ids []string) (int, error)

	// GetOpenedAlarmsByAlarmIDs gets ongoing alarms related the provided alarm ids
	GetOpenedAlarmsByAlarmIDs(ids []string, alarms *[]types.Alarm) error
	GetOpenedAlarmsWithEntityByAlarmIDs(ids []string, alarms *[]types.AlarmWithEntity) error

	// MassUpdate updates alarms with a mgo.Bulk object, by slices of max 1000 alarms.
	MassUpdate(alarms []types.Alarm, notUpdateResolved bool) error

	// MassUpdateWithEntity updates alarms with a mgo.Bulk object, by slices of max 1000 alarms.
	MassUpdateWithEntity(alarms []types.AlarmWithEntity) error

	// GetOpenedAlarmsWithLastEventDateBefore gets ongoing alarms which don't receive events after "date".
	GetOpenedAlarmsWithLastEventDateBefore(date time.Time) ([]types.AlarmWithEntity, error)

	// GetOpenedAlarmsWithLastUpdateDateBefore gets ongoing alarms which don't change after "date".
	GetOpenedAlarmsWithLastUpdateDateBefore(date time.Time) ([]types.AlarmWithEntity, error)

	CountResolvedAlarm(alarmList []string) (int, error)
}

type EventProcessor interface {
	// ProcessAlarmEvent processes an event and updates the corresponding
	// alarm. It enriches the event with this alarm, and returns an AlarmChange
	// representing the change that occured on this alarm and its previous
	// state.
	Process(ctx context.Context, event *types.Event) (types.AlarmChange, error)
}

type Service interface {
	// ResolveAlarms that have v.resolved to null
	ResolveAlarms(ctx context.Context, baggotTime time.Duration) ([]types.Alarm, error)

	// ResolveCancels close canceled alarms when time has expired
	ResolveCancels(ctx context.Context, cancelAutosolveDelay time.Duration) ([]types.Alarm, error)

	// ResolveDone close one alarms when time has expired
	ResolveDone(ctx context.Context) ([]types.Alarm, error)

	// ResolveSnoozes remove snooze state when snooze time has expired; DisableActionSnoozeDelayOnPbh as second argument
	ResolveSnoozes(context.Context, bool) ([]types.Alarm, error)

	// UpdateFlappingAlarms updates the status of the flapping alarms, removing
	// the flapping status if needed.
	UpdateFlappingAlarms(ctx context.Context) ([]types.Alarm, error)
}
