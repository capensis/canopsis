package pbehavior_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/pbehavior"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	mock_redis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/redis"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for suiteName, suiteData := range dataSetsForService() {
		t.Run(suiteName, func(t *testing.T) {
			for caseName, data := range suiteData.cases {
				t.Run(caseName, func(t *testing.T) {
					mockDbClient := mock_mongo.NewMockDbClient(ctrl)
					mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
					mockCursor := mock_mongo.NewMockCursor(ctrl)
					mockCursor.EXPECT().All(gomock.Any(), gomock.Any()).Return(nil).
						MinTimes(0).
						MaxTimes(len(suiteData.pbehaviors))
					mockDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).
						Return(mockCursor, nil).
						MinTimes(0).
						MaxTimes(len(suiteData.pbehaviors))
					mockDbClient.EXPECT().Collection(gomock.Eq(mongo.EntityMongoCollection)).Return(mockDbCollection)
					mockProvider := newMockModelProvider(ctrl, suiteData)
					mockDecoder := mock_encoding.NewMockDecoder(ctrl)
					typeComputer := pbehavior.NewTypeComputer(mockProvider, mockDecoder)
					mockStore := mock_pbehavior.NewMockStore(ctrl)
					mockLockClient := mock_redis.NewMockLockClient(ctrl)
					mockLock := mock_redis.NewMockLock(ctrl)
					mockLockClient.EXPECT().Obtain(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(mockLock, nil)
					mockLock.EXPECT().Release(gomock.Any()).Return(nil)
					mockStore.EXPECT().GetSpan(gomock.Any()).Return(timespan.Span{}, pbehavior.ErrNoComputed)
					mockStore.EXPECT().SetSpan(gomock.Any(), gomock.Any()).Return(nil)
					mockStore.EXPECT().SetComputed(gomock.Any(), gomock.Any()).Return(nil)
					service := pbehavior.NewService(mockDbClient, typeComputer, mockStore, mockLockClient, zerolog.Nop())
					resolver, _, err := service.Compute(ctx, data.date)
					if err != nil {
						t.Errorf("expected no error but got %v", err)
						return
					}

					r, err := resolver.Resolve(ctx, suiteData.entity, data.t)
					if err != nil {
						t.Errorf("expected no error but got %v", err)
						return
					}

					if r.Type != data.expected {
						t.Errorf("expected output: %q, but got %q", data.expected.ID, r.Type.ID)
					}
				})
			}
		})
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
	entity     types.Entity
	cases      map[string]serviceCaseDataSet
}

type serviceCaseDataSet struct {
	date     timespan.Span
	t        time.Time
	expected pbehavior.Type
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
	entity := types.Entity{ID: "entity1"}
	entityPattern := pattern.Entity{{
		{
			Field:     "_id",
			Condition: pattern.NewStringCondition(pattern.ConditionEqual, entity.ID),
		},
	}}

	return map[string]serviceSuiteDataSet{
		"Given single maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default active type": {
					date: timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:    genTime("01-06-2020 01:00"),
				},
				"and date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: maintenanceType,
				},
				"and next date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: maintenanceType,
				},
				"and next date during behavior and time after behavior Should return default active type": {
					date: timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:    genTime("02-06-2020 20:00"),
				},
				"and date before behavior Should return default active type": {
					date: timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:    genTime("31-05-2020 01:00"),
				},
				"and date after behavior Should return default active type": {
					date: timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:    genTime("03-06-2020 01:00"),
				},
			},
		},
		"Given single pause behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  nil,
					Type:  pauseType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default active type": {
					date: timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:    genTime("01-06-2020 01:00"),
				},
				"and date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: pauseType,
				},
				"and date before behavior Should return default active type": {
					date: timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:    genTime("31-05-2020 01:00"),
				},
			},
		},
		"Given every day in for 7 times maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;COUNT=7",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default active type": {
					date: timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:    genTime("02-06-2020 01:00"),
				},
				"and date during behavior and time during behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: maintenanceType,
				},
				"and date during behavior and time after behavior Should return default active type": {
					date: timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:    genTime("02-06-2020 20:00"),
				},
				"and date before behavior Should return default active type": {
					date: timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:    genTime("31-05-2020 01:00"),
				},
				"and date after behavior Should return default active type": {
					date: timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:    genTime("08-06-2020 20:00"),
				},
			},
		},
		"Given every 7 days in for 4 times maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=7;COUNT=4",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and date during activity of behavior Should return behavior maintenance type": {
					date:     timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:        genTime("08-06-2020 11:00"),
					expected: maintenanceType,
				},
				"and date during inactivity of behavior Should return default active type": {
					date: timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:    genTime("03-06-2020 01:00"),
				},
				"and date before behavior Should return default active type": {
					date: timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:    genTime("31-05-2020 01:00"),
				},
				"and date after behavior Should return default active type": {
					date: timespan.New(genTime("01-07-2020 00:00"), genTime("02-07-2020 00:00")),
					t:    genTime("01-07-2020 20:00"),
				},
			},
		},
		"Given single active behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  activeType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default inactive type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 01:00"),
					expected: defaultInactiveType,
				},
				"and date during behavior and time during behavior Should return behavior active type": {
					date:     timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: activeType,
				},
				"and next date during behavior and time during behavior Should return behavior active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: activeType,
				},
				"and next date during behavior and time after behavior Should return default inactive type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 20:00"),
					expected: defaultInactiveType,
				},
				"and date before behavior Should return default active type": {
					date: timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:    genTime("31-05-2020 01:00"),
				},
				"and date after behavior Should return default active type": {
					date: timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:    genTime("03-06-2020 01:00"),
				},
				"and timespan from date before to date after and time before on first date Should return default inactive type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("01-06-2020 09:30"),
					expected: defaultInactiveType,
				},
				"and timespan from date before to date after and time during behavior on first date Should return behavior active type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("01-06-2020 11:00"),
					expected: activeType,
				},
				"and timespan from date before to date after and time during behavior on second date Should return behavior active type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("02-06-2020 01:00"),
					expected: activeType,
				},
				"and timespan from date before to date after and time after behavior on second date Should return default inactive type": {
					date:     timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:        genTime("02-06-2020 14:00"),
					expected: defaultInactiveType,
				},
				"and timespan from date before to date and time after behavior after Should return default active type": {
					date: timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
					t:    genTime("03-06-2020 01:00"),
				},
				"and 2 hours timespan on first date and time before behavior after Should return default inactive type": {
					date:     timespan.New(genTime("01-06-2020 05:00"), genTime("01-06-2020 06:00")),
					t:        genTime("01-06-2020 05:10"),
					expected: defaultInactiveType,
				},
				"and 2 hours during behavior timespan Should return behavior active type": {
					date:     timespan.New(genTime("01-06-2020 11:00"), genTime("01-06-2020 13:00")),
					t:        genTime("01-06-2020 11:10"),
					expected: activeType,
				},
				"and 2 hours timespan on after date after Should return default active type": {
					date: timespan.New(genTime("03-06-2020 05:00"), genTime("03-06-2020 06:00")),
					t:    genTime("03-06-2020 05:10"),
				},
			},
		},
		"Given every day in for 7 times active behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;COUNT=7",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  activeType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and date during behavior and time before behavior Should return default inactive type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 01:00"),
					expected: defaultInactiveType,
				},
				"and date during behavior and time during behavior Should return behavior active type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 11:00"),
					expected: activeType,
				},
				"and date during behavior and time after behavior Should return default inactive type": {
					date:     timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:        genTime("02-06-2020 20:00"),
					expected: defaultInactiveType,
				},
				"and date before behavior Should return default active type": {
					date: timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:    genTime("31-05-2020 01:00"),
				},
				"and date after behavior Should return default active type": {
					date: timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:    genTime("08-06-2020 20:00"),
				},
			},
		},
		"Given every 7 days in for 4 times active behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=7;COUNT=4",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
					Type:  activeType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and date during activity of behavior Should return behavior active type": {
					date:     timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
					t:        genTime("08-06-2020 11:00"),
					expected: activeType,
				},
				"and date during inactivity of behavior Should return default active type": {
					date: timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:    genTime("03-06-2020 01:00"),
				},
				"and date before behavior Should return default active type": {
					date: timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
					t:    genTime("31-05-2020 01:00"),
				},
				"and date after behavior Should return default active type": {
					date: timespan.New(genTime("01-07-2020 00:00"), genTime("02-07-2020 00:00")),
					t:    genTime("01-07-2020 20:00"),
				},
			},
		},
		"Given every 2 days active behavior and every day maintenance behavior": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=2",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  mostPriorityActiveType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 09:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh1 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: mostPriorityActiveType,
				},
				"and time in pbh2 and default maintenance timestamp Should return pbh2 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 09:10"),
					expected: maintenanceType,
				},
				"and time in default maintenance timestamp Should return default inactive type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 07:10"),
					expected: defaultInactiveType,
				},
			},
		},
		"Given every 2 days active behavior with exdate and every day maintenance behavior with more priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=2",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  activeType.ID,
					Exdates: []pbehavior.Exdate{
						{
							Exdate: types.Exdate{
								Begin: datetime.CpsTime{Time: genTime("02-06-2020 00:00")},
								End:   datetime.CpsTime{Time: genTime("04-06-2020 00:00")},
							},
							Type: mostPriorityMaintenanceType.ID,
						},
					},

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: maintenanceType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: mostPriorityMaintenanceType,
				},
				"and time in exdate and in pbh2 Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: mostPriorityMaintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date: timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:    genTime("02-06-2020 10:10"),
				},
			},
		},
		"Given every 2 days active behavior with exception and every day maintenance behavior with more priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:         "pbh1",
					RRule:      "FREQ=DAILY;INTERVAL=2",
					Start:      &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:       &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:       activeType.ID,
					Exceptions: []string{"exception1"},

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types: pbhTypes,
			exceptions: []pbehavior.Exception{
				{
					ID: "exception1",
					Exdates: []pbehavior.Exdate{
						{
							Exdate: types.Exdate{
								Begin: datetime.CpsTime{Time: genTime("02-06-2020 00:00")},
								End:   datetime.CpsTime{Time: genTime("04-06-2020 00:00")},
							},
							Type: mostPriorityMaintenanceType.ID,
						},
					},
				},
			},
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: maintenanceType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: mostPriorityMaintenanceType,
				},
				"and time in exdate and in pbh2 Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: mostPriorityMaintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date: timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:    genTime("02-06-2020 10:10"),
				},
			},
		},
		"Given every 2 days active behavior with exdate and every day maintenance behavior with less priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					RRule: "FREQ=DAILY;INTERVAL=2",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  mostPriorityActiveType.ID,
					Exdates: []pbehavior.Exdate{
						{
							Exdate: types.Exdate{
								Begin: datetime.CpsTime{Time: genTime("02-06-2020 00:00")},
								End:   datetime.CpsTime{Time: genTime("04-06-2020 00:00")},
							},
							Type: defaultInactiveType.ID,
						},
					},

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh1 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: mostPriorityActiveType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: defaultInactiveType,
				},
				"and time in exdate and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: maintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date: timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:    genTime("02-06-2020 10:10"),
				},
			},
		},
		"Given every 2 days active behavior with exception and every day maintenance behavior with less priority": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:         "pbh1",
					RRule:      "FREQ=DAILY;INTERVAL=2",
					Start:      &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:       &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:       mostPriorityActiveType.ID,
					Exceptions: []string{"exception1"},

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					RRule: "FREQ=DAILY",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types: pbhTypes,
			exceptions: []pbehavior.Exception{
				{
					ID: "exception1",
					Exdates: []pbehavior.Exdate{
						{
							Exdate: types.Exdate{
								Begin: datetime.CpsTime{Time: genTime("02-06-2020 00:00")},
								End:   datetime.CpsTime{Time: genTime("04-06-2020 00:00")},
							},
							Type: defaultInactiveType.ID,
						},
					},
				},
			},
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"and time in pbh1 and in pbh2 Should return pbh1 type": {
					date:     timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
					t:        genTime("05-06-2020 10:58"),
					expected: mostPriorityActiveType,
				},
				"and time in exdate Should return exdate type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:10"),
					expected: defaultInactiveType,
				},
				"and time in exdate and in pbh2 Should return pbh2 type": {
					date:     timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
					t:        genTime("03-06-2020 10:58"),
					expected: maintenanceType,
				},
				"and time in exdate and no in pbh1 Should return default active type": {
					date: timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
					t:    genTime("02-06-2020 10:10"),
				},
			},
		},
		"Given 2 intersected pbh with the same priorities": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  anotherMaintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 12:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"should return the newest pbh's type": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: anotherMaintenanceType,
				},
			},
		},
		"Given 2 intersected pbh with the same priorities (inversed test)": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 12:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  anotherMaintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"should return the newest pbh's type": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: maintenanceType,
				},
			},
		},
		"Given 2 intersected pbh with the same priorities and created date": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh1",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh2",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  anotherMaintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"should return pbehavior type by greatest pbehavior id": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: anotherMaintenanceType,
				},
			},
		},
		"Given 2 intersected pbh with the same priorities and created date (inversed test)": {
			pbehaviors: []pbehavior.PBehavior{
				{
					ID:    "pbh2",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
					Type:  maintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
				{
					ID:    "pbh1",
					Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:55")},
					Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 11:00")},
					Type:  anotherMaintenanceType.ID,

					Created:             &datetime.CpsTime{Time: genTime("01-01-2020 10:00")},
					EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
				},
			},
			types:  pbhTypes,
			entity: entity,
			cases: map[string]serviceCaseDataSet{
				"should return pbehavior type by greatest pbehavior id": {
					date:     timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
					t:        genTime("01-06-2020 10:58"),
					expected: maintenanceType,
				},
			},
		},
	}
}
