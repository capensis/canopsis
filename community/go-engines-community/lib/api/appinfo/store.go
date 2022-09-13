package appinfo

import (
	"context"
	"errors"
	"sort"

	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPopupInterval = 3 //seconds

type Store interface {
	RetrieveLoginConfig(ctx context.Context) (LoginConf, error)
	RetrieveUserInterfaceConfig(ctx context.Context) (UserInterfaceConf, error)
	RetrieveVersionConfig(ctx context.Context) (VersionConf, error)
	RetrieveGlobalConfig(ctx context.Context) (GlobalConf, error)
	RetrieveRemediationConfig(ctx context.Context) (RemediationConf, error)
	UpdateUserInterfaceConfig(ctx context.Context, conf *UserInterfaceConf) error
	DeleteUserInterfaceConfig(ctx context.Context) error
	RetrieveApiSecurityConfig(ctx context.Context) (map[string]apisecurity.AuthMethodConf, error)
	UpdateApiSecurityConfig(ctx context.Context, conf map[string]apisecurity.AuthMethodConf) (map[string]apisecurity.AuthMethodConf, error)
}

type store struct {
	dbClient         mongo.DbClient
	objectCollection mongo.DbCollection
	configCollection mongo.DbCollection
	authProviders    []string
}

// NewStore instantiates configuration store.
func NewStore(db mongo.DbClient, authProviders []string) Store {
	return &store{
		dbClient:         db,
		objectCollection: db.Collection(mongo.ObjectMongoCollection),
		configCollection: db.Collection(mongo.ConfigurationMongoCollection),
		authProviders:    authProviders,
	}
}

func (s *store) RetrieveLoginConfig(ctx context.Context) (LoginConf, error) {
	var login = LoginConf{}
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

func (s *store) RetrieveGlobalConfig(ctx context.Context) (GlobalConf, error) {
	conf := config.CanopsisConf{}
	err := s.configCollection.FindOne(ctx, bson.M{"_id": config.ConfigKeyName}).Decode(&conf)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return GlobalConf{}, nil
		}

		return GlobalConf{}, err
	}

	return GlobalConf{
		Timezone:          conf.Timezone.Timezone,
		FileUploadMaxSize: conf.File.UploadMaxSize,
	}, nil
}

func (s *store) RetrieveRemediationConfig(ctx context.Context) (RemediationConf, error) {
	conf := config.RemediationConf{}
	result := RemediationConf{}
	err := s.configCollection.FindOne(ctx, bson.M{"_id": config.RemediationKeyName}).Decode(&conf)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return result, nil
		}

		return result, err
	}

	result.JobConfigTypes = make([]JobConfigType, len(conf.ExternalAPI))
	i := 0
	for name, apiConfig := range conf.ExternalAPI {
		result.JobConfigTypes[i] = JobConfigType{
			Name:     name,
			AuthType: apiConfig.Auth.Type,
		}
		i++
	}

	sort.Slice(result.JobConfigTypes, func(i, j int) bool {
		return result.JobConfigTypes[i].Name < result.JobConfigTypes[j].Name
	})

	return result, nil
}

func (s *store) UpdateUserInterfaceConfig(ctx context.Context, model *UserInterfaceConf) error {
	defaultInterval := types.DurationWithUnit{
		Value: defaultPopupInterval,
		Unit:  "s",
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

	var updatedModel UserInterfaceConf
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedModel = UserInterfaceConf{}
		_, err := s.configCollection.UpdateOne(ctx, bson.M{"_id": config.UserInterfaceKeyName},
			bson.M{"$set": model}, options.Update().SetUpsert(true))

		if err != nil {
			return err
		}

		updatedModel, err = s.RetrieveUserInterfaceConfig(ctx)
		return err
	})
	if err != nil {
		return err
	}
	*model = updatedModel

	return nil
}

func (s *store) DeleteUserInterfaceConfig(ctx context.Context) error {
	_, err := s.configCollection.DeleteOne(ctx, bson.M{"_id": config.UserInterfaceKeyName})
	return err
}

func (s *store) RetrieveApiSecurityConfig(ctx context.Context) (map[string]apisecurity.AuthMethodConf, error) {
	conf := make(map[string]apisecurity.AuthMethodConf)
	err := s.configCollection.
		FindOne(ctx, bson.M{"_id": config.ApiSecurityKeyName}, options.FindOne().SetProjection(bson.M{"_id": 0})).
		Decode(&conf)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, err
	}

	for _, provider := range s.authProviders {
		if provider == security.AuthMethodApiKey {
			continue
		}
		providerConf := conf[provider]
		if providerConf.InactivityInterval == nil {
			providerConf.InactivityInterval = &types.DurationWithUnit{
				Value: config.ApiSecurityInactivityInterval,
				Unit:  "h",
			}
		}
		conf[provider] = providerConf
	}

	return conf, nil
}

func (s *store) UpdateApiSecurityConfig(ctx context.Context, conf map[string]apisecurity.AuthMethodConf) (map[string]apisecurity.AuthMethodConf, error) {
	set := bson.M{}
	unset := bson.M{}
	for _, provider := range s.authProviders {
		if provider == security.AuthMethodApiKey {
			continue
		}
		providerConf := conf[provider]
		if providerConf.ExpirationInterval != nil || providerConf.InactivityInterval != nil {
			set[provider] = providerConf
		} else {
			unset[provider] = ""
		}
	}
	update := bson.M{}
	if len(set) > 0 {
		update["$set"] = set
	}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	var result map[string]apisecurity.AuthMethodConf
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		result = nil

		_, err := s.configCollection.UpdateOne(ctx, bson.M{"_id": config.ApiSecurityKeyName}, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}

		result, err = s.RetrieveApiSecurityConfig(ctx)
		return err
	})

	return result, err
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
