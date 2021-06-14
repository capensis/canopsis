package config

import (
	"context"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter to the concrete database
type Adapter interface {
	// GetConfig return the whole config document
	GetConfig() (CanopsisConf, error)
	// UpsertConfig upsert a config in mongo.
	UpsertConfig(conf CanopsisConf) error
}

type adapter struct {
	collection mongo.DbCollection
}

func NewAdapter(client mongo.DbClient) Adapter {
	return &adapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (c *adapter) GetConfig() (CanopsisConf, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf := CanopsisConf{}
	err := c.collection.FindOne(ctx, bson.M{"_id": ConfigKeyName}).Decode(&conf)

	return conf, err
}

func (c *adapter) UpsertConfig(conf CanopsisConf) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf.ID = ConfigKeyName
	_, err := c.collection.UpdateOne(ctx, bson.M{"_id": ConfigKeyName},
		bson.M{"$set": conf}, options.Update().SetUpsert(true))

	return err
}

type UserInterfaceAdapter interface {
	GetConfig() (UserInterfaceConf, error)
}

type userInterfaceAdapter struct {
	collection mongo.DbCollection
}

func NewUserInterfaceAdapter(client mongo.DbClient) UserInterfaceAdapter {
	return &userInterfaceAdapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *userInterfaceAdapter) GetConfig() (UserInterfaceConf, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := UserInterfaceConf{}
	err := a.collection.FindOne(ctx, bson.M{"_id": UserInterfaceKeyName}).Decode(&conf)

	return conf, err
}
