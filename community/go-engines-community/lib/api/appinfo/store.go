package appinfo

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPopupInterval = 3 //seconds

type Store interface {
	RetrieveLoginConfig(ctx context.Context) (LoginConfig, error)
	RetrieveUserInterfaceConfig(ctx context.Context) (UserInterfaceConf, error)
	RetrieveVersionConfig(ctx context.Context) (VersionConf, error)
	RetrieveTimezoneConf(ctx context.Context) (TimezoneConf, error)
	RetrieveRemediationConf(ctx context.Context) (RemediationConf, error)
	UpdateUserInterfaceConfig(ctx context.Context, conf *UserInterfaceConf) error
	DeleteUserInterfaceConfig(ctx context.Context) error
}

type store struct {
	objectCollection mongo.DbCollection
	configCollection mongo.DbCollection
	authProviders    []string
}

// NewStore instantiates configuration store.
func NewStore(db mongo.DbClient, authProviders []string) Store {
	return &store{
		objectCollection: db.Collection(mongo.ObjectMongoCollection),
		configCollection: db.Collection(mongo.ConfigurationMongoCollection),
		authProviders:    authProviders,
	}
}

func (s *store) RetrieveLoginConfig(ctx context.Context) (LoginConfig, error) {
	var login = LoginConfig{}
	for _, p := range s.authProviders {
		switch p {
		case security.AuthMethodLdap:
			login.LdapConfig.Enable = true
		case security.AuthMethodCas:
			login.CasConfig.Enable = true
			var err error
			login.CasConfig.Title, err = s.findAuthMethodTitle(ctx, security.CasConfigID)
			if err != nil {
				return login, err
			}
		case security.AuthMethodSaml:
			login.SamlConfig.Enable = true
			var err error
			login.SamlConfig.Title, err = s.findAuthMethodTitle(ctx, security.SamlConfigID)
			if err != nil {
				return login, err
			}
		}
	}

	return login, nil
}

func (s *store) RetrieveUserInterfaceConfig(ctx context.Context) (UserInterfaceConf, error) {
	filter := bson.M{
		"_id": config.UserInterfaceKeyName,
	}

	var conf UserInterfaceConf
	err := s.configCollection.FindOne(ctx, filter).Decode(&conf)
	if err == mongodriver.ErrNoDocuments {
		return conf, nil
	}

	return conf, err
}

func (s *store) RetrieveVersionConfig(ctx context.Context) (VersionConf, error) {
	filter := bson.M{
		"_id": config.VersionKeyName,
	}

	var version VersionConf
	err := s.configCollection.FindOne(ctx, filter).Decode(&version)
	if err == mongodriver.ErrNoDocuments {
		return version, nil
	}

	return version, err
}

func (s *store) RetrieveTimezoneConf(ctx context.Context) (TimezoneConf, error) {
	var tz TimezoneConf
	conf := config.CanopsisConf{}
	err := s.configCollection.FindOne(ctx, bson.M{"_id": config.ConfigKeyName}).Decode(&conf)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return tz, nil
		}

		return tz, err
	}

	tz.Timezone = conf.Timezone.Timezone
	return tz, nil
}

func (s *store) RetrieveRemediationConf(ctx context.Context) (RemediationConf, error) {
	var remediation RemediationConf
	conf := config.CanopsisConf{}
	err := s.configCollection.FindOne(ctx, bson.M{"_id": config.ConfigKeyName}).Decode(&conf)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return remediation, nil
		}

		return remediation, err
	}

	remediation.JobExecutorFetchTimeoutSeconds = conf.Remediation.JobExecutorFetchTimeoutSeconds
	return remediation, nil
}

func (s *store) UpdateUserInterfaceConfig(ctx context.Context, model *UserInterfaceConf) error {
	defaultInterval := IntervalUnit{
		Interval: defaultPopupInterval,
		Unit:     "s",
	}

	if model.PopupTimeout == nil {
		model.PopupTimeout = &PopupTimeout{
			Info:  &defaultInterval,
			Error: &defaultInterval,
		}
	} else if model.PopupTimeout.Error == nil {
		model.PopupTimeout.Error = &defaultInterval
	} else if model.PopupTimeout.Info == nil {
		model.PopupTimeout.Info = &defaultInterval
	}

	_, err := s.configCollection.UpdateOne(ctx, bson.M{"_id": config.UserInterfaceKeyName},
		bson.M{"$set": model}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	updatedModel, err := s.RetrieveUserInterfaceConfig(ctx)
	if err != nil {
		return err
	}
	*model = updatedModel

	return err
}

func (s *store) DeleteUserInterfaceConfig(ctx context.Context) error {
	_, err := s.configCollection.DeleteOne(ctx, bson.M{"_id": config.UserInterfaceKeyName})
	return err
}

func (s *store) findAuthMethodTitle(ctx context.Context, id string) (string, error) {
	cfg := &struct {
		Title string `bson:"title"`
	}{}
	err := s.objectCollection.
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&cfg)
	if err != nil && err != mongodriver.ErrNoDocuments {
		return "", err
	}

	return cfg.Title, nil
}
