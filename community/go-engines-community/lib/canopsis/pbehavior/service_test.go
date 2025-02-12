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
	customTZName := "Australia/Sydney"
	customTZ, err := time.LoadLocation(customTZName)
	if err != nil {
		t.Fatalf("failed to load timezone %v", err)
	}

	f := func(
		pbehaviors []pbehavior.PBehavior,
		exceptions []pbehavior.Exception,
		span timespan.Span,
		date time.Time,
		expected pbehavior.Type,
	) {
		t.Helper()

		mockDbClient := mock_mongo.NewMockDbClient(ctrl)
		mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
		mockCursor := mock_mongo.NewMockCursor(ctrl)
		mockCursor.EXPECT().All(gomock.Any(), gomock.Any()).Return(nil).
			MinTimes(0).
			MaxTimes(len(pbehaviors))
		mockDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).
			Return(mockCursor, nil).
			MinTimes(0).
			MaxTimes(len(pbehaviors))
		mockDbClient.EXPECT().Collection(gomock.Eq(mongo.EntityMongoCollection)).Return(mockDbCollection)
		mockProvider := newMockModelProvider(ctrl, pbehaviors, pbhTypes, nil, exceptions)
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
		resolver, _, err := service.Compute(ctx, span, time.UTC)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
			return
		}

		r, err := resolver.Resolve(ctx, entity, date)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
			return
		}

		if r.Type != expected {
			t.Errorf("expected output: %q, but got %q", expected.ID, r.Type.ID)
		}
	}

	// single maintenance behavior
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:    "pbh1",
				Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
				Type:  maintenanceType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during behavior and time before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
			genTime("01-06-2020 01:00"),
			pbehavior.Type{},
		)
		// span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
			genTime("01-06-2020 11:00"),
			maintenanceType,
		)
		// next span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 11:00"),
			maintenanceType,
		)
		// next span during behavior and time after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 20:00"),
			pbehavior.Type{},
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
			genTime("31-05-2020 01:00"),
			pbehavior.Type{},
		)
		// span after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 01:00"),
			pbehavior.Type{},
		)
	}

	// single pause behavior
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:    "pbh1",
				Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				Stop:  nil,
				Type:  pauseType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during behavior and time before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
			genTime("01-06-2020 01:00"),
			pbehavior.Type{},
		)
		// span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
			genTime("01-06-2020 11:00"),
			pauseType,
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
			genTime("31-05-2020 01:00"),
			pbehavior.Type{},
		)
	}

	// every day in for 7 times maintenance behavior
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:    "pbh1",
				RRule: "FREQ=DAILY;COUNT=7",
				Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
				Type:  maintenanceType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during behavior and time before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 01:00"),
			pbehavior.Type{},
		)
		// span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 11:00"),
			maintenanceType,
		)
		// span during behavior and time after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 20:00"),
			pbehavior.Type{},
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
			genTime("31-05-2020 01:00"),
			pbehavior.Type{},
		)
		// span after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
			genTime("08-06-2020 20:00"),
			pbehavior.Type{},
		)
	}

	// every 7 days in for 4 times maintenance behavior
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:    "pbh1",
				RRule: "FREQ=DAILY;INTERVAL=7;COUNT=4",
				Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
				Type:  maintenanceType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during activity of behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
			genTime("08-06-2020 11:00"),
			maintenanceType,
		)
		// span during inactivity of behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 01:00"),
			pbehavior.Type{},
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
			genTime("31-05-2020 01:00"),
			pbehavior.Type{},
		)
		// span after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-07-2020 00:00"), genTime("02-07-2020 00:00")),
			genTime("01-07-2020 20:00"),
			pbehavior.Type{},
		)
	}

	// single active behavior
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:    "pbh1",
				Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
				Type:  activeType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during behavior and time before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
			genTime("01-06-2020 01:00"),
			defaultInactiveType,
		)
		// span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00"), genTime("02-06-2020 00:00")),
			genTime("01-06-2020 11:00"),
			activeType,
		)
		// next span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 11:00"),
			activeType,
		)
		// next span during behavior and time after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 20:00"),
			defaultInactiveType,
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
			genTime("31-05-2020 01:00"),
			pbehavior.Type{},
		)
		// span after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 01:00"),
			pbehavior.Type{},
		)
		// timespan from date before to date after and time before on first date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
			genTime("01-06-2020 09:30"),
			defaultInactiveType,
		)
		// timespan from date before to date after and time during behavior on first date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
			genTime("01-06-2020 11:00"),
			activeType,
		)
		// timespan from date before to date after and time during behavior on second date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
			genTime("02-06-2020 01:00"),
			activeType,
		)
		// timespan from date before to date after and time after behavior on second date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
			genTime("02-06-2020 14:00"),
			defaultInactiveType,
		)
		// timespan from date before to date and time after behavior after
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00"), genTime("03-06-2020 09:00")),
			genTime("03-06-2020 01:00"),
			pbehavior.Type{},
		)
		// 2 hours timespan on first date and time before behavior after
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 05:00"), genTime("01-06-2020 06:00")),
			genTime("01-06-2020 05:10"),
			defaultInactiveType,
		)
		// 2 hours during behavior timespan
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 11:00"), genTime("01-06-2020 13:00")),
			genTime("01-06-2020 11:10"),
			activeType,
		)
		// 2 hours timespan on after date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 05:00"), genTime("03-06-2020 06:00")),
			genTime("03-06-2020 05:10"),
			pbehavior.Type{},
		)
	}

	// single active behavior in custom timezone
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:       "pbh1",
				Timezone: customTZName,
				Start:    &datetime.CpsTime{Time: genTime("01-06-2020 10:00", customTZ).UTC()},
				Stop:     &datetime.CpsTime{Time: genTime("02-06-2020 12:00", customTZ).UTC()},
				Type:     activeType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00", customTZ).UTC()},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during behavior and time before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00", customTZ).UTC(), genTime("02-06-2020 00:00", customTZ).UTC()),
			genTime("01-06-2020 01:00", customTZ).UTC(),
			defaultInactiveType,
		)
		// span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 00:00", customTZ).UTC(), genTime("02-06-2020 00:00", customTZ).UTC()),
			genTime("01-06-2020 11:00", customTZ).UTC(),
			activeType,
		)
		// next span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00", customTZ).UTC(), genTime("03-06-2020 00:00", customTZ).UTC()),
			genTime("02-06-2020 11:00", customTZ).UTC(),
			activeType,
		)
		// next span during behavior and time after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00", customTZ).UTC(), genTime("03-06-2020 00:00", customTZ).UTC()),
			genTime("02-06-2020 20:00", customTZ).UTC(),
			defaultInactiveType,
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00", customTZ).UTC(), genTime("01-06-2020 00:00", customTZ).UTC()),
			genTime("31-05-2020 01:00", customTZ).UTC(),
			pbehavior.Type{},
		)
		// span after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00", customTZ).UTC(), genTime("04-06-2020 00:00", customTZ).UTC()),
			genTime("03-06-2020 01:00", customTZ).UTC(),
			pbehavior.Type{},
		)
		// timespan from date before to date after and time before on first date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00", customTZ).UTC(), genTime("03-06-2020 09:00", customTZ).UTC()),
			genTime("01-06-2020 09:30", customTZ).UTC(),
			defaultInactiveType,
		)
		// timespan from date before to date after and time during behavior on first date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00", customTZ).UTC(), genTime("03-06-2020 09:00", customTZ).UTC()),
			genTime("01-06-2020 11:00", customTZ).UTC(),
			activeType,
		)
		// timespan from date before to date after and time during behavior on second date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00", customTZ).UTC(), genTime("03-06-2020 09:00", customTZ).UTC()),
			genTime("02-06-2020 01:00", customTZ).UTC(),
			activeType,
		)
		// timespan from date before to date after and time after behavior on second date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00", customTZ).UTC(), genTime("03-06-2020 09:00", customTZ).UTC()),
			genTime("02-06-2020 14:00", customTZ).UTC(),
			defaultInactiveType,
		)
		// timespan from date before to date and time after behavior after
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 09:00", customTZ).UTC(), genTime("03-06-2020 09:00", customTZ).UTC()),
			genTime("03-06-2020 01:00", customTZ).UTC(),
			pbehavior.Type{},
		)
		// 2 hours timespan on first date and time before behavior after
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 05:00", customTZ).UTC(), genTime("01-06-2020 06:00", customTZ).UTC()),
			genTime("01-06-2020 05:10", customTZ).UTC(),
			defaultInactiveType,
		)
		// 2 hours during behavior timespan
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 11:00", customTZ).UTC(), genTime("01-06-2020 13:00", customTZ).UTC()),
			genTime("01-06-2020 11:10", customTZ).UTC(),
			activeType,
		)
		// 2 hours timespan on after date
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 05:00", customTZ).UTC(), genTime("03-06-2020 06:00", customTZ).UTC()),
			genTime("03-06-2020 05:10", customTZ).UTC(),
			pbehavior.Type{},
		)
	}

	// every day in for 7 times active behavior
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:    "pbh1",
				RRule: "FREQ=DAILY;COUNT=7",
				Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				Stop:  &datetime.CpsTime{Time: genTime("01-06-2020 12:00")},
				Type:  activeType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during behavior and time before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 01:00"),
			defaultInactiveType,
		)
		// span during behavior and time during behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 11:00"),
			activeType,
		)
		// span during behavior and time after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 20:00"),
			defaultInactiveType,
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
			genTime("31-05-2020 01:00"),
			pbehavior.Type{},
		)
		// span after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
			genTime("08-06-2020 20:00"),
			pbehavior.Type{},
		)
	}

	// every 7 days in for 4 times active behavior
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:    "pbh1",
				RRule: "FREQ=DAILY;INTERVAL=7;COUNT=4",
				Start: &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				Stop:  &datetime.CpsTime{Time: genTime("02-06-2020 12:00")},
				Type:  activeType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-06-2020 10:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during activity of behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("08-06-2020 00:00"), genTime("09-06-2020 00:00")),
			genTime("08-06-2020 11:00"),
			activeType,
		)
		// span during inactivity of behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 01:00"),
			pbehavior.Type{},
		)
		// span before behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("31-05-2020 00:00"), genTime("01-06-2020 00:00")),
			genTime("31-05-2020 01:00"),
			pbehavior.Type{},
		)
		// span after behavior
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-07-2020 00:00"), genTime("02-07-2020 00:00")),
			genTime("01-07-2020 20:00"),
			pbehavior.Type{},
		)
	}

	// every 2 days active behavior and every day maintenance behavior
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		// time in pbh1 and in pbh2
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
			genTime("05-06-2020 10:58"),
			mostPriorityActiveType,
		)
		// time in pbh2 and default maintenance timestamp
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
			genTime("05-06-2020 09:10"),
			maintenanceType,
		)
		// time in default maintenance timestamp
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
			genTime("05-06-2020 07:10"),
			defaultInactiveType,
		)
	}

	// every 2 days active behavior with exdate and every day maintenance behavior with more priority
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		// time in pbh1 and in pbh2
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
			genTime("05-06-2020 10:58"),
			maintenanceType,
		)
		// time in exdate
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:10"),
			mostPriorityMaintenanceType,
		)
		// time in exdate and in pbh2
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:58"),
			mostPriorityMaintenanceType,
		)
		// time in exdate and no in pbh1
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 10:10"),
			pbehavior.Type{},
		)
	}

	// every 2 days active behavior with exception and every day maintenance behavior with more priority
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		exceptions := []pbehavior.Exception{
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
		}
		// time in pbh1 and in pbh2
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
			genTime("05-06-2020 10:58"),
			maintenanceType,
		)
		// time in exdate
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:10"),
			mostPriorityMaintenanceType,
		)
		// time in exdate and in pbh2
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:58"),
			mostPriorityMaintenanceType,
		)
		// time in exdate and no in pbh1
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 10:10"),
			pbehavior.Type{},
		)
	}

	// every 2 days active behavior with exdate and every day maintenance behavior with less priority
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		// time in pbh1 and in pbh2
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
			genTime("05-06-2020 10:58"),
			mostPriorityActiveType,
		)
		// time in exdate
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:10"),
			defaultInactiveType,
		)
		// time in exdate and in pbh2
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:58"),
			maintenanceType,
		)
		// time in exdate and no in pbh1
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 10:10"),
			pbehavior.Type{},
		)
	}

	// every 2 days active behavior with exception and every day maintenance behavior with less priority
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		exceptions := []pbehavior.Exception{
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
		}
		// time in pbh1 and in pbh2
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("05-06-2020 00:00"), genTime("06-06-2020 00:00")),
			genTime("05-06-2020 10:58"),
			mostPriorityActiveType,
		)
		// time in exdate
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:10"),
			defaultInactiveType,
		)
		// time in exdate and in pbh2
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("03-06-2020 00:00"), genTime("04-06-2020 00:00")),
			genTime("03-06-2020 10:58"),
			maintenanceType,
		)
		// time in exdate and no in pbh1
		f(
			pbehaviors,
			exceptions,
			timespan.New(genTime("02-06-2020 00:00"), genTime("03-06-2020 00:00")),
			genTime("02-06-2020 10:10"),
			pbehavior.Type{},
		)
	}

	// 2 intersected pbh with the same priorities
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		// should return the newest pbh's type
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
			genTime("01-06-2020 10:58"),
			anotherMaintenanceType,
		)
	}

	// 2 intersected pbh with the same priorities (inversed test)
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		// should return the newest pbh's type
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
			genTime("01-06-2020 10:58"),
			maintenanceType,
		)
	}

	// 2 intersected pbh with the same priorities and created date
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		// should return pbehavior type by greatest pbehavior id
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
			genTime("01-06-2020 10:58"),
			anotherMaintenanceType,
		)
	}

	// 2 intersected pbh with the same priorities and created date (inversed test)
	{
		pbehaviors := []pbehavior.PBehavior{
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
		}
		// should return pbehavior type by greatest pbehavior id
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("01-06-2020 09:00"), genTime("01-06-2020 12:00")),
			genTime("01-06-2020 10:58"),
			maintenanceType,
		)
	}

	// maintenance behavior during summer/winter time change
	{
		pbehaviors := []pbehavior.PBehavior{
			{
				ID:       "pbh1",
				RRule:    "FREQ=DAILY",
				Start:    &datetime.CpsTime{Time: genTime("01-01-2020 10:00", customTZ).UTC()},
				Stop:     &datetime.CpsTime{Time: genTime("01-01-2020 11:00", customTZ).UTC()},
				Timezone: customTZName,
				Type:     maintenanceType.ID,

				Created:             &datetime.CpsTime{Time: genTime("01-01-2020 00:00")},
				EntityPatternFields: savedpattern.EntityPatternFields{EntityPattern: entityPattern},
			},
		}
		// span during summer change
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("26-03-2020 00:00", customTZ).UTC(), genTime("02-04-2020 00:00", customTZ).UTC()),
			genTime("28-03-2020 10:30", customTZ).UTC(),
			maintenanceType,
		)
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("26-03-2020 00:00", customTZ).UTC(), genTime("02-04-2020 00:00", customTZ).UTC()),
			genTime("29-03-2020 10:30", customTZ).UTC(),
			maintenanceType,
		)
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("26-03-2020 00:00", customTZ).UTC(), genTime("02-04-2020 00:00", customTZ).UTC()),
			genTime("30-03-2020 10:30", customTZ).UTC(),
			maintenanceType,
		)
		// span during winter change
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("22-10-2020 00:00", customTZ).UTC(), genTime("29-10-2020 00:00", customTZ).UTC()),
			genTime("24-10-2020 10:30", customTZ).UTC(),
			maintenanceType,
		)
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("22-10-2020 00:00", customTZ).UTC(), genTime("29-10-2020 00:00", customTZ).UTC()),
			genTime("25-10-2020 10:30", customTZ).UTC(),
			maintenanceType,
		)
		f(
			pbehaviors,
			nil,
			timespan.New(genTime("22-10-2020 00:00", customTZ).UTC(), genTime("29-10-2020 00:00", customTZ).UTC()),
			genTime("26-10-2020 10:30", customTZ).UTC(),
			maintenanceType,
		)
	}
}

func newMockModelProvider(
	ctrl *gomock.Controller,
	pbehaviors []pbehavior.PBehavior,
	types []pbehavior.Type,
	reasons []pbehavior.Reason,
	exceptions []pbehavior.Exception,
) *mock_pbehavior.MockModelProvider {
	mock := mock_pbehavior.NewMockModelProvider(ctrl)
	typesByID := make(map[string]pbehavior.Type)
	pbehaviorsByID := make(map[string]pbehavior.PBehavior)
	reasonsByID := make(map[string]pbehavior.Reason)
	exceptionsByID := make(map[string]pbehavior.Exception)
	for k := range types {
		typesByID[types[k].ID] = types[k]
	}
	for k := range pbehaviors {
		pbehaviorsByID[pbehaviors[k].ID] = pbehaviors[k]
	}
	for k := range reasons {
		reasonsByID[reasons[k].ID] = reasons[k]
	}
	for k := range exceptions {
		exceptionsByID[exceptions[k].ID] = exceptions[k]
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
