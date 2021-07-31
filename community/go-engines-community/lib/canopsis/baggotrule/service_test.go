package baggotrule_test

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/baggotrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_baggotrule "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/baggotrule"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func newClosedAlarm(time types.CpsTime) types.AlarmWithEntity {
	return types.AlarmWithEntity{
		Alarm: types.Alarm{
			Value: types.AlarmValue{
				State: &types.AlarmStep{
					Value: types.AlarmStateOK,
				},
				Steps: []types.AlarmStep{
					{
						Timestamp: time,
					},
				},
				Resource: "test_resource",
			},
		},
	}
}

func baggotRule() []baggotrule.Rule {
	return []baggotrule.Rule{
		{
			ID: "baggot_rule",
			Duration: types.DurationWithUnit{
				Seconds: 10,
				Unit:    "s",
			},
			AlarmPatterns: pattern.AlarmPatternList{
				Patterns: []pattern.AlarmPattern{
					{
						AlarmFields: pattern.AlarmFields{
							Value: pattern.AlarmValuePattern{
								AlarmValueFields: pattern.AlarmValueFields{
									Resource: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "test_resource",
											},
										},
									},
								},
							},
						},
					},
				},
				Set:   true,
				Valid: true,
			},
			EntityPatterns: pattern.EntityPatternList{},
		},
	}
}

func TestService_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lastStepTime := 61 * time.Second

	var dataSets = []struct {
		testName       string
		findAlarms     []types.AlarmWithEntity
		findError      error
		expectedClosed int
	}{
		{
			"given no alarms should return empty result",
			[]types.AlarmWithEntity{},
			nil,
			0,
		},
		{
			"given done alarms with last step time < baggotTime should return empty result",
			[]types.AlarmWithEntity{
				newClosedAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now(),
				}),
			},
			nil,
			0,
		},
		{
			"given done alarms and done alarms with last steps time > baggotTime should return count of alarms with time > baggotTime",
			[]types.AlarmWithEntity{
				newClosedAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
			},
			nil,
			2,
		},
		{
			"given done alarms with valid time should return count of alarms",
			[]types.AlarmWithEntity{
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
			},
			nil,
			3,
		},
		{
			"given find error should return error",
			[]types.AlarmWithEntity{
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
				newClosedAlarm(types.CpsTime{
					Time: time.Now().Add(-lastStepTime),
				}),
			},
			fmt.Errorf("not found"),
			0,
		},
	}

	baggotRuleAdapterMock := mock_baggotrule.NewMockAdapter(ctrl)
	baggotRuleAdapterMock.
		EXPECT().
		Get(context.Background()).
		Return(baggotRule(), nil).AnyTimes()

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
			alarmAdapterMock.
				EXPECT().
				GetOpenedAlarmsWithEntity(context.Background()).
				Return(createMockCursor(ctrl, dataset.findAlarms, dataset.findError), dataset.findError)

			service := baggotrule.NewService(
				baggotRuleAdapterMock,
				alarmAdapterMock,
				zerolog.Logger{},
			)

			closedAlarms, err := service.Process(context.Background())
			if err != nil {
				expectedErr := fmt.Sprintf("cannot fetch open alarms: %v", dataset.findError.Error())
				if err.Error() != expectedErr {
					t.Errorf("expected err %v but got %v", expectedErr, err)
				}
			}

			if len(closedAlarms) != dataset.expectedClosed {
				t.Errorf("%s expected %d done alarms but got %d", dataset.testName, dataset.expectedClosed, len(closedAlarms))
			}
		})
	}
}

/**
---------------------------
Benchmarking
*/
func generateListOfClosedAlarm(n int, time types.CpsTime) []types.AlarmWithEntity {
	var alarms = make([]types.AlarmWithEntity, n)
	for i := 0; i < n; i++ {
		alarms[i] = types.AlarmWithEntity{
			Alarm: types.Alarm{
				ID:       uuid.New().String(),
				EntityID: fmt.Sprintf("entity_%d", i),
				Value: types.AlarmValue{
					State: &types.AlarmStep{
						Value: types.AlarmStateOK,
					},
					Steps: []types.AlarmStep{
						{
							Timestamp: time,
						},
					},
					Resource:       fmt.Sprintf("resource_%d", i),
					Component:      fmt.Sprintf("component_%d", i),
					ConnectorName:  fmt.Sprintf("component_name_%d", i),
					Connector:      fmt.Sprintf("connector_%d", i),
					CreationDate:   time,
					LastUpdateDate: time,
					LastEventDate:  time,
					DisplayName:    fmt.Sprintf("display_name_%d", i),
					Output:         fmt.Sprintf("output_%d", i),
				},
			},
			Entity: types.Entity{
				ID:      fmt.Sprintf("entity_%d", i),
				Name:    fmt.Sprintf("entity_name_%d", i),
				Enabled: true,
				Infos: map[string]types.Info{
					"name": {
						Name:        fmt.Sprintf("test_info_%d", i),
						Description: "this is test value name",
						Value:       fmt.Sprintf("value_%d", i),
					},
				},
			},
		}

	}
	return alarms
}

func generateListOfBaggotRules(n int, updated types.CpsTime) []baggotrule.Rule {
	var rules = make([]baggotrule.Rule, n)
	var defaultRule = baggotrule.Rule{
		ID: "default_rule",
		Duration: types.DurationWithUnit{
			Seconds: 60,
			Unit:    "s",
		},
		AlarmPatterns:  pattern.AlarmPatternList{},
		EntityPatterns: pattern.EntityPatternList{},
		Updated:        &updated,
	}

	for i := 0; i < n-1; i++ {
		rand.Seed(time.Now().UnixNano())
		rules[i] = baggotrule.Rule{
			ID: fmt.Sprintf("baggot_ruule_%d", i),
			Duration: types.DurationWithUnit{
				Seconds: rand.Int63n(300),
				Unit:    "s",
			},
			AlarmPatterns: pattern.AlarmPatternList{
				Patterns: []pattern.AlarmPattern{
					{
						AlarmFields: pattern.AlarmFields{
							Value: pattern.AlarmValuePattern{
								AlarmValueFields: pattern.AlarmValueFields{
									Resource: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: fmt.Sprintf("resource_%d", i),
											},
										},
									},
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: fmt.Sprintf("component_%d", i),
											},
										},
									},
									Connector: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: fmt.Sprintf("connector_%d", i),
											},
										},
									},
									Output: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: fmt.Sprintf("output_%d", i),
											},
										},
									},
								},
							},
						},
					},
				},
				Set:   true,
				Valid: true,
			},
			EntityPatterns: pattern.EntityPatternList{
				Patterns: []pattern.EntityPattern{
					{
						EntityFields: pattern.EntityFields{
							ID: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: types.OptionalString{
										Set:   true,
										Value: fmt.Sprintf("entity_%d", i),
									},
								},
							},
							Infos: map[string]pattern.InfoPattern{
								"name": pattern.InfoPattern{
									InfoFields: pattern.InfoFields{
										Value: pattern.StringPattern{
											StringConditions: pattern.StringConditions{
												RegexMatch: types.OptionalRegexp{
													Set: true,
													Value: utils.WrapperBuiltInRegex{
														Regexp: regexp.MustCompile("test_[0-9]{1,4}"),
													},
												},
											}},
									},
								},
							},
						},
					},
				},
				Set:   true,
				Valid: true,
			},
		}
	}
	rules[n-1] = defaultRule
	return rules
}

func BenchmarkService_Process(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	alrsNum, ruleNum := 100000, 5000

	pivotTime := types.CpsTime{
		Time: time.Now().Add(-360 * time.Second),
	}
	alrs := createMockCursor(ctrl, generateListOfClosedAlarm(alrsNum, pivotTime), nil)
	rules := generateListOfBaggotRules(ruleNum, pivotTime)

	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)
	alarmAdapterMock.
		EXPECT().
		GetOpenedAlarmsWithEntity(context.Background()).
		Return(alrs, nil)
	baggotRuleAdapterMock := mock_baggotrule.NewMockAdapter(ctrl)
	baggotRuleAdapterMock.
		EXPECT().
		Get(context.Background()).
		Return(rules, nil)

	service := baggotrule.NewService(
		baggotRuleAdapterMock,
		alarmAdapterMock,
		zerolog.Logger{},
	)
	b.ResetTimer()

	closedAlarms, err := service.Process(context.Background())
	if err != nil {
		b.Error(err)
	}

	if len(closedAlarms) != alrsNum {
		b.Errorf("expected %d done alarms but got %d", alrsNum, len(closedAlarms))
	}
}

func createMockCursor(ctrl *gomock.Controller, models []types.AlarmWithEntity, err error) libmongo.Cursor {
	if err != nil {
		return nil
	}
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	if len(models) > 0 {
		//calls := make([]*gomock.Call, 2*len(models))
		for i := range models {
			model := models[i]
			mockCursor.EXPECT().Next(gomock.Any()).Return(true)
			mockCursor.EXPECT().Decode(gomock.Any()).Do(func(m *types.AlarmWithEntity) {
				*m = model
			}).Return(nil)
		}
	}
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	return mockCursor
}
