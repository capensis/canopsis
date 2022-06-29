package main

import (
	"flag"
)

const (
	DefaultCfgFile       = "/opt/canopsis/etc/canopsis.toml"
	DefaultMongoConfPath = "/opt/canopsis/share/config/mongo"
)

func (f *flags) Parse() {
	flag.StringVar(&f.confFile, "conf", DefaultCfgFile, "The configuration file used to initialize Canopsis")
	flag.StringVar(&f.overrideConfFile, "override", "", "The configuration file used to override default Canopsis configurations, for example /opt/canopsis/etc/conf.d/canopsis-override.toml")
	flag.StringVar(&f.mongoConfPath, "mongoConf", DefaultMongoConfPath, "The configuration file path is used to create mongo indexes")
	flag.BoolVar(&f.modeDebug, "d", false, "debug mode")
	flag.BoolVar(&f.modeMigrateMongo, "migrate-mongo", false, "If true, it will execute mongo migration scripts")
	flag.BoolVar(&f.modeMigratePostgres, "migrate-postgres", false, "If true, it will execute postgres migration scripts")
	flag.StringVar(&f.mongoMigrationDirectory, "mongo-migration-directory", "", "The directory with migration scripts")
	flag.StringVar(&f.postgresMigrationDirectory, "postgres-migration-directory", "", "The directory with migration scripts")
	flag.StringVar(&f.postgresMigrationMode, "postgres-migration-mode", "", "should be up or down")
	flag.IntVar(&f.postgresMigrationSteps, "postgres-migration-steps", 0, "number of migration steps, will execute all migrations if empty or 0")
	flag.BoolVar(&f.modeMigrateOnly, "migrate-only", false, "If true, it will only execute migration scripts")
	flag.StringVar(&f.mongoContainer, "mongo-container", "", "Should contain docker container_id. If set, it will execute migration scripts inside the container")
	flag.StringVar(&f.mongoURL, "mongo-url", "", "mongo url")
	flag.Parse()
}

type flags struct {
	confFile                   string
	overrideConfFile           string
	mongoConfPath              string
	mongoMigrationDirectory    string
	modeDebug                  bool
	mongoContainer             string
	mongoURL                   string
	modeMigrateOnly            bool
	modeMigrateMongo           bool
	modeMigratePostgres        bool
	postgresMigrationDirectory string
	postgresMigrationMode      string
	postgresMigrationSteps     int
}
