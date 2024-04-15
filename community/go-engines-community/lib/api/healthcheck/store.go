package healthcheck

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/jackc/pgx/v5"
	libredis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Load(ctx context.Context)
	GetInfo() Info
	GetStatus() Status
	GetEnginesOrder() Graph
	GetParameters(ctx context.Context) (config.HealthCheckParameters, error)
	UpdateParameters(ctx context.Context, params config.HealthCheckParameters) (config.HealthCheckParameters, error)
}

func NewStore(
	dbClient mongo.DbClient,
	manager engine.RunInfoManager,
	configAdapter config.HealthCheckAdapter,
	configProvider *config.BaseHealthCheckConfigProvider,
	logger zerolog.Logger,
	websocketHub websocket.Hub,
) Store {
	return &store{
		collection:     dbClient.Collection(mongo.ConfigurationMongoCollection),
		manager:        manager,
		logger:         logger,
		configAdapter:  configAdapter,
		configProvider: configProvider,
		websocketHub:   websocketHub,
	}
}

type store struct {
	collection     mongo.DbCollection
	manager        engine.RunInfoManager
	mxEngines      sync.RWMutex
	engines        Engines
	mxServices     sync.RWMutex
	services       []Service
	logger         zerolog.Logger
	configAdapter  config.HealthCheckAdapter
	configProvider *config.BaseHealthCheckConfigProvider
	websocketHub   websocket.Hub

	mongoClient    mongo.DbClient
	redisClient    libredis.UniversalClient
	amqpConnection amqp.Connection
	pgConn         *pgx.Conn

	mxChangedParamsCh sync.RWMutex
	changedParamsCh   chan<- bool
}

func (s *store) GetParameters(ctx context.Context) (config.HealthCheckParameters, error) {
	conf := config.HealthCheckConf{}
	err := s.collection.FindOne(ctx, bson.M{"_id": config.HealthCheckName}).Decode(&conf)

	return conf.Parameters, err
}

func (s *store) UpdateParameters(ctx context.Context, params config.HealthCheckParameters) (config.HealthCheckParameters, error) {
	conf := config.HealthCheckConf{}
	err := s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": config.HealthCheckName},
		bson.M{"$set": bson.M{
			"parameters": params,
		}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&conf)

	s.mxChangedParamsCh.RLock()
	defer s.mxChangedParamsCh.RUnlock()

	if s.changedParamsCh != nil {
		select {
		case s.changedParamsCh <- true:
		default:
			s.logger.Warn().Msg("cannot send message to channel")
		}
	}

	return conf.Parameters, err
}

func (s *store) GetEnginesOrder() Graph {
	engineOrder := s.configProvider.Get().EngineOrder
	graph := Graph{
		Nodes: make([]string, 0, len(engineOrder)),
		Edges: make([]Edge, len(engineOrder)),
	}
	nodeExists := make(map[string]struct{}, len(engineOrder))

	for idx, pair := range engineOrder {
		if _, ok := nodeExists[pair.From]; !ok {
			graph.Nodes = append(graph.Nodes, pair.From)
			nodeExists[pair.From] = struct{}{}
		}

		if _, ok := nodeExists[pair.To]; !ok {
			graph.Nodes = append(graph.Nodes, pair.To)
			nodeExists[pair.To] = struct{}{}
		}

		graph.Edges[idx].From = pair.From
		graph.Edges[idx].To = pair.To
	}

	return graph
}

func (s *store) Load(ctx context.Context) {
	ch := make(chan bool)
	s.mxChangedParamsCh.Lock()
	s.changedParamsCh = ch
	s.mxChangedParamsCh.Unlock()

	ticker := time.NewTicker(s.configProvider.Get().ParseUpdateInterval(s.logger))

	defer func() {
		s.mxChangedParamsCh.Lock()
		s.changedParamsCh = nil
		s.mxChangedParamsCh.Unlock()
		close(ch)

		ticker.Stop()

		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		s.closeConnections(closeCtx)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ch:
			s.doLoad(ctx)
		case <-ticker.C:
			s.doLoad(ctx)
		}
	}
}

func (s *store) GetInfo() Info {
	engines := s.getEngines()
	enginesOrder := s.GetEnginesOrder()
	queueParameters := s.configProvider.Get().Parameters.Queue
	maxQueueLength := queueParameters.Limit
	if !queueParameters.Enabled {
		maxQueueLength = 0
	}
	messagesParameters := s.configProvider.Get().Parameters.Messages
	maxMessagesLength := messagesParameters.Limit
	if !messagesParameters.Enabled {
		maxMessagesLength = 0
	}

	return Info{
		Services:          s.getServices(),
		Engines:           engines,
		MaxQueueLength:    maxQueueLength,
		MaxMessagesLength: maxMessagesLength,

		HasInvalidEnginesOrder: !equalEdges(engines.Graph.Edges, enginesOrder.Edges),
	}
}

func (s *store) GetStatus() Status {
	services := s.getServices()
	engines := s.getEngines()
	enginesOrder := s.GetEnginesOrder()
	status := Status{
		Services: make([]Service, 0, len(services)),
		Engines:  make([]Engine, 0, len(engines.Parameters)),

		HasInvalidEnginesOrder: !equalEdges(engines.Graph.Edges, enginesOrder.Edges),
	}

	for _, service := range services {
		if !service.IsRunning {
			status.Services = append(status.Services, service)
		}
	}

	for name, e := range engines.Parameters {
		if !e.IsRunning || e.IsQueueOverflown || e.IsTooFewInstances || e.IsDiffInstancesConfig {
			status.Engines = append(status.Engines, Engine{
				Name:                  name,
				IsRunning:             e.IsRunning,
				IsQueueOverflown:      e.IsQueueOverflown,
				IsTooFewInstances:     e.IsTooFewInstances,
				IsDiffInstancesConfig: e.IsDiffInstancesConfig,
			})
		}
	}

	return status
}

func (s *store) doLoad(ctx context.Context) {
	s.loadConfig(ctx)
	s.loadServices(ctx)
	s.loadEngines(ctx)
	s.websocketHub.Send(websocket.RoomHealthcheck, s.GetInfo())
	s.websocketHub.Send(websocket.RoomHealthcheckStatus, s.GetStatus())
}

func (s *store) loadEngines(ctx context.Context) {
	engines, err := s.manager.GetEngines(ctx)
	enginesByName := make(map[string]engine.RunInfo, len(engines))
	for _, info := range engines {
		enginesByName[info.Name] = info
	}

	s.mxEngines.Lock()
	defer s.mxEngines.Unlock()

	if err != nil {
		s.engines = Engines{}
		s.logger.Err(err).Msg("cannot load engines run info")
		return
	}

	canonicalOrder := s.configProvider.Get().EngineOrder
	order := make([]string, 0)
	has := make(map[string]bool)

	for _, pair := range canonicalOrder {
		if !has[pair.From] {
			order = append(order, pair.From)
			has[pair.From] = true
		}
		if !has[pair.To] {
			order = append(order, pair.To)
			has[pair.To] = true
		}
	}

	s.engines = transformEngineInfoToGraph(enginesByName, order, s.configProvider.Get().Parameters)
}

func (s *store) loadConfig(ctx context.Context) {
	cfg, err := s.configAdapter.GetConfig(ctx)
	if err != nil {
		s.logger.Err(err).Msg("fail to load healthcheck config")
		return
	}

	err = s.configProvider.Update(cfg)
	if err != nil {
		s.logger.Err(err).Msg("fail to update healthcheck config")
	}
}

func (s *store) loadServices(ctx context.Context) {
	services := []Service{
		{
			Name:      ServiceMongoDB,
			IsRunning: s.checkMongoDB(ctx),
		},
		{
			Name:      ServiceRedis,
			IsRunning: s.checkRedis(ctx),
		},
		{
			Name:      ServiceRabbitMQ,
			IsRunning: s.checkRabbitMQ(),
		},
		{
			Name:      ServiceTimescaleDB,
			IsRunning: s.checkTimescaleDB(ctx),
		},
	}

	s.mxServices.Lock()
	defer s.mxServices.Unlock()
	s.services = services
}

func (s *store) closeConnections(ctx context.Context) {
	if s.mongoClient != nil {
		err := s.mongoClient.Disconnect(ctx)
		if err != nil {
			s.logger.Err(err).Msg("cannot close mongo")
		}
		s.mongoClient = nil
	}
	if s.redisClient != nil {
		err := s.redisClient.Close()
		if err != nil {
			s.logger.Err(err).Msg("cannot close redis")
		}
		s.redisClient = nil
	}
	if s.amqpConnection != nil {
		if !s.amqpConnection.IsClosed() {
			err := s.amqpConnection.Close()
			if err != nil {
				s.logger.Err(err).Msg("cannot close amqp")
			}
		}
		s.amqpConnection = nil
	}
	if s.pgConn != nil {
		if !s.pgConn.IsClosed() {
			err := s.pgConn.Close(ctx)
			if err != nil {
				s.logger.Err(err).Msg("cannot close postgres")
			}
		}
		s.pgConn = nil
	}
}

func (s *store) getEngines() Engines {
	s.mxEngines.RLock()
	defer s.mxEngines.RUnlock()

	return s.engines
}

func (s *store) getServices() []Service {
	s.mxServices.RLock()
	defer s.mxServices.RUnlock()

	return s.services
}

func (s *store) checkMongoDB(ctx context.Context) bool {
	if s.mongoClient == nil {
		var err error
		s.mongoClient, err = mongo.NewClientWithOptions(ctx, 0, 0, time.Second, time.Second, zerolog.Nop())
		if err != nil {
			s.logger.Err(err).Msg("cannot connect to mongo")
			return false
		}
	}

	if err := s.mongoClient.Ping(ctx, nil); err != nil {
		s.logger.Err(err).Msg("cannot ping mongo")
		return false
	}

	return true
}

func (s *store) checkRedis(ctx context.Context) bool {
	if s.redisClient == nil {
		var err error
		s.redisClient, err = redis.NewSession(ctx, 0, s.logger, -1, -1)
		if err != nil {
			s.logger.Err(err).Msg("cannot connect to redis")
			return false
		}
	}

	if err := s.redisClient.Ping(ctx).Err(); err != nil {
		s.logger.Err(err).Msg("cannot ping redis")
		return false
	}

	return true
}

func (s *store) checkRabbitMQ() bool {
	if s.amqpConnection == nil || s.amqpConnection.IsClosed() {
		var err error
		s.amqpConnection, err = amqp.NewConnection(s.logger, 0, 0)
		if err != nil {
			s.logger.Err(err).Msg("cannot connect to rabbitmq")
			return false
		}
	}

	return true
}

func (s *store) checkTimescaleDB(ctx context.Context) bool {
	if s.pgConn == nil || s.pgConn.IsClosed() {
		connStr, err := postgres.GetConnStr()
		if err != nil {
			s.logger.Err(err).Msg("cannot connect to postgres")
			return false
		}

		s.pgConn, err = pgx.Connect(ctx, connStr)
		if err != nil {
			s.logger.Err(err).Msg("cannot connect to postgres")
			return false
		}
	}

	if err := s.pgConn.Ping(ctx); err != nil {
		s.logger.Err(err).Msg("cannot ping postgres")
		return false
	}

	return true
}

func transformEngineInfoToGraph(engines map[string]engine.RunInfo, order []string, parameters config.HealthCheckParameters) Engines {
	graph := Graph{}
	graph.Nodes = make([]string, len(order))
	enginesParams := make(map[string]Engine, len(order))
	for i, name := range order {
		if info, ok := engines[name]; ok {
			engineParams := parameters.GetEngineParameters(name)

			minInstances := engineParams.Minimal
			optimalInstances := engineParams.Optimal
			if !engineParams.Enabled {
				minInstances = 1
				optimalInstances = 1
			}

			graph.Nodes[i] = name
			enginesParams[name] = Engine{
				Instances:             &info.Instances,
				MinInstances:          &minInstances,
				OptimalInstances:      &optimalInstances,
				QueueLength:           &info.QueueLength,
				Time:                  &info.Time,
				IsRunning:             true,
				IsQueueOverflown:      parameters.Queue.Enabled && info.QueueLength > parameters.Queue.Limit,
				IsTooFewInstances:     info.Instances < minInstances,
				IsDiffInstancesConfig: info.HasDiffConfig,
			}
		} else {
			graph.Nodes[i] = name
			enginesParams[name] = Engine{
				IsRunning: false,
			}
		}
	}

	graph.Edges = make([]Edge, 0)
	for i := range engines {
		for j := range engines {
			if engines[i].PublishQueue != "" && engines[i].PublishQueue == engines[j].ConsumeQueue ||
				isIntersected(engines[i].RpcPublishQueues, engines[j].RpcConsumeQueues) {
				graph.Edges = append(graph.Edges, Edge{
					From: engines[i].Name,
					To:   engines[j].Name,
				})
			}
		}
	}

	return Engines{
		Graph:      graph,
		Parameters: enginesParams,
	}
}

func isIntersected(l, r []string) bool {
	for _, lv := range l {
		for _, rv := range r {
			if lv == rv {
				return true
			}
		}
	}

	return false
}

func equalEdges(l, r []Edge) bool {
	sort.Slice(l, func(i, j int) bool {
		return l[i].To < l[j].To || l[i].To == l[j].To && l[i].From == l[j].From
	})
	sort.Slice(r, func(i, j int) bool {
		return r[i].To < r[j].To || r[i].To == r[j].To && r[i].From == r[j].From
	})

	return reflect.DeepEqual(l, r)
}
