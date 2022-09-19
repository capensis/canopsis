package che

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"

	"github.com/teambition/rrule-go"

	"go.mongodb.org/mongo-driver/bson"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"

	"github.com/rs/zerolog"
)

func NewEventfilterIntervalsWorker(
	client mongo.DbClient,
	provider config.TimezoneConfigProvider,
	periodicalInterval time.Duration,
	logger zerolog.Logger,
) engine.PeriodicalWorker {
	return &eventfilterIntervalsWorker{
		collection:             client.Collection(mongo.EventFilterRulesMongoCollection),
		timezoneConfigProvider: provider,
		periodicalInterval:     periodicalInterval,
		logger:                 logger,
	}
}

type eventfilterIntervalsWorker struct {
	collection             mongo.DbCollection
	timezoneConfigProvider config.TimezoneConfigProvider
	periodicalInterval     time.Duration
	logger                 zerolog.Logger
}

func (w *eventfilterIntervalsWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *eventfilterIntervalsWorker) Work(ctx context.Context) {
	location := w.timezoneConfigProvider.Get().Location
	now := time.Now().In(location)

	cursor, err := w.collection.
		Aggregate(
			ctx,
			[]bson.M{
				{
					"$match": bson.M{
						"start":   bson.M{"$nin": bson.A{0, nil}},
						"stop":    bson.M{"$nin": bson.A{0, nil}},
						"enabled": bson.M{"$nin": bson.A{false, nil}},
					},
				},
				{
					"$lookup": bson.M{
						"from":         mongo.PbehaviorExceptionMongoCollection,
						"localField":   "exceptions",
						"foreignField": "_id",
						"as":           "exceptions",
					},
				},
				{
					"$unwind": bson.M{
						"path":                       "$exceptions",
						"preserveNullAndEmptyArrays": true,
					},
				},
				{
					"$unwind": bson.M{
						"path":                       "$exceptions.exdates",
						"preserveNullAndEmptyArrays": true,
					},
				},
				{
					"$group": bson.M{
						"_id":                 "$_id",
						"exdates":             bson.M{"$first": "$exdates"},
						"rrule":               bson.M{"$first": "$rrule"},
						"start":               bson.M{"$first": "$start"},
						"stop":                bson.M{"$first": "$stop"},
						"resolved_start":      bson.M{"$first": "$resolved_start"},
						"resolved_stop":       bson.M{"$first": "$resolved_stop"},
						"next_resolved_start": bson.M{"$first": "$next_resolved_start"},
						"next_resolved_stop":  bson.M{"$first": "$next_resolved_stop"},
						"in_exdate":           bson.M{"$first": "$in_exdate"},
						"exceptions_exdates": bson.M{
							"$push": "$exceptions.exdates.exdate",
						},
					},
				},
				{
					"$addFields": bson.M{
						"resolved_exdates": bson.M{
							"$concatArrays": []bson.M{
								{
									"$ifNull": bson.A{"$exceptions_exdates", []bson.M{}},
								},
								{
									"$ifNull": bson.A{"$exdates", []bson.M{}},
								},
							},
						},
					},
				},
				{
					"$unwind": bson.M{
						"path":                       "$resolved_exdates",
						"preserveNullAndEmptyArrays": true,
					},
				},
				{
					"$group": bson.M{
						"_id":                 "$_id",
						"rrule":               bson.M{"$first": "$rrule"},
						"start":               bson.M{"$first": "$start"},
						"stop":                bson.M{"$first": "$stop"},
						"resolved_start":      bson.M{"$first": "$resolved_start"},
						"resolved_stop":       bson.M{"$first": "$resolved_stop"},
						"next_resolved_start": bson.M{"$first": "$next_resolved_start"},
						"next_resolved_stop":  bson.M{"$first": "$next_resolved_stop"},
						"in_exdate":           bson.M{"$first": "$in_exdate"},
						"resolved_exdates": bson.M{
							"$push": bson.M{
								"$cond": bson.A{
									bson.M{"$and": bson.A{
										bson.M{"gt": bson.A{time.Now().Unix(), "$resolved_exdates.end"}},
										bson.M{"lt": bson.A{"$resolved_exdates.start", time.Now().Add(w.periodicalInterval * 2).Unix()}},
									}},
									"$resolved_exdates",
									"$$REMOVE",
								},
							},
						},
					},
				},
			},
		)

	if err != nil {
		w.logger.Error().Err(err).Msg("unable to load eventfilter rules with rrule")
	}

	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	bulkBytesSize := 0

	for cursor.Next(ctx) {
		var ef eventfilter.Rule
		err = cursor.Decode(&ef)
		if err != nil {
			w.logger.Error().Err(err).Str("rule_id", ef.ID).Msg("failed to decode the eventfilter with rrule")
			continue
		}

		if ef.Start.IsZero() {
			w.logger.Error().Str("rule_id", ef.ID).Msg("failed to calculate the rrule interval, start date can't be zero")
			continue
		}

		if ef.Stop.IsZero() {
			w.logger.Error().Str("rule_id", ef.ID).Msg("failed to calculate the rrule interval, stop date can't be zero")
			continue
		}

		if ef.ResolvedStart.IsZero() {
			w.logger.Error().Str("rule_id", ef.ID).Msg("failed to calculate the rrule interval, resolved start date can't be zero")
			continue
		}

		if ef.ResolvedStop.IsZero() {
			w.logger.Error().Str("rule_id", ef.ID).Msg("failed to calculate the rrule interval, resolved stop date can't be zero")
			continue
		}

		var r *rrule.RRule
		if ef.RRule != "" {
			opt, err := rrule.StrToROption(ef.RRule)
			if err != nil {
				w.logger.Error().Err(err).Str("rule_id", ef.ID).Str("rrule", ef.RRule).Msg("failed to parse the rrule string")
				continue
			}

			r, err = rrule.NewRRule(*opt)
			if err != nil {
				w.logger.Error().Err(err).Str("rule_id", ef.ID).Str("rrule", ef.RRule).Msg("failed to create the rrule")
				continue
			}

			if opt.Count != 0 {
				r.DTStart(ef.Start.Time.In(location))
			} else {
				r.DTStart(ef.ResolvedStart.Time.In(location))
			}
		}

		eventfilter.ResolveIntervals(&ef, r, now, location)

		set := bson.M{"resolved_start": ef.ResolvedStart, "resolved_stop": ef.ResolvedStop, "resolved_exdates": ef.ResolvedExdates}
		unset := bson.M{}
		if ef.NextResolvedStart != nil && ef.NextResolvedStop != nil {
			set["next_resolved_start"] = ef.NextResolvedStart
			set["next_resolved_stop"] = ef.NextResolvedStop
		} else {
			unset["next_resolved_start"] = ""
			unset["next_resolved_stop"] = ""
		}

		newModel := mongodriver.
			NewUpdateOneModel().
			SetFilter(bson.M{"_id": ef.ID}).
			SetUpdate(bson.M{"$set": set, "$unset": unset})

		b, err := bson.Marshal(newModel)
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to marshal eventfilter update model")
			continue
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err := w.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.logger.Error().Err(err).Msg("unable to bulk write eventfilters")
				return
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, newModel)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := w.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				w.logger.Error().Err(err).Msg("unable to bulk write eventfilters")
				return
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		_, err = w.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			w.logger.Error().Err(err).Msg("unable to bulk write eventfilters")
			return
		}
	}
}
