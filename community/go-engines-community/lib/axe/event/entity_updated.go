package event

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewEntityUpdatedProcessor(
	dbClient mongo.DbClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	externalTagUpdater alarmtag.ExternalUpdater,
) Processor {
	return &entityUpdatedProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityCollection:                dbClient.Collection(mongo.EntityMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		externalTagUpdater:              externalTagUpdater,
	}
}

type entityUpdatedProcessor struct {
	dbClient                        mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityCollection                mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	externalTagUpdater              alarmtag.ExternalUpdater
}

func (p *entityUpdatedProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	entity := *event.Entity
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo

	var componentStateChanged bool
	var newComponentState int

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedServiceStates = nil
		result = Result{}

		var update []bson.M
		if event.Parameters.ImportSource != "" {
			importTags := types.TransformEventTags(event.Parameters.ImportTags)
			var setTags bson.M
			if len(importTags) > 0 {
				setTags = bson.M{"$concatArrays": bson.A{
					bson.M{"$cond": bson.M{"if": "$etags", "then": "$etags", "else": bson.A{}}},
					bson.M{"$cond": bson.M{"if": "$itags", "then": "$itags", "else": bson.A{}}},
					importTags,
				}}
			} else {
				setTags = bson.M{"$cond": bson.M{
					"if": bson.M{"$and": bson.A{
						"$imtags",
						bson.M{"$ne": bson.A{"$imtags", nil}},
						bson.M{"$ne": bson.A{"$imtags", bson.A{}}},
					}},
					"then": bson.M{"$concatArrays": []bson.M{
						{"$cond": bson.M{"if": "$etags", "then": "$etags", "else": bson.A{}}},
						{"$cond": bson.M{"if": "$itags", "then": "$itags", "else": bson.A{}}},
					}},
					"else": "$tags",
				}}
			}

			update = []bson.M{
				{"$set": bson.M{"tags": setTags}},
				{"$set": bson.M{"imtags": importTags}},
			}
		}

		alarm := types.Alarm{}
		var err error
		if len(update) > 0 {
			err = p.alarmCollection.FindOneAndUpdate(ctx, getOpenAlarmMatch(event), update).Decode(&alarm)
		} else {
			err = p.alarmCollection.FindOne(ctx, getOpenAlarmMatch(event)).Decode(&alarm)
		}

		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		result.IsCountersUpdated, updatedServiceStates, componentStateChanged, newComponentState, err = processComponentAndServiceCounters(
			ctx,
			p.entityServiceCountersCalculator,
			p.componentCountersCalculator,
			&alarm,
			&entity,
			result.AlarmChange,
		)
		if err != nil {
			return err
		}

		if entity.Type == types.EntityTypeComponent {
			// force update
			componentStateChanged = true

			newComponentState, err = p.componentCountersCalculator.RecomputeCounters(ctx, &entity)
			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		return result, err
	}

	for servID, servInfo := range updatedServiceStates {
		err = p.eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			return result, fmt.Errorf("failed to update service state: %w", err)
		}
	}

	if componentStateChanged {
		err = p.eventsSender.UpdateComponentState(ctx, event.Entity.Component, newComponentState)
		if err != nil {
			return result, fmt.Errorf("failed to update component state: %w", err)
		}
	}

	if event.Parameters.ImportSource != "" {
		p.externalTagUpdater.Add(event.Parameters.ImportTags)
	}

	return result, nil
}
