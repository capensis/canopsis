package messageratestats

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const MinutePointsNumber = 61 // [0;60] = 61

type RatePerMinute struct {
	Time int64 `bson:"_id" json:"time"`
	Rate int64 `bson:"received" json:"rate"`
}

func GetWebsocketHandler(
	hub websocket.Hub,
	dbClient mongo.DbClient,
	logger zerolog.Logger,
	tickDuration time.Duration,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		ticker := time.NewTicker(tickDuration)
		defer ticker.Stop()

		var rates [MinutePointsNumber]RatePerMinute

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				nowTrunc := time.Now().Truncate(time.Minute)

				end := nowTrunc.Unix()
				begin := nowTrunc.Add(-time.Hour).Unix()

				for i := int64(0); i < MinutePointsNumber; i++ {
					rates[i].Time = i*60 + begin
					rates[i].Rate = 0
				}

				cursor, err := dbClient.Collection(mongo.MessageRateStatsMinuteCollectionName).Find(
					ctx,
					bson.M{"_id": bson.M{"$gte": begin, "$lte": end}},
				)
				if err != nil {
					logger.Err(err).Msg("Failed to find message rates from mongo")
					continue
				}

				for cursor.Next(ctx) {
					var rate RatePerMinute
					err = cursor.Decode(&rate)
					if err != nil {
						logger.Err(err).Msg("Failed to decode RatePerMinute")
						break
					}

					rates[(rate.Time-begin)/60].Rate = rate.Rate
				}

				err = cursor.Close(ctx)
				if err != nil {
					logger.Err(err).Msg("Failed to close cursor")
					continue
				}

				hub.Send(websocket.RoomMessageRates, rates)
			}
		}
	}
}
