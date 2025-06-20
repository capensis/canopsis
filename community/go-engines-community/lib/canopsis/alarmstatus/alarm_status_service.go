package alarmstatus

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/alarmstatus/alarmstatus.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus Service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service interface {
	Load(ctx context.Context) error
	ComputeStatusOnStatusChange(ctx context.Context, alarm types.Alarm, entity types.Entity) (status types.CpsNumber, msg string, err error)
	ComputeStatusOnStateChange(alarm types.Alarm, entity types.Entity) (status types.CpsNumber, msg string)
}

func NewService(
	dbClient mongo.DbClient,
	flappingRuleAdapter flappingrule.Adapter,
	configProvider config.AlarmConfigProvider,
	logger zerolog.Logger,
) Service {
	return &service{
		flappingRuleAdapter: flappingRuleAdapter,
		configProvider:      configProvider,
		alarmCollection:     dbClient.Collection(mongo.AlarmMongoCollection),
		logger:              logger,
	}
}

type service struct {
	flappingRuleAdapter flappingrule.Adapter
	configProvider      config.AlarmConfigProvider
	alarmCollection     mongo.DbCollection
	logger              zerolog.Logger

	flappingRulesMx sync.RWMutex
	flappingRules   []flappingrule.Rule
}

func (s *service) Load(ctx context.Context) error {
	s.flappingRulesMx.Lock()
	defer s.flappingRulesMx.Unlock()

	rules, err := s.flappingRuleAdapter.Get(ctx)
	if err != nil {
		return err
	}

	ids := make([]string, len(rules))
	for i, rule := range rules {
		ids[i] = rule.ID
	}
	s.logger.Debug().Strs("rules", ids).Msg("load flapping rules")

	s.flappingRules = rules
	return nil
}

func (s *service) ComputeStatusOnStatusChange(ctx context.Context, alarm types.Alarm, entity types.Entity) (types.CpsNumber, string, error) {
	if alarm.Value.Canceled != nil {
		return types.AlarmStatusCancelled, "", nil
	}

	isUpstreamOK, err := s.isUpstreamOK(ctx, entity.Upstream)
	if err != nil {
		return 0, "", err
	}

	if !isUpstreamOK {
		return types.AlarmStatusUnknown, types.OutputUpstreamPrefix + entity.Upstream, nil
	}

	if alarm.Value.NoEventsDate != nil {
		return types.AlarmStatusNoEvents, "", nil
	}

	if isFlapping, msg := s.isFlapping(alarm, entity); isFlapping {
		return types.AlarmStatusFlapping, msg, nil
	}

	if s.isStealthy(alarm) {
		return types.AlarmStatusStealthy, "", nil
	}

	if alarm.Value.State != nil && alarm.Value.State.Value != types.AlarmStateOK {
		return types.AlarmStatusOngoing, "", nil
	}

	return types.AlarmStatusOff, "", nil
}

func (s *service) ComputeStatusOnStateChange(alarm types.Alarm, entity types.Entity) (types.CpsNumber, string) {
	if alarm.Value.Status != nil && alarm.Value.Status.Value == types.AlarmStatusCancelled {
		return types.AlarmStatusCancelled, ""
	}

	if alarm.Value.Status != nil && alarm.Value.Status.Value == types.AlarmStatusUnknown {
		return types.AlarmStatusUnknown, ""
	}

	if alarm.Value.NoEventsDate != nil {
		return types.AlarmStatusNoEvents, ""
	}

	if isFlapping, msg := s.isFlapping(alarm, entity); isFlapping {
		return types.AlarmStatusFlapping, msg
	}

	if s.isStealthy(alarm) {
		return types.AlarmStatusStealthy, ""
	}

	if alarm.Value.State != nil && alarm.Value.State.Value != types.AlarmStateOK {
		return types.AlarmStatusOngoing, ""
	}

	return types.AlarmStatusOff, ""
}

func (s *service) isFlapping(alarm types.Alarm, entity types.Entity) (bool, string) {
	s.flappingRulesMx.RLock()
	defer s.flappingRulesMx.RUnlock()

	now := datetime.NewCpsTime()
	alarmWithEntity := types.AlarmWithEntity{
		Alarm:  alarm,
		Entity: entity,
	}
	lastStepType := ""
	freq := 0
	for _, rule := range s.flappingRules {
		matched, err := rule.Matches(alarmWithEntity)
		if err != nil {
			s.logger.Error().Err(err).Str("flapping_rule", rule.ID).Msg("match flapping rule returned error, skip")
			continue
		}

		if matched {
			before := rule.Duration.SubFrom(now)

			for i := len(alarm.Value.Steps) - 1; i >= 0; i-- {
				step := alarm.Value.Steps[i]

				if step.Timestamp.Before(before) {
					break
				}

				if step.Type != lastStepType {
					switch step.Type {
					case types.AlarmStepStateIncrease, types.AlarmStepStateDecrease:
						lastStepType = step.Type
						freq++
					}
				}

				if freq > rule.FreqLimit {
					return true, types.RuleNameRulePrefix + rule.Name
				}
			}

			break
		}
	}

	return false, ""
}

func (s *service) isStealthy(alarm types.Alarm) bool {
	interval := s.configProvider.Get().StealthyInterval

	for i := len(alarm.Value.Steps) - 1; i >= 0; i-- {
		step := alarm.Value.Steps[i]
		if time.Since(step.Timestamp.Time) >= interval {
			break
		}

		if step.Value != types.AlarmStateOK {
			switch step.Type {
			case types.AlarmStepStatusIncrease, types.AlarmStepStateDecrease:
				return true
			default:
				break
			}
		}
	}

	return false
}

func (s *service) isUpstreamOK(ctx context.Context, upstream string) (bool, error) {
	if upstream == "" {
		return true, nil
	}

	cursor, err := s.alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"d":          upstream,
			"v.resolved": nil,
			"v.meta":     nil,
			"$and": []bson.M{
				{
					"v.status.val": bson.M{"$ne": types.AlarmStatusOff},
				},
				{"$or": []bson.M{
					{"v.state.val": bson.M{"$ne": types.AlarmStateOK}},
					{"v.status.val": bson.M{"$ne": types.AlarmStatusCancelled}},
				}},
			},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
			"pipeline": []bson.M{
				{"$match": bson.M{
					"enabled": true,
					"type":    bson.M{"$in": bson.A{types.EntityTypeResource, types.EntityTypeComponent}},
				}},
				{"$limit": 1},
			},
		}},
		{"$unwind": "$entity"},
		{"$project": bson.M{
			"_id": 1,
		}},
	})
	if err != nil {
		return false, fmt.Errorf("cannot find upstream alarm: %w", err)
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		return false, nil
	}

	if err = cursor.Err(); err != nil {
		return false, fmt.Errorf("cannot fetch upstream alarm: %w", err)
	}

	return true, nil
}
