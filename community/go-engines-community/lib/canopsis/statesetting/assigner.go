package statesetting

//go:generate mockgen -destination=../../../mocks/lib/statesetting/assigner.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting Assigner

import (
	"context"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Assigner interface {
	AssignStateSetting(ctx context.Context, entity *types.Entity, commRegister mongo.CommandsRegister) (bool, error)
	LoadRules(ctx context.Context) error
}

type assigner struct {
	dbClient                 mongo.DbClient
	dbCollection             mongo.DbCollection
	entityCountersCollection mongo.DbCollection

	rulesMutex     sync.RWMutex
	componentRules []StateSetting
	serviceRules   []StateSetting

	ruleQuery        bson.M
	ruleQueryOptions *options.FindOptions

	logger zerolog.Logger
}

func NewService(dbClient mongo.DbClient, logger zerolog.Logger) Assigner {
	return &assigner{
		dbClient:                 dbClient,
		dbCollection:             dbClient.Collection(mongo.StateSettingsMongoCollection),
		entityCountersCollection: dbClient.Collection(mongo.EntityCountersCollection),

		componentRules: make([]StateSetting, 0),
		serviceRules:   make([]StateSetting, 0),
		rulesMutex:     sync.RWMutex{},

		ruleQuery: bson.M{
			"enabled": true,
			"method":  bson.M{"$in": bson.A{MethodInherited, MethodDependencies}},
		},
		ruleQueryOptions: options.Find().SetSort(bson.D{{Key: "priority", Value: 1}, {Key: "_id", Value: 1}}),

		logger: logger,
	}
}

func (a *assigner) LoadRules(ctx context.Context) error {
	cursor, err := a.dbCollection.Find(ctx, a.ruleQuery, a.ruleQueryOptions)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	componentRules := make([]StateSetting, 0, len(a.componentRules))
	serviceRules := make([]StateSetting, 0, len(a.serviceRules))

	componentRulesIDs := make([]string, 0, len(a.componentRules))
	serviceRulesIDs := make([]string, 0, len(a.serviceRules))

	for cursor.Next(ctx) {
		var rule StateSetting

		err = cursor.Decode(&rule)
		if err != nil {
			return err
		}

		if rule.Type == RuleTypeComponent {
			componentRules = append(componentRules, rule)
			componentRulesIDs = append(componentRulesIDs, rule.ID)
		} else if rule.Type == RuleTypeService {
			serviceRules = append(serviceRules, rule)
			serviceRulesIDs = append(serviceRulesIDs, rule.ID)
		}
	}

	a.rulesMutex.Lock()

	a.componentRules = componentRules
	a.serviceRules = serviceRules

	a.rulesMutex.Unlock()

	a.logger.
		Debug().
		Strs("component_rules", componentRulesIDs).
		Strs("service_rules", serviceRulesIDs).
		Msg("Loading state settings rules")

	return nil
}

func (a *assigner) AssignStateSetting(ctx context.Context, entity *types.Entity, commRegister mongo.CommandsRegister) (bool, error) {
	if entity.Type != types.EntityTypeComponent && entity.Type != types.EntityTypeService {
		return false, nil
	}

	prevStateMethodID := ""
	if entity.StateInfo != nil {
		prevStateMethodID = entity.StateInfo.ID
	}

	a.rulesMutex.RLock()
	defer a.rulesMutex.RUnlock()

	if entity.Type == types.EntityTypeComponent {
		return a.assignToComponent(ctx, entity, prevStateMethodID, commRegister)
	} else if entity.Type == types.EntityTypeService {
		return a.assignToService(ctx, entity, prevStateMethodID, commRegister)
	}

	return false, nil
}

func (a *assigner) assignToComponent(ctx context.Context, entity *types.Entity, prevStateMethodID string, commRegister mongo.CommandsRegister) (bool, error) {
	for idx := range a.componentRules {
		if a.componentRules[idx].EntityPattern == nil {
			continue
		}

		matched, err := match.MatchEntityPattern(*a.componentRules[idx].EntityPattern, entity)
		if err != nil {
			return false, err
		}

		if matched {
			// for component, save inherited pattern to a state info to match resources easily.
			entity.StateInfo = &types.StateInfo{
				ID:               a.componentRules[idx].ID,
				InheritedPattern: a.componentRules[idx].InheritedEntityPattern,
			}

			// save rule to a corresponding state counter document to get it in the engine-axe on state calculation.
			_, err := a.entityCountersCollection.UpdateOne(
				ctx,
				bson.M{"_id": entity.ID},
				bson.M{"$set": bson.M{"rule": a.componentRules[idx]}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return false, err
			}

			commRegister.RegisterUpdate(entity.ID, bson.M{"state_info": entity.StateInfo})

			return true, nil
		}
	}

	// if we're here then no rule was matched, set it to nil.
	entity.StateInfo = nil
	if prevStateMethodID != "" {
		// if a component doesn't have state setting anymore, we should delete its counters, because we don't need them anymore.
		_, err := a.entityCountersCollection.DeleteOne(ctx, bson.M{"_id": entity.ID})
		if err != nil {
			return false, err
		}

		commRegister.RegisterUpdate(entity.ID, bson.M{"state_info": nil})
	}

	return false, nil
}

func (a *assigner) assignToService(ctx context.Context, entity *types.Entity, prevStateMethodID string, commRegister mongo.CommandsRegister) (bool, error) {
	for idx := range a.serviceRules {
		if a.serviceRules[idx].EntityPattern == nil {
			continue
		}

		matched, err := match.MatchEntityPattern(*a.serviceRules[idx].EntityPattern, entity)
		if err != nil {
			return false, err
		}

		if matched {
			// for service, save only rule's id, there is no need to save inherited pattern,
			// because resources are matched to a service with service's pattern and inherited pattern
			// will be used in axe for state calculation.
			entity.StateInfo = &types.StateInfo{ID: a.serviceRules[idx].ID}

			// save rule to a corresponding state counter document to get it in the engine-axe on state calculation.
			_, err := a.entityCountersCollection.UpdateOne(
				ctx,
				bson.M{"_id": entity.ID},
				bson.M{"$set": bson.M{"rule": a.serviceRules[idx]}},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return false, err
			}

			commRegister.RegisterUpdate(entity.ID, bson.M{"state_info": entity.StateInfo})

			return true, nil
		}
	}

	// if we're here then no rule was matched, set it to nil.
	entity.StateInfo = nil
	if prevStateMethodID != "" {
		// if a service doesn't have state setting anymore, we should delete only rule field, other counters are needed.
		_, err := a.entityCountersCollection.UpdateOne(
			ctx,
			bson.M{"_id": entity.ID},
			bson.M{"$unset": bson.M{"rule": ""}},
		)

		if err != nil {
			return false, err
		}

		commRegister.RegisterUpdate(entity.ID, bson.M{"state_info": nil})
	}

	return false, nil
}
