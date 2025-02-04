package alarm_test

import (
	"context"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCleaner_Clean(t *testing.T) {
	ctx := t.Context()
	client, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		t.Fatalf("cannot connect to mongodb: %v", err)
	}

	now := datetime.NewCpsTime()
	resolvedColl := client.Collection(mongo.ResolvedAlarmMongoCollection)
	archivedColl := client.Collection(mongo.ArchivedAlarmMongoCollection)
	f := func(
		conf datastorage.Config,
		limit int,
		resolvedAlarms, archivedAlarms []any,
		expectedResolvedIDs, expectedArchivedIDs []string,
		expectedDeleted, expectedArchived int64,
	) {
		t.Helper()

		err = cleanCollections(ctx, resolvedColl, archivedColl, resolvedAlarms, archivedAlarms)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}

		c := alarm.NewCleaner(zerolog.Nop())
		res, err := c.Clean(ctx, client, conf, now, limit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		resolvedIDs, err := findIDs(ctx, resolvedColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedResolvedIDs, resolvedIDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", resolvedColl.Name(), diff)
		}

		archivedIDs, err := findIDs(ctx, archivedColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedArchivedIDs, archivedIDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", archivedColl.Name(), diff)
		}

		if res.Archived != expectedArchived {
			t.Fatalf("unexpected archived count, expected %d, got %d", expectedArchived, res.Archived)
		}

		if res.Deleted != expectedDeleted {
			t.Fatalf("unexpected deleted count, expected %d, got %d", expectedDeleted, res.Deleted)
		}

		err = cleanCollections(ctx, resolvedColl, archivedColl, nil, nil)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}
	}

	dayAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -1)}
	weekAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -7)}
	f(
		newDSConfig("1d", "7d"),
		0,
		[]any{
			newTestAlarm("r1", nil),
			newTestAlarm("r2", &now),
			newTestAlarm("r3", &dayAgo),
		},
		[]any{
			newTestAlarm("a1", &dayAgo),
			newTestAlarm("a2", &weekAgo),
		},
		[]string{"r1", "r2"},
		[]string{"a1", "r3"},
		1,
		1,
	)
	f(
		newDSConfig("1d", "7d"),
		3,
		[]any{
			newTestAlarm("r1", &dayAgo),
			newTestAlarm("r2", &dayAgo),
			newTestAlarm("r3", &dayAgo),
			newTestAlarm("r4", &dayAgo),
		},
		[]any{
			newTestAlarm("a1", &weekAgo),
			newTestAlarm("a2", &weekAgo),
			newTestAlarm("a3", &weekAgo),
			newTestAlarm("a4", &weekAgo),
		},
		[]string{"r4"},
		[]string{"a4", "r1", "r2", "r3"},
		3,
		3,
	)
}

func newDSConfig(archivedAfterStr, deletedAfterStr string) datastorage.Config {
	c := datastorage.Config{}
	enabled := true
	d, err := datetime.ParseDurationWithUnit(archivedAfterStr)
	if err != nil {
		panic(err)
	}

	c.Alarm.ArchiveAfter = &datetime.DurationWithEnabled{
		DurationWithUnit: d,
		Enabled:          &enabled,
	}

	d, err = datetime.ParseDurationWithUnit(deletedAfterStr)
	if err != nil {
		panic(err)
	}

	c.Alarm.DeleteAfter = &datetime.DurationWithEnabled{
		DurationWithUnit: d,
		Enabled:          &enabled,
	}

	return c
}

func newTestAlarm(id string, resolved *datetime.CpsTime) types.Alarm {
	return types.Alarm{
		ID: id,
		Value: types.AlarmValue{
			Resolved: resolved,
		},
	}
}

func cleanCollections(ctx context.Context, resolvedColl, archivedColl mongo.DbCollection, resolvedAlarms, archivedAlarms []any) error {
	_, err := resolvedColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", resolvedColl.Name(), err)
	}

	_, err = archivedColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", archivedColl.Name(), err)
	}

	if len(resolvedAlarms) > 0 {
		_, err = resolvedColl.InsertMany(ctx, resolvedAlarms)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", resolvedColl.Name(), err)
		}
	}

	if len(archivedAlarms) > 0 {
		_, err = archivedColl.InsertMany(ctx, archivedAlarms)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", archivedColl.Name(), err)
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
