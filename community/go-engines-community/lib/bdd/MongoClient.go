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
func NewMongoClient(db libmongo.DbClient) *MongoClient {
	return &MongoClient{
		client: db,
	}
}

/*
*
Step example:

	And an alarm test_post_resource/test_post_component should be in the db
*/
func (c *MongoClient) AlarmShouldBeInTheDb(ctx context.Context, eid string) error {
	var expectedAlarm types.Alarm
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

/*
*
Step example:

	And an entity test_post_resource/test_post_component should be in the db
*/
func (c *MongoClient) EntityShouldBeInTheDb(ctx context.Context, eid string) error {
	var expectedEntity types.Entity
	err := c.client.Collection(libmongo.EntityMongoCollection).FindOne(ctx, bson.M{"_id": eid}).Decode(&expectedEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("couldn't find an entity for eid = %s", eid)
		}

		return err
	}

	return nil
}

/*
*
Step example:

	And an entity test_post_resource/test_post_component should not be in the db
*/
func (c *MongoClient) EntityShouldNotBeInTheDb(ctx context.Context, eid string) error {
	var expectedEntity types.Entity
	err := c.client.Collection(libmongo.EntityMongoCollection).FindOne(ctx, bson.M{"_id": eid}).Decode(&expectedEntity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		return err
	}

	return fmt.Errorf("could find an entity for eid = %s", eid)
}
