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

	flag.BoolVar(&f.modeMigrateMongo, "migrate-mongo", true, "If true, it will execute Mongo migration scripts")
	flag.StringVar(&f.mongoMigrationDirectory, "mongo-migration-directory", DefaultMongoMigrationsPath, "The directory with Mongo migration scripts")
	flag.StringVar(&f.mongoFixtureDirectory, "mongo-fixture-directory", DefaultMongoFixturesPath, "The directory with Mongo fixtures")

	flag.BoolVar(&f.modeMigratePostgres, "migrate-postgres", false, "If true, it will execute Postgres migration scripts")
	flag.StringVar(&f.postgresMigrationDirectory, "postgres-migration-directory", DefaultPostgresMigrationsPath, "The directory with Postgres migration scripts")
	flag.StringVar(&f.postgresMigrationMode, "postgres-migration-mode", "up", "Should be up or down")
	flag.IntVar(&f.postgresMigrationSteps, "postgres-migration-steps", 0, "Number of migration steps, will execute all migrations if empty or 0")

	flag.Parse()
}
