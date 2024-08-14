package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/migration/cli"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo/goja"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo/mongosh"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pgxdriver "github.com/jackc/pgx/v5"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	f := flags{}
	f.Parse()

	if f.version {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(f.modeDebug)
	conf, err := parseConfig(f, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	err = GracefulStart(ctx, f.modeMigratePostgres, f.modeMigrateTechPostgres, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open one of required sessions")
	}

	err = initRabbitMQ(conf, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize rabbitmq")
	}

	client, err := mongo.NewClient(ctx, 0, 0, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to mongo")
	}
	defer func() {
		err = client.Disconnect(context.Background())
		if err != nil {
			logger.Err(err).Msg("failed to close mongo")
		}
	}()

	err = applyMongoFixtures(ctx, f, client, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to run mongo fixtures")
	}

	err = updateMongoConfig(ctx, conf, client)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to update config in mongo")
	}

	err = migrateMongo(ctx, f, client, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to run mongo migrations")
	}

	err = migratePostgres(f, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to run postgres migrations")
	}

	err = migrateTechPostgres(f, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to run tech postgres migrations")
	}

	// keep it in the end for cmd/ready
	err = updateVersionConfig(ctx, f, client)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to update config in mongo")
	}

	err = generateSerialName(ctx, f.forceGenerateSerialName)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to generate serial name")
	}
}

type Conf struct {
	RabbitMQ    config.RabbitMQConf    `toml:"RabbitMQ"`
	Canopsis    config.CanopsisConf    `toml:"Canopsis"`
	Remediation config.RemediationConf `toml:"Remediation"`
	HealthCheck config.HealthCheckConf `toml:"HealthCheck"`
}

func initRabbitMQ(conf Conf, logger zerolog.Logger) error {
	amqpConn, err := amqp.NewConnection(logger, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to open amqp: %w", err)
	}

	defer func() {
		err = amqpConn.Close()
		if err != nil {
			logger.Err(err).Msg("cannot close amqp session")
		}
	}()

	ch, err := amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open amqp channel: %w", err)
	}

	logger.Info().Msg("initialising rabbitmq exchanges")

	for _, exchange := range conf.RabbitMQ.Exchanges {
		err := ch.ExchangeDeclare(
			exchange.Name,
			exchange.Kind,
			exchange.Durable,
			exchange.AutoDelete,
			exchange.Internal,
			exchange.NoWait,
			exchange.Args,
		)
		if err != nil {
			return fmt.Errorf("cannot initialise exchange %q: %w", exchange.Name, err)
		}
		logger.Info().Msgf("exchange %q created", exchange.Name)
	}

	logger.Info().Msg("initialising rabbitmq queues")
	for _, queue := range conf.RabbitMQ.Queues {
		_, err := ch.QueueDeclare(
			queue.Name,
			queue.Durable,
			queue.AutoDelete,
			queue.Exclusive,
			queue.NoWait,
			queue.Args,
		)
		if err != nil {
			return fmt.Errorf("cannot initialise queue %q: %w", queue.Name, err)
		}
		logger.Info().Msgf("queue %q created", queue.Name)

		if queue.Bind != nil {
			err := ch.QueueBind(
				queue.Name,
				queue.Bind.Key,
				queue.Bind.Exchange,
				queue.Bind.NoWait,
				queue.Bind.Args,
			)
			if err != nil {
				return fmt.Errorf("cannot bind queue %q to exchange %q: %w", queue.Name, queue.Bind.Exchange, err)
			}
			logger.Info().Msgf("%q bind to %q exchange with %q routing key", queue.Name, queue.Bind.Exchange, queue.Bind.Key)
		}
	}

	return nil
}

func applyMongoFixtures(ctx context.Context, f flags, dbClient mongo.DbClient, logger zerolog.Logger) error {
	if f.mongoFixtureDirectory == "" {
		return errors.New("-mongo-fixture-directory is not set")
	}

	collections, err := dbClient.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return err
	}
	if len(collections) > 0 {
		return nil
	}

	logger.Info().Msg("start mongo fixtures")
	loader := fixtures.NewLoader(dbClient, []string{f.mongoFixtureDirectory},
		fixtures.NewParser(fixtures.NewFaker(password.NewBcryptEncoder())), logger)
	err = loader.Load(ctx)
	if err != nil {
		return err
	}

	if f.mongoFixtureMigrations {
		cmd := cli.NewSkipCmd(f.mongoMigrationDirectory, f.mongoFixtureMigrationsVersion, dbClient, logger)
		err = cmd.Exec(ctx)
		if err != nil {
			return err
		}
	}

	logger.Info().Msg("finish mongo fixtures")
	return nil
}

func updateMongoConfig(ctx context.Context, conf Conf, dbClient mongo.DbClient) error {
	globalConfAdapter := config.NewAdapter(dbClient)
	prevGlobalConf, err := globalConfAdapter.GetConfig(ctx)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return fmt.Errorf("failed to fetch global config: %w", err)
	}

	//todo: fix it with config refactoring
	conf.Canopsis.Metrics.EnabledInstructions = prevGlobalConf.Metrics.EnabledInstructions
	conf.Canopsis.Metrics.EnabledNotAckedMetrics = prevGlobalConf.Metrics.EnabledNotAckedMetrics
	err = globalConfAdapter.UpsertConfig(ctx, conf.Canopsis)
	if err != nil {
		return fmt.Errorf("failed to update global config: %w", err)
	}

	err = config.NewRemediationAdapter(dbClient).UpsertConfig(ctx, conf.Remediation)
	if err != nil {
		return fmt.Errorf("failed to update remediation config: %w", err)
	}

	err = config.NewHealthCheckAdapter(dbClient).UpsertConfig(ctx, conf.HealthCheck)
	if err != nil {
		return fmt.Errorf("failed to update healthcheck config: %w", err)
	}

	return nil
}

func updateVersionConfig(ctx context.Context, f flags, dbClient mongo.DbClient) error {
	adapter := config.NewVersionAdapter(dbClient)
	prevConf, err := adapter.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch version config: %w", err)
	}

	buildInfo := canopsis.GetBuildInfo()
	conf := config.VersionConf{
		Version: buildInfo.Version,
		Edition: f.edition,
		Stack:   "go",
	}

	if prevConf.Version != conf.Version {
		versionUpdated := datetime.NewCpsTime()
		conf.VersionUpdated = &versionUpdated
	}

	err = adapter.UpsertConfig(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to update version config: %w", err)
	}

	return nil
}

func migrateMongo(ctx context.Context, f flags, dbClient mongo.DbClient, logger zerolog.Logger) error {
	if !f.modeMigrateMongo {
		return nil
	}

	if f.mongoMigrationDirectory == "" {
		return errors.New("-mongo-migration-directory is not set")
	}

	var executor mongo.ScriptExecutor
	switch f.mongoMigrationExec {
	case MongoMigrationExecGoja:
		executor = goja.NewScriptExecutor(dbClient)
	case MongoMigrationExecMongosh:
		executor = mongosh.NewScriptExecutor()
	default:
		return errors.New("-mongo-migration-exec is invalid")
	}

	logger.Info().Msg("start mongo migrations")
	cmd := cli.NewUpCmd(f.mongoMigrationDirectory, "", dbClient, executor, logger)
	err := cmd.Exec(ctx)
	if err != nil {
		return err
	}
	logger.Info().Msg("finish mongo migrations")
	return nil
}

func migratePostgres(f flags, logger zerolog.Logger) error {
	if !f.modeMigratePostgres {
		return nil
	}
	if f.postgresMigrationDirectory == "" {
		return errors.New("-postgres-migration-directory is not set")
	}

	logger.Info().Msg("start postgres migrations")

	connStr, err := postgres.GetConnStr()
	if err != nil {
		return err
	}

	err = runPostgresMigrations(f.postgresMigrationDirectory, f.postgresMigrationMode, f.postgresMigrationSteps, connStr)
	if err != nil {
		return err
	}

	logger.Info().Msg("finish postgres migrations")
	return nil
}

func migrateTechPostgres(f flags, logger zerolog.Logger) error {
	if !f.modeMigrateTechPostgres {
		return nil
	}
	if f.techPostgresMigrationDirectory == "" {
		return errors.New("-tech-postgres-migration-directory is not set")
	}

	logger.Info().Msg("start tech postgres migrations")

	connStr, err := postgres.GetTechConnStr()
	if err != nil {
		return err
	}

	err = runPostgresMigrations(f.techPostgresMigrationDirectory, f.techPostgresMigrationMode, f.techPostgresMigrationSteps, connStr)
	if err != nil {
		return err
	}

	logger.Info().Msg("finish tech postgres migrations")
	return nil
}

func generateSerialName(ctx context.Context, force bool) error {
	if os.Getenv(postgres.EnvURL) == "" {
		return nil
	}

	pool, err := postgres.NewPool(ctx, 0, 0)
	if err != nil {
		return err
	}

	res := pool.QueryRow(ctx, "SELECT id FROM serial_name LIMIT 1")

	serialName := ""
	err = res.Scan(&serialName)
	if errors.Is(err, pgxdriver.ErrNoRows) {
		_, err = pool.Exec(ctx, "INSERT INTO serial_name (id) VALUES ($1)", petname.Generate(2, "-"))
		return err
	}

	if err != nil {
		return err
	}

	if force {
		_, err = pool.Exec(ctx, "UPDATE serial_name SET id = $1 WHERE id = $2", petname.Generate(2, "-"), serialName)
	}

	return err
}

func runPostgresMigrations(migrationDirectory, mode string, steps int, connStr string) error {
	p := &pgx.Postgres{}
	driver, err := p.Open(connStr)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationDirectory, "pgx", driver)
	if err != nil {
		return err
	}

	if steps < 0 {
		return errors.New("postgres migration steps should be >= 0")
	}

	switch mode {
	case "up":
		if steps != 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps != 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
	default:
		return errors.New("postgres migration mode should be up or down")
	}

	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func parseConfig(f flags, logger zerolog.Logger) (Conf, error) {
	data, err := os.ReadFile(f.confFile)
	if err != nil {
		return Conf{}, err
	}
	if f.overrideConfFile != "" {
		overrideData, err := os.ReadFile(f.overrideConfFile)
		if err == nil {
			data, err = mergeConfigs(data, overrideData)
		}

		if err != nil {
			logger.Warn().Err(err).Msgf("skip configuration override")
		}
	}

	var conf Conf
	err = toml.Unmarshal(data, &conf)
	return conf, err
}

func mergeConfigs(configs ...[]byte) ([]byte, error) {
	var res map[string]interface{}
	for _, b := range configs {
		v := make(map[string]interface{})
		err := toml.Unmarshal(b, &v)
		if err != nil {
			return nil, err
		}
		if len(res) == 0 {
			res = v
		} else {
			mergeMaps(res, v)
		}
	}

	return toml.Marshal(res)
}

func mergeMaps(l, r map[string]interface{}) {
	for k, v := range r {
		if rm, ok := v.(map[string]interface{}); ok {
			if lm, ok := l[k].(map[string]interface{}); ok {
				mergeMaps(lm, rm)
				continue
			}
		}
		l[k] = v
	}
}
