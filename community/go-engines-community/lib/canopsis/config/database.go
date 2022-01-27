package config

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
)

func SetDbClientRetry(dbClient mongo.DbClient, cfg CanopsisConf) {
	dbClient.SetRetry(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
}

func SetPgPoolRetry(pgPool postgres.Pool, cfg CanopsisConf) {
	pgPool.SetRetry(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
}
