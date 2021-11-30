package main

import (
	"flag"
)

func (f *Flags) ParseArgs() {
	flag.Int64Var(&f.Port, "port", 8082, "Server port")
	flag.StringVar(&f.ConfigDir, "c", "/opt/canopsis/share/config", "Configuration files directory")
	flag.BoolVar(&f.Debug, "d", false, "debug")
	flag.BoolVar(&f.SecureSession, "secure", false, "Secure session")
	flag.BoolVar(&f.Test, "test", false, "Set for functional tests")
	flag.BoolVar(&f.EnableSameServiceNames, "enableSameServiceNames", false, "Enable same service names, services have unique names by default")
	flag.Parse()
}

type Flags struct {
	Port          int64
	ConfigDir     string
	Debug         bool
	SecureSession bool
	Test          bool
	// EnableSameServiceNames affects entityservice Create/Update payload validation
	EnableSameServiceNames bool
}
