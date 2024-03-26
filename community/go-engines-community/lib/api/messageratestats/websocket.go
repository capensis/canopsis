package messageratestats

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
)

const MinutePointsNumber = 61 // [0;60] = 61

type RatePerMinute struct {
	Time int64 `json:"time"`
	Rate int64 `json:"rate"`
}

func GetWebsocketHandler(
	hub websocket.Hub,
	pgPoolProvider postgres.PoolProvider,
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
				end := nowTrunc
				begin := nowTrunc.Add(-time.Hour)
				pgPool, err := pgPoolProvider.Get(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot connect to postgres")
					continue
				}

				bucket := "1 minute"
				table := metrics.MessageRate
				rows, err := pgPool.Query(ctx, "SELECT time_bucket_gapfill('"+bucket+"', time), count(*) FROM "+table+
					" WHERE time >= $1 AND time <= $2 GROUP BY time_bucket_gapfill('"+bucket+"', time)", begin.UTC(), end.UTC())
				if err != nil {
					logger.Err(err).Msg("cannot find message rates")
					continue
				}

				i := 0
				for rows.Next() {
					var t time.Time
					var rateColumn *int64
					err = rows.Scan(&t, &rateColumn)
					if err != nil {
						logger.Err(err).Msg("cannot scan result")
						break
					}

					var rate int64
					if rateColumn != nil {
						rate = *rateColumn
					}

					if i >= len(rates) {
						logger.Error().Msg("invalid postgres query")
						break
					}

					rates[i] = RatePerMinute{
						Time: t.Unix(),
						Rate: rate,
					}
					i++
				}

				rows.Close()
				if i == len(rates) {
					hub.Send(ctx, websocket.RoomMessageRates, rates)
				}
			}
		}
	}
}
