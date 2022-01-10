package main

import (
	"flag"
)

type flags struct {
	confFile  string
	modeDebug bool

	modeMigrateMongo        bool
	mongoMigrationDirectory string
	mongoFixtureDirectory   string

	modeMigratePostgres        bool
	postgresMigrationDirectory string
	postgresMigrationMode      string
	postgresMigrationSteps     int
}

func (f *flags) Parse() {
	flag.StringVar(&f.confFile, "conf", DefaultCfgFile, FlagUsageConf)
	flag.BoolVar(&f.modeDebug, "d", false, "debug mode")

	flag.BoolVar(&f.modeMigrateMongo, "migrate-mongo", false, "If true, it will execute mongo migration scripts")
	flag.StringVar(&f.mongoMigrationDirectory, "mongo-migration-directory", DefaultMigrationsPath, "The directory with migration scripts")
	flag.StringVar(&f.mongoFixtureDirectory, "mongo-fixture-directory", DefaultFixturesPath, "The directory with fixtures")

	flag.BoolVar(&f.modeMigratePostgres, "migrate-postgres", false, "If true, it will execute postgres migration scripts")
	flag.StringVar(&f.postgresMigrationDirectory, "postgres-migration-directory", "", "The directory with migration scripts")
	flag.StringVar(&f.postgresMigrationMode, "postgres-migration-mode", "", "should be up or down")
	flag.IntVar(&f.postgresMigrationSteps, "postgres-migration-steps", 0, "number of migration steps, will execute all migrations if empty or 0")

	flag.Parse()
}
