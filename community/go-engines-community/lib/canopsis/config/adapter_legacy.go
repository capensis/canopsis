package config

import (
	cps "git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/globalsign/mgo"
)

type legacyAdapter struct {
	Collection mongo.Collection
}

// NewLegacyAdapter give the correct config adapter
func NewLegacyAdapter(session *mgo.Session) Adapter {
	return &legacyAdapter{
		Collection: mongo.FromMgo(session.DB(cps.DbName).C(mongo.ConfigurationMongoCollection)),
	}
}

func (c *legacyAdapter) GetConfig() (CanopsisConf, error) {
	conf := CanopsisConf{}
	err := c.Collection.GetByID(ConfigKeyName, &conf)

	return conf, err
}

func (c *legacyAdapter) UpsertConfig(conf CanopsisConf) error {
	conf.ID = ConfigKeyName

	return c.Collection.Upsert(ConfigKeyName, conf)
}
