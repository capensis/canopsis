package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/migration/cli"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ErrGeneral    = 1
	ErrRabbitInit = 2
)

type Conf struct {
	RabbitMQ    config.RabbitMQConf    `toml:"RabbitMQ"`
	Canopsis    config.CanopsisConf    `toml:"Canopsis"`
	Remediation config.RemediationConf `toml:"Remediation"`
	HealthCheck config.HealthCheckConf `toml:"HealthCheck"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := flags{}
	f.Parse()

	if f.version {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(f.modeDebug)
	conf, err := parseConfig(f, logger)
	utils.FailOnError(err, "Failed to parse config")

	err = GracefullStart(ctx, logger)
	utils.FailOnError(err, "Failed to open one of required sessions")

	if f.modeMigratePostgres {
		if f.postgresMigrationDirectory == "" {
			logger.Error().Msg("-postgres-migration-directory is not set")
			os.Exit(ErrGeneral)
		}

		logger.Info().Msg("Start postgres migrations")

		err = runPostgresMigrations(f.postgresMigrationDirectory, f.postgresMigrationMode, f.postgresMigrationSteps)
		if err != nil {
			utils.FailOnError(err, "Failed to migrate")
		}

		logger.Info().Msg("Finish postgres migrations")
	}

	amqpConn, err := amqp.NewConnection(logger, 0, 0)
	utils.FailOnError(err, "Failed to open amqp")
	defer func() {
		err = amqpConn.Close()
		if err != nil {
			logger.Err(err).Msg("cannot close rmq session")
		}
	}()

	ch, err := amqpConn.Channel()
	utils.FailOnError(err, "Failed to open amqp channel")

	logger.Info().Msg("Initialising RabbitMQ exchanges")
	for _, exchange := range conf.RabbitMQ.Exchanges {

		err := ch.ExchangeDeclare(exchange.Name,
			exchange.Kind,
			exchange.Durable,
			exchange.AutoDelete,
			exchange.Internal,
			exchange.NoWait,
			exchange.Args)
		if err != nil {
			logger.Error().Err(err).Int("exit status", 2).Msgf("Can not initialise exchange %s", exchange.Name)
			os.Exit(ErrRabbitInit)
		}
		logger.Info().Msgf("Exchange \"%s\" created.", exchange.Name)
	}

	logger.Info().Msg("Initialising RabbitMQ queues")
	for _, queue := range conf.RabbitMQ.Queues {

		_, err := ch.QueueDeclare(queue.Name,
			queue.Durable,
			queue.AutoDelete,
			queue.Exclusive,
			queue.NoWait,
			queue.Args)

		if err != nil {
			logger.Error().Err(err).Int("exit status", 2).Msgf("Can not initialise queue %s", queue.Name)
			os.Exit(ErrRabbitInit)
		}

		logger.Info().Msgf("Queue \"%s\" created.", queue.Name)

		if queue.Bind != nil {

			err := ch.QueueBind(queue.Name,
				queue.Bind.Key,
				queue.Bind.Exchange,
				queue.Bind.NoWait,
				queue.Bind.Args)

			if err != nil {
				logger.Error().Err(err).Int("exit status", 2).Msgf("Can not bind queue %s to exchange %s %v", queue.Name, queue.Bind.Exchange, err)
				os.Exit(ErrRabbitInit)
			}
			logger.Info().Msgf("\"%s\" bind to \"%s\" exchange with \"%s\" routing key.\n",
				queue.Name,
				queue.Bind.Exchange,
				queue.Bind.Key)
		}

	}

	client, err := mongo.NewClient(ctx, 0, 0, logger)
	utils.FailOnError(err, "Failed to create mongo session")
	defer func() {
		err = client.Disconnect(context.Background())
		if err != nil {
			logger.Err(err).Msg("cannot close mongo session")
		}
	}()

	collections, err := client.ListCollectionNames(ctx, bson.M{})
	utils.FailOnError(err, "Failed to apply fixtures")
	if len(collections) == 0 {
		logger.Info().Msg("Start fixtures")
		loader := fixtures.NewLoader(client, []string{f.mongoFixtureDirectory},
			fixtures.NewParser(fixtures.NewFaker(password.NewSha1Encoder())), logger)
		err = loader.Load(ctx)
		utils.FailOnError(err, "Failed to apply fixtures")
		logger.Info().Msg("Finish fixtures")
	}

	err = updateMongoConfig(ctx, f, conf, client)
	utils.FailOnError(err, "Failed to save config into mongo")

	if f.modeMigrateMongo {
		if f.mongoMigrationDirectory == "" {
			logger.Error().Msg("-mongo-migration-directory is not set")
			os.Exit(ErrGeneral)
		}

		logger.Info().Msg("Start migrations")
		cmd := cli.NewUpCmd(f.mongoMigrationDirectory, "", client, mongo.NewScriptExecutor(), logger)
		err = cmd.Exec(ctx)
		utils.FailOnError(err, "Failed to migrate")
		logger.Info().Msg("Finish migrations")
	}
}

func updateMongoConfig(ctx context.Context, f flags, conf Conf, dbClient mongo.DbClient) error {
	versionConfAdapter := config.NewVersionAdapter(dbClient)
	prevVersionConf, err := versionConfAdapter.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch version config: %w", err)
	}
	buildInfo := canopsis.GetBuildInfo()
	versionConf := config.VersionConf{
		Version: buildInfo.Version,
		Edition: f.edition,
		Stack:   "go",
	}
	if prevVersionConf.Version != versionConf.Version {
		versionUpdated := types.NewCpsTime()
		versionConf.VersionUpdated = &versionUpdated
	}
	err = versionConfAdapter.UpsertConfig(ctx, versionConf)
	if err != nil {
		return fmt.Errorf("failed to update version config: %w", err)
	}
	err = config.NewAdapter(dbClient).UpsertConfig(ctx, conf.Canopsis)
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

func runPostgresMigrations(migrationDirectory, mode string, steps int) error {
	connStr, err := postgres.GetConnStr()
	if err != nil {
		return err
	}

	p := &pgx.Postgres{}
	driver, err := p.Open(connStr)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationDirectory), "pgx", driver)
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

	if err != migrate.ErrNoChange {
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
