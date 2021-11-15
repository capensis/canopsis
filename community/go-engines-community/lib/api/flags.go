package api

import "flag"

const (
	defaultPort      = 8082
	defaultConfigDir = "/opt/canopsis/share/config"
)

func (f *Flags) ParseArgs() {
	flag.Int64Var(&f.Port, "port", defaultPort, "Server port")
	flag.StringVar(&f.ConfigDir, "c", defaultConfigDir, "Configuration files directory")
	flag.BoolVar(&f.Debug, "d", false, "debug")
	flag.BoolVar(&f.SecureSession, "secure", false, "Secure session")
	flag.BoolVar(&f.Test, "test", false, "Set for functional tests")
	flag.BoolVar(&f.EnableDocs, "docs", false, "Set to enable Swagger docs")
	flag.Parse()
}

type Flags struct {
	Port          int64
	ConfigDir     string
	Debug         bool
	SecureSession bool
	Test          bool
	EnableDocs    bool
}
