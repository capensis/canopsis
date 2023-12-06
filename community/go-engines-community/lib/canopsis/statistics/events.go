package statistics

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventStatisticsSender interface {
	Send(ctx context.Context, entityID string, stats EventStatistics)
}

type EventStatistics struct {
	OK        int               `json:"ok" bson:"ok"`
	KO        int               `json:"ko" bson:"ko"`
	LastEvent *datetime.CpsTime `json:"last_event,omitempty" bson:"last_event,omitempty" swaggertype:"integer"`
	LastKO    *datetime.CpsTime `json:"last_ko,omitempty" bson:"last_ko,omitempty" swaggertype:"integer"`
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
	if stats.LastEvent.IsZero() {
		return
	}

	set := bson.M{"last_event": stats.LastEvent}
	if stats.KO > 0 {
		set["last_ko"] = stats.LastEvent
		stats.LastKO = stats.LastEvent
	}

	res := s.eventStatisticsCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": entityID},
		bson.M{
			"$set": set,
			"$inc": bson.M{"ok": stats.OK, "ko": stats.KO},
		},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.Before),
	)
	if err := res.Err(); err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		s.logger.Err(err).Msg("failed to send event statistics")
	}

	var prev = EventStatistics{}
	if err := res.Decode(&prev); err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		s.logger.Err(err).Msg("failed to decode event statistics")
	}

	if prev.LastEvent == nil {
		return
	}

	location := s.timezoneConfigProvider.Get().Location

	year, month, day := stats.LastEvent.In(location).Date()
	truncatedInLocation := time.Date(year, month, day, 0, 0, 0, 0, location)

	//basically if it's the next day then start a new statistics
	if truncatedInLocation.Unix() > prev.LastEvent.Unix() {
		set["ok"] = stats.OK
		set["ko"] = stats.KO

		_, err := s.eventStatisticsCollection.UpdateOne(
			ctx,
			bson.M{"_id": entityID},
			bson.M{"$set": set},
		)
		if err != nil {
			s.logger.Err(err).Msg("failed to send event statistics")
		}
	}
}
