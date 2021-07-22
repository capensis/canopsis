package config

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Adapter to the concrete database
type Adapter interface {
	// GetConfig return the whole config document
	GetConfig(ctx context.Context) (CanopsisConf, error)
	// UpsertConfig upsert a config in mongo.
	UpsertConfig(ctx context.Context, conf CanopsisConf) error
}

type adapter struct {
	collection mongo.DbCollection
}

func NewAdapter(client mongo.DbClient) Adapter {
	return &adapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (c *adapter) GetConfig(ctx context.Context) (CanopsisConf, error) {
	conf := CanopsisConf{}
	err := c.collection.FindOne(ctx, bson.M{"_id": ConfigKeyName}).Decode(&conf)

	return conf, err
}

func (c *adapter) UpsertConfig(ctx context.Context, conf CanopsisConf) error {
	conf.ID = ConfigKeyName
	_, err := c.collection.UpdateOne(ctx, bson.M{"_id": ConfigKeyName},
		bson.M{"$set": conf}, options.Update().SetUpsert(true))

	return err
}

type UserInterfaceAdapter interface {
	GetConfig(ctx context.Context) (UserInterfaceConf, error)
}

type userInterfaceAdapter struct {
	collection mongo.DbCollection
}

func NewUserInterfaceAdapter(client mongo.DbClient) UserInterfaceAdapter {
	return &userInterfaceAdapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *userInterfaceAdapter) GetConfig(ctx context.Context) (UserInterfaceConf, error) {
	conf := UserInterfaceConf{}
	err := a.collection.FindOne(ctx, bson.M{"_id": UserInterfaceKeyName}).Decode(&conf)

	return conf, err
}

type RemediationAdapter interface {
	GetConfig(ctx context.Context) (RemediationConf, error)
	UpsertConfig(ctx context.Context, conf RemediationConf) error
}

type remediationAdapter struct {
	collection mongo.DbCollection
}

func NewRemediationAdapter(client mongo.DbClient) RemediationAdapter {
	return &remediationAdapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *remediationAdapter) GetConfig(ctx context.Context) (RemediationConf, error) {
	conf := RemediationConf{}
	err := a.collection.FindOne(ctx, bson.M{"_id": RemediationKeyName}).Decode(&conf)

	return conf, err
}

func (a *remediationAdapter) UpsertConfig(ctx context.Context, conf RemediationConf) error {
	_, err := a.collection.UpdateOne(ctx, bson.M{"_id": RemediationKeyName},
		bson.M{"$set": conf}, options.Update().SetUpsert(true))

	return err
}
