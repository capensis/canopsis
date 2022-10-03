package config

type Exchange struct {
	Name       string                 `toml:"name"`
	Kind       string                 `toml:"kind"`
	Durable    bool                   `toml:"durable"`
	AutoDelete bool                   `toml:"autodelete"`
	Internal   bool                   `toml:"internal"`
	NoWait     bool                   `toml:"noWait"`
	Args       map[string]interface{} `toml:"args"`
}

type QueueBinding struct {
	Key      string `toml:"key"`
	Exchange string `toml:"exchange"`
	NoWait   bool   `toml:"noWait"`
	Args     map[string]interface{}
}

type Queue struct {
	Name       string                 `toml:"name"`
	Durable    bool                   `toml:"durable"`
	AutoDelete bool                   `toml:"autoDelete"`
	Exclusive  bool                   `toml:"exclusive"`
	NoWait     bool                   `toml:"noWait"`
	Bind       *QueueBinding          `toml:"bind"`
	Args       map[string]interface{} `toml:"args"`
}

type RabbitMQConf struct {
	Exchanges []Exchange `toml:"exchanges"`
	Queues    []Queue    `toml:"queues"`
}
