package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const shutdownTimout = 5 * time.Second

//  Router is used to implement adding new routes to API.
type Router func(gin.IRouter)

//  Worker is used to implement adding new worker to API.
type Worker func(context.Context)

// API is used to implement API http server.
type API interface {
	// Run starts http server.
	Run(context.Context) error
	// AddRouter adds new router.
	AddRouter(Router)
	// AddWorker adds new worker.
	AddWorker(string, Worker)
	// AddNoRoute adds handlers for no roure.
	AddNoRoute([]gin.HandlerFunc)
}

type api struct {
	addr            string
	logger          zerolog.Logger
	routers         []Router
	noRouteHandlers []gin.HandlerFunc
	workers         map[string]Worker
	waitGroup       sync.WaitGroup
}

// New creates new api.
func New(
	addr string,
	logger zerolog.Logger,
) API {
	return &api{
		addr:    addr,
		logger:  logger,
		routers: make([]Router, 0),
		workers: make(map[string]Worker),
	}
}

func (a *api) AddWorker(key string, worker Worker) {
	a.workers[key] = worker
}

func (a *api) AddRouter(router Router) {
	a.routers = append(a.routers, router)
}

func (a *api) AddNoRoute(handlers []gin.HandlerFunc) {
	a.noRouteHandlers = handlers
}

func (a *api) Run(ctx context.Context) error {
	handler := a.registerRoutes()
	a.runWorkers(ctx)

	// Start server.
	server := &http.Server{
		Addr:    a.addr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			a.logger.Err(err).Msg("server forced to shutdown")
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Err(err).Msg("server fail to start")

		return err
	}

	a.waitGroup.Wait()

	return nil
}

func (a *api) registerRoutes() http.Handler {
	ginRouter := gin.New()
	ginRouter.Use(gin.Logger())
	ginRouter.Use(middleware.Recovery(a.logger))

	for _, router := range a.routers {
		router(ginRouter)
	}

	if len(a.noRouteHandlers) > 0 {
		ginRouter.NoRoute(a.noRouteHandlers...)
	}

	return ginRouter
}

func (a *api) runWorkers(ctx context.Context) {
	for key := range a.workers {
		a.waitGroup.Add(1)
		f := a.workers[key]

		go RestartGoroutine(fmt.Sprintf("worker %s", key), func() {
			f(ctx)
		}, &a.waitGroup, a.logger)
	}
}

// RestartGoroutine starts goroutine with panic recovery. RestartGoroutine logs
// recovery and restarts goroutine on panic.
func RestartGoroutine(
	key string,
	f func(),
	wg *sync.WaitGroup,
	logger zerolog.Logger,
) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%s have been recovered from panic and restarted", key)

			if err, ok := r.(error); ok {
				logger.Err(err).Msg(msg)
			} else {
				logger.Error().Interface("recover", r).Msg(msg)
			}

			go RestartGoroutine(key, f, wg, logger)
		} else {
			wg.Done()
		}
	}()

	f()
}
