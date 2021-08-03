package config

type HealthCheckConf struct {
	EngineOrder []EngineOrder         `toml:"engine_order" bson:"engine_order" json:"engine_order"`
	Parameters  HealthCheckParameters `toml:"parameters" bson:"parameters" json:"parameters"`
}

type EngineOrder struct {
	From string `toml:"from" bson:"from" json:"from"`
	To   string `toml:"to" bson:"to" json:"to"`
}

type EngineParameters struct {
	Minimal int `toml:"minimal" bson:"minimal" json:"minimal" binding:"required,gt=0"`
	Optimal int `toml:"optimal" bson:"optimal" json:"optimal" binding:"required,gtefield=Minimal"`
}

type HealthCheckParameters struct {
	QueueLimit int                         `toml:"queue_limit" bson:"queue_limit" json:"queue_limit" binding:"gt=0"`
	Engines    map[string]EngineParameters `toml:"engines" bson:"engines" json:"engines" binding:"dive"`
}
