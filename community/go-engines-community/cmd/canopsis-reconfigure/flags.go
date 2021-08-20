package main

import (
	"flag"
)

type flags struct {
	confFile           string
	migrationDirectory string
	modeDebug          bool
}

func (f *flags) Parse() {
	flag.StringVar(&f.confFile, "conf", DefaultCfgFile, FlagUsageConf)
	flag.BoolVar(&f.modeDebug, "d", false, "debug mode")
	flag.StringVar(&f.migrationDirectory, "migration-directory", DefaultMigrationsPath, "The directory with migration scripts")

	flag.Parse()
}
