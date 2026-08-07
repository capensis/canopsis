package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimout    = 5 * time.Second
	readHeaderTimeout = 30 * time.Second
)

// Router is used to implement adding new routes to API.
type Router func(*gin.Engine) error

// Worker is used to implement adding new worker to API.
type Worker func(context.Context) error

type DeferFunc func(ctx context.Context)

type ReloadWorkerError struct {
	err error
}

func NewReloadWorkerError(err error) error {
	return ReloadWorkerError{err: err}
}

func (e ReloadWorkerError) Error() string {
	return e.err.Error()
}

func (e ReloadWorkerError) Unwrap() error {
	return e.err
}

// API is used to implement API http server.
type API interface {
	// Run starts http server.
	Run(context.Context) error
	// AddRouter adds new routes.
	AddRouter(Router)
	// AddWorker adds new worker.
	AddWorker(string, Worker)
	// AddNoRoute adds handlers for no route.
	AddNoRoute(...gin.HandlerFunc)
	// AddNoMethod adds handlers for no method.
	AddNoMethod(...gin.HandlerFunc)
}

type api struct {
	addr      string
	deferFunc DeferFunc
	logger    zerolog.Logger
	routers   []Router
	workers   map[string]Worker

	noRouteHandlers  []gin.HandlerFunc
	noMethodHandlers []gin.HandlerFunc
}

// New creates new api.
func New(
	addr string,
	deferFunc DeferFunc,
	logger zerolog.Logger,
) API {
	return &api{
		addr:      addr,
		deferFunc: deferFunc,
		logger:    logger,
		routers:   make([]Router, 0),
		workers:   make(map[string]Worker),
	}
}

func (a *api) AddWorker(key string, worker Worker) {
	if _, ok := a.workers[key]; ok {
		panic(fmt.Errorf("%q worker already exists", key))
	}

	a.workers[key] = worker
}

func (a *api) AddRouter(router Router) {
	a.routers = append(a.routers, router)
}

func (a *api) AddNoRoute(handlers ...gin.HandlerFunc) {
	a.noRouteHandlers = handlers
}

func (a *api) AddNoMethod(handlers ...gin.HandlerFunc) {
	a.noMethodHandlers = handlers
}

func (a *api) Run(ctx context.Context) error {
	apiErrGroup, ctx := errgroup.WithContext(ctx)

	// Start server.
	h, err := a.registerRoutes()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              a.addr,
		Handler:           h,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() { // nolint:gosec
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			a.logger.Err(err).Msg("server forced to shutdown")
		}
	}()

	defer func() { // nolint:contextcheck
		if a.deferFunc != nil {
			deferCtx, deferCancel := context.WithTimeout(context.Background(), shutdownTimout)
			defer deferCancel()
			a.deferFunc(deferCtx)
		}
	}()

	apiErrGroup.Go(func() error {
		return a.runWorkers(ctx).Wait()
	})
	apiErrGroup.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Err(err).Msg("server fail to start")

			return err
		}

		return nil
	})

	return apiErrGroup.Wait()
}

func (a *api) registerRoutes() (http.Handler, error) {
	ginRouter := gin.New()
	ginRouter.HandleMethodNotAllowed = true
	ginRouter.ContextWithFallback = true

	for _, router := range a.routers {
		err := router(ginRouter)
		if err != nil {
			return nil, fmt.Errorf("cannot register routes: %w", err)
		}
	}

	if len(a.noRouteHandlers) > 0 {
		ginRouter.NoRoute(a.noRouteHandlers...)
	}

	if len(a.noMethodHandlers) > 0 {
		ginRouter.NoMethod(a.noMethodHandlers...)
	}

	return ginRouter, nil
}

func (a *api) runWorkers(ctx context.Context) *errgroup.Group {
	g, ctx := errgroup.WithContext(ctx)

	for key := range a.workers {
		f := a.workers[key]

		restartGoroutine(g, "worker "+key, func() error {
			return f(ctx)
		}, a.logger)
	}

	return g
}

// restartGoroutine starts goroutine with panic recovery. RestartGoroutine logs
// recovery and restarts goroutine on panic.
func restartGoroutine(
	g *errgroup.Group,
	key string,
	f func() error,
	logger zerolog.Logger,
) {
	g.Go(func() (gErr error) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				var ok bool
				if err, ok = r.(error); !ok {
					err = fmt.Errorf("%v", r)
				}

				logger.Err(err).Str("worker", key).Msgf("panic recovered\n%s\n", debug.Stack())

				restartGoroutine(g, key, f, logger)
			}
		}()

		for {
			err := f()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					logger.Debug().Err(err).Str("worker", key).Msg("worker is canceled")

					return nil
				}

				if _, ok := errors.AsType[ReloadWorkerError](err); ok {
					logger.Err(err).Str("worker", key).Msgf("worker restart")

					continue
				}

				return err
			}

			return nil
		}
	})
}
