package datastorage

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PeriodicalWorker interface {
	engine.PeriodicalWorker
	AddCleaner(k string, c Cleaner)
	// OnSchedule defines if data cleaning should be executed only on schedule.
	OnSchedule(v bool)
}

type Cleaner interface {
	IsEnabled(conf Config) bool
	Clean(ctx context.Context, dbClient mongo.DbClient, conf Config, t datetime.CpsTime, limit int) (CleanResult, error)
}

type CleanResult struct {
	Deleted  int64
	Archived int64
}

type NewDbClient func(ctx context.Context, clientTimeout time.Duration) (mongo.DbClient, error)

func NewPeriodicalWorker(
	newDbClient NewDbClient,
	periodicalInterval time.Duration,
	timezoneConfigProvider config.TimezoneConfigProvider,
	scheduleConfigProvider config.DataStorageConfigProvider,
	logger zerolog.Logger,
) PeriodicalWorker {
	return &worker{
		newDbClient:            newDbClient,
		periodicalInterval:     periodicalInterval,
		timezoneConfigProvider: timezoneConfigProvider,
		scheduleConfigProvider: scheduleConfigProvider,
		onSchedule:             true,
		cleaners:               make([]orderedCleaner, 0),
		logger:                 logger,
	}
}

type worker struct {
	newDbClient            NewDbClient
	periodicalInterval     time.Duration
	timezoneConfigProvider config.TimezoneConfigProvider
	scheduleConfigProvider config.DataStorageConfigProvider
	onSchedule             bool
	cleaners               []orderedCleaner
	logger                 zerolog.Logger
}

type orderedCleaner struct {
	Key     string
	Cleaner Cleaner
}

func (w *worker) AddCleaner(k string, c Cleaner) {
	w.cleaners = append(w.cleaners, orderedCleaner{k, c})
}

func (w *worker) OnSchedule(v bool) {
	w.onSchedule = v
}

func (w *worker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *worker) Work(ctx context.Context) {
	if len(w.cleaners) == 0 {
		return
	}

	now := datetime.NewCpsTime().In(w.timezoneConfigProvider.Get().Location)
	schConf := w.scheduleConfigProvider.Get()
	if w.onSchedule && !schConf.TimeToExecute.IsScheduledTime(now) {
		return
	}

	limit := schConf.MaxUpdates
	dbClient, err := w.newDbClient(ctx, schConf.MongoClientTimeout)
	if err != nil {
		w.logger.Err(err).Msg("cannot connect to mongo")

		return
	}

	defer func() {
		err = dbClient.Disconnect(context.WithoutCancel(ctx))
		if err != nil {
			w.logger.Err(err).Msg("cannot disconnect from mongo")
		}
	}()

	conf := DataStorage{}
	confCollection := dbClient.Collection(mongo.ConfigurationMongoCollection)
	err = confCollection.FindOne(ctx, bson.M{"_id": ID}).Decode(&conf)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		w.logger.Err(err).Msg("cannot find configuration")

		return
	}

	cleanCtx := ctx
	deadline, ok := ctx.Deadline()
	if schConf.Timeout > 0 && (!ok || schConf.Timeout < time.Until(deadline)) {
		var cancel context.CancelFunc
		cleanCtx, cancel = context.WithTimeout(ctx, schConf.Timeout)
		defer cancel()
	}

	// start from the oldest executed worker to avoid timeout run out on the same workers
	startIdx := w.getStartIdx(conf)
	if startIdx < 0 {
		return
	}

	i := startIdx
	for {
		c := w.cleaners[i]
		if c.Cleaner.IsEnabled(conf.Config) && (!w.onSchedule || !w.isAlreadyExecuted(now, conf.History, c.Key)) {
			res, err := c.Cleaner.Clean(cleanCtx, dbClient, conf.Config, now, limit)
			h := HistoryWithCount{
				Time:     now,
				Deleted:  res.Deleted,
				Archived: res.Archived,
			}
			if err != nil {
				if schConf.Timeout > 0 && errors.Is(err, context.DeadlineExceeded) &&
					ctx.Err() == nil && cleanCtx.Err() != nil {
					w.logger.Warn().Msg("cannot finish data cleaning for " + schConf.Timeout.String())
					_, err = confCollection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": bson.M{
						"history." + c.Key: h,
					}}, options.UpdateOne().SetUpsert(true))
					if err != nil {
						w.logger.Err(err).Msg("cannot update config history")
					}
				} else {
					w.logger.Err(err).Msg("cannot clean " + c.Key)
				}

				return
			}

			_, err = confCollection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": bson.M{
				"history." + c.Key: h,
			}}, options.UpdateOne().SetUpsert(true))
			if err != nil {
				w.logger.Err(err).Msg("cannot update config history")
			}
		}

		i++
		if i == len(w.cleaners) {
			i = 0
		}

		// stop after loops reaches the starting worker
		if i == startIdx {
			break
		}
	}
}

func (w *worker) isAlreadyExecuted(t datetime.CpsTime, history map[string]HistoryWithCount, key string) bool {
	if h, ok := history[key]; ok {
		lastExecuted := h.Time.In(w.timezoneConfigProvider.Get().Location)
		if t.EqualDay(lastExecuted) && t.Hour() == lastExecuted.Hour() {
			return true
		}
	}

	return false
}

func (w *worker) getStartIdx(conf DataStorage) int {
	oldestExecCleanerIdx, neverExecCleanerIdx := -1, -1
	var oldestExecTime datetime.CpsTime
	for i, c := range w.cleaners {
		if !c.Cleaner.IsEnabled(conf.Config) {
			continue
		}

		if h, ok := conf.History[c.Key]; ok {
			if oldestExecCleanerIdx < 0 || h.Time.Before(oldestExecTime) {
				oldestExecCleanerIdx = i
				oldestExecTime = h.Time
			}
		} else if neverExecCleanerIdx < 0 {
			neverExecCleanerIdx = i
		}
	}

	if neverExecCleanerIdx >= 0 {
		return neverExecCleanerIdx
	}

	return oldestExecCleanerIdx
}
