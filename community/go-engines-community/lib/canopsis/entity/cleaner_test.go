package entity_test

import (
	"context"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCleaner_Clean(t *testing.T) {
	ctx := t.Context()
	client, err := mongo.NewClient(ctx)
	if err != nil {
		t.Fatalf("cannot connect to mongodb: %v", err)
	}

	now := datetime.NewCpsTime()
	entityCollection := client.Collection(mongo.EntityMongoCollection)
	archivedCollection := client.Collection(mongo.ArchivedEntitiesMongoCollection)

	redisSession, err := redis.NewSession(ctx, redis.EngineLockStorage, zerolog.Nop(), 0, 0)
	if err != nil {
		t.Fatalf("cannot create redis session: %v", err)
	}

	redlockClient := redis.NewLockClient(redisSession)

	dataStorageAdapter := datastorage.NewAdapter(client)

	configAdapter := config.NewAdapter(client)
	cfg, err := configAdapter.GetConfig(ctx)
	if err != nil {
		t.Fatalf("cannot load config: %v", err)
	}

	dataStorageCfgProvider := config.NewDataStorageConfigProvider(cfg, zerolog.Nop())

	f := func(
		conf datastorage.Config,
		limit int,
		entities []any,
		expectedIDs, expectedArchivedIDs []string, expectedArchived int64,
	) {
		t.Helper()

		err = cleanCollections(ctx, entityCollection, archivedCollection, entities)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}

		c := entity.NewCleaner(redlockClient, dataStorageAdapter, dataStorageCfgProvider, metrics.NewNullMetaUpdater(), zerolog.Nop())
		res, err := c.Clean(ctx, client, conf, now, limit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ids, err := findIDs(ctx, entityCollection)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedIDs, ids); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", entityCollection.Name(), diff)
		}

		archivedIDs, err := findIDs(ctx, archivedCollection)
		if err != nil {
			t.Fatalf("cannot find: %v", err)
		}

		if diff := pretty.Compare(expectedArchivedIDs, archivedIDs); diff != "" {
			t.Fatalf("unexpected result in %s: (-want +got):\n%s", archivedCollection.Name(), diff)
		}

		if res.Archived != expectedArchived {
			t.Fatalf("unexpected archived count, expected %d, got %d", expectedArchived, res.Archived)
		}

		err = cleanCollections(ctx, entityCollection, archivedCollection, nil)
		if err != nil {
			t.Fatalf("cannot clean: %v", err)
		}
	}

	dayAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -1)}
	weekAgo := datetime.CpsTime{Time: now.AddDate(0, 0, -7)}
	f(
		newDSConfig("1d"),
		0,
		[]any{
			newTestEntity("id1", weekAgo, &dayAgo, types.EntityTypeResource),
			newTestEntity("id2", dayAgo, nil, types.EntityTypeResource),
			newTestEntity("id3", now, nil, types.EntityTypeResource),
			newTestEntity("id4", weekAgo, &now, types.EntityTypeResource),
			newTestEntity("id5", weekAgo, &dayAgo, types.EntityTypeComponent),
			newTestEntity("id6", dayAgo, nil, types.EntityTypeComponent),
			newTestEntity("id7", now, nil, types.EntityTypeComponent),
			newTestEntity("id8", weekAgo, &now, types.EntityTypeComponent),
			newTestEntity("id9", weekAgo, &dayAgo, types.EntityTypeConnector),
			newTestEntity("id10", dayAgo, nil, types.EntityTypeConnector),
			newTestEntity("id11", now, nil, types.EntityTypeConnector),
			newTestEntity("id12", weekAgo, &now, types.EntityTypeConnector),
		},
		[]string{"id3", "id4", "id7", "id8", "id11", "id12"},
		[]string{"id1", "id2", "id5", "id6", "id9", "id10"},
		6,
	)
}

func newDSConfig(archivedAfterStr string) datastorage.Config {
	c := datastorage.Config{}
	enabled := true
	d, err := datetime.ParseDurationWithUnit(archivedAfterStr)
	if err != nil {
		panic(err)
	}

	c.Entity.ArchiveAfter = &datetime.DurationWithEnabled{
		DurationWithUnit: d,
		Enabled:          &enabled,
	}

	return c
}

func newTestEntity(id string, created datetime.CpsTime, lastEventDate *datetime.CpsTime, t string) types.Entity {
	return types.Entity{
		ID:            id,
		Created:       created,
		LastEventDate: lastEventDate,
		Type:          t,
	}
}

func cleanCollections(ctx context.Context, coll, archivedColl mongo.DbCollection, entities []any) error {
	_, err := coll.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", coll.Name(), err)
	}

	_, err = archivedColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("cannot clean %s: %w", archivedColl.Name(), err)
	}

	if len(entities) > 0 {
		_, err = coll.InsertMany(ctx, entities)
		if err != nil {
			return fmt.Errorf("cannot insert into %s: %w", coll.Name(), err)
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
