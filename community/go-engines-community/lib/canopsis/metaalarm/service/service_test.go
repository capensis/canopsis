package service_test

import (
	"context"
	"encoding/json"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAddChildToMetaAlarmWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	const expectedState = types.AlarmStateCritical

	rulesAdapter, er := getRuleAdapter()
	if er != nil {
		t.Error(er)
	}

	s := service.NewMetaAlarmService(alarmAdapterMock, rulesAdapter,
		alarmConfigProviderMock, log.NewTestLogger())

	alarmJSON := []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 2
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue alarm component",
					"val" : 2
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [ 
				"testvaluecomponent922-110", 
				"testvaluecomponent922-112", 
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
	var metaAlarm types.Alarm
	err := json.Unmarshal(alarmJSON, &metaAlarm)
	if err != nil {
		t.Fatalf("MetaAlarm unmarshal error %s", err)
	}

	alarmJSON = []byte(`{
		"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
		"t" : 1594204048,
		"d" : "testvalueressource311/testvaluecomponent922-A311",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 3
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 3
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-A311",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594204048,
			"display_name" : "VX-LV-EA",
			"initial_output" : "testvalue alarm",
			"output" : "testvalue alarm",
			"initial_long_output" : "",
			"last_update_date" : 1594204048,
			"last_event_date" : 1594204048,
			"resource" : "testvalueressource311",
			"parents" : [ 
				"testvaluecomponent922-A311"
			],
			"children" : [],
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	var child types.AlarmWithEntity
	err = json.Unmarshal(alarmJSON, &child.Alarm)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	rule := metaalarm.Rule{}
	event := &types.Event{}

	var updatedAlarms *[]types.Alarm

	alarmAdapterMock.
		EXPECT().
		MassUpdate(gomock.Any(), gomock.Any(), true).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm, _ bool) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any(), gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddChildToMetaAlarm(event, metaAlarm, child, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	if state := (*updatedAlarms)[0].CurrentState(); state != expectedState {
		t.Errorf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

func TestAddMultipleChildsToMetaAlarmWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	const expectedState = types.AlarmStateCritical

	rulesAdapter, er := getRuleAdapter()
	if er != nil {
		t.Error(er)
	}
	s := service.NewMetaAlarmService(alarmAdapterMock, rulesAdapter,
		alarmConfigProviderMock, log.NewTestLogger())

	alarmJSON := []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue alarm component",
					"val" : 2
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [ 
				"testvaluecomponent922-112", 
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
	var metaAlarm types.Alarm
	err := json.Unmarshal(alarmJSON, &metaAlarm)
	if err != nil {
		t.Fatalf("MetaAlarm unmarshal error %s", err)
	}

	alarmJSON = []byte(`{
		"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
		"t" : 1594204048,
		"d" : "testvalueressource311/testvaluecomponent922-A311",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 2
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 2
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-A311",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594204048,
			"display_name" : "VX-LV-EA",
			"initial_output" : "testvalue alarm",
			"output" : "testvalue alarm",
			"initial_long_output" : "",
			"last_update_date" : 1594204048,
			"last_event_date" : 1594204048,
			"resource" : "testvalueressource311",
			"parents" : [ 
				"testvaluecomponent922-A311"
			],
			"children" : [],
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	var child types.AlarmWithEntity
	err = json.Unmarshal(alarmJSON, &child.Alarm)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	children := []types.AlarmWithEntity{child}

	alarmJSON = []byte(`{
		"_id" : "9f327dd2-763f-4997-9a09-8f3a4b94e155",
		"t" : 1594203725,
		"d" : "testvaluecomponent922-110",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm component",
				"val" : 3
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594203725,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm component",
					"val" : 3
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm component",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-110",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594203725,
			"display_name" : "PW-DK-XZ",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	err = json.Unmarshal(alarmJSON, &child.Alarm)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	children = append(children, child)
	rule := metaalarm.Rule{}
	event := &types.Event{}

	var updatedAlarms *[]types.Alarm

	alarmAdapterMock.
		EXPECT().
		MassUpdate(gomock.Any(), gomock.Any(), true).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm, _ bool) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any(), gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddMultipleChildsToMetaAlarm(event, metaAlarm, children, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	if state := (*updatedAlarms)[0].CurrentState(); state != expectedState {
		t.Errorf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

func TestUpdateChildToMetaAlarmIncreaseWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	const expectedState = types.AlarmStateCritical

	rulesAdapter, er := getRuleAdapter()
	if er != nil {
		t.Error(er)
	}
	s := service.NewMetaAlarmService(alarmAdapterMock, rulesAdapter,
		alarmConfigProviderMock, log.NewTestLogger())

	alarmJSON := []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 2
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue alarm component",
					"val" : 2
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [ 
				"testvaluecomponent922-110", 
				"testvaluecomponent922-112", 
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
	var metaAlarm types.Alarm
	err := json.Unmarshal(alarmJSON, &metaAlarm)
	if err != nil {
		t.Fatalf("MetaAlarm unmarshal error %s", err)
	}

	alarmJSON = []byte(`{
		"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
		"t" : 1594204048,
		"d" : "testvaluecomponent922-110",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 3
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 3
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-A311",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594204048,
			"display_name" : "VX-LV-EA",
			"initial_output" : "testvalue alarm",
			"output" : "testvalue alarm",
			"initial_long_output" : "",
			"last_update_date" : 1594204048,
			"last_event_date" : 1594204048,
			"resource" : "testvalueressource311",
			"parents" : [ 
				"testvaluecomponent922-A311"
			],
			"children" : [],
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	var child types.AlarmWithEntity
	err = json.Unmarshal(alarmJSON, &child.Alarm)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	rule := metaalarm.Rule{}
	event := &types.Event{}
	var updatedAlarms *[]types.Alarm

	alarmAdapterMock.
		EXPECT().
		MassUpdate(gomock.Any(), gomock.Any(), true).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm, _ bool) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any(), gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddChildToMetaAlarm(event, metaAlarm, child, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	// AddChildToMetaAlarm has metaAlarm as first item of updatedAlarms
	validatedAlarm := (*updatedAlarms)[0]

	if !validatedAlarm.IsMetaAlarm() {
		t.Fatalf("validatedAlarm is not meta-alarm %+v", validatedAlarm)
	}

	if state := validatedAlarm.CurrentState(); state != expectedState {
		t.Errorf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

func TestUpdateChildToMetaAlarmDecreaseWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)
	alarmConfigProviderMock.EXPECT().Get().Return(config.AlarmConfig{})

	const expectedState = types.AlarmStateMajor

	rulesAdapter, er := getRuleAdapter()
	if er != nil {
		t.Error(er)
	}
	s := service.NewMetaAlarmService(alarmAdapterMock, rulesAdapter,
		alarmConfigProviderMock, log.NewTestLogger())

	alarmJSON := []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 3
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue alarm component",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue alarm component",
					"val" : 2
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue alarm component",
			"output" : "testvalue alarm component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [ 
				"testvaluecomponent922-110", 
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
	var metaAlarm types.Alarm
	err := json.Unmarshal(alarmJSON, &metaAlarm)
	if err != nil {
		t.Fatalf("MetaAlarm unmarshal error %s", err)
	}

	alarmJSON = []byte(`{
		"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
		"t" : 1594204048,
		"d" : "testvaluecomponent922-110",
		"v" : {
			"state" : {
				"_t" : "stateinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 2
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594204048,
				"a" : "testvalueconnector.testvalueconnectorname",
				"m" : "testvalue alarm",
				"val" : 1
			},
			"steps" : [ 
				{
					"_t" : "stateinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 2
				}, 
				{
					"_t" : "statusinc",
					"t" : 1594204048,
					"a" : "testvalueconnector.testvalueconnectorname",
					"m" : "testvalue alarm",
					"val" : 1
				}
			],
			"component" : "testvaluecomponent922-A311",
			"connector" : "testvalueconnector",
			"connector_name" : "testvalueconnectorname",
			"creation_date" : 1594204048,
			"display_name" : "VX-LV-EA",
			"initial_output" : "testvalue alarm",
			"output" : "testvalue alarm",
			"initial_long_output" : "",
			"last_update_date" : 1594204048,
			"last_event_date" : 1594204048,
			"resource" : "testvalueressource311",
			"parents" : [ 
				"meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm"
			],
			"children" : [],
			"total_state_changes" : 1,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)

	var child types.Alarm
	err = json.Unmarshal(alarmJSON, &child)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}

	var otherChild types.Alarm
	err = json.Unmarshal(alarmJSON, &otherChild)
	if err != nil {
		t.Fatalf("Child Alarm unmarshal error %s", err)
	}
	otherChild.EntityID = "testvaluecomponent922-A311"
	rule := metaalarm.Rule{}
	event := &types.Event{}

	var (
		updatedAlarms *[]types.Alarm
		openedAlarms  []types.Alarm
	)
	childrenEIDs := make([]string, 2)
	copy(childrenEIDs, metaAlarm.Value.Children)

	alarmAdapterMock.
		EXPECT().
		GetOpenedAlarmsByIDs(gomock.Any(), childrenEIDs, &openedAlarms).
		DoAndReturn(func(_ context.Context, _ []string, c *[]types.Alarm) error {
			*c = []types.Alarm{child, otherChild}
			return nil
		})

	alarmAdapterMock.
		EXPECT().
		MassUpdate(gomock.Any(), gomock.Any(), true).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm, _ bool) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any(), gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddChildToMetaAlarm(event, metaAlarm, types.AlarmWithEntity{Alarm: child}, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	if state := (*updatedAlarms)[0].CurrentState(); state != expectedState {
		t.Errorf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

// All meta-alarm's children to update
func TestUpdateAllChildrenToMetaAlarmDecreaseWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)
	alarmConfigProviderMock.EXPECT().Get().Return(config.AlarmConfig{})

	const (
		expectedState   = types.AlarmStateMajor
		metaAlarmEntity = "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm"
	)

	rulesAdapter, er := getRuleAdapter()
	if er != nil {
		t.Error(er)
	}
	s := service.NewMetaAlarmService(alarmAdapterMock, rulesAdapter,
		alarmConfigProviderMock, log.NewTestLogger())

	childrenEIDs := []string{"a/b", "a/c"}

	metaAlarm := types.Alarm{
		EntityID: metaAlarmEntity,
		Value: types.AlarmValue{
			Meta:     " ",
			Children: childrenEIDs,
			State: &types.AlarmStep{
				Value: 3,
			},
		},
	}

	children := []types.AlarmWithEntity{
		{
			Alarm: types.Alarm{
				EntityID: childrenEIDs[0],
				Value: types.AlarmValue{
					State: &types.AlarmStep{
						Value: 1,
					},
					Parents: []string{metaAlarmEntity},
				},
			},
		},
		{
			Alarm: types.Alarm{
				EntityID: childrenEIDs[1],
				Value: types.AlarmValue{
					State: &types.AlarmStep{
						Value: 2,
					},
					Parents: []string{metaAlarmEntity},
				},
			},
		},
	}

	rule := metaalarm.Rule{}
	event := &types.Event{}

	var updatedAlarms *[]types.Alarm

	alarmAdapterMock.
		EXPECT().
		MassUpdate(gomock.Any(), gomock.Any(), true).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm, _ bool) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any(), gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddMultipleChildsToMetaAlarm(event, metaAlarm, children, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	if state := (*updatedAlarms)[len(*updatedAlarms)-1].CurrentState(); state != expectedState {
		t.Errorf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

func TestUpdateSomeChildrenToMetaAlarmDecreaseWorstState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)
	alarmConfigProviderMock.EXPECT().Get().Return(config.AlarmConfig{})

	const (
		expectedState   = types.AlarmStateMajor
		metaAlarmEntity = "meta-alarm-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm"
	)

	rulesAdapter, er := getRuleAdapter()
	if er != nil {
		t.Error(er)
	}
	s := service.NewMetaAlarmService(alarmAdapterMock, rulesAdapter,
		alarmConfigProviderMock, log.NewTestLogger())

	childrenEIDs := []string{"a/b", "a/c"}

	childrenToUpdate := []types.AlarmWithEntity{
		{
			Alarm: types.Alarm{
				ID:       "b-1",
				EntityID: childrenEIDs[0],
				Value: types.AlarmValue{
					State: &types.AlarmStep{
						Value: types.AlarmStateMajor,
					},
					Parents: []string{metaAlarmEntity},
				},
			},
		},
	}

	children := []types.Alarm{
		{
			ID:       "b-1",
			EntityID: childrenEIDs[0],
			Value: types.AlarmValue{
				State: &types.AlarmStep{
					Value: types.AlarmStateCritical,
				},
				Parents: []string{metaAlarmEntity},
			},
		},
		{
			ID:       "b-2",
			EntityID: childrenEIDs[1],
			Value: types.AlarmValue{
				State: &types.AlarmStep{
					Value: types.AlarmStateMinor,
				},
				Parents: []string{metaAlarmEntity},
			},
		},
	}

	metaAlarm := types.Alarm{
		EntityID: metaAlarmEntity,
		Value: types.AlarmValue{
			Meta:     " ",
			Children: childrenEIDs,
			State: &types.AlarmStep{
				Value: children[0].Value.State.Value,
			},
		},
	}

	rule := metaalarm.Rule{}
	event := &types.Event{}

	var (
		updatedAlarms *[]types.Alarm
		openedAlarms  []types.Alarm
	)

	alarmAdapterMock.
		EXPECT().
		GetOpenedAlarmsByIDs(gomock.Any(), childrenEIDs, &openedAlarms).
		DoAndReturn(func(_ context.Context, _ []string, c *[]types.Alarm) error {
			*c = children
			return nil
		})

	alarmAdapterMock.
		EXPECT().
		MassUpdate(gomock.Any(), gomock.Any(), true).
		DoAndReturn(func(_ context.Context, alarms []types.Alarm, _ bool) error {
			updatedAlarms = &alarms
			return nil
		})
	alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any(), gomock.Any()).Return(int64(1), nil)

	updateEvent, err := s.AddMultipleChildsToMetaAlarm(event, metaAlarm, childrenToUpdate, rule)
	if err != nil {
		t.Fatalf("AddChildToMetaAlarm error %s", err)
	}

	if state := (*updatedAlarms)[len(*updatedAlarms)-1].CurrentState(); state != expectedState {
		t.Errorf("wrong state value %d", state)
	}

	if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
		t.Fatalf("AddChildToMetaAlarm event %v", updateEvent)
	}
}

func TestChildInheritMetaAlarmActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmConfigProviderMock := mock_config.NewMockAlarmConfigProvider(ctrl)

	rulesAdapter, er := getRuleAdapter()
	if er != nil {
		t.Error(er)
	}

	s := service.NewMetaAlarmService(alarmAdapterMock, rulesAdapter,
		alarmConfigProviderMock, log.NewTestLogger())

	var dataSets = []struct {
		testName      string
		metaalarm     []byte
		child         []byte
		inheritAck    bool
		inheritTicket bool
		inheritSnooze bool
	}{
		{
			testName:  "test inherit all",
			metaalarm: getMetaalarm(),
			child: []byte(`{
				"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
				"t" : 1594204048,
				"d" : "testvalueressource311/testvaluecomponent922-A311",
				"v" : {
					"state" : {
						"_t" : "stateinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 2
					},
					"status" : {
						"_t" : "statusinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 1
					},
					"steps" : [ 
						{
							"_t" : "stateinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 2
						}, 
						{
							"_t" : "statusinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 1
						}
					],
					"component" : "testvaluecomponent922-A311",
					"connector" : "testvalueconnector",
					"connector_name" : "testvalueconnectorname",
					"creation_date" : 1594204048,
					"display_name" : "VX-LV-EA",
					"initial_output" : "testvalue child",
					"output" : "testvalue child",
					"initial_long_output" : "",
					"last_update_date" : 1594204048,
					"last_event_date" : 1594204048,
					"resource" : "testvalueressource311",
					"parents" : [ 
						"testvaluecomponent922-A311"
					],
					"children" : [],
					"total_state_changes" : 1,
					"extra" : {},
					"infos" : {},
					"infos_rule_version" : {}
				}
			}`),
			inheritAck:    true,
			inheritSnooze: true,
			inheritTicket: true,
		},
		{
			testName:  "test inherit none",
			metaalarm: getMetaalarm(),
			child: []byte(`{
				"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
				"t" : 1594204048,
				"d" : "testvalueressource311/testvaluecomponent922-A311",
				"v" : {
					"ack": {
						"_t": "ack",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Ack original alarm",
						"role": "admin",
						"val" : 0
					},
					"ticket": {
						"_t": "assocticket",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Ticket original alarm",
						"role": "admin",
						"val" : "Ticket original alarm"
					},
					"state" : {
						"_t" : "stateinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 2
					},
					"status" : {
						"_t" : "statusinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 1
					},
					"snooze": {
						"_t": "snooze",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Snooze original alarm",
						"role": "admin",
						"val" : 2594204000
					},
					"steps" : [ 
						{
							"_t" : "stateinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 2
						}, 
						{
							"_t" : "statusinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 1
						},
						{
							"_t": "assocticket",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Ticket original alarm",
							"role": "admin",
							"val" : 0
						},
						{
							"_t": "ack",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Ack original alarm",
							"role": "admin",
							"val" : 0
						},
						{
							"_t": "snooze",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Snooze original alarm",
							"role": "admin",
							"val" : 2594204000
						}
					],
					"component" : "testvaluecomponent922-A311",
					"connector" : "testvalueconnector",
					"connector_name" : "testvalueconnectorname",
					"creation_date" : 1594204048,
					"display_name" : "VX-LV-EA",
					"initial_output" : "testvalue child",
					"output" : "testvalue child",
					"initial_long_output" : "",
					"last_update_date" : 1594204048,
					"last_event_date" : 1594204048,
					"resource" : "testvalueressource311",
					"parents" : [ 
						"testvaluecomponent922-A311"
					],
					"children" : [],
					"total_state_changes" : 1,
					"extra" : {},
					"infos" : {},
					"infos_rule_version" : {}
				}
			}`),
			inheritAck:    false,
			inheritSnooze: false,
			inheritTicket: false,
		},
		{
			testName:  "test inherit some steps",
			metaalarm: getMetaalarm(),
			child: []byte(`{
				"_id" : "1e6e41ce-3823-4a44-9ee5-c31b7a9fd812",
				"t" : 1594204048,
				"d" : "testvalueressource311/testvaluecomponent922-A311",
				"v" : {
					"ack": {
						"_t": "ack",
						"t" : 1594204000,
						"a" : "root",
						"m" : "Ack original alarm",
						"role": "admin",
						"val" : 0
					},
					"state" : {
						"_t" : "stateinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 2
					},
					"status" : {
						"_t" : "statusinc",
						"t" : 1594204048,
						"a" : "testvalueconnector.testvalueconnectorname",
						"m" : "testvalue child",
						"val" : 1
					},
					"steps" : [ 
						{
							"_t" : "stateinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 2
						}, 
						{
							"_t" : "statusinc",
							"t" : 1594204048,
							"a" : "testvalueconnector.testvalueconnectorname",
							"m" : "testvalue child",
							"val" : 1
						},
						{
							"_t": "ack",
							"t" : 1594204000,
							"a" : "root",
							"m" : "Ack original alarm",
							"role": "admin",
							"val" : 0
						}
					],
					"component" : "testvaluecomponent922-A311",
					"connector" : "testvalueconnector",
					"connector_name" : "testvalueconnectorname",
					"creation_date" : 1594204048,
					"display_name" : "VX-LV-EA",
					"initial_output" : "testvalue child",
					"output" : "testvalue child",
					"initial_long_output" : "",
					"last_update_date" : 1594204048,
					"last_event_date" : 1594204048,
					"resource" : "testvalueressource311",
					"parents" : [ 
						"testvaluecomponent922-A311"
					],
					"children" : [],
					"total_state_changes" : 1,
					"extra" : {},
					"infos" : {},
					"infos_rule_version" : {}
				}
			}`),
			inheritAck:    false,
			inheritSnooze: true,
			inheritTicket: true,
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			var metaAlarm types.Alarm
			err := json.Unmarshal(dataset.metaalarm, &metaAlarm)
			if err != nil {
				t.Fatalf("MetaAlarm unmarshal error %s", err)
			}

			var child types.Alarm
			err = json.Unmarshal(dataset.child, &child)
			if err != nil {
				t.Fatalf("Child Alarm unmarshal error %s", err)
			}

			children := []types.AlarmWithEntity{
				{Alarm: child},
			}

			var updatedAlarms *[]types.Alarm

			alarmAdapterMock.
				EXPECT().
				MassUpdate(gomock.Any(), gomock.Any(), true).
				DoAndReturn(func(_ context.Context, alarms []types.Alarm, _ bool) error {
					updatedAlarms = &alarms
					return nil

				})
			alarmAdapterMock.EXPECT().GetCountOpenedAlarmsByIDs(gomock.Any(), gomock.Any()).Return(int64(1), nil)

			event := &types.Event{}
			rule := metaalarm.Rule{}
			updateEvent, err := s.AddMultipleChildsToMetaAlarm(event, metaAlarm, children, rule)
			if err != nil {
				t.Fatalf("AddMultipleChildsToMetaAlarm error %s", err)
			}

			if updateEvent.EventType != types.EventTypeMetaAlarmUpdated {
				t.Fatalf("AddMultipleChildsToMetaAlarm event %v", updateEvent)
			}

			for _, alarm := range *updatedAlarms {
				if alarm.EntityID != metaAlarm.EntityID {
					alarmAck := alarm.Value.ACK
					if alarmAck == nil {
						t.Fatalf("invalid child's ACK %s", alarm.EntityID)
					}

					if dataset.inheritAck {
						if alarmAck.Message != metaAlarm.Value.ACK.Message {
							t.Errorf("child's ACK message mismatch %s != %s", alarmAck.Message, metaAlarm.Value.ACK.Message)
						}
					} else {
						if alarmAck.Message == metaAlarm.Value.ACK.Message {
							t.Errorf("child's ACK message shouldn't be inherited")
						}
					}

					alarmTicket := alarm.Value.Ticket
					if alarmTicket == nil {
						t.Fatalf("invalid child's ACK %s", alarm.EntityID)
					}

					if dataset.inheritTicket {
						if alarmTicket.Message != metaAlarm.Value.Ticket.Message {
							t.Errorf("child's assocticket message mismatch %s != %s", alarmTicket.Message, metaAlarm.Value.Ticket.Message)
						}
					} else {
						if alarmTicket.Message == metaAlarm.Value.Ticket.Message {
							t.Errorf("child's assocticket message shouldn't be inherited")
						}
					}

					alarmSnooze := alarm.Value.Snooze
					if alarmSnooze == nil {
						t.Fatalf("invalid child's snooze %s", alarm.EntityID)
					}
					if dataset.inheritSnooze {
						if alarmSnooze.Message != metaAlarm.Value.Snooze.Message {
							t.Errorf("child's snooze message mismatch %s != %s", alarmSnooze.Message, metaAlarm.Value.Snooze.Message)
						}
					} else {
						if alarmSnooze.Message == metaAlarm.Value.Snooze.Message {
							t.Errorf("child's snooze message shouldn't be inherited")
						}
					}

					for _, step := range alarm.Value.Steps {
						switch step.Type {
						case alarmAck.Type:
							if dataset.inheritAck {
								if step.Message != metaAlarm.Value.ACK.Message {
									t.Errorf("step's ACK message mismatch %s != %s", step.Message, metaAlarm.Value.ACK.Message)
								}
								if step.Author != "correlation" {
									t.Errorf("invalid author value %s", step.Author)
								}
							} else {
								if step.Message == metaAlarm.Value.ACK.Message {
									t.Errorf("step's ACK message shouldn't be inherited")
								}
								if step.Author == "correlation" {
									t.Errorf("ack step's author shouldn't be %s", step.Author)
								}
							}
						case alarmTicket.Type:
							if dataset.inheritTicket {
								if step.Message != metaAlarm.Value.Ticket.Message {
									t.Errorf("step's assocticket message mismatch %s != %s", step.Message, metaAlarm.Value.Ticket.Message)
								}
								if step.Author != "correlation" {
									t.Errorf("invalid author value %s", step.Author)
								}
							} else {
								if step.Message == metaAlarm.Value.Ticket.Message {
									t.Errorf("step's assocticket message shouldn't be inherited")
								}
								if step.Author == "correlation" {
									t.Errorf("assocticket step's author shouldn't be %s", step.Author)
								}
							}
						case alarmSnooze.Type:
							if dataset.inheritSnooze {
								if step.Message != metaAlarm.Value.Snooze.Message {
									t.Errorf("step's snooze message mismatch %s != %s", step.Message, metaAlarm.Value.Snooze.Message)
								}
								if step.Author != "correlation" {
									t.Errorf("invalid author value %s", step.Author)
								}
							} else {
								if step.Message == metaAlarm.Value.Snooze.Message {
									t.Errorf("step's snooze message shouldn't be inherited")
								}
								if step.Author == "correlation" {
									t.Errorf("snooze step's author shouldn't be %s", step.Author)
								}
							}
						}
					}
				}
			}
		})
	}
}

func getMetaalarm() []byte {
	return []byte(`{
		"_id" : "40b8aaef-a975-491a-a143-2dfdabe059c8",
		"t" : 1594203725,
		"d" : "meta-child-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9/metaalarm",
		"v" : {
			"ack": {
				"_t": "ack",
				"t" : 1594204000,
				"a" : "root",
				"m" : "Ack meta-child",
				"role": "admin",
				"val" : 0
			},
			"ticket": {
				"_t": "assocticket",
				"t" : 1594204000,
				"a" : "root",
				"m" : "Ticket meta-child",
				"role": "admin",
				"val" : "Ticket meta-child"
			},
			"state" : {
				"_t" : "stateinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue child component",
				"val" : 1
			},
			"status" : {
				"_t" : "statusinc",
				"t" : 1594203725,
				"a" : "engine.correlation",
				"m" : "testvalue child component",
				"val" : 1
			},
			"snooze": {
				"_t": "snooze",
				"t" : 1594204000,
				"a" : "root",
				"m" : "Snooze meta-child",
				"role": "admin",
				"val" : 2594204000
			},
			"steps" : [ 
				{
					"_t" : "statusinc",
					"t" : 1594203725,
					"a" : "engine.correlation",
					"m" : "testvalue child component",
					"val" : 2
				},
				{
					"_t": "assocticket",
					"t" : 1594204000,
					"a" : "root",
					"m" : "Ticket meta-child",
					"role": "admin",
					"val" : 0
				},
				{
					"_t": "ack",
					"t" : 1594204000,
					"a" : "root",
					"m" : "Ack meta-child",
					"role": "admin",
					"val" : 0
				},
				{
					"_t": "snooze",
					"t" : 1594204000,
					"a" : "root",
					"m" : "Snooze meta-child",
					"role": "admin",
					"val" : 2594204000
				}
			],
			"component" : "metaalarm",
			"connector" : "engine",
			"connector_name" : "correlation",
			"creation_date" : 1594203725,
			"display_name" : "WP-NZ-KV",
			"initial_output" : "testvalue child component",
			"output" : "testvalue child component",
			"last_update_date" : 1594203725,
			"last_event_date" : 1594203725,
			"resource" : "meta-child-entity-24cc5fff-6d0e-4f64-ae91-7bd4a0cdc9b9",
			"meta" : "testRule",
			"parents" : [],
			"children" : [
				"testvaluecomponent922-A311"
			],
			"total_state_changes" : 0,
			"extra" : {},
			"infos" : {},
			"infos_rule_version" : {}
		}
	}`)
}

func getRuleAdapter() (metaalarm.RulesAdapter, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbClient, err := mongo.NewClient(ctx, 0, 0)
	if err != nil {
		panic(err)
	}

	rulesCollection := dbClient.Collection(mongo.MetaAlarmRulesMongoCollection)
	_, err = rulesCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	rulesAdapter := metaalarm.NewRuleAdapter(dbClient)
	return rulesAdapter, nil
}
