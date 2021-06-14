package metaalarm

import (
	"context"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"time"
)

type valueGroupEntityCounter struct {
	redisClient *redis.Client
	dbClient    mongo.DbClient
	logger      zerolog.Logger
}

func NewValueGroupEntityCounter(client mongo.DbClient, redisClient *redis.Client, logger zerolog.Logger) ValueGroupEntityCounter {
	return &valueGroupEntityCounter{
		dbClient:    client,
		redisClient: redisClient,
		logger:      logger,
	}
}

func (c *valueGroupEntityCounter) CountTotalEntitiesAmount(ctx context.Context, rule Rule) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if rule.Type == RuleValueGroup && rule.Config.ThresholdRate != nil && len(rule.Config.ValuePaths) != 0 {
		var concatValues []string
		var andConditions []bson.M
		for _, valuePath := range rule.Config.ValuePaths {
			idx := strings.Index(valuePath, ".")
			// use only entity valuepath
			if idx != -1 && valuePath[:idx] == "entity" {
				andConditions = append(andConditions, bson.M{valuePath[idx+1:]: bson.M{"$ne": nil}})
				concatValues = append(concatValues, ".", "$"+valuePath[idx+1:])
			}
		}

		if len(concatValues) == 0 || len(andConditions) == 0 {
			return c.logAndReturn(errors.New("there are no entity valuePaths"), "")
		}

		var filter bson.M
		if rule.Config.TotalEntityPatterns == nil {
			filter = rule.Config.EntityPatterns.AsMongoDriverQuery()
		} else {
			filter = rule.Config.TotalEntityPatterns.AsMongoDriverQuery()
		}

		andConditions = append(andConditions, filter)

		cursor, err := c.dbClient.Collection(mongo.EntityMongoCollection).
			Aggregate(ctx, []bson.M{
				{
					"$match": bson.M{
						"$and": andConditions,
					},
				},
				{
					"$addFields": bson.M{
						"valuepath": bson.M{
							"$concat": concatValues[1:],
						},
					},
				},
				{
					"$group": bson.M{
						"_id": "$valuepath",
						"total": bson.M{
							"$sum": 1,
						},
					},
				},
			})

		if err != nil {
			return c.logAndReturn(err, "ValueGroupEntityCounter: Mongo aggregate operation failed")
		}

		var result struct {
			Name  string `bson:"_id"`
			Total int64  `bson:"total"`
		}

		pipe := c.redisClient.Pipeline()
		for i := 0; cursor.Next(ctx); i++ {
			err := cursor.Decode(&result)
			if err != nil {
				return c.logAndReturn(err, "ValueGroupEntityCounter: Failed to decode value group total result")
			}

			pipe.Set(ctx, fmt.Sprintf("%s-%s", rule.ID, result.Name), result.Total, time.Hour)

			if i%1000 == 0 && i != 0 {
				_, err := pipe.Exec(ctx)
				if err != nil {
					c.logger.Err(err).Msgf("Failed to count entities for the %s rule", rule.ID)

					return err
				}

				pipe = c.redisClient.Pipeline()
			}
		}

		_, err = pipe.Exec(ctx)
		if err != nil {
			c.logger.Err(err).Msgf("Failed to count entities for the %s rule", rule.ID)

			return err
		}
	}

	return nil
}

func (c *valueGroupEntityCounter) CountTotalEntitiesAmountForValuePaths(ctx context.Context, rule Rule, valuePathsMap map[string]string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if len(valuePathsMap) == 0 {
		return c.logAndReturn(errors.New("ValueGroupEntityCounter: valueMap is empty"), "")
	}

	var andConditions []bson.M
	var values []string

	for valuePath, value := range valuePathsMap {
		idx := strings.Index(valuePath, ".")
		if idx != -1 && valuePath[:idx] == "entity" {
			andConditions = append(andConditions, bson.M{valuePath[idx+1:]: value})
			values = append(values, value)
		}
	}

	if len(andConditions) == 0 || len(values) == 0 {
		return c.logAndReturn(errors.New("there are no entity valuePaths"), "")
	}

	var filter bson.M
	if rule.Config.TotalEntityPatterns == nil {
		filter = rule.Config.EntityPatterns.AsMongoDriverQuery()
	} else {
		filter = rule.Config.TotalEntityPatterns.AsMongoDriverQuery()
	}

	andConditions = append(andConditions, filter)

	cursor, err := c.dbClient.Collection(mongo.EntityMongoCollection).
		Aggregate(ctx, []bson.M{
			{
				"$match": bson.M{
					"$and": andConditions,
				},
			},
			{"$count": "total"},
		})

	if err != nil {
		return c.logAndReturn(err, "ValueGroupEntityCounter: Mongo aggregate operation failed")
	}

	var totalResult struct {
		Total int64 `bson:"total"`
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&totalResult)
		if err != nil {
			return c.logAndReturn(err, "ValueGroupEntityCounter: Failed to decode value group total result")
		}
	}

	setResult := c.redisClient.Set(ctx, fmt.Sprintf("%s-%s", rule.ID, strings.Join(values, ".")), totalResult.Total, time.Hour)
	if err := setResult.Err(); err != nil {
		c.logger.Err(err).Msgf("Failed to count entities for the %s rule and %s valuePath", rule.ID, strings.Join(values, "."))

		return err
	}

	return nil
}

func (c *valueGroupEntityCounter) GetTotalEntitiesAmount(ctx context.Context, ruleId string, valueGroup string) (int, error) {
	key := fmt.Sprintf("%s-%s", ruleId, valueGroup)
	result := c.redisClient.Get(ctx, key)
	if result.Err() != nil {
		return 0, result.Err()
	}

	return strconv.Atoi(result.Val())
}

func (c *valueGroupEntityCounter) logAndReturn(err error, message string) error {
	c.logger.Err(err).Msg(message)

	return err
}
