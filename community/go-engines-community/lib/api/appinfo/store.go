package appinfo

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ObjectCollection          = "object"
	UserInterfaceDocumentID   = "user_interface"
	CanopsisVersionDocumentId = "canopsis_version"
)

type Store interface {
	RetrieveObjectConfig() ([]LoginServiceConf, error)
	RetrieveUserInterfaceConfig() (UserInterfaceConf, error)
	RetrieveCanopsisVersionConfig() (CanopsisVersionConf, error)
	RetrieveTimezoneConf() (TimezoneConf, error)
	UpdateUserInterfaceConfig(*UserInterfaceConf) error
	DeleteUserInterfaceConfig() error
}

type store struct {
	db         mongo.DbClient
	configAdpt config.Adapter
}

// NewStore instantiates configuration store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		db:         db,
		configAdpt: config.NewAdapter(db),
	}
}

func (s *store) RetrieveObjectConfig() ([]LoginServiceConf, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{
		CrecordName: bson.M{"$in": bson.A{Casconfig, Ldapconfig}},
	}

	records := make([]LoginServiceConf, 0)
	cursor, err := s.db.Collection(ObjectCollection).Find(ctx, filter)
	if err == mongodriver.ErrNoDocuments {
		return records, nil
	}

	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var object LoginServiceConf
		err = cursor.Decode(&object)
		if err != nil {
			return nil, err
		}
		records = append(records, object)
	}

	return records, nil
}

func (s *store) RetrieveUserInterfaceConfig() (UserInterfaceConf, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{
		"_id": UserInterfaceDocumentID,
	}

	var conf UserInterfaceConf
	err := s.db.Collection(mongo.ConfigurationMongoCollection).FindOne(ctx, filter).Decode(&conf)
	if err == mongodriver.ErrNoDocuments {
		return conf, nil
	}

	return conf, err
}

func (s *store) RetrieveCanopsisVersionConfig() (CanopsisVersionConf, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filter := bson.M{
		"_id": CanopsisVersionDocumentId,
	}

	var version CanopsisVersionConf
	err := s.db.Collection(mongo.ConfigurationMongoCollection).FindOne(ctx, filter).Decode(&version)
	if err == mongodriver.ErrNoDocuments {
		return version, nil
	}

	return version, err
}

func (s *store) RetrieveTimezoneConf() (TimezoneConf, error) {
	var tz TimezoneConf
	conf, err := s.configAdpt.GetConfig()
	if err == mongodriver.ErrNoDocuments {
		return tz, nil
	}
	if err != nil {
		return tz, err
	}
	tz.Timezone = conf.Timezone.Timezone
	return tz, nil
}

func (s *store) UpdateUserInterfaceConfig(model *UserInterfaceConf) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var upsert = true
	_, err := s.db.Collection(mongo.ConfigurationMongoCollection).UpdateOne(ctx, bson.M{"_id": UserInterfaceDocumentID},
		bson.M{"$set": model}, &options.UpdateOptions{Upsert: &upsert})

	if err != nil {
		return err
	}

	updatedModel, err := s.RetrieveUserInterfaceConfig()
	if err != nil {
		return err
	}
	*model = updatedModel

	return err
}

func (s *store) DeleteUserInterfaceConfig() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := s.db.Collection(mongo.ConfigurationMongoCollection).DeleteOne(ctx, bson.M{"_id": UserInterfaceDocumentID})
	return err
}
