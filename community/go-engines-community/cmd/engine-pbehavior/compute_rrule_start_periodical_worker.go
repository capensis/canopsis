package main

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	librrule "github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type computeRruleStartPeriodicalWorker struct {
	PeriodicalInterval     time.Duration
	PbhCollection          libmongo.DbCollection
	TimezoneConfigProvider config.TimezoneConfigProvider
	Logger                 zerolog.Logger
}

func (w *computeRruleStartPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *computeRruleStartPeriodicalWorker) Work(ctx context.Context) {
	err := w.updatePbehaviorComputedStarts(ctx)
	if err != nil {
		w.Logger.Err(err).Msgf("cannot update pbehaviors")
	}
}

func (w *computeRruleStartPeriodicalWorker) updatePbehaviorComputedStarts(ctx context.Context) error {
	now := time.Now().In(w.TimezoneConfigProvider.Get().Location)
	currMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// Find computed start between the first day of the previous month and the last week of the previous month.
	// Extra week interval is required for Calendar view in case if the first day of the current month is not Monday.
	after := currMonth.AddDate(0, -1, 0)
	before := currMonth.AddDate(0, 0, -7)
	cursor, err := w.PbhCollection.Find(ctx, bson.M{
		"rrule": bson.M{"$nin": bson.A{"", nil}},
		"$or": []bson.M{
			{"rrule_cstart": bson.M{
				"$ne": nil,
				"$lt": after.Unix(),
			}},
			{
				"rrule_cstart": nil,
				"tstart":       bson.M{"$lt": after.Unix()},
			},
		},
	})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	writeModels := make([]mongo.WriteModel, 0)
	locs := make(map[string]*time.Location)
	for cursor.Next(ctx) {
		pbh := pbehavior.PBehavior{}
		err := cursor.Decode(&pbh)
		if err != nil {
			return err
		}

		rOption, err := librrule.StrToROption(pbh.RRule)
		if err != nil {
			w.Logger.Err(err).Str("pbehavior", pbh.ID).Msgf("cannot parse rrule")
			continue
		}

		location := w.TimezoneConfigProvider.Get().Location
		if pbh.Timezone != "" {
			var ok bool
			if location, ok = locs[pbh.Timezone]; !ok {
				location, err = time.LoadLocation(pbh.Timezone)
				if err != nil {
					w.Logger.Err(err).Str("pbehavior", pbh.ID).Str("timezone", pbh.Timezone).Msgf("invalid timezone")
					continue
				}

				locs[pbh.Timezone] = location
			}
		}

		if pbh.RRuleComputedStart != nil {
			rOption.Dtstart = pbh.RRuleComputedStart.Time.In(location)
		} else {
			rOption.Dtstart = pbh.Start.Time.In(location)
		}

		r, err := librrule.NewRRule(*rOption)
		if err != nil {
			w.Logger.Err(err).Str("pbehavior", pbh.ID).Msgf("cannot create rrule")
			continue
		}

		dates := r.Between(after, before, true)
		if len(dates) == 0 {
			continue
		}

		computedStart := dates[len(dates)-1]
		writeModels = append(writeModels, mongo.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id":    pbh.ID,
				"tstart": pbh.Start,
				"rrule":  pbh.RRule,
				"$or": []bson.M{
					{"rrule_cstart": nil},
					{"rrule_cstart": bson.M{"$lt": computedStart.Unix()}},
				},
			}).
			SetUpdate(bson.M{"$set": bson.M{
				"rrule_cstart": computedStart.Unix(),
			}}))

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := w.PbhCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err := w.PbhCollection.BulkWrite(ctx, writeModels)
		return err
	}

	return nil
}
