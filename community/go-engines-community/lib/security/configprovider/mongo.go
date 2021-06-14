package configprovider

import (
	"context"
	"errors"
	"fmt"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoProvider provides config from mongodb.
type mongoProvider struct {
	dbClient libmongo.DbClient
}

// NewMongoProvider creates new config provider.
func NewMongoProvider(dbClient libmongo.DbClient) security.ConfigProvider {
	return &mongoProvider{dbClient: dbClient}
}

const ObjectCollection = "object"
const ldapConfigID = "cservice.ldapconfig"
const casConfigID = "cservice.casconfig"
const defaultLdapPort = 389

func (p *mongoProvider) LoadLdapConfig() (*security.LdapConfig, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	col := p.dbClient.Collection(ObjectCollection)
	res := col.FindOne(ctx, bson.M{"_id": ldapConfigID})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	config := security.LdapConfig{}
	err := res.Decode(&config)
	if err != nil {
		return nil, err
	}

	if config.Url == "" {
		if config.Host == "" {
			return nil, errors.New("invalid config, url or host must be defined")
		}

		if config.Port == 0 {
			config.Port = defaultLdapPort
		}

		config.Url = fmt.Sprintf("ldap://%s:%d", config.Host, config.Port)
	}

	return &config, nil
}

func (p *mongoProvider) LoadCasConfig() (*security.CasConfig, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	col := p.dbClient.Collection(ObjectCollection)
	res := col.FindOne(ctx, bson.M{"_id": casConfigID})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	config := security.CasConfig{}
	err := res.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
