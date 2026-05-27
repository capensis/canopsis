package main

import (
	"flag"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

const (
	DefaultCfgFile         = "/opt/canopsis/etc/canopsis.toml"
	DefaultOverrideCfgFile = "/opt/canopsis/etc/conf.d/canopsis-override.toml"

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
	log.Options

	confFile         string
	overrideConfFile string

	version bool
	edition string

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

	postgresUnsafeMigrations     bool
	techPostgresUnsafeMigrations bool
	diagnoseMigrations           bool

	mongoFixtureMigrations        bool
	mongoFixtureMigrationsVersion string

	forceGenerateSerialName bool
}

func (f *flags) Parse() {
	const (
		dbMongoName            = "MongoDB"
		dbPostgresName         = "Postgres (canopsis)"
		dbTechPostgresName     = "Tech Postgres (tech_metrics)"
		migrateHelp            = "If true, it will execute %s migration scripts"
		migrationDirectoryHelp = "The directory with %s migration scripts"
		migrationModeHelp      = "The migration mode for %s migrations: up or down"
		migrationStepsHelp     = "Number of migration steps, will execute all migrations if empty or 0"
		unsafeMigrationsHelp   = "If true, bypass cross-line transition checks for %s postgres migrations (recovery only)"
	)

	log.BindCmdFlags(&f.Options)

	flag.StringVar(&f.confFile, "conf", DefaultCfgFile, "The configuration file used to initialize Canopsis")
	flag.StringVar(&f.overrideConfFile, "override", DefaultOverrideCfgFile, "The configuration file used to override default Canopsis configurations, for example /opt/canopsis/etc/conf.d/canopsis-override.toml")

	flag.BoolVar(&f.version, "version", false, "Show the version information")
	flag.StringVar(&f.edition, "edition", "", fmt.Sprintf("Canopsis edition: %s or %s", EditionCommunity, EditionPro))

	flag.BoolVar(&f.modeMigrateMongo, "migrate-mongo", true, fmt.Sprintf(migrateHelp, dbMongoName))
	flag.StringVar(&f.mongoMigrationDirectory, "mongo-migration-directory", DefaultMongoMigrationsPath, fmt.Sprintf(migrationDirectoryHelp, dbMongoName))
	flag.StringVar(&f.mongoMigrationExec, "mongo-migration-exec", MongoMigrationExecGoja, "The execution of Mongo migration scripts: "+MongoMigrationExecGoja+" or "+MongoMigrationExecMongosh)

	flag.StringVar(&f.mongoFixtureDirectory, "mongo-fixture-directory", DefaultMongoFixturesPath, "The directory with Mongo fixtures")
	flag.BoolVar(&f.mongoFixtureMigrations, "mongo-fixture-migrations", false, "If true, it will fill migration collection with migration versions without executing them during mongo fixtures loading")
	flag.StringVar(&f.mongoFixtureMigrationsVersion, "mongo-fixture-migrations-version", "", "The max migration version to be inserted to migration collection during mongo fixtures loading")

	flag.BoolVar(&f.modeMigratePostgres, "migrate-postgres", false, fmt.Sprintf(migrateHelp, dbPostgresName))
	flag.StringVar(&f.postgresMigrationDirectory, "postgres-migration-directory", DefaultPostgresMigrationsPath, fmt.Sprintf(migrationDirectoryHelp, dbPostgresName))
	flag.StringVar(&f.postgresMigrationMode, "postgres-migration-mode", "up", fmt.Sprintf(migrationModeHelp, dbPostgresName))
	flag.IntVar(&f.postgresMigrationSteps, "postgres-migration-steps", 0, migrationStepsHelp)

	flag.BoolVar(&f.modeMigrateTechPostgres, "migrate-tech-postgres", false, fmt.Sprintf(migrateHelp, dbTechPostgresName))
	flag.StringVar(&f.techPostgresMigrationDirectory, "tech-postgres-migration-directory", DefaultTechPostgresMigrationsPath, fmt.Sprintf(migrationDirectoryHelp, dbTechPostgresName))
	flag.StringVar(&f.techPostgresMigrationMode, "tech-postgres-migration-mode", "up", fmt.Sprintf(migrationModeHelp, dbTechPostgresName))
	flag.IntVar(&f.techPostgresMigrationSteps, "tech-postgres-migration-steps", 0, migrationStepsHelp)

	flag.BoolVar(&f.postgresUnsafeMigrations, "postgres-migration-unsafe", false, fmt.Sprintf(unsafeMigrationsHelp, dbPostgresName))
	flag.BoolVar(&f.techPostgresUnsafeMigrations, "tech-postgres-migration-unsafe", false, fmt.Sprintf(unsafeMigrationsHelp, dbTechPostgresName))

	flag.BoolVar(&f.forceGenerateSerialName, "force-generate-serial-name", false, "If true, it will regenerate serial name even if it exists")
	flag.BoolVar(&f.diagnoseMigrations, "diagnose-migrations", false, "If true, it will print diagnostics about postgres migrations, print and exit when -migrate-postgres or -migrate-tech-postgres is set, and migration mode is not set")

	flag.Parse()
}
