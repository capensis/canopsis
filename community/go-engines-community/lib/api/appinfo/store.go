package appinfo

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/colortheme"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPopupInterval = 3 //seconds

type Store interface {
	RetrieveLoginConfig() LoginConf
	RetrieveUserInterfaceConfig(ctx context.Context) (UserInterfaceConf, error)
	RetrieveVersionConfig(ctx context.Context) (VersionConf, error)
	RetrieveGlobalConfig(ctx context.Context) (GlobalConf, error)
	RetrieveRemediationConfig(ctx context.Context) (RemediationConf, error)
	UpdateUserInterfaceConfig(ctx context.Context, conf *UserInterfaceConf) error
	DeleteUserInterfaceConfig(ctx context.Context) error
	RetrieveMaintenanceState(ctx context.Context) (bool, error)
	RetrieveDefaultColorTheme(ctx context.Context) (colortheme.Theme, error)
	RetrieveSerialName(ctx context.Context) (string, error)
}

type store struct {
	dbClient             mongo.DbClient
	configCollection     mongo.DbCollection
	colorThemeCollection mongo.DbCollection
	maintenanceAdapter   config.MaintenanceAdapter
	pgPoolProvider       postgres.PoolProvider
	securityConfig       security.Config
}

// NewStore instantiates configuration store.
func NewStore(db mongo.DbClient, maintenanceAdapter config.MaintenanceAdapter, pgPoolProvider postgres.PoolProvider, securityConfig security.Config) Store {
	return &store{
		dbClient:             db,
		configCollection:     db.Collection(mongo.ConfigurationMongoCollection),
		colorThemeCollection: db.Collection(mongo.ColorThemeCollection),
		maintenanceAdapter:   maintenanceAdapter,
		pgPoolProvider:       pgPoolProvider,
		securityConfig:       securityConfig,
	}
}

func (s *store) RetrieveLoginConfig() LoginConf {
	var login = LoginConf{}
	for _, p := range s.securityConfig.Security.AuthProviders {
		switch p {
		case security.AuthMethodBasic:
			login.BasicConfig.Enable = true
		case security.AuthMethodLdap:
			login.LdapConfig.Enable = true
		case security.AuthMethodCas:
			login.CasConfig.Enable = true
			login.CasConfig.Title = s.securityConfig.Security.Cas.Title
		case security.AuthMethodSaml:
			login.SamlConfig.Enable = true
			login.SamlConfig.Title = s.securityConfig.Security.Saml.Title
		case security.AuthMethodOAuth2:
			providersLen := len(s.securityConfig.Security.OAuth2.Providers)
			if providersLen == 0 {
				continue
			}

			login.OAuth2Config.Enable = true
			login.OAuth2Config.Providers = make([]string, 0, providersLen)
			for name := range s.securityConfig.Security.OAuth2.Providers {
				login.OAuth2Config.Providers = append(login.OAuth2Config.Providers, name)
			}
		}
	}

	return login
}

func (s *store) RetrieveUserInterfaceConfig(ctx context.Context) (UserInterfaceConf, error) {
	filter := bson.M{
		"_id": config.UserInterfaceKeyName,
	}

	var conf UserInterfaceConf
	err := s.configCollection.FindOne(ctx, filter).Decode(&conf)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
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
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return version, nil
	}

	return version, err
}

func (s *store) RetrieveGlobalConfig(ctx context.Context) (GlobalConf, error) {
	conf := config.CanopsisConf{
		Global: config.SectionGlobal{
			EventsCountTriggerDefaultThreshold: config.DefaultEventsCountThreshold,
		},
	}
	err := s.configCollection.FindOne(ctx, bson.M{"_id": config.ConfigKeyName}).Decode(&conf)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return GlobalConf{}, nil
		}

		return GlobalConf{}, err
	}

	return GlobalConf{
		Timezone:                           conf.Timezone.Timezone,
		FileUploadMaxSize:                  conf.File.UploadMaxSize,
		EventsCountTriggerDefaultThreshold: conf.Global.EventsCountTriggerDefaultThreshold,
	}, nil
}

func (s *store) RetrieveRemediationConfig(ctx context.Context) (RemediationConf, error) {
	conf := config.RemediationConf{}
	result := RemediationConf{}
	err := s.configCollection.FindOne(ctx, bson.M{"_id": config.RemediationKeyName}).Decode(&conf)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	result.JobConfigTypes = make([]JobConfigType, len(conf.ExternalAPI))
	i := 0
	for name, apiConfig := range conf.ExternalAPI {
		result.JobConfigTypes[i] = JobConfigType{
			Name:      name,
			AuthType:  apiConfig.Auth.Type,
			WithBody:  apiConfig.LaunchEndpoint.WithBody,
			WithQuery: apiConfig.LaunchEndpoint.WithUrlQuery,
		}
		i++
	}

	sort.Slice(result.JobConfigTypes, func(i, j int) bool {
		return result.JobConfigTypes[i].Name < result.JobConfigTypes[j].Name
	})

	return result, nil
}

func (s *store) RetrieveMaintenanceState(ctx context.Context) (bool, error) {
	maintenanceConf, err := s.maintenanceAdapter.GetConfig(ctx)
	if err != nil {
		return false, err
	}

	return maintenanceConf.Enabled, nil
}

func (s *store) UpdateUserInterfaceConfig(ctx context.Context, model *UserInterfaceConf) error {
	defaultInterval := datetime.NewDurationWithUnit(defaultPopupInterval, datetime.DurationUnitSecond)

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

func (s *store) RetrieveDefaultColorTheme(ctx context.Context) (colortheme.Theme, error) {
	var t colortheme.Theme

	return t, s.colorThemeCollection.FindOne(ctx, bson.M{"_id": colortheme.Canopsis}).Decode(&t)
}

func (s *store) RetrieveSerialName(ctx context.Context) (string, error) {
	pool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get postgres pool: %w", err)
	}

	var serialName string
	err = pool.QueryRow(ctx, "SELECT id FROM serial_name LIMIT 1").Scan(&serialName)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("failed to get serial name: %w", err)
	}

	return serialName, nil
}
