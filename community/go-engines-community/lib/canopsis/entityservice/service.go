package entityservice

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"runtime/trace"
	"sync"
)

const maxWorkersCount = 10

type service struct {
	adapter       Adapter
	entityAdapter entity.Adapter
	logger        zerolog.Logger
}

func NewService(
	adapter Adapter,
	entityAdapter entity.Adapter,
	logger zerolog.Logger,
) IdleSinceService {
	service := service{
		adapter:       adapter,
		entityAdapter: entityAdapter,
		logger:        logger,
	}
	return &service
}

func (s *service) markServices(parentCtx context.Context, idleSinceMap *ServicesIdleSinceMap, services []EntityService, impacts []string, timestamp int64) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	wg := sync.WaitGroup{}
	workerCh := make(chan string)
	defer close(workerCh)

	wg.Add(len(impacts))
	go func() {
		for _, impact := range impacts {
			workerCh <- impact
		}
	}()

	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			for impact := range workerCh {
				func() {
					defer wg.Done()

					select {
					case <-ctx.Done():
						return
					default:
					}

					if !idleSinceMap.Mark(impact, timestamp) {
						return
					}

					for _, service := range services {
						if service.ID == impact && len(service.Impacts) > 0 {
							wg.Add(len(service.Impacts))
							go func(service EntityService) {
								for _, impact := range service.Impacts {
									select {
									case <-ctx.Done():
										return
									case workerCh <- impact:
									}
								}
							}(service)

							return
						}
					}
				}()
			}
		}()
	}

	wg.Wait()
}

func (s *service) RecomputeIdleSince(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	defer trace.StartRegion(ctx, "service.RecomputeIdleSince").End()

	services, err := s.adapter.GetValid(ctx)
	if err != nil {
		return err
	}

	if len(services) == 0 {
		return nil
	}

	idleSinceMap := NewServicesIdleSinceMap()
	cursor, err := s.entityAdapter.GetWithIdleSince(ctx)
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	wg := sync.WaitGroup{}
	workerCh := make(chan types.Entity)
	go func() {
		defer close(workerCh)
		for cursor.Next(ctx) {
			var ent types.Entity
			err := cursor.Decode(&ent)
			if err != nil {
				s.logger.Err(err).Msg("Can't decode entity")
			}

			select {
			case <-ctx.Done():
				return
			case workerCh <- ent:
			}
		}
	}()

	errCh := make(chan error, maxWorkersCount)
	defer close(errCh)

	for i := 0; i < maxWorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case ent, ok := <-workerCh:
					if !ok {
						return
					}

					if ent.IdleSince == nil {
						continue
					}

					if ent.Type == types.EntityTypeResource || ent.Type == types.EntityTypeComponent {
						s.markServices(ctx, &idleSinceMap, services, ent.Impacts, ent.IdleSince.Unix())
					}

					if ent.Type == types.EntityTypeConnector {
						s.markServices(ctx, &idleSinceMap, services, ent.ImpactedServicesFromDependencies, ent.IdleSince.Unix())
					}
				}
			}
		}()
	}

	wg.Wait()

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	var newModel mongodriver.WriteModel
	bulkBytesSize := 0

	for _, service := range services {
		idleSince := idleSinceMap.idleMap[service.ID]
		if idleSince > 0 {
			newModel = mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": service.ID}).
				SetUpdate(bson.M{"$set": bson.M{"idle_since": idleSince}})
		} else {
			newModel = mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": service.ID}).
				SetUpdate(bson.M{"$unset": bson.M{"idle_since": ""}})
		}

		b, err := bson.Marshal(newModel)
		if err != nil {
			return err
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			err := s.adapter.UpdateBulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(
			writeModels,
			newModel,
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			err := s.adapter.UpdateBulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		err = s.adapter.UpdateBulk(ctx, writeModels)
	}

	return err
}
