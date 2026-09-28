package alarmtimeline

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store interface {
	Exists(ctx context.Context, id string) (bool, error)
	Export(ctx context.Context, t export.Task) (export.DataCursor, error)
}

func NewStore(
	dbClient mongo.DbClient,
	dbExportClient mongo.DbClient,
	decoder encoding.Decoder,
	timezoneConfigProvider config.TimezoneConfigProvider,
	userInterfaceAdapter config.UserInterfaceAdapter,
) Store {
	return &store{
		dbCollection:           dbClient.Collection(mongo.AlarmMongoCollection),
		dbResolvedCollection:   dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		dbExportCollection:     dbExportClient.Collection(mongo.AlarmMongoCollection),
		decoder:                decoder,
		timezoneConfigProvider: timezoneConfigProvider,
		userInterfaceAdapter:   userInterfaceAdapter,
	}
}

type store struct {
	dbCollection           mongo.DbCollection
	dbResolvedCollection   mongo.DbCollection
	dbExportCollection     mongo.DbCollection
	decoder                encoding.Decoder
	timezoneConfigProvider config.TimezoneConfigProvider
	userInterfaceAdapter   config.UserInterfaceAdapter
}

func (s *store) Exists(ctx context.Context, id string) (bool, error) {
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}, opts).Err()
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		err = s.dbResolvedCollection.FindOne(ctx, bson.M{"_id": id}, opts).Err()
	}

	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, fmt.Errorf("fetch alarm %q: %w", id, err)
	}

	return true, nil
}

func (s *store) Export(ctx context.Context, t export.Task) (export.DataCursor, error) {
	var r ExportFetchParameters

	if err := s.decoder.Decode([]byte(t.Parameters), &r); err != nil {
		return nil, fmt.Errorf("decode timeline export parameters: %w", err)
	}

	exists, err := s.Exists(ctx, r.AlarmID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, httperror.ErrNotFound
	}

	// The alarm may have moved to resolved storage since the export was requested.
	pipeline := []bson.M{
		{"$match": bson.M{"_id": r.AlarmID}},
		{"$project": bson.M{"v.steps": 1}},
		{"$unionWith": bson.M{
			"coll": mongo.ResolvedAlarmMongoCollection,
			"pipeline": []bson.M{
				{"$match": bson.M{"_id": r.AlarmID}},
				{"$project": bson.M{"v.steps": 1}},
			},
		}},
		{"$limit": 1},
		{"$unwind": bson.M{"path": "$v.steps", "includeArrayIndex": "step_index"}},
		{"$sort": bson.D{{Key: "v.steps.t", Value: 1}, {Key: "step_index", Value: 1}}},
		{"$replaceRoot": bson.M{"newRoot": "$v.steps"}},
	}

	cursor, err := s.dbExportCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("fetch timeline of alarm %q: %w", r.AlarmID, err)
	}

	location := s.timezoneConfigProvider.Get().Location
	timeFormat, err := export.ResolveTimeFormat(ctx, t.TimeFormat, s.userInterfaceAdapter)
	if err != nil {
		return nil, fmt.Errorf("resolve time format: %w", err)
	}

	return export.NewMongoCursor(cursor, t.Fields.Fields(), func(k string, v any) (any, error) {
		if k == "t" {
			var timestamp int64
			switch value := v.(type) {
			case int64:
				timestamp = value
			case int32:
				timestamp = int64(value)
			default:
				return v, nil
			}

			return datetime.NewCpsTime(timestamp).In(location).Time.Format(timeFormat), nil
		}

		return v, nil
	}), nil
}
