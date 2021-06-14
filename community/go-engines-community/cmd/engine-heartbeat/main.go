package main

import (
	"context"
	"flag"
	"os"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/debug"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
)

func main() {
	flagVersion := flag.Bool("version", false, "version infos")
	flag.Parse()

	if *flagVersion {
		canopsis.PrintVersionExit()
	}

	logger := log.NewLogger(false)

	trace := debug.Start(logger)

	ctx := context.Background()

	depMaker := DependencyMaker{}
	references := depMaker.GetDefaultReferences(ctx, logger)

	engine := NewEngineHeartBeat(references)
	engine.SendAlarmFunc = engine.sendalarm

	logger.Info().Msg("starting heartbeat")
	exitStatus, err := canopsis.StartEngine(ctx, engine, nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
	}

	trace.Stop()
	os.Exit(exitStatus)
}
