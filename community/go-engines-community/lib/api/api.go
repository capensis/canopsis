package api

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimout    = 5 * time.Second
	readHeaderTimeout = 30 * time.Second
)

// Router is used to implement adding new routes to API.
type Router func(*gin.Engine)

// Worker is used to implement adding new worker to API.
type Worker func(context.Context)

type DeferFunc func(ctx context.Context)

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
	// SetWebsocketHub sets websocket hub.
	SetWebsocketHub(websocket.Hub)
	// GetWebsocketHub gets websocket hub.
	GetWebsocketHub() websocket.Hub
}

type api struct {
	addr      string
	deferFunc DeferFunc
	logger    zerolog.Logger
	routers   []Router
	workers   map[string]Worker

	noRouteHandlers  []gin.HandlerFunc
	noMethodHandlers []gin.HandlerFunc

	websocketHub websocket.Hub
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

func (a *api) Run(ctx context.Context) (resErr error) {
	handler := a.registerRoutes()
	workersGroup := a.runWorkers(ctx)

	// Start server.
	server := &http.Server{
		Addr:              a.addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
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

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Err(err).Msg("server fail to start")

		return err
	}

	return workersGroup.Wait()
}

func (a *api) SetWebsocketHub(v websocket.Hub) {
	a.websocketHub = v
}
func (a *api) GetWebsocketHub() websocket.Hub {
	return a.websocketHub
}

func (a *api) registerRoutes() http.Handler {
	ginRouter := gin.New()
	ginRouter.Use(gin.Logger())
	ginRouter.Use(middleware.Recovery(a.logger))
	ginRouter.HandleMethodNotAllowed = true
	ginRouter.ContextWithFallback = true

	for _, router := range a.routers {
		router(ginRouter)
	}

	if len(a.noRouteHandlers) > 0 {
		ginRouter.NoRoute(a.noRouteHandlers...)
	}

	if len(a.noMethodHandlers) > 0 {
		ginRouter.NoMethod(a.noMethodHandlers...)
	}

	return ginRouter
}

func (a *api) runWorkers(ctx context.Context) *errgroup.Group {
	g, ctx := errgroup.WithContext(ctx)

	for key := range a.workers {
		f := a.workers[key]

		restartGoroutine(g, "worker "+key, func() error {
			f(ctx)

			return nil
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
	g.Go(func() error {
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

		return f()
	})
}
