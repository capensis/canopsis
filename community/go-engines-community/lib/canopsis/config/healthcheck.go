package config

type HealthCheckConf struct {
	EngineOrder []EngineOrder `toml:"engine_order" bson:"engine_order"`
}

type EngineOrder struct {
	From string `toml:"from" bson:"from"`
	To   string `toml:"to" bson:"to"`
}
