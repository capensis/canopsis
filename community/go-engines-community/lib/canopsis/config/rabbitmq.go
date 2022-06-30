package config

type Exchange struct {
	Name       string                 `toml:"name" mapstructure:"name"`
	Kind       string                 `toml:"kind" mapstructure:"kind"`
	Durable    bool                   `toml:"durable" mapstructure:"durable"`
	AutoDelete bool                   `toml:"autodelete" mapstructure:"autodelete"`
	Internal   bool                   `toml:"internal" mapstructure:"internal"`
	NoWait     bool                   `toml:"noWait" mapstructure:"noWait"`
	Args       map[string]interface{} `toml:"args" mapstructure:"args"`
}

type QueueBinding struct {
	Key      string `toml:"key"`
	Exchange string `toml:"exchange"`
	NoWait   bool   `toml:"noWait"`
	Args     map[string]interface{}
}

type Queue struct {
	Name       string                 `toml:"name" mapstructure:"name"`
	Durable    bool                   `toml:"durable" mapstructure:"durable"`
	AutoDelete bool                   `toml:"autoDelete" mapstructure:"autoDelete"`
	Exclusive  bool                   `toml:"exclusive" mapstructure:"exclusive"`
	NoWait     bool                   `toml:"noWait" mapstructure:"noWait"`
	Bind       *QueueBinding          `toml:"bind" mapstructure:"bind"`
	Args       map[string]interface{} `toml:"args" mapstructure:"args"`
}

type RabbitMQConf struct {
	Exchanges []Exchange `toml:"exchanges"`
	Queues    []Queue    `toml:"queues"`
}
