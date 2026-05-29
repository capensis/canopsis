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

	dbMongoName        = "MongoDB"
	dbPostgresName     = "Postgres (canopsis)"
	dbTechPostgresName = "Tech Postgres (tech_metrics)"

	flagMigratePostgres                      = "migrate-postgres"
	flagPostgresMigrationDirectory           = "postgres-migration-directory"
	flagPostgresMigrationMode                = "postgres-migration-mode"
	flagPostgresMigrationSteps               = "postgres-migration-steps"
	flagPostgresMigrationForceVersion        = "postgres-migration-force-version"
	flagPostgresMigrationForceClearDirty     = "postgres-migration-force-clear-dirty"
	flagPostgresMigrationUnsafe              = "postgres-migration-unsafe"
	flagMigrateTechPostgres                  = "migrate-tech-postgres"
	flagTechPostgresMigrationDirectory       = "tech-postgres-migration-directory"
	flagTechPostgresMigrationMode            = "tech-postgres-migration-mode"
	flagTechPostgresMigrationSteps           = "tech-postgres-migration-steps"
	flagTechPostgresMigrationForceVersion    = "tech-postgres-migration-force-version"
	flagTechPostgresMigrationForceClearDirty = "tech-postgres-migration-force-clear-dirty"
	flagTechPostgresMigrationUnsafe          = "tech-postgres-migration-unsafe"
	flagDiagnoseMigrations                   = "diagnose-migrations"
)

func flagRef(name string) string {
	return "-" + name
}

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

	modeMigratePostgres              bool
	postgresMigrationDirectory       string
	postgresMigrationMode            string
	postgresMigrationSteps           int
	postgresMigrationForceVersion    int
	postgresMigrationForceClearDirty bool

	modeMigrateTechPostgres              bool
	techPostgresMigrationDirectory       string
	techPostgresMigrationMode            string
	techPostgresMigrationSteps           int
	techPostgresMigrationForceVersion    int
	techPostgresMigrationForceClearDirty bool

	postgresUnsafeMigrations     bool
	techPostgresUnsafeMigrations bool
	diagnoseMigrations           bool

	mongoFixtureMigrations        bool
	mongoFixtureMigrationsVersion string

	forceGenerateSerialName bool
}

func (f *flags) Parse() error {
	const (
		migrateHelp            = "If true, it will execute %s migration scripts"
		migrationDirectoryHelp = "The directory with %s migration scripts"
		migrationModeHelp      = "The migration mode for %s migrations: up or down"
		migrationStepsHelp     = "Number of migration steps, will execute all migrations if empty or 0"
		unsafeMigrationsHelp   = "If true, bypass cross-line transition checks for %s postgres migrations (recovery only)"
		recoveryVersionHelp    = "Recovery only: force-set the migration version for %s and clear dirty state"
		recoveryClearDirtyHelp = "Recovery only: required together with %s to clear dirty state"
	)

	log.BindCmdFlags(&f.Options)

	flag.StringVar(&f.confFile, "conf", DefaultCfgFile, "The configuration file used to initialize Canopsis")
	flag.StringVar(&f.overrideConfFile, "override", DefaultOverrideCfgFile, "The configuration file used to override default Canopsis configurations, for example "+DefaultOverrideCfgFile)

	flag.BoolVar(&f.version, "version", false, "Show the version information")
	flag.StringVar(&f.edition, "edition", "", fmt.Sprintf("Canopsis edition: %s or %s", EditionCommunity, EditionPro))

	flag.BoolVar(&f.modeMigrateMongo, "migrate-mongo", true, fmt.Sprintf(migrateHelp, dbMongoName))
	flag.StringVar(&f.mongoMigrationDirectory, "mongo-migration-directory", DefaultMongoMigrationsPath, fmt.Sprintf(migrationDirectoryHelp, dbMongoName))
	flag.StringVar(&f.mongoMigrationExec, "mongo-migration-exec", MongoMigrationExecGoja, "The execution of Mongo migration scripts: "+MongoMigrationExecGoja+" or "+MongoMigrationExecMongosh)

	flag.StringVar(&f.mongoFixtureDirectory, "mongo-fixture-directory", DefaultMongoFixturesPath, fmt.Sprintf(migrationDirectoryHelp, dbMongoName))
	flag.BoolVar(&f.mongoFixtureMigrations, "mongo-fixture-migrations", false, fmt.Sprintf(migrateHelp, dbMongoName))
	flag.StringVar(&f.mongoFixtureMigrationsVersion, "mongo-fixture-migrations-version", "", "The max migration version to be inserted to migration collection during mongo fixtures loading")

	flag.BoolVar(&f.modeMigratePostgres, flagMigratePostgres, false, fmt.Sprintf(migrateHelp, dbPostgresName))
	flag.StringVar(&f.postgresMigrationDirectory, flagPostgresMigrationDirectory, DefaultPostgresMigrationsPath, fmt.Sprintf(migrationDirectoryHelp, dbPostgresName))
	flag.StringVar(&f.postgresMigrationMode, flagPostgresMigrationMode, "up", fmt.Sprintf(migrationModeHelp, dbPostgresName))
	flag.IntVar(&f.postgresMigrationSteps, flagPostgresMigrationSteps, 0, migrationStepsHelp)
	flag.IntVar(&f.postgresMigrationForceVersion, flagPostgresMigrationForceVersion, -1, fmt.Sprintf(recoveryVersionHelp, dbPostgresName))
	flag.BoolVar(&f.postgresMigrationForceClearDirty, flagPostgresMigrationForceClearDirty, false, fmt.Sprintf(recoveryClearDirtyHelp, flagRef(flagPostgresMigrationForceVersion)))

	flag.BoolVar(&f.modeMigrateTechPostgres, flagMigrateTechPostgres, false, fmt.Sprintf(migrateHelp, dbTechPostgresName))
	flag.StringVar(&f.techPostgresMigrationDirectory, flagTechPostgresMigrationDirectory, DefaultTechPostgresMigrationsPath, fmt.Sprintf(migrationDirectoryHelp, dbTechPostgresName))
	flag.StringVar(&f.techPostgresMigrationMode, flagTechPostgresMigrationMode, "up", fmt.Sprintf(migrationModeHelp, dbTechPostgresName))
	flag.IntVar(&f.techPostgresMigrationSteps, flagTechPostgresMigrationSteps, 0, migrationStepsHelp)
	flag.IntVar(&f.techPostgresMigrationForceVersion, flagTechPostgresMigrationForceVersion, -1, fmt.Sprintf(recoveryVersionHelp, dbTechPostgresName))
	flag.BoolVar(&f.techPostgresMigrationForceClearDirty, flagTechPostgresMigrationForceClearDirty, false, fmt.Sprintf(recoveryClearDirtyHelp, flagRef(flagTechPostgresMigrationForceVersion)))

	flag.BoolVar(&f.postgresUnsafeMigrations, flagPostgresMigrationUnsafe, false, fmt.Sprintf(unsafeMigrationsHelp, dbPostgresName))
	flag.BoolVar(&f.techPostgresUnsafeMigrations, flagTechPostgresMigrationUnsafe, false, fmt.Sprintf(unsafeMigrationsHelp, dbTechPostgresName))

	flag.BoolVar(&f.forceGenerateSerialName, "force-generate-serial-name", false, "If true, it will regenerate serial name even if it exists")
	flag.BoolVar(&f.diagnoseMigrations, flagDiagnoseMigrations, false, "If true, it will print non-mutating diagnostics about postgres migrations and exit when "+flagRef(flagMigratePostgres)+" or "+flagRef(flagMigrateTechPostgres)+" is set")

	flag.Parse()

	visited := map[string]bool{}
	flag.Visit(func(fl *flag.Flag) {
		visited[fl.Name] = true
	})

	if err := validateRecoveryRequest(
		dbPostgresName+" recovery",
		f.modeMigratePostgres,
		recoveryRequest{forceVersion: f.postgresMigrationForceVersion, forceClearDirty: f.postgresMigrationForceClearDirty},
		visited[flagPostgresMigrationMode],
		visited[flagPostgresMigrationSteps],
		f.diagnoseMigrations,
		f.postgresUnsafeMigrations,
		flagMigratePostgres,
		flagPostgresMigrationForceVersion,
		flagPostgresMigrationForceClearDirty,
		flagPostgresMigrationMode,
		flagPostgresMigrationSteps,
		flagPostgresMigrationUnsafe,
	); err != nil {
		return err
	}

	if err := validateRecoveryRequest(
		dbTechPostgresName+" recovery",
		f.modeMigrateTechPostgres,
		recoveryRequest{forceVersion: f.techPostgresMigrationForceVersion, forceClearDirty: f.techPostgresMigrationForceClearDirty},
		visited[flagTechPostgresMigrationMode],
		visited[flagTechPostgresMigrationSteps],
		f.diagnoseMigrations,
		f.techPostgresUnsafeMigrations,
		flagMigrateTechPostgres,
		flagTechPostgresMigrationForceVersion,
		flagTechPostgresMigrationForceClearDirty,
		flagTechPostgresMigrationMode,
		flagTechPostgresMigrationSteps,
		flagTechPostgresMigrationUnsafe,
	); err != nil {
		return err
	}

	return nil
}
