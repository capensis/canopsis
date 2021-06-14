package metaalarm

import (
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog"
	"strconv"
	"strings"
	"time"
)

type valueGroupEntityCounter struct {
	redisClient *redis.Client
	dbClient 	*mgo.Session
	logger	 	zerolog.Logger
}

func NewValueGroupEntityCounter(client *mgo.Session, redisClient *redis.Client, logger zerolog.Logger) ValueGroupEntityCounter {
	return &valueGroupEntityCounter{
		dbClient: client,
		redisClient: redisClient,
		logger: logger,
	}
}

func (c *valueGroupEntityCounter) CountTotalEntitiesAmount(rule Rule) error {
	if rule.Type == RuleValueGroup && rule.Config.ThresholdRate != nil && len(rule.Config.ValuePaths) != 0 {
		var concatValues []string
		var andConditions []bson.M
		for _, valuePath := range rule.Config.ValuePaths {
			idx := strings.Index(valuePath, ".")
			// use only entity valuepath
			if idx != -1 && valuePath[:idx] == "entity" {
				andConditions = append(andConditions, bson.M{valuePath[idx + 1:]: bson.M{"$ne": nil}})
				concatValues = append(concatValues, ".", "$" + valuePath[idx + 1:])
			}
		}

		if len(concatValues) == 0 || len(andConditions) == 0 {
			return c.logAndReturn(errors.New("there are no entity valuePaths"), "")
		}

		var filter bson.M
		if rule.Config.TotalEntityPatterns == nil {
			filter = rule.Config.EntityPatterns.AsMongoQuery()
		} else {
			filter = rule.Config.TotalEntityPatterns.AsMongoQuery()
		}

		andConditions = append(andConditions, filter)

		cursor := c.
			dbClient.DB(canopsis.DbName).
			C(entity.EntityCollectionName).
			Pipe([]bson.M{
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
			},
		)

		iter := cursor.Iter()
		if iter.Err() != nil {
			return c.logAndReturn(iter.Err(), "ValueGroupEntityCounter: Mongo aggregate operation failed")
		}

		var result struct{
			Name	string	`bson:"_id"`
			Total	int64	`bson:"total"`
		}

		pipe := c.redisClient.Pipeline()
		for i := 0; iter.Next(&result); i++ {
			err := iter.Err()
			if err != nil {
				return c.logAndReturn(err, "ValueGroupEntityCounter: Failed to decode value group total result")
			}

			pipe.Set(fmt.Sprintf("%s-%s", rule.ID, result.Name), result.Total, time.Hour)

			if i % 1000 == 0 && i != 0 {
				_, err := pipe.Exec()
				if err != nil {
					c.logger.Err(err).Msgf("Failed to count entities for the %s rule", rule.ID)

					return err
				}

				pipe = c.redisClient.Pipeline()
			}
		}

		_, err := pipe.Exec()
		if err != nil {
			c.logger.Err(err).Msgf("Failed to count entities for the %s rule", rule.ID)

			return err
		}
	}

	return nil
}

func (c *valueGroupEntityCounter) CountTotalEntitiesAmountForValuePaths(rule Rule, valuePathsMap map[string]string) error {
	if len(valuePathsMap) == 0 {
		return c.logAndReturn(errors.New("ValueGroupEntityCounter: valueMap is empty"), "")
	}

	var andConditions []bson.M
	var values []string

	for valuePath, value := range valuePathsMap {
		idx := strings.Index(valuePath, ".")
		if idx != -1 && valuePath[:idx] == "entity" {
			andConditions = append(andConditions, bson.M{valuePath[idx + 1:]: value})
			values = append(values, value)
		}
	}

	if len(andConditions) == 0 || len(values) == 0 {
		return c.logAndReturn(errors.New("there are no entity valuePaths"), "")
	}

	var filter bson.M
	if rule.Config.TotalEntityPatterns == nil {
		filter = rule.Config.EntityPatterns.AsMongoQuery()
	} else {
		filter = rule.Config.TotalEntityPatterns.AsMongoQuery()
	}

	andConditions = append(andConditions, filter)

	cursor := c.
		dbClient.
		DB(canopsis.DbName).
		C(entity.EntityCollectionName).
		Pipe([]bson.M{
			{
				"$match": bson.M{
					"$and": andConditions,
				},
			},
			{"$count": "total"},
		},
	)

	iter := cursor.Iter()
	if iter.Err() != nil {
		return c.logAndReturn(iter.Err(), "ValueGroupEntityCounter: Mongo aggregate operation failed")
	}

	var totalResult struct{
		Total	int64	`bson:"total"`
	}

	iter.Next(&totalResult)
	err := iter.Err()
	if err != nil {
		return c.logAndReturn(err, "ValueGroupEntityCounter: Failed to decode value group total result")
	}

	setResult := c.redisClient.Set(fmt.Sprintf("%s-%s", rule.ID, strings.Join(values, ".")), totalResult.Total, time.Hour)
	if err := setResult.Err(); err != nil {
		c.logger.Err(err).Msgf("Failed to count entities for the %s rule and %s valuePath", rule.ID, strings.Join(values, "."))

		return err
	}

	return nil
}

func (c *valueGroupEntityCounter) GetTotalEntitiesAmount(ruleId string, valueGroup string) (int, error) {
	key := fmt.Sprintf("%s-%s", ruleId, valueGroup)
	result := c.redisClient.Get(key)
	if result.Err() != nil {
		return 0, result.Err()
	}

	return strconv.Atoi(result.Val())
}

func (c *valueGroupEntityCounter) logAndReturn(err error, message string) error {
	c.logger.Err(err).Msg(message)

	return err
}
