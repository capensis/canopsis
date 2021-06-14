package context_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/cache"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	sconnector     = "CONN"
	sconnectorName = "CONN_I"
	scomponent     = "COMP_"
	sresource      = "RES_"
)

func getDirtyEvent(state, Ci, ci, ri int) (types.Event, []byte) {
	evt := types.Event{
		Connector:     sconnector,
		ConnectorName: sconnectorName + strconv.Itoa(Ci),
		Component:     scomponent + strconv.Itoa(ci),
		Resource:      sresource + strconv.Itoa(ri),
		State:         types.CpsNumber(state),
		SourceType:    types.SourceTypeResource,
		EventType:     types.EventTypeCheck,
	}

	source, err := json.Marshal(evt)
	if err != nil {
		panic(fmt.Errorf("getDirtyEvent: %v", err))
	}

	freeEvt := make(map[string]interface{})

	if err = json.Unmarshal(source, &freeEvt); err != nil {
		panic(fmt.Errorf("getDirtyEvent: %v", err))
	}

	freeEvt["waha"] = "ça défonce"

	source, err = json.Marshal(freeEvt)

	if err != nil {
		panic(fmt.Errorf("gitDirtyEvent: %v", err))
	}

	evt.InjectExtraInfos(source)

	return evt, source
}

type testEntityAdapter struct {
	istore  cache.Cache
	ustore  cache.Cache
	stored  map[string]context.EntityState
	flushed map[string]types.Entity
}

func (e *testEntityAdapter) Insert(entity types.Entity) error {
	//e.stored[entity.ID] = entity
	return nil
}

func (e *testEntityAdapter) Update(entity types.Entity) error {
	//e.stored[entity.ID] = entity
	return nil
}

func (e *testEntityAdapter) Remove(entity types.Entity) error {
	//delete(e.stored, entity.ID)
	return nil
}

func (e *testEntityAdapter) Get(id string) (types.Entity, bool) {
	entityState, exists := e.stored[id]
	return entityState.Entity, exists
}

func (e *testEntityAdapter) GetIDs(filter bson.M, ids *[]interface{}) error {
	panic("not implemented")
}

func (e *testEntityAdapter) GetEntityByID(id string) (types.Entity, error) {
	return e.stored[id].Entity, nil
}

func (e *testEntityAdapter) FlushBulk() error {
	for eid, entityState := range e.stored {
		e.flushed[eid] = entityState.Entity
	}
	_ = e.FlushBulkInsert()
	_ = e.FlushBulkUpdate()
	return nil
}

func (e *testEntityAdapter) FlushBulkInsert() error {
	e.istore.Flush()
	return nil
}

func (e *testEntityAdapter) FlushBulkUpdate() error {
	e.ustore.Flush()
	return nil
}

func (e *testEntityAdapter) BulkInsert(entity types.Entity) error {
	e.stored[entity.CacheID()] = context.EntityState{
		Entity: entity,
		State:  context.PendingInsert,
	}
	return e.istore.Set(entity)
}

func (e *testEntityAdapter) Count() (int, error) {
	return len(e.flushed), nil
}

func (e *testEntityAdapter) BulkUpdate(entity types.Entity) error {
	e.stored[entity.CacheID()] = context.EntityState{
		Entity: entity,
		State:  context.PendingUpdate,
	}
	return e.ustore.Set(entity)
}

func (e *testEntityAdapter) AddToBulkUpdate(entityId string, data bson.M) error {
	return nil
}

func (e *testEntityAdapter) BulkUpsert(entity types.Entity, newImpacts []string, newDepends []string) error {
	//simulate $addToSet
	storedEntityState, ok := e.stored[entity.CacheID()]
	if ok {
		for _, storedEntityImpact := range storedEntityState.Entity.Impacts {
			dup := false

			for _, entityImpact := range entity.Impacts {
				if entityImpact == storedEntityImpact {
					dup = true
					break
				}
			}

			if !dup {
				entity.Impacts = append(entity.Impacts, storedEntityImpact)
			}
		}

		for _, storedEntityDepend := range storedEntityState.Entity.Depends {
			dup := false

			for _, entityDepend := range entity.Depends {
				if entityDepend == storedEntityDepend {
					dup = true
					break
				}
			}

			if !dup {
				entity.Depends = append(entity.Depends, storedEntityDepend)
			}
		}
	}

	e.stored[entity.CacheID()] = context.EntityState{
		Entity:     entity,
		State:      context.PendingUpdate,
		NewImpacts: newImpacts,
		NewDepends: newDepends,
	}
	return e.ustore.Set(entity)
}

func (e *testEntityAdapter) RemoveAll() error {
	e.stored = make(map[string]context.EntityState)
	e.ustore.Flush()
	e.istore.Flush()
	return nil
}

func newTestEntityAdapter() entity.Adapter {
	return &testEntityAdapter{
		istore:  cache.NewKV(),
		ustore:  cache.NewKV(),
		stored:  make(map[string]context.EntityState),
		flushed: make(map[string]types.Entity),
	}
}

type testWatcherAdapter struct {
	Watchers map[string]watcher.Watcher
	Entities []types.Entity
}

func (e *testWatcherAdapter) Insert(aWatcher watcher.Watcher) error {
	e.Watchers[aWatcher.ID] = aWatcher
	return nil
}

func (e *testWatcherAdapter) GetAll(watchers *[]watcher.Watcher) error {
	values := []watcher.Watcher{}
	for _, value := range e.Watchers {
		values = append(values, value)
	}
	*watchers = values
	return nil
}

func (e *testWatcherAdapter) GetAllValidWatchers(watchers *[]watcher.Watcher) error {
	return e.GetAll(watchers)
}

func (e testWatcherAdapter) GetByID(id string, aWatcher *watcher.Watcher) error {
	_, ok := e.Watchers[id]
	if ok {
		*aWatcher = e.Watchers[id]
	}
	return nil
}

func (e *testWatcherAdapter) GetEntities(aWatcher watcher.Watcher, entities *[]types.Entity) error {
	for _, watchedEntity := range e.Entities {
		if aWatcher.CheckEntityInWatcher(watchedEntity) {
			*entities = append(*entities, watchedEntity)
		}
	}

	return nil
}

func (e *testWatcherAdapter) GetAllAnnotatedEntities() ([]watcher.AnnotatedEntity, error) {
	return []watcher.AnnotatedEntity{}, nil
}

func (e *testWatcherAdapter) GetAnnotatedEntitiesIter() *mgo.Iter {
	return nil
}

func (e *testWatcherAdapter) GetAnnotatedDependencies(watcherID string) ([]watcher.AnnotatedEntity, error) {
	return []watcher.AnnotatedEntity{}, nil
}

func (e *testWatcherAdapter) Update(_ string, _ interface{}) error {
	return nil
}

func (e *testWatcherAdapter) GetAnnotatedDependenciesIter(_ string) *mgo.Iter {
	return nil
}

func newTestWatcherAdapter() *testWatcherAdapter {
	return &testWatcherAdapter{Watchers: make(map[string]watcher.Watcher), Entities: nil}
}

func newTestWatcherAdapterWithEntities(entities []types.Entity) *testWatcherAdapter {
	return &testWatcherAdapter{Watchers: make(map[string]watcher.Watcher), Entities: entities}
}

// Simulate a full enrichment job with only local caches to avoid
// problems with external libraries
func TestEnrichmentHandle(t *testing.T) {
	Convey("Setup", t, func() {
		entityAdapter := newTestEntityAdapter()
		watcherAdapter := newTestWatcherAdapter()
		ec := context.NewEnrichmentCenter(1000, true, entityAdapter, watcherAdapter, log.NewTestLogger())
		ef := context.NewEnrichFields("", "")

		nConnectors := 3
		nComponents := 4
		nResources := 5

		msg := fmt.Sprintf(
			"Generate %dConnectors * %dComponents * %dResources = %dEvents to generate context",
			nConnectors,
			nComponents,
			nResources,
			nConnectors*nComponents*nResources,
		)
		Convey(msg, func() {
			for Ci := 0; Ci < nConnectors; Ci++ {
				for ci := 0; ci < nComponents; ci++ {
					for ri := 0; ri < nResources; ri++ {
						evt, _ := getDirtyEvent(3, Ci, ci, ri)
						So(ec.Handle(evt, ef), ShouldNotBeNil)
					}
				}
			}

			tea, ok := entityAdapter.(*testEntityAdapter)
			So(ok, ShouldBeTrue)

			Convey("Check for expected data", func() {
				// Check the total amount of entities
				So(len(tea.stored), ShouldEqual, nConnectors+(nComponents*nResources)+nComponents)

				// Split connectors, components and resources entities to check their content
				connectors := make(map[string]context.EntityState)
				components := make(map[string]context.EntityState)
				resources := make(map[string]context.EntityState)

				for eid, entityState := range tea.stored {
					switch entityState.Entity.Type {
					case types.EntityTypeConnector:
						connectors[eid] = entityState
					case types.EntityTypeComponent:
						components[eid] = entityState
					case types.EntityTypeResource:
						resources[eid] = entityState
					}
				}

				So(len(connectors), ShouldEqual, nConnectors)
				So(len(components), ShouldEqual, nComponents)
				So(len(resources), ShouldEqual, nResources*nComponents)

				for _, connector := range connectors {
					So(connector.Entity.Type, ShouldEqual, types.EntityTypeConnector)

					impacts := connector.Entity.Impacts
					depends := connector.Entity.Depends

					So(len(impacts), ShouldEqual, nResources*nComponents)
					So(len(depends), ShouldEqual, nComponents)

					for nComponent := 0; nComponent < nComponents; nComponent++ {
						So(depends, ShouldContain, fmt.Sprintf("%s%d", scomponent, nComponent))
						for nResource := 0; nResource < nResources; nResource++ {
							So(impacts, ShouldContain, fmt.Sprintf("%s%d/%s%d", sresource, nResource, scomponent, nComponent))
						}
					}
				}

				for _, component := range components {
					So(component.Entity.Type, ShouldEqual, types.EntityTypeComponent)

					impacts := component.Entity.Impacts
					depends := component.Entity.Depends

					So(len(impacts), ShouldEqual, 3)
					So(len(depends), ShouldEqual, 5)

					for nResource := 0; nResource < nResources; nResource++ {
						So(depends, ShouldContain, fmt.Sprintf("%s%d/%s", sresource, nResource, component.Entity.ID))
					}

					for nConnector := 0; nConnector < nConnectors; nConnector++ {
						So(impacts, ShouldContain, fmt.Sprintf("%s/%s%d", sconnector, sconnectorName, nConnector))
					}
				}

				for _, resource := range resources {
					So(resource.Entity.Type, ShouldEqual, types.EntityTypeResource)

					impacts := resource.Entity.Impacts
					depends := resource.Entity.Depends

					So(len(impacts), ShouldEqual, 1)
					componentID := strings.Split(resource.Entity.ID, "/")[1]
					So(impacts[0], ShouldEqual, componentID)

					So(resource.Entity.Infos, ShouldContainKey, "waha")
					So(len(resource.Entity.Infos), ShouldEqual, 1)

					for nConnector := 0; nConnector < nConnectors; nConnector++ {
						So(depends, ShouldContain, fmt.Sprintf("%s/%s%d", sconnector, sconnectorName, nConnector))
					}

				}
			})
		})
	})
}
