package main

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func computeRruleEnd(ctx context.Context, logger zerolog.Logger) error {
	dbClient, err := mongo.NewClient(ctx, 0, 0, logger)
	if err != nil {
		return err
	}

	defer dbClient.Disconnect(ctx)
	cfg, err := config.NewAdapter(dbClient).GetConfig(ctx)
	if err != nil {
		return err
	}
	timezoneConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	loc := timezoneConfigProvider.Get().Location
	dbCollection := dbClient.Collection(mongo.PbehaviorMongoCollection)
	cursor, err := dbCollection.Find(ctx, bson.M{
		"rrule":     bson.M{"$nin": bson.A{"", nil}},
		"rrule_end": nil,
	})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	writeModels := make([]mongodriver.WriteModel, 0)
	for cursor.Next(ctx) {
		pbh := pbehavior.PBehavior{}
		err := cursor.Decode(&pbh)
		if err != nil {
			return err
		}

		if pbh.Start == nil {
			continue
		}

		rruleEnd, err := pbehavior.GetRruleEnd(*pbh.Start, pbh.RRule, loc)
		if err != nil {
			logger.Err(err).Str("pbehavior", pbh.ID).Msgf("cannot update pbehavior rrule end")
			continue
		}

		if rruleEnd == nil {
			continue
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id":       pbh.ID,
				"rrule":     pbh.RRule,
				"tstart":    pbh.Start,
				"rrule_end": nil,
			}).
			SetUpdate(bson.M{"$set": bson.M{"rrule_end": rruleEnd}}),
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = dbCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err = dbCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}

	return nil
}
