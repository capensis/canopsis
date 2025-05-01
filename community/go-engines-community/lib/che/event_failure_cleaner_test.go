package che_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/usernotification"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	"github.com/kylelemons/godebug/pretty"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/mock/gomock"
)

func TestEventFailureCleaner_Clean(t *testing.T) {
	ctx := t.Context()
	client, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		t.Fatalf("cannot connect to mongodb: %v", err)
	}

	now := datetime.NewCpsTime()
	failColl := client.Collection(mongo.EventFilterFailureCollection)
	ruleColl := client.Collection(mongo.EventFilterRuleCollection)
	notifColl := client.Collection(mongo.UserNotificationCollection)
	decoder := json.NewDecoder()
	f := func(
		conf datastorage.Config,
		limit int,
		fails, rules, notifs []any,
		expectedFailIDs []string,
		expectedRules []eventfilter.Rule,
		expectedDeleted int64,
		expectedNotifIDs []string,
		expectedEvent *rpc.ApiNotificationEvent,
	) {
		t.Helper()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		err = cleanCollections(ctx, failColl, ruleColl, notifColl, fails, rules, notifs)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}

		amqpChannel := mock_amqp.NewMockPublisher(ctrl)
		if expectedEvent != nil {
			amqpChannel.EXPECT().
				PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(_ context.Context, _, _ string, _, _ bool, msg amqp.Publishing) {
					received := rpc.ApiNotificationEvent{}
					err := decoder.Decode(msg.Body, &received)
					if err != nil {
						t.Fatalf("cannot decode event: %v", err)
					}

					slices.Sort(received.Roles)
					if diff := pretty.Compare(expectedEvent, received); diff != "" {
						t.Fatalf("unexpected event: (-want +got):\n%s", diff)
					}
				})
		}
		c := che.NewEventFailureCleaner(amqpChannel, json.NewEncoder(),
			canopsis.ApiNotificationExchangeName, "", canopsis.JsonContentType, zerolog.Nop())
		res, err := c.Clean(ctx, client, conf, now, limit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		failDs, err := findIDs(ctx, failColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedFailIDs, failDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", failColl.Name(), diff)
		}

		updatedRules, err := findRules(ctx, ruleColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedRules, updatedRules); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", ruleColl.Name(), diff)
		}

		if res.Deleted != expectedDeleted {
			t.Fatalf("unexpected deleted count, expected %d, got %d", expectedDeleted, res.Deleted)
		}

		notifIDs, err := findIDs(ctx, notifColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedNotifIDs, notifIDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", notifColl.Name(), diff)
		}

		err = cleanCollections(ctx, failColl, ruleColl, notifColl, nil, nil, nil)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}
	}

	weekAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -7)}
	f(
		newDSConfig("7d"),
		0,
		[]any{
			newTestFail("f1", "r1", now, true),
			newTestFail("f2", "r1", now, true),
			newTestFail("f3", "r1", now, false),
			newTestFail("f4", "r1", now, false),
			newTestFail("f5", "r1", weekAgo, true),
			newTestFail("f6", "r1", weekAgo, true),
			newTestFail("f7", "r1", weekAgo, false),
			newTestFail("f8", "r1", weekAgo, false),
			newTestFail("f9", "r2", now, true),
			newTestFail("f10", "r3", weekAgo, true),
			newTestFail("f11", "r4", weekAgo, true),
		},
		[]any{
			newTestRule("r1", 8, 4),
			newTestRule("r2", 1, 1),
			newTestRule("r3", 1, 1),
			newTestRule("r4", 1, 1),
			newTestRule("r5", 0, 0),
		},
		[]any{
			newTestNotif("n1", "r1", weekAgo, []string{"ro1"}),
			newTestNotif("n2", "r2", now, []string{"ro2"}),
			newTestNotif("n3", "r3", weekAgo, []string{"ro3", "ro4"}),
			newTestNotif("n4", "r4", weekAgo, []string{"ro4", "ro5"}),
		},
		[]string{"f1", "f2", "f3", "f4", "f9"},
		[]eventfilter.Rule{
			newTestRule("r1", 4, 2),
			newTestRule("r2", 1, 1),
			newTestRule("r3", 0, 0),
			newTestRule("r4", 0, 0),
			newTestRule("r5", 0, 0),
		},
		6,
		[]string{"n1", "n2"},
		&rpc.ApiNotificationEvent{Roles: []string{"ro3", "ro4", "ro5"}},
	)
	f(
		newDSConfig("7d"),
		3,
		[]any{
			newTestFail("f1", "r1", weekAgo, false),
			newTestFail("f2", "r1", weekAgo, false),
			newTestFail("f3", "r1", weekAgo, false),
			newTestFail("f4", "r1", weekAgo, false),
		},
		[]any{
			newTestRule("r1", 4, 0),
		},
		[]any{},
		[]string{"f4"},
		[]eventfilter.Rule{
			newTestRule("r1", 1, 0),
		},
		3,
		[]string{},
		nil,
	)
}

func newDSConfig(deletedAfterStr string) datastorage.Config {
	c := datastorage.Config{}
	enabled := true
	d, err := datetime.ParseDurationWithUnit(deletedAfterStr)
	if err != nil {
		panic(err)
	}

	c.EventFilterFailure.DeleteAfter = &datetime.DurationWithEnabled{
		DurationWithUnit: d,
		Enabled:          &enabled,
	}

	return c
}

func newTestRule(id string, f, u int64) eventfilter.Rule {
	return eventfilter.Rule{
		ID:                  id,
		FailuresCount:       f,
		UnreadFailuresCount: u,
	}
}

func newTestFail(id, r string, t datetime.CpsTime, u bool) eventfilter.Failure {
	return eventfilter.Failure{
		ID:        id,
		Rule:      r,
		Timestamp: t,
		Unread:    u,
	}
}

func newTestNotif(id, r string, t datetime.CpsTime, roles []string) usernotification.Notification {
	return usernotification.Notification{
		ID:      id,
		Type:    usernotification.TypeEventFilterFailure,
		Roles:   roles,
		Time:    t,
		Comment: "test",
		Rule: &usernotification.Rule{
			ID:   r,
			Name: "test",
		},
	}
}

func cleanCollections(ctx context.Context, failColl, ruleColl, notifColl mongo.DbCollection, fails, rules, notifs []any) error {
	_, err := failColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", failColl.Name(), err)
	}

	_, err = ruleColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", ruleColl.Name(), err)
	}

	_, err = notifColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", notifColl.Name(), err)
	}

	if len(fails) > 0 {
		_, err = failColl.InsertMany(ctx, fails)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", failColl.Name(), err)
		}
	}

	if len(rules) > 0 {
		_, err = ruleColl.InsertMany(ctx, rules)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", ruleColl.Name(), err)
		}
	}

	if len(notifs) > 0 {
		_, err = notifColl.InsertMany(ctx, notifs)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", notifColl.Name(), err)
		}
	}

	return nil
}

func findIDs(ctx context.Context, coll mongo.DbCollection) ([]string, error) {
	c, err := coll.Aggregate(ctx, []bson.M{
		{"$group": bson.M{
			"_id": nil,
			"ids": bson.M{"$push": "$_id"},
		}},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot find in %s: %w", coll.Name(), err)
	}

	defer c.Close(ctx)
	r := struct {
		IDs []string `bson:"ids"`
	}{}
	if c.Next(ctx) {
		err = c.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("cannot decode result from %s: %w", coll.Name(), err)
		}
	}

	if err = c.Err(); err != nil {
		return nil, fmt.Errorf("cannot fetch from %s: %w", coll.Name(), err)
	}

	return r.IDs, nil
}

func findRules(ctx context.Context, coll mongo.DbCollection) ([]eventfilter.Rule, error) {
	c, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("cannot find in %s: %w", coll.Name(), err)
	}

	defer c.Close(ctx)
	var rules []eventfilter.Rule
	err = c.All(ctx, &rules)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch from %s: %w", coll.Name(), err)
	}

	return rules, nil
}
