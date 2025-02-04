package pbehavior_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
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
	pbhColl := client.Collection(mongo.PbehaviorMongoCollection)
	f := func(
		conf datastorage.Config,
		limit int,
		pbhs []any,
		expectedPbhIDs []string,
		expectedDeleted int64,
	) {
		t.Helper()

		err = cleanCollections(ctx, pbhColl, pbhs)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}

		c := libpbehavior.NewCleaner(zerolog.Nop())
		res, err := c.Clean(ctx, client, conf, now, limit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		pbhIDs, err := findIDs(ctx, pbhColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedPbhIDs, pbhIDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", pbhColl.Name(), diff)
		}

		if res.Deleted != expectedDeleted {
			t.Fatalf("unexpected deleted count, expected %d, got %d", expectedDeleted, res.Deleted)
		}

		err = cleanCollections(ctx, pbhColl, nil)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}
	}

	loc, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		panic(err)
	}

	weekAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -7)}
	weekAndHourAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -7).Add(-time.Hour)}
	weekWithoutHourAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -7).Add(time.Hour)}
	twoWeekAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -14)}
	twoWeekAndHourAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -14).Add(-time.Hour)}
	f(
		newDSConfig("7d"),
		0,
		[]any{
			newTestPbh("p1", weekAndHourAgo, weekAgo, loc, ""),
			newTestPbh("p2", weekAndHourAgo, weekWithoutHourAgo, loc, ""),
			newTestPbh("p3", weekAgo, weekWithoutHourAgo, loc, ""),
			newTestPbh("p4", weekAndHourAgo, datetime.CpsTime{}, loc, ""),
			newTestPbh("p5", twoWeekAndHourAgo, twoWeekAgo, loc, "FREQ=DAILY;COUNT=8"),
			newTestPbh("p6", twoWeekAndHourAgo, twoWeekAgo, loc, "FREQ=DAILY;COUNT=9"),
			newTestPbh("p7", weekAndHourAgo, weekAgo, loc, "FREQ=DAILY;COUNT=7"),
			newTestPbh("p8", twoWeekAndHourAgo, twoWeekAgo, loc, "FREQ=DAILY"),
		},
		[]string{"p2", "p3", "p4", "p6", "p7", "p8"},
		2,
	)
	f(
		newDSConfig("7d"),
		3,
		[]any{
			newTestPbh("p1", weekAndHourAgo, weekAgo, loc, ""),
			newTestPbh("p2", weekAndHourAgo, weekAgo, loc, ""),
			newTestPbh("p3", weekAndHourAgo, weekAgo, loc, ""),
			newTestPbh("p4", weekAndHourAgo, weekAgo, loc, ""),
		},
		[]string{"p4"},
		3,
	)
}

func newDSConfig(deletedAfterStr string) datastorage.Config {
	c := datastorage.Config{}
	enabled := true
	d, err := datetime.ParseDurationWithUnit(deletedAfterStr)
	if err != nil {
		panic(err)
	}

	c.Pbehavior.DeleteAfter = &datetime.DurationWithEnabled{
		DurationWithUnit: d,
		Enabled:          &enabled,
	}

	return c
}

func newTestPbh(id string, b, e datetime.CpsTime, l *time.Location, rrule string) libpbehavior.PBehavior {
	p := libpbehavior.PBehavior{
		ID:       id,
		Name:     "name" + id,
		Timezone: l.String(),
		Start:    &b,
		RRule:    rrule,
	}
	if !e.IsZero() {
		p.Stop = &e
	}

	var err error
	p.RRuleEnd, err = libpbehavior.GetRruleEnd(b, rrule, l)
	if err != nil {
		panic(err)
	}

	return p
}

func cleanCollections(ctx context.Context, pbhColl mongo.DbCollection, pbhs []any) error {
	_, err := pbhColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", pbhColl.Name(), err)
	}

	if len(pbhs) > 0 {
		_, err = pbhColl.InsertMany(ctx, pbhs)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", pbhColl.Name(), err)
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
