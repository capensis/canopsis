package config

type MaintenanceConf struct {
	Enabled     bool   `bson:"enabled"`
	BroadcastID string `bson:"broadcast_id"`
}
