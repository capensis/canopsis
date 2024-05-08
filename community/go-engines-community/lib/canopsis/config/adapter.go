package config

//go:generate mockgen -destination=../../../mocks/lib/canopsis/config/adapter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config MaintenanceAdapter

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
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

type HealthCheckAdapter interface {
	GetConfig(ctx context.Context) (HealthCheckConf, error)
	UpsertConfig(ctx context.Context, conf HealthCheckConf) error
}

type healthCheckAdapter struct {
	collection mongo.DbCollection
}

func NewHealthCheckAdapter(client mongo.DbClient) HealthCheckAdapter {
	return &healthCheckAdapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *healthCheckAdapter) GetConfig(ctx context.Context) (HealthCheckConf, error) {
	conf := HealthCheckConf{}
	err := a.collection.FindOne(ctx, bson.M{"_id": HealthCheckName}).Decode(&conf)

	return conf, err
}

func (a *healthCheckAdapter) UpsertConfig(ctx context.Context, conf HealthCheckConf) error {
	_, err := a.collection.UpdateOne(ctx, bson.M{"_id": HealthCheckName},
		bson.M{"$set": bson.M{
			"engine_order":    conf.EngineOrder,
			"update_interval": conf.UpdateInterval,
		}}, options.Update().SetUpsert(true))

	return err
}

type VersionAdapter interface {
	GetConfig(ctx context.Context) (VersionConf, error)
	UpsertConfig(ctx context.Context, conf VersionConf) error
}

type versionAdapter struct {
	collection mongo.DbCollection
}

func NewVersionAdapter(client mongo.DbClient) VersionAdapter {
	return &versionAdapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *versionAdapter) GetConfig(ctx context.Context) (VersionConf, error) {
	conf := VersionConf{}
	err := a.collection.FindOne(ctx, bson.M{"_id": VersionKeyName}).Decode(&conf)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return VersionConf{}, nil
		}
		return VersionConf{}, err
	}

	return conf, nil
}

func (a *versionAdapter) UpsertConfig(ctx context.Context, conf VersionConf) error {
	err := a.collection.FindOneAndUpdate(ctx, bson.M{"_id": VersionKeyName}, bson.M{"$set": conf},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)).Decode(&conf)
	if err != nil {
		return err
	}
	if conf.Edition == "" {
		return errors.New("edition is undefined")
	}
	return nil
}

type MaintenanceAdapter interface {
	GetConfig(ctx context.Context) (MaintenanceConf, error)
}

type maintenanceAdapter struct {
	collection mongo.DbCollection
}

func NewMaintenanceAdapter(client mongo.DbClient) MaintenanceAdapter {
	return &maintenanceAdapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *maintenanceAdapter) GetConfig(ctx context.Context) (MaintenanceConf, error) {
	conf := MaintenanceConf{}
	err := a.collection.FindOne(ctx, bson.M{"_id": MaintenanceKeyName}).Decode(&conf)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return conf, nil
		}

		return conf, err
	}

	return conf, nil
}
