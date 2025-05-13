package axe_test

import (
	"context"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestExternalTagCleaner_Clean(t *testing.T) {
	ctx := t.Context()
	client, err := mongo.NewClient(ctx, mongo.ClientOptions{})
	if err != nil {
		t.Fatalf("cannot connect to mongodb: %v", err)
	}

	now := datetime.NewCpsTime()
	tagColl := client.Collection(mongo.AlarmTagCollection)
	colorColl := client.Collection(mongo.AlarmTagColorCollection)
	f := func(
		conf datastorage.Config,
		limit int,
		tags, colors []any,
		expectedTagIDs, expectedColorIDs []string,
		expectedDeleted int64,
	) {
		t.Helper()

		err = cleanCollections(ctx, tagColl, colorColl, tags, colors)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}

		c := axe.NewExternalTagCleaner(zerolog.Nop())
		res, err := c.Clean(ctx, client, conf, now, limit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tagIDs, err := findIDs(ctx, tagColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedTagIDs, tagIDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", tagColl.Name(), diff)
		}

		colorIDs, err := findIDs(ctx, colorColl)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedColorIDs, colorIDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", colorColl.Name(), diff)
		}

		if res.Deleted != expectedDeleted {
			t.Fatalf("unexpected deleted count: expected %d, got %d", expectedDeleted, res.Deleted)
		}

		err = cleanCollections(ctx, tagColl, colorColl, nil, nil)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}
	}

	weekAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -7)}
	f(
		newDSConfig("7d"),
		0,
		[]any{
			newTestExTag("t1", "l1", "c1", now),
			newTestExTag("t2", "l1", "c1", weekAgo),
			newTestExTag("t3", "l3", "c3", weekAgo),
			newTestInTag("t4", "l4", "c4"),
		},
		[]any{
			newTestColor("l1", "c1"),
			newTestColor("l2", "c2"),
			newTestColor("l3", "c3"),
			newTestColor("l4", "c4"),
		},
		[]string{"t1", "t4"},
		[]string{"l1", "l4"},
		2,
	)
	f(
		newDSConfig("7d"),
		3,
		[]any{
			newTestExTag("t1", "l1", "c1", weekAgo),
			newTestExTag("t2", "l1", "c1", weekAgo),
			newTestExTag("t3", "l1", "c1", weekAgo),
			newTestExTag("t4", "l1", "c1", weekAgo),
		},
		[]any{
			newTestColor("l1", "c1"),
		},
		[]string{"t4"},
		[]string{"l1"},
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

	c.AlarmExternalTag.DeleteAfter = &datetime.DurationWithEnabled{
		DurationWithUnit: d,
		Enabled:          &enabled,
	}

	return c
}

func newTestExTag(id, l, c string, d datetime.CpsTime) alarmtag.AlarmTag {
	return alarmtag.AlarmTag{
		ID:            id,
		Label:         l,
		Value:         l + ": val" + id,
		Type:          alarmtag.TypeExternal,
		Color:         c,
		LastEventDate: d,
	}
}

func newTestInTag(id, l, c string) alarmtag.AlarmTag {
	return alarmtag.AlarmTag{
		ID:    id,
		Label: l,
		Value: l + ": val" + id,
		Type:  alarmtag.TypeInternal,
		Color: c,
	}
}

func newTestColor(l, c string) map[string]string {
	return map[string]string{
		"_id":   l,
		"color": c,
	}
}

func cleanCollections(ctx context.Context, tagColl, colorColl mongo.DbCollection, tags, colors []any) error {
	_, err := tagColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", tagColl.Name(), err)
	}

	_, err = colorColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", colorColl.Name(), err)
	}

	if len(tags) > 0 {
		_, err = tagColl.InsertMany(ctx, tags)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", tagColl.Name(), err)
		}
	}

	if len(colors) > 0 {
		_, err = colorColl.InsertMany(ctx, colors)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", colorColl.Name(), err)
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
