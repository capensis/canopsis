package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewContextUpdateProcessor(
	dbClient mongo.DbClient,
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator,
	componentCountersCalculator calculator.ComponentCountersCalculator,
	eventsSender entitycounters.EventsSender,
	logger zerolog.Logger,
) Processor {
	return &contextUpdateProcessor{
		dbClient:                        dbClient,
		alarmCollection:                 dbClient.Collection(mongo.AlarmMongoCollection),
		entityServiceCountersCalculator: entityServiceCountersCalculator,
		componentCountersCalculator:     componentCountersCalculator,
		eventsSender:                    eventsSender,
		logger:                          logger,
	}
}

type contextUpdateProcessor struct {
	dbClient                        mongo.DbClient
	alarmCollection                 mongo.DbCollection
	entityServiceCountersCalculator calculator.EntityServiceCountersCalculator
	componentCountersCalculator     calculator.ComponentCountersCalculator
	eventsSender                    entitycounters.EventsSender
	logger                          zerolog.Logger
}

func (p *contextUpdateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	entity := *event.Entity
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

	go p.postProcess(context.WithoutCancel(ctx), event, updatedServiceStates, componentStateChanged, newComponentState)

	return result, nil
}

func (p *contextUpdateProcessor) postProcess(
	ctx context.Context,
	event rpc.AxeEvent,
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
}
