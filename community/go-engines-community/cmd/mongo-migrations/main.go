package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/migration/cli"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo/goja"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo/mongosh"
	"github.com/rs/zerolog"
)

const helpFlag = "-h"

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show the version information")
	flag.Parse()

	if version {
		canopsis.PrintVersionInfo()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := log.NewLogger(false)

	err := execCmd(ctx, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}

func execCmd(ctx context.Context, logger zerolog.Logger) error {
	var cmdName string
	var args []string
	if len(os.Args) >= 2 {
		cmdName = os.Args[1]
		args = os.Args[2:]
	}

	var cmd cli.Cmd
	switch cmdName {
	case "create":
		flags := createFlags{}
		err := flags.Parse(args)
		handleFlagErr(err)
		if flags.name == "" {
			_ = flags.Parse([]string{helpFlag})
			return nil
		}

		cmd = cli.NewCreateCmd(flags.path, flags.name, logger)
	case "up":
		flags := upFlags{}
		err := flags.Parse(args)
		handleFlagErr(err)
		client, err := mongo.NewClient(ctx, 0, 0, logger)
		if err != nil {
			return err
		}

		scriptExecutor, err := newScriptExecutor(flags.migrationExec, client)
		if err != nil {
			return err
		}

		cmd = cli.NewUpCmd(flags.path, flags.to, client, scriptExecutor, logger)
	case "down":
		flags := downFlags{}
		err := flags.Parse(args)
		handleFlagErr(err)
		client, err := mongo.NewClient(ctx, 0, 0, logger)
		if err != nil {
			return err
		}

		scriptExecutor, err := newScriptExecutor(flags.migrationExec, client)
		if err != nil {
			return err
		}

		cmd = cli.NewDownCmd(flags.path, flags.to, client, scriptExecutor, logger)
	case "status":
		flags := statusFlags{}
		err := flags.Parse(args)
		handleFlagErr(err)
		client, err := mongo.NewClient(ctx, 0, 0, logger)
		if err != nil {
			return err
		}
		cmd = cli.NewStatusCmd(flags.path, client)
	case "skip":
		flags := skipFlags{}
		err := flags.Parse(args)
		handleFlagErr(err)
		client, err := mongo.NewClient(ctx, 0, 0, logger)
		if err != nil {
			return err
		}
		cmd = cli.NewSkipCmd(flags.path, flags.version, client, logger)
	case "help", "":
		args := []string{helpFlag}
		_ = (&createFlags{}).Parse(args)
		_ = (&upFlags{}).Parse(args)
		_ = (&downFlags{}).Parse(args)
		_ = (&statusFlags{}).Parse(args)
		_ = (&skipFlags{}).Parse(args)
		return nil
	default:
		return fmt.Errorf("unknown command %q", cmdName)
	}

	return cmd.Exec(ctx)
}

func handleFlagErr(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	}

	os.Exit(2)
}

func newScriptExecutor(execName string, dbClient mongo.DbClient) (mongo.ScriptExecutor, error) {
	switch execName {
	case MigrationExecGoja:
		return goja.NewScriptExecutor(dbClient), nil
	case MigrationExecMongosh:
		return mongosh.NewScriptExecutor(), nil
	default:
		return nil, errors.New("-migration-exec is invalid")
	}
}
