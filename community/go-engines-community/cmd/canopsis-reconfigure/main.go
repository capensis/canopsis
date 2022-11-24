package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/pelletier/go-toml"
	"github.com/rs/zerolog"
)

const (
	ErrGeneral    = 1
	ErrRabbitInit = 2

	DefaultCfgFile = "/opt/canopsis/etc/canopsis.toml"
	FlagUsageConf  = "The configuration file used to initialize Canopsis."

	DefaultMongoConfPath = "/opt/canopsis/share/config/mongo"
)

type Conf struct {
	RabbitMQ config.RabbitMQConf `toml:"RabbitMQ"`
	Canopsis config.CanopsisConf `toml:"Canopsis"`
}

func main() {
	var confFile string
	var mongoConfPath string
	var migrationDirectory string
	var modeDebug bool
	var mongoContainer string
	var mongoURL string
	var modeMigrateOnly bool
	var modeMigrate bool

	flag.StringVar(&confFile, "conf", DefaultCfgFile, FlagUsageConf)
	flag.StringVar(&mongoConfPath, "mongoConf", DefaultMongoConfPath, "The configuration file path is used to create mongo indexes.")
	flag.BoolVar(&modeDebug, "d", false, "debug mode")
	flag.BoolVar(&modeMigrate, "migrate", false, "If true, it will execute migration scripts")
	flag.BoolVar(&modeMigrateOnly, "migrate-only", false, "If true, it will only execute migration scripts")
	flag.StringVar(&migrationDirectory, "migration-directory", "", "The directory with migration scripts")
	flag.StringVar(&mongoContainer, "mongo-container", "", "Should contain docker container_id. If set, it will execute migration scripts inside the container")
	flag.StringVar(&mongoURL, "mongo-url", "", "mongo url")

	flag.Parse()

	logger := log.NewLogger(modeDebug)

	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		logger.Error().Err(err).Int("status", 1).Msg("exit")
		os.Exit(1)
	}

	var conf Conf
	if err := toml.Unmarshal(data, &conf); err != nil {
		logger.Error().Err(err).Int("status", 2).Msg("exit")
		os.Exit(2)
	}

	err = GracefullStart(logger)
	utils.FailOnError(err, "Failed to open one of required sessions")

	if modeMigrate || modeMigrateOnly {
		if migrationDirectory == "" {
			logger.Error().Msg("-migration-directory is not set")
			os.Exit(ErrGeneral)
		}

		logger.Info().Msg("Start migrations")

		err = executeMigrations(logger, migrationDirectory, mongoURL, mongoContainer)
		if err != nil {
			utils.FailOnError(err, "Failed to migrate")
		}

		logger.Info().Msg("Finish migrations")

		if modeMigrateOnly {
			return
		}
	}

	amqpConn, err := amqp.NewConnection(logger, 0, 0)
	utils.FailOnError(err, "Failed to open amqp")

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
			logger.Error().Err(err).Int("status", 2).Msgf("Can not initialise exchange %s", exchange.Name)
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
			logger.Error().Err(err).Int("status", 2).Msgf("Can not initialise queue %s", queue.Name)
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
				logger.Error().Err(err).Int("status", 2).Msgf("Can not bind queue %s to exchange %s %v", queue.Name, queue.Bind.Exchange, err)
				os.Exit(ErrRabbitInit)
			}
			logger.Info().Msgf("\"%s\" bind to \"%s\" exchange with \"%s\" routing key.\n",
				queue.Name,
				queue.Bind.Exchange,
				queue.Bind.Key)
		}

	}

	client, err := mongo.NewClient(0, 0)
	utils.FailOnError(err, "Failed to create mongo session")

	configAdapter := config.NewAdapter(client)
	err = configAdapter.UpsertConfig(conf.Canopsis)
	utils.FailOnError(err, "Failed to save config into mongo")

	logger.Info().Msg("Initialising Mongo indexes")
	err = createMongoIndexes(mongoConfPath, logger)
	utils.FailOnError(err, "Failed to create Mongo indexes")
}

func createMongoIndexes(mongoConfPath string, logger zerolog.Logger) error {
	client, err := mongo.NewClient(0, 0)

	if err != nil {
		logger.
			Error().
			Err(err).
			Msg("failed to open MongoDB session")

		return err
	}

	service := mongo.NewIndexService(
		client,
		mongoConfPath,
		&logger,
	)

	return service.Create()
}

func executeMigrations(logger zerolog.Logger, migrationDirectory, mongoURL, mongoContainer string) error {
	return filepath.Walk(migrationDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*.js", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			logger.Info().Msg(fmt.Sprintf("Start migration %s", path))

			command := fmt.Sprintf("mongo %s %s", mongoURL, path)
			if mongoContainer != "" {
				command = fmt.Sprintf("sudo docker exec -i %s mongo %s < %s", mongoContainer, mongoURL, path)
			}

			result := exec.Command("bash", "-c", command)

			output, err := result.CombinedOutput()
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + string(output))
				return err
			}

			fmt.Println(string(output))
			logger.Info().Msg(fmt.Sprintf("Finish migration %s", path))
		}

		return nil
	})
}
