package metaalarm

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"sync"
	"time"
)

type ruleEntityCounter struct {
	countMutex    sync.RWMutex
	redisClient   *redis.Client
	entityAdapter entity.Adapter
	logger        zerolog.Logger
}

func (c *ruleEntityCounter) CountTotalEntitiesAmount(ctx context.Context, rule Rule) error {
	if rule.Type == RuleTypeComplex && rule.Config.ThresholdRate != nil {
		c.countMutex.Lock()
		defer c.countMutex.Unlock()

		var ids []interface{}

		var filter bson.M
		if rule.Config.TotalEntityPatterns == nil {
			filter = rule.Config.EntityPatterns.AsMongoDriverQuery()
		} else {
			filter = rule.Config.TotalEntityPatterns.AsMongoDriverQuery()
		}

		err := c.entityAdapter.GetIDs(filter, &ids)
		if err != nil {
			c.logger.Err(err).Msgf("Failed to count entities for the %s rule", rule.ID)

			return err
		}

		result := c.redisClient.Set(ctx, rule.ID, len(ids), time.Hour)
		if err := result.Err(); err != nil {
			c.logger.Err(err).Msgf("Failed to count entities for the %s rule", rule.ID)

			return err
		}

		c.logger.
			Debug().
			Msgf("Total entities matched the %s rule = %d", rule.ID, len(ids))
	}

	return nil
}

func (c *ruleEntityCounter) GetTotalEntitiesAmount(ctx context.Context, rule Rule) (int, error) {
	c.countMutex.Lock()
	defer c.countMutex.Unlock()

	result := c.redisClient.Get(ctx, rule.ID)
	if result.Err() != nil {
		return 0, result.Err()
	}

	return strconv.Atoi(result.Val())
}

func NewRuleEntityCounter(entityAdapter entity.Adapter, redisClient *redis.Client, logger zerolog.Logger) RuleEntityCounter {
	return &ruleEntityCounter{
		countMutex:    sync.RWMutex{},
		entityAdapter: entityAdapter,
		redisClient:   redisClient,
		logger:        logger,
	}
}
