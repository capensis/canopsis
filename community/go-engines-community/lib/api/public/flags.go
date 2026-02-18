package api

import (
	"flag"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

const (
	defaultPort      = 8083
	defaultConfigDir = "/opt/canopsis/share/config"
)

func BindCmdFlags(o *Options) {
	log.BindCmdFlags(&o.Options)
	flag.BoolVar(&o.Version, "version", false, "Show the version information")
	flag.IntVar(&o.Port, "port", defaultPort, "Server port")
	flag.StringVar(&o.ConfigDir, "c", defaultConfigDir, "Configuration files directory")
	flag.BoolVar(&o.EnableDocs, "docs", false, "Set to enable Swagger docs")
	flag.BoolVar(&o.LogBody, "logBody", false, "Set to enable logging response and request bodies")
	flag.BoolVar(&o.LogBodyOnError, "logBodyOnError", false, "Set to enable logging response and request bodies in case of error")
}

type Options struct {
	log.Options

	Version    bool
	Port       int
	ConfigDir  string
	EnableDocs bool

	LogBody        bool
	LogBodyOnError bool
}
