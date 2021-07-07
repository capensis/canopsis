package bdd

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoClient represents utility struct which implements db steps to feature context.
type MongoClient struct {
	client libmongo.DbClient
}

// NewMongoClient creates new mongo client.
func NewMongoClient(db libmongo.DbClient) (*MongoClient, error) {
	var mongoClient MongoClient
	// TODO: change database to test
	mongoClient.client = db

	return &mongoClient, nil
}

/**
Step example:
	And an alarm test_post_resource/test_post_component should be in the db
*/
func (c *MongoClient) AlarmShouldBeInTheDb(eid string) error {
	var expectedAlarm types.Alarm
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res := c.client.Collection(alarm.AlarmCollectionName).FindOne(ctx, bson.M{"d": eid})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("couldn't find an alarm for eid = %s", eid)
		}

		return err
	}

	err := res.Decode(&expectedAlarm)
	if err != nil {
		return err
	}

	return nil
}

/**
Step example:
	And an entity test_post_resource/test_post_component should be in the db
*/
func (c *MongoClient) EntityShouldBeInTheDb(eid string) error {
	var expectedEntity types.Entity
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res := c.client.Collection(libmongo.EntityMongoCollection).FindOne(ctx, bson.M{"_id": eid})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("couldn't find an entity for eid = %s", eid)
		}

		return err
	}

	err := res.Decode(&expectedEntity)
	if err != nil {
		return err
	}

	return nil
}
