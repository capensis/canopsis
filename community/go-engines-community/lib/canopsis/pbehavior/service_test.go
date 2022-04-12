package pbehavior_test

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	mock_pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/pbehavior"
	mock_redis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/redis"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for suiteName, suiteData := range dataSetsForService() {
		for caseName, data := range suiteData.cases {
			mockProvider := newMockModelProvider(ctrl, suiteData)
			mockEntityMatcher := mock_pbehavior.NewMockComputedEntityMatcher(ctrl)
			mockStore := mock_pbehavior.NewMockStore(ctrl)
			mockLockClient := mock_redis.NewMockLockClient(ctrl)
			mockLock := mock_redis.NewMockLock(ctrl)
			mockLockClient.EXPECT().Obtain(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(mockLock, nil)
			mockLock.EXPECT().Release(gomock.Any()).Return(nil)
			mockStore.EXPECT().GetSpan(gomock.Any()).Return(timespan.Span{}, pbehavior.ErrNoComputed)
			mockStore.EXPECT().SetSpan(gomock.Any(), gomock.Any()).Return(nil)
			mockStore.EXPECT().SetComputed(gomock.Any(), gomock.Any()).Return(nil)
			mockEntityMatcher.EXPECT().LoadAll(gomock.Any(), gomock.Any()).Return(nil)
			mockEntityMatcher.
				EXPECT().
				Match(gomock.Any(), gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string) ([]string, error) {
					ids := make([]string, len(suiteData.pbehaviors))
					for i, v := range suiteData.pbehaviors {
						ids[i] = v.ID
					}

					return ids, nil
				}).MinTimes(0).MaxTimes(1)

			service := pbehavior.NewService(mockProvider, mockEntityMatcher, mockStore, mockLockClient)
			_, err := service.Compute(ctx, data.date)
			if err != nil {
				t.Errorf("%s %s: expected no error but got %v", suiteName, caseName, err)
				continue
			}

			r, err := service.Resolve(ctx, suiteData.entityID, data.t)
			if err != nil {
				t.Errorf("[Resolve] %s %s: expected no error but got %v", suiteName, caseName, err)
				continue
			}

			if data.expected == nil && r.ResolvedType != data.expected ||
				data.expected != nil && (r.ResolvedType == nil || *r.ResolvedType != *data.expected) {
				var expectedOutput interface{} = nil
				var resOutput interface{} = nil

				if data.expected != nil {
					expectedOutput = data.expected.ID
				}

				if r.ResolvedType != nil {
					resOutput = r.ResolvedType.ID
				}

				t.Errorf(
					"[Resolve] %s %s: expected output: %v, but got %v",
					suiteName,
					caseName,
					expectedOutput,
					resOutput,
				)
			}
		}
	}
}

func newMockModelProvider(ctrl *gomock.Controller, suiteData serviceSuiteDataSet) *mock_pbehavior.MockModelProvider {
	mock := mock_pbehavior.NewMockModelProvider(ctrl)
	typesByID := make(map[string]pbehavior.Type)
	pbehaviorsByID := make(map[string]pbehavior.PBehavior)
	reasonsByID := make(map[string]pbehavior.Reason)
	exceptionsByID := make(map[string]pbehavior.Exception)
	for k := range suiteData.types {
		typesByID[suiteData.types[k].ID] = suiteData.types[k]
	}
	for k := range suiteData.pbehaviors {
		pbehaviorsByID[suiteData.pbehaviors[k].ID] = suiteData.pbehaviors[k]
	}
	for k := range suiteData.reasons {
		reasonsByID[suiteData.reasons[k].ID] = suiteData.reasons[k]
	}
	for k := range suiteData.exceptions {
		exceptionsByID[suiteData.exceptions[k].ID] = suiteData.exceptions[k]
	}
	mock.
		EXPECT().
		GetTypes(gomock.Any()).
		Return(typesByID, nil)
	mock.
		EXPECT().
		GetEnabledPbehaviors(gomock.Any(), gomock.Any()).
		Return(pbehaviorsByID, nil)
	mock.
		EXPECT().
		GetReasons(gomock.Any()).
		Return(reasonsByID, nil)
	mock.
		EXPECT().
		GetExceptions(gomock.Any()).
		Return(exceptionsByID, nil)

	return mock
}

type serviceSuiteDataSet struct {
	pbehaviors []pbehavior.PBehavior
	types      []pbehavior.Type
	exceptions []pbehavior.Exception
	reasons    []pbehavior.Reason
	entityID   string
	cases      map[string]serviceCaseDataSet
}

type serviceCaseDataSet struct {
	date     timespan.Span
	t        time.Time
	expected *pbehavior.Type
}

func dataSetsForService() map[string]serviceSuiteDataSet {
	defaultActiveType := pbehavior.Type{
		ID:       "type1",
		Type:     pbehavior.TypeActive,
		Priority: 1,
	}
	defaultInactiveType := pbehavior.Type{
		ID:       "type2",
		Type:     pbehavior.TypeInactive,
		Priority: 2,
	}
	activeType := pbehavior.Type{
		ID:       "type3",
		Type:     pbehavior.TypeActive,
		Priority: 3,
	}
	maintenanceType := pbehavior.Type{
		ID:       "type4",
		Type:     pbehavior.TypeMaintenance,
		Priority: 4,
	}
	//it's not possible to have 2 types with the same priority, but we allow it in tests
	//in order to easily differentiate the results in according tests
	anotherMaintenanceType := pbehavior.Type{
		ID:       "type4-another",
		Type:     pbehavior.TypeMaintenance,
		Priority: 4,
	}
	mostPriorityActiveType := pbehavior.Type{
		ID:       "type5",
		Type:     pbehavior.TypeActive,
		Priority: 5,
	}
	mostPriorityMaintenanceType := pbehavior.Type{
		ID:       "type6",
		Type:     pbehavior.TypeMaintenance,
		Priority: 6,
	}
	pauseType := pbehavior.Type{
		ID:       "type7",
		Type:     pbehavior.TypeMaintenance,
		Priority: 7,
	}
	pbhTypes := []pbehavior.Type{
		defaultActiveType,
		defaultInactiveType,
		activeType,
		maintenanceType,
		anotherMaintenanceType,
		mostPriorityActiveType,
		mostPriorityMaintenanceType,
		pauseType,
	}
	entityID := "entity1"

	return map[string]serviceSuiteDataSet{
		"Given single maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  maintenanceType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default active type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 01:00"),
					expected: nil,
				},
				"and date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: &maintenanceType,
				},
				"and next date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: &maintenanceType,
				},
				"and next date during behavior and time after behavior Should return default active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 20:00"),
					expected: nil,
				},
				"and date before behavior Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:        genTime("31-05-2020 01:00"),
					expected: nil,
				},
				"and date after behavior Should return default active type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 01:00"),
					expected: nil,
				},
			},
		},
		"Given single pause behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  nil,
					Type:  pauseType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default active type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 01:00"),
					expected: nil,
				},
				"and date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: &pauseType,
				},
				"and date before behavior Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:        genTime("31-05-2020 01:00"),
					expected: nil,
				},
			},
		},
		"Given every day in for 7 times maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;COUNT=7",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  maintenanceType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 01:00"),
					expected: nil,
				},
				"and date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: &maintenanceType,
				},
				"and date during behavior and time after behavior Should return default active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 20:00"),
					expected: nil,
				},
				"and date before behavior Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:        genTime("31-05-2020 01:00"),
					expected: nil,
				},
				"and date after behavior Should return default active type": {
					date:     timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:        genTime("08-06-2020 20:00"),
					expected: nil,
				},
			},
		},
		"Given every 7 days in for 4 times maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=7;COUNT=4",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  maintenanceType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and date during activity of behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:        genTime("08-06-2020 11:00"),
					expected: &maintenanceType,
				},
				"and date during inactivity of behavior Should return default active type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 01:00"),
					expected: nil,
				},
				"and date before behavior Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:        genTime("31-05-2020 01:00"),
					expected: nil,
				},
				"and date after behavior Should return default active type": {
					date:     timespan.New(genTime("01-07-2020 00:00"), genTime("02-07-2020 00:00")),
					t:        genTime("01-07-2020 20:00"),
					expected: nil,
				},
			},
		},
		"Given single active behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  activeType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default inactive type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 01:00"),
					expected: &defaultInactiveType,
				},
				"and date during behavior and time during behavior Should return behavior active type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: &activeType,
				},
				"and next date during behavior and time during behavior Should return behavior active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: &activeType,
				},
				"and next date during behavior and time after behavior Should return default inactive type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 20:00"),
					expected: &defaultInactiveType,
				},
				"and date before behavior Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:        genTime("31-05-2020 01:00"),
					expected: nil,
				},
				"and date after behavior Should return default active type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 01:00"),
					expected: nil,
				},
				"and timespan from date before to date after and time before on first date Should return default inactive type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("01-06-2020 09:30"),
					expected: &defaultInactiveType,
				},
				"and timespan from date before to date after and time during behavior on first date Should return behavior active type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: &activeType,
				},
				"and timespan from date before to date after and time during behavior on second date Should return behavior active type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("02-06-2020 01:00"),
					expected: &activeType,
				},
				"and timespan from date before to date after and time after behavior on second date Should return default inactive type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("02-06-2020 14:00"),
					expected: &defaultInactiveType,
				},
				"and timespan from date before to date and time after behavior after Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("03-06-2020 01:00"),
					expected: nil,
				},
				"and 2 hours timespan on first date and time before behavior after Should return default inactive type": {
					date:     timespan.New(genTime("01-06-2020 05:00"), genTime("01-06-2020 06:00")),
					t:        genTime("01-06-2020 05:10"),
					expected: &defaultInactiveType,
				},
				"and 2 hours during behavior timespan Should return behavior active type": {
					date:     timespan.New(genTime("01-06-2020 11:00"), genTime("01-06-2020 13:00")),
					t:        genTime("01-06-2020 11:10"),
					expected: &activeType,
				},
				"and 2 hours timespan on after date after Should return default active type": {
					date:     timespan.New(genTime("03-06-2020 05:00"), genTime("03-06-2020 06:00")),
					t:        genTime("03-06-2020 05:10"),
					expected: nil,
				},
			},
		},
		"Given every day in for 7 times active behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;COUNT=7",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  activeType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default inactive type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 01:00"),
					expected: &defaultInactiveType,
				},
				"and date during behavior and time during behavior Should return behavior active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: &activeType,
				},
				"and date during behavior and time after behavior Should return default inactive type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 20:00"),
					expected: &defaultInactiveType,
				},
				"and date before behavior Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:        genTime("31-05-2020 01:00"),
					expected: nil,
				},
				"and date after behavior Should return default active type": {
					date:     timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:        genTime("08-06-2020 20:00"),
					expected: nil,
				},
			},
		},
		"Given every 7 days in for 4 times active behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=7;COUNT=4",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  activeType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and date during activity of behavior Should return behavior active type": {
					date:     timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:        genTime("08-06-2020 11:00"),
					expected: &activeType,
				},
				"and date during inactivity of behavior Should return default active type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 01:00"),
					expected: nil,
				},
				"and date before behavior Should return default active type": {
					date:     timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:        genTime("31-05-2020 01:00"),
					expected: nil,
				},
				"and date after behavior Should return default active type": {
					date:     timespan.New(genTime("01-07-2020 00:00"), genTime("02-07-2020 00:00")),
					t:        genTime("01-07-2020 20:00"),
					expected: nil,
				},
			},
		},
		"Given every 2 days active behavior and every day maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=2",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  mostPriorityActiveType.ID,
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &types.CpsTime{Time: genTime("01-06-2020 09:00")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh1 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: &mostPriorityActiveType,
				},
				"and time in pbh2 and default maintenance timestamp Should return pbh2 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 09:10"),
					expected: &maintenanceType,
				},
				"and time in default maintenance timestamp Should return default inactive type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 07:10"),
					expected: &defaultInactiveType,
				},
			},
		},
		"Given every 2 days active behavior with exdate and every day maintenance behavior with more priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=2",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  activeType.ID,
					Exdates: []pbehavior.Exdate{
						{
							Begin: types.CpsTime{Time: genTime("02-06-2020 00:00")},
							End:   types.CpsTime{Time: genTime("04-06-2020 00:00")},
							Type:  mostPriorityMaintenanceType.ID,
						},
					},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: &maintenanceType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: &mostPriorityMaintenanceType,
				},
				"and time in exdate and in pbh2 Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: &mostPriorityMaintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 10:10"),
					expected: nil,
				},
			},
		},
		"Given every 2 days active behavior with exception and every day maintenance behavior with more priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:         "pbh1",
					RRule:      "FREQ=DAILY;INTERVAL=2",
					Start:      &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:       &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:       activeType.ID,
					Exceptions: []string{"exception1"},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,
				},
			},
			types: pbhTypes,
			exceptions: []pbehavior.Exception{
				{
					ID: "exception1",
					Exdates: []pbehavior.Exdate{
						{
							Begin: types.CpsTime{Time: genTime("02-06-2020 00:00")},
							End:   types.CpsTime{Time: genTime("04-06-2020 00:00")},
							Type:  mostPriorityMaintenanceType.ID,
						},
					},
				},
			},
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: &maintenanceType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: &mostPriorityMaintenanceType,
				},
				"and time in exdate and in pbh2 Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: &mostPriorityMaintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 10:10"),
					expected: nil,
				},
			},
		},
		"Given every 2 days active behavior with exdate and every day maintenance behavior with less priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=2",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  mostPriorityActiveType.ID,
					Exdates: []pbehavior.Exdate{
						{
							Begin: types.CpsTime{Time: genTime("02-06-2020 00:00")},
							End:   types.CpsTime{Time: genTime("04-06-2020 00:00")},
							Type:  defaultInactiveType.ID,
						},
					},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh1 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: &mostPriorityActiveType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: &defaultInactiveType,
				},
				"and time in exdate and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: &maintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 10:10"),
					expected: nil,
				},
			},
		},
		"Given every 2 days active behavior with exception and every day maintenance behavior with less priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:         "pbh1",
					RRule:      "FREQ=DAILY;INTERVAL=2",
					Start:      &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:       &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:       mostPriorityActiveType.ID,
					Exceptions: []string{"exception1"},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,
				},
			},
			types: pbhTypes,
			exceptions: []pbehavior.Exception{
				{
					ID: "exception1",
					Exdates: []pbehavior.Exdate{
						{
							Begin: types.CpsTime{Time: genTime("02-06-2020 00:00")},
							End:   types.CpsTime{Time: genTime("04-06-2020 00:00")},
							Type:  defaultInactiveType.ID,
						},
					},
				},
			},
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh1 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: &mostPriorityActiveType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: &defaultInactiveType,
				},
				"and time in exdate and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: &maintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 10:10"),
					expected: nil,
				},
			},
		},
		"Given 2 intersected pbh with the same priorities": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:      "pbh1",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:    maintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 10:00")},
				},
				{
					ID:      "pbh2",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:    anotherMaintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 12:00")},
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"should return the newest pbh's type": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: &anotherMaintenanceType,
				},
			},
		},
		"Given 2 intersected pbh with the same priorities (inversed test)": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:      "pbh1",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:    maintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 12:00")},
				},
				{
					ID:      "pbh2",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:    anotherMaintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 10:00")},
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"should return the newest pbh's type": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: &maintenanceType,
				},
			},
		},
		"Given 2 intersected pbh with the same priorities and created date": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:      "pbh1",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:    maintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 10:00")},
				},
				{
					ID:      "pbh2",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:    anotherMaintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 10:00")},
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"should return pbehavior type by greatest pbehavior id": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: &anotherMaintenanceType,
				},
			},
		},
		"Given 2 intersected pbh with the same priorities and created date (inversed test)": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:      "pbh2",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:    maintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 10:00")},
				},
				{
					ID:      "pbh1",
					Start:   &types.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:    &types.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:    anotherMaintenanceType.ID,
					Created: types.CpsTime{Time: genTime("01-01-2020 10:00")},
				},
			},
			types:    pbhTypes,
			entityID: entityID,
			cases: map[string]serviceCaseDataSet{
				"should return pbehavior type by greatest pbehavior id": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: &maintenanceType,
				},
			},
		},
	}
}
