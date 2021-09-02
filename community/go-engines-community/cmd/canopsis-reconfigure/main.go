package main

import (
	"context"
	"io/ioutil"
	"os"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/migration/cli"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/pelletier/go-toml"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ErrRabbitInit = 2

	DefaultCfgFile = "/opt/canopsis/etc/canopsis.toml"
	FlagUsageConf  = "The configuration file used to initialize Canopsis."

	DefaultMigrationsPath = "/opt/canopsis/share/database/migrations"
	DefaultFixturesPath   = "/opt/canopsis/share/database/fixtures/prod"
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

	logger := log.NewLogger(f.modeDebug)
	data, err := ioutil.ReadFile(f.confFile)
	if err != nil {
		logger.Error().Err(err).Int("exit status", 1).Msg("")
		os.Exit(1)
	}

	var conf Conf
	if err := toml.Unmarshal(data, &conf); err != nil {
		logger.Error().Err(err).Int("exit status", 2).Msg("")
		os.Exit(2)
	}

	err = GracefullStart(ctx, logger)
	utils.FailOnError(err, "Failed to open one of required sessions")

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

	client, err := mongo.NewClient(ctx, 0, 0)
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
		loader := fixtures.NewLoader(client, []string{f.fixtureDirectory}, true,
			fixtures.NewParser(password.NewSha1Encoder()), logger)
		err = loader.Load(ctx)
		utils.FailOnError(err, "Failed to apply fixtures")
		logger.Info().Msg("Finish fixtures")
	}

	err = config.NewAdapter(client).UpsertConfig(ctx, conf.Canopsis)
	utils.FailOnError(err, "Failed to save config into mongo")
	err = config.NewRemediationAdapter(client).UpsertConfig(ctx, conf.Remediation)
	utils.FailOnError(err, "Failed to save config into mongo")
	err = config.NewHealthCheckAdapter(client).UpsertConfig(ctx, conf.HealthCheck)
	utils.FailOnError(err, "Failed to save config into mongo")

	logger.Info().Msg("Start migrations")
	cmd := cli.NewUpCmd(f.migrationDirectory, "", client, mongo.NewScriptExecutor(), logger)
	err = cmd.Exec(ctx)
	utils.FailOnError(err, "Failed to migrate")
	logger.Info().Msg("Finish migrations")
}
