package main

import "flag"

const defaultConfigPath = "/opt/canopsis/share/config/connector-centreon-downtime-adapter/config.yml"

type Flags struct {
	ConfigPath string
	Version    bool
	Debug      bool
}

func (f *Flags) ParseArgs() {
	flag.StringVar(&f.ConfigPath, "c", defaultConfigPath, "Configuration file path")
	flag.BoolVar(&f.Debug, "d", false, "Enable debug mode")
	flag.BoolVar(&f.Version, "version", false, "Show the version information")
	flag.Parse()
}
