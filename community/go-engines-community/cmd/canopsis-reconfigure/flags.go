package main

import (
	"flag"
	"fmt"
)

const (
	DefaultCfgFile = "/opt/canopsis/etc/canopsis.toml"

	DefaultMongoMigrationsPath        = "/opt/canopsis/share/database/migrations"
	DefaultMongoFixturesPath          = "/opt/canopsis/share/database/fixtures"
	DefaultPostgresMigrationsPath     = "/opt/canopsis/share/database/postgres_migrations"
	DefaultTechPostgresMigrationsPath = "/opt/canopsis/share/database/tech_postgres_migrations"

	EditionCommunity = "community"
	EditionPro       = "pro"

	MongoMigrationExecGoja    = "goja"
	MongoMigrationExecMongosh = "mongosh"
)

type flags struct {
	confFile         string
	overrideConfFile string

	version bool
	edition string

	modeDebug bool

	modeMigrateMongo        bool
	mongoMigrationDirectory string
	mongoMigrationExec      string
	mongoFixtureDirectory   string

	modeMigratePostgres        bool
	postgresMigrationDirectory string
	postgresMigrationMode      string
	postgresMigrationSteps     int

	modeMigrateTechPostgres        bool
	techPostgresMigrationDirectory string
	techPostgresMigrationMode      string
	techPostgresMigrationSteps     int

	mongoFixtureMigrations        bool
	mongoFixtureMigrationsVersion string
}

func (f *flags) Parse() {
	flag.StringVar(&f.confFile, "conf", DefaultCfgFile, "The configuration file used to initialize Canopsis")
	flag.StringVar(&f.overrideConfFile, "override", "", "The configuration file used to override default Canopsis configurations, for example /opt/canopsis/etc/conf.d/canopsis-override.toml")

	flag.BoolVar(&f.version, "version", false, "Show the version information")
	flag.StringVar(&f.edition, "edition", EditionCommunity, fmt.Sprintf("Canopsis edition: %s or %s", EditionCommunity, EditionPro))

	flag.BoolVar(&f.modeDebug, "d", false, "debug mode")

	flag.BoolVar(&f.modeMigrateMongo, "migrate-mongo", true, "If true, it will execute Mongo migration scripts")
	flag.StringVar(&f.mongoMigrationDirectory, "mongo-migration-directory", DefaultMongoMigrationsPath, "The directory with Mongo migration scripts")
	flag.StringVar(&f.mongoMigrationExec, "mongo-migration-exec", MongoMigrationExecGoja, "The execution of Mongo migration scripts: "+MongoMigrationExecGoja+" or "+MongoMigrationExecMongosh)

	flag.StringVar(&f.mongoFixtureDirectory, "mongo-fixture-directory", DefaultMongoFixturesPath, "The directory with Mongo fixtures")
	flag.BoolVar(&f.mongoFixtureMigrations, "mongo-fixture-migrations", false, "If true, it will fill migration collection with migration versions without executing them during mongo fixtures loading")
	flag.StringVar(&f.mongoFixtureMigrationsVersion, "mongo-fixture-migrations-version", "", "The max migration version to be inserted to migration collection during mongo fixtures loading")

	flag.BoolVar(&f.modeMigratePostgres, "migrate-postgres", false, "If true, it will execute Postgres migration scripts")
	flag.StringVar(&f.postgresMigrationDirectory, "postgres-migration-directory", DefaultPostgresMigrationsPath, "The directory with Postgres migration scripts")
	flag.StringVar(&f.postgresMigrationMode, "postgres-migration-mode", "up", "Should be up or down")
	flag.IntVar(&f.postgresMigrationSteps, "postgres-migration-steps", 0, "Number of migration steps, will execute all migrations if empty or 0")

	flag.BoolVar(&f.modeMigrateTechPostgres, "migrate-tech-postgres", false, "If true, it will execute Tech Postgres migration scripts")
	flag.StringVar(&f.techPostgresMigrationDirectory, "tech-postgres-migration-directory", DefaultTechPostgresMigrationsPath, "The directory with Tech Postgres migration scripts")
	flag.StringVar(&f.techPostgresMigrationMode, "tech-postgres-migration-mode", "up", "Should be up or down")
	flag.IntVar(&f.techPostgresMigrationSteps, "tech-postgres-migration-steps", 0, "Number of migration steps, will execute all migrations if empty or 0")

	flag.Parse()
}
