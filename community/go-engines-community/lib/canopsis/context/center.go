package context

import (
	"github.com/globalsign/mgo/bson"
	"strings"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
	"github.com/rs/zerolog"
)

type center struct {
	entityAdapter entity.Adapter
	builder       Builder
	enableEnrich  bool
	updatedAfter  time.Duration
	maxOldObjects uint64
	logger        zerolog.Logger
}

func (e *center) Handle(event types.Event, ef EnrichFields) *types.Entity {
	eventEntity := e.builder.UpdateLinkedEntities(event)
	if e.enableEnrich {
		eventEntity = e.builder.Enrich(event, ef)
	}
	eventEntity = e.builder.UpdateWatchersLinks(event, eventEntity)

	err := e.Flush()
	if err != nil {
		e.logger.Warn().Err(err).Msg("entity flush")
	}

	return eventEntity
}

func (e *center) Update(entity types.Entity) types.Entity {
	entity = e.builder.Update(entity)

	err := e.Flush()
	if err != nil {
		e.logger.Warn().Err(err).Msg("entity flush")
	}

	return entity
}

func (e *center) Get(event types.Event) (*types.Entity, error) {
	eid := event.GetEID()

	entity, err := e.entityAdapter.GetEntityByID(eid)
	if err != nil {
		_, ok := err.(errt.NotFound)
		if ok {
			return nil, nil
		}
	}

	return &entity, err
}

func (e *center) Flush() error {
	var err error

	updates := e.builder.Extract()

	//fix for duplicated upsert, caused by several instances
	for {
		for _, entityState := range updates {
			if err := e.entityAdapter.BulkUpsert(entityState.Entity, entityState.NewImpacts, entityState.NewDepends); err != nil {
				e.logger.Warn().Err(err).Msg("entity update")
			}
		}

		e.logger.Debug().Int("updated", len(updates)).Msg("")

		err = e.entityAdapter.FlushBulk()
		//TODO: should use mgo.IsDup(), but the error type is changed by mongo_adapter
		if err == nil || !strings.Contains(err.Error(), " E11000 ") {
			break
		}
	}

	return err
}

func (e *center) EnrichResourceInfoWithComponentInfo(event *types.Event, entity *types.Entity) error {
	if event.Connector == "watcher" {
		return nil
	}

	if event.SourceType == types.SourceTypeResource && entity.IsNew {
		component, _ := e.entityAdapter.Get(event.Component)
		entity.ComponentInfos = component.Infos
		err := e.entityAdapter.Update(*entity)
		if err != nil {
			return err
		}
	}

	if event.SourceType == types.SourceTypeComponent {
		for _, v := range entity.Depends {
			err := e.entityAdapter.AddToBulkUpdate(v, bson.M{"$set": bson.M{"component_infos": entity.Infos}})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *center) LoadWatchers() error {
	return e.builder.LoadWatchers()
}

// NewEnrichmentCenter creates an instance of the high level context creation api.
// param maxOldObjects: each Handle() call will check for objects that has not received updates since 10 seconds and
//    handle flush through entityAdapter
// param enableEnrich: enable extra infos enrichment
func NewEnrichmentCenter(maxOldObjects uint64, enableEnrich bool, entityAdapter entity.Adapter, watcherAdapter watcher.Adapter, logger zerolog.Logger) EnrichmentCenter {
	builder := NewBuilder(entityAdapter, watcherAdapter, logger)
	// getting the list of watchers currently in database
	err := builder.LoadWatchers()
	if err != nil {
		logger.Warn().Err(err).Msg("creating enrichment center builder, updating the watcher list")
	}
	return &center{
		builder:       builder,
		entityAdapter: entityAdapter,
		enableEnrich:  enableEnrich,
		updatedAfter:  time.Second * 10,
		maxOldObjects: maxOldObjects,
		logger:        logger,
	}
}
