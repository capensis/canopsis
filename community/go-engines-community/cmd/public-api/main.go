package main

//go:generate go tool github.com/swaggo/swag/cmd/swag init -parseDependency -d ../../lib/api/public -g ../../../cmd/public-api/$GOFILE -o ../../lib/api/public/docs --outputTypes yaml --instanceName schemas

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	api "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"golang.org/x/sync/errgroup"
)

const (
	okExitCode  = 0
	errExitCode = 1
)

func main() {
	var opts api.Options
	api.BindCmdFlags(&opts)
	flag.Parse()
	if opts.Version {
		canopsis.PrintVersionInfo()

		return
	}

	exitCode := run(opts)
	os.Exit(exitCode)
}

func run(opts api.Options) int {
	// Graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := log.NewLogger(ctx, opts.LogOptions)
	s, cleanup, err := api.NewServices(ctx, opts, logger)
	defer cleanup()
	if err != nil {
		logger.Err(err).Msg("cannot create services")

		return errExitCode
	}

	r, err := api.RegisterRoutes(s, opts)
	if err != nil {
		logger.Err(err).Msg("failed to register routes")

		return errExitCode
	}

	jobHandlers, jobMaintainers := api.AsyncJobs(s)
	g, gctx := errgroup.WithContext(ctx)
	api.AddBackground(gctx, g, s, jobHandlers, jobMaintainers, cancel)
	hs := serve(r, opts.Port, s, cancel)

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), api.ShutdownTimeout)
	defer shutdownCancel()
	err = hs.Shutdown(shutdownCtx)
	if err != nil {
		logger.Err(err).Msg("failed to shutdown http server")

		return errExitCode
	}

	logger.Info().Msg("http server stopped")

	// Let the worker finish the running jobs before connections are closed.
	err = g.Wait()
	if err != nil {
		logger.Err(err).Msg("workers stopped with error")

		return errExitCode
	}

	return okExitCode
}

// serve starts the http server and returns it, so it can be shut down gracefully.
func serve(r http.Handler, port int, s api.Services, cancel context.CancelFunc) *http.Server {
	hs := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           r,
		ReadHeaderTimeout: api.ReadHeaderTimeout,
	}

	go func() {
		s.Logger.Info().Msg("http server started")
		err := hs.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Err(err).Msg("failed to run http server")
			cancel()
		}
	}()

	return hs
}
