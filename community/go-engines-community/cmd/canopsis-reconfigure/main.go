package main

import (
	"context"
	"flag"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
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
	RabbitMQ    config.RabbitMQConf    `toml:"RabbitMQ"`
	Canopsis    config.CanopsisConf    `toml:"Canopsis"`
	Remediation config.RemediationConf `toml:"Remediation"`
	HealthCheck config.HealthCheckConf `toml:"HealthCheck"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	err = config.NewAdapter(client).UpsertConfig(ctx, conf.Canopsis)
	utils.FailOnError(err, "Failed to save config into mongo")
	err = config.NewRemediationAdapter(client).UpsertConfig(ctx, conf.Remediation)
	utils.FailOnError(err, "Failed to save config into mongo")
	err = config.NewHealthCheckAdapter(client).UpsertConfig(ctx, conf.HealthCheck)
	utils.FailOnError(err, "Failed to save config into mongo")

	logger.Info().Msg("Initialise TimescaleDB")
	err = createTimescaleDBTables(ctx)
	if os.Getenv(postgres.EnvURL) != "" && err != nil {
		utils.FailOnError(err, "Failed to create timescaleDB tables")
	}

	logger.Info().Msg("Initialising Mongo indexes")
	err = createMongoIndexes(ctx, client, mongoConfPath, logger)
	utils.FailOnError(err, "Failed to create Mongo indexes")
}

func createTimescaleDBTables(ctx context.Context) error {
	postgresPool, err := postgres.NewPool(ctx)
	if err != nil {
		return err
	}

	defer postgresPool.Close()

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS total_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
		   	value INT);
		   	SELECT create_hypertable('total_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS non_displayed_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
		   	value INT);
		   	SELECT create_hypertable('non_displayed_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS pbh_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
		   	value INT);
		   	SELECT create_hypertable('pbh_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS instruction_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
		   	value INT);
		   	SELECT create_hypertable('instruction_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS correlation_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
		   	value INT);
		   	SELECT create_hypertable('correlation_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS ticket_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
			user_id VARCHAR(255),
		   	value INT);
		   	SELECT create_hypertable('ticket_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS ack_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
			user_id VARCHAR(255),
		   	value INT);
		   	SELECT create_hypertable('ack_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS cancel_ack_alarm_number (
		   	time TIMESTAMP NOT NULL,
		   	entity_id INT,
			user_id VARCHAR(255),
		   	value INT);
		   	SELECT create_hypertable('cancel_ack_alarm_number', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS ack_duration (
			time TIMESTAMP NOT NULL,
			entity_id INT,
			user_id VARCHAR(255),
			value INT);
			SELECT create_hypertable('ack_duration', 'time', if_not_exists => TRUE);
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS resolve_duration (
			time TIMESTAMP NOT NULL,
			entity_id INT,
			value INT);
			SELECT create_hypertable('resolve_duration', 'time', if_not_exists => TRUE);
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS user_logins (
		   	time TIMESTAMP NOT NULL,
		   	user_id VARCHAR(255),
		   	value INT);
		   	SELECT create_hypertable('user_logins', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS user_activity (
		   	time TIMESTAMP NOT NULL,
		   	user_id VARCHAR(255),
		   	value INT);
		   	SELECT create_hypertable('user_activity', 'time', if_not_exists => TRUE);   
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS sli_duration (
			time TIMESTAMP NOT NULL,
			entity_id INT,
			type SMALLINT,
			value INT);
			SELECT create_hypertable('sli_duration', 'time', if_not_exists => TRUE);
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS entities (
			id SERIAL PRIMARY KEY,
			custom_id VARCHAR(500),
			name VARCHAR(500),
		   	category VARCHAR(255),
		   	impact_level INT,
		   	type VARCHAR(255),
			enabled BOOLEAN,
			infos JSONB,
			component_infos JSONB,
			component VARCHAR(500),
			UNIQUE(custom_id)
			);
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			username VARCHAR(255),
		   	role VARCHAR(255)
			);
       	`,
	)
	if err != nil {
		return err
	}

	_, err = postgresPool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS metrics_criteria (
			id INT PRIMARY KEY,
			type INT,
		   	name VARCHAR(255)
			);
       	`,
	)
	if err != nil {
		return err
	}

	for _, c := range defaultCriteria() {
		_, err := postgresPool.Exec(
			ctx,
			fmt.Sprintf(
				`
				INSERT INTO %s (id, type, name) VALUES($1, $2, $3)
				ON CONFLICT ON CONSTRAINT metrics_criteria_pkey DO UPDATE SET type = $2, name = $3
			`,
				postgres.MetricsCriteria,
			),
			c.ID,
			c.Type,
			c.Name,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

type CriteriaConfig struct {
	ID   int
	Type int
	Name string
}

func defaultCriteria() []CriteriaConfig {
	return []CriteriaConfig{
		{
			ID:   1,
			Type: metrics.EntityCriteriaType,
			Name: "name",
		},
		{
			ID:   2,
			Type: metrics.EntityCriteriaType,
			Name: "category",
		},
		{
			ID:   3,
			Type: metrics.EntityCriteriaType,
			Name: "impact_level",
		},
		{
			ID:   4,
			Type: metrics.EntityCriteriaType,
			Name: "type",
		},
		{
			ID:   5,
			Type: metrics.UserCriteriaType,
			Name: "username",
		},
		{
			ID:   6,
			Type: metrics.UserCriteriaType,
			Name: "role",
		},
	}
}

func createMongoIndexes(ctx context.Context, client mongo.DbClient, mongoConfPath string, logger zerolog.Logger) error {
	service := mongo.NewIndexService(
		client,
		mongoConfPath,
		&logger,
	)

	return service.Create(ctx)
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
