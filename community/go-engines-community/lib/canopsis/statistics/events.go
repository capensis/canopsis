package statistics

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type EventStatisticsSender interface {
	Send(ctx context.Context, entityID string, stats EventStatistics)
}

type EventStatistics struct {
	OK        int            `bson:"ok"`
	KO        int            `bson:"ko"`
	Timestamp *types.CpsTime `bson:"timestamp,omitempty"`
}

type eventStatisticsSender struct {
	eventStatisticsCollection mongo.DbCollection
	logger                    zerolog.Logger
	timezoneConfigProvider    config.TimezoneConfigProvider
}

func NewEventStatisticsSender(dbClient mongo.DbClient, logger zerolog.Logger, timezoneConfigProvider config.TimezoneConfigProvider) EventStatisticsSender {
	return &eventStatisticsSender{
		eventStatisticsCollection: dbClient.Collection(mongo.EventStatistics),
		logger:                    logger,
		timezoneConfigProvider:    timezoneConfigProvider,
	}
}

func (s *eventStatisticsSender) Send(ctx context.Context, entityID string, stats EventStatistics) {
	if stats.Timestamp.IsZero() {
		return
	}

	res := s.eventStatisticsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": entityID},
		bson.M{
			"$set": bson.M{"timestamp": stats.Timestamp},
			"$inc": bson.M{"ok": stats.OK, "ko": stats.KO},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.Before),
	)
	if err := res.Err(); err != nil && err != mongodriver.ErrNoDocuments {
		s.logger.Err(err).Msg("failed to send event statistics")
	}

	var prev = EventStatistics{}
	if err := res.Decode(&prev); err != nil && err != mongodriver.ErrNoDocuments {
		s.logger.Err(err).Msg("failed to decode event statistics")
	}

	if prev.Timestamp == nil {
		return
	}

	location := s.timezoneConfigProvider.Get().Location

	year, month, day := stats.Timestamp.In(location).Date()
	truncatedInLocation := time.Date(year, month, day, 0, 0, 0, 0, location)

	//basically if it's the next day then start a new statistics
	if truncatedInLocation.Unix() > prev.Timestamp.Unix() {
		_, err := s.eventStatisticsCollection.UpdateOne(
			ctx,
			bson.M{"_id": entityID},
			bson.M{"$set": bson.M{"timestamp": stats.Timestamp, "ok": stats.OK, "ko": stats.KO}},
		)
		if err != nil {
			s.logger.Err(err).Msg("failed to send event statistics")
		}
	}
}
