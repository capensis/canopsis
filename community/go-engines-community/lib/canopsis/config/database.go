package config

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

func SetDbClientRetry(dbClient mongo.DbClient, cfg CanopsisConf) {
	dbClient.SetRetry(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
}
