package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewEntityUpdatedProcessor(
	dbClient mongo.DbClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	externalTagUpdater alarmtag.ExternalUpdater,
	metaAlarmPostProcessor MetaAlarmPostProcessor,
	logger zerolog.Logger,
) Processor {
	return &entityUpdatedProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		externalTagUpdater:              externalTagUpdater,
		metaAlarmPostProcessor:          metaAlarmPostProcessor,
		logger:                          logger,
	}
}

type entityUpdatedProcessor struct {
	dbClient                        mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	externalTagUpdater              alarmtag.ExternalUpdater
	metaAlarmPostProcessor          MetaAlarmPostProcessor
	logger                          zerolog.Logger
}

func (p *entityUpdatedProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	entity := *event.Entity
	importTags := types.TransformEventTags(event.Parameters.ImportTags)
	var updatedServiceStates map[string]entitycounters.UpdatedServicesInfo
	var componentStateChanged bool
	var newComponentState int
	var alarm types.Alarm
	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedServiceStates = nil
		result = Result{}
		alarm = types.Alarm{}
		err := p.alarmCollection.FindOne(ctx, getOpenAlarmMatch(event)).Decode(&alarm)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		if event.Parameters.ImportSource != "" && (len(importTags) > 0 || len(alarm.ImportTags) > 0) {
			var setTags bson.M
			if len(importTags) > 0 {
				hasTagsInAlarm := make(map[string]bool, len(alarm.ImportTags))
				for _, tag := range alarm.ImportTags {
					hasTagsInAlarm[tag] = true
				}

				addedTags := make([]string, 0)
				removedTags := make([]string, 0)
				hasInImportTags := make(map[string]bool, len(importTags))
				for _, tag := range importTags {
					hasInImportTags[tag] = true
					if !hasTagsInAlarm[tag] {
						addedTags = append(addedTags, tag)
					}
				}

				for _, tag := range alarm.ImportTags {
					if !hasInImportTags[tag] {
						removedTags = append(removedTags, tag)
					}
				}

				result.AddedExternalTags = addedTags
				result.RemovedExternalTags = removedTags
				setTags = bson.M{"$concatArrays": bson.A{
					bson.M{"$cond": bson.M{"if": "$etags", "then": "$etags", "else": bson.A{}}},
					bson.M{"$cond": bson.M{"if": "$itags", "then": "$itags", "else": bson.A{}}},
					bson.M{"$literal": importTags},
				}}
			} else {
				result.RemovedExternalTags = alarm.ImportTags
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

			_, err = p.alarmCollection.UpdateOne(ctx, getOpenAlarmMatch(event), []bson.M{
				{"$set": bson.M{"tags": setTags}},
				{"$set": bson.M{"imtags": bson.M{"$literal": importTags}}},
			})
			if err != nil {
				return err
			}
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

	// to dynamic-infos
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = types.NewAlarmChange()

	go p.postProcess(context.WithoutCancel(ctx), event, result, updatedServiceStates, componentStateChanged, newComponentState)

	return result, nil
}

func (p *entityUpdatedProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
	result Result,
	updatedServiceStates map[string]entitycounters.UpdatedServicesInfo,
	componentStateChanged bool,
	newComponentState int,
) {
	for servID, servInfo := range updatedServiceStates {
		err := p.eventsSender.UpdateServiceState(ctx, servID, servInfo)
		if err != nil {
			p.logger.Err(err).Msg("failed to update service state")
		}
	}

	if componentStateChanged {
		err := p.eventsSender.UpdateComponentState(ctx, event.Entity.Component, newComponentState)
		if err != nil {
			p.logger.Err(err).Msg("failed to update component state")
		}
	}

	if event.Parameters.ImportSource != "" {
		p.externalTagUpdater.Add(event.Parameters.ImportTags)
	}

	err := p.metaAlarmPostProcessor.Process(ctx, event, rpc.AxeResultEvent{
		Alarm:               &result.Alarm,
		AlarmChangeType:     result.AlarmChange.Type,
		AddedExternalTags:   result.AddedExternalTags,
		RemovedExternalTags: result.RemovedExternalTags,
	})
	if err != nil {
		p.logger.Err(err).Msg("cannot process meta alarm")
	}
}
