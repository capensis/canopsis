package eventfilter_test

import (
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

type MockEnrichmentCenter struct {
	HandleCalls int
	Entity      *types.Entity
}

func (e *MockEnrichmentCenter) Handle(event types.Event, ef context.EnrichFields) *types.Entity {
	e.HandleCalls++
	return e.Entity
}

func (e *MockEnrichmentCenter) Update(entity types.Entity) types.Entity {
	return entity
}

func (e *MockEnrichmentCenter) Flush() error {
	return nil
}

func (e *MockEnrichmentCenter) LoadWatchers() error {
	return nil
}

func (e *MockEnrichmentCenter) Get(event types.Event) (*types.Entity, error) {
	return nil, nil
}

func (e *MockEnrichmentCenter) EnrichResourceInfoWithComponentInfo(event *types.Event, entity *types.Entity) error {
	return nil
}

func TestEntityDataSourceCreate(t *testing.T) {
	Convey("Given an entity data source factory", t, func() {
		entity := types.NewEntity("component", "name", types.EntityTypeComponent, nil, nil, nil)
		enrichmentCenter := MockEnrichmentCenter{
			Entity: &entity,
		}
		enrichFields := context.NewEnrichFields("", "")

		factory := eventfilter.NewEntityDataSourceFactory(&enrichmentCenter, enrichFields)

		Convey("Creating an entity data source with parameters returns an error", func() {
			_, err := factory.Create(map[string]interface{}{
				"unexpected_parameters": "test",
			})
			So(err, ShouldNotBeNil)
		})
	})
}

func TestEntityDataSourceGet(t *testing.T) {
	Convey("Given an entity data source factory", t, func() {
		entity := types.NewEntity("component", "name", types.EntityTypeComponent, nil, nil, nil)
		enrichmentCenter := MockEnrichmentCenter{
			Entity: &entity,
		}
		enrichFields := context.NewEnrichFields("", "")

		factory := eventfilter.NewEntityDataSourceFactory(&enrichmentCenter, enrichFields)

		Convey("Creating an entity data source without parameters succeeds", func() {
			source, err := factory.Create(map[string]interface{}{})
			So(err, ShouldBeNil)
			So(source, ShouldNotBeNil)

			Convey("Getting an event's entity succeeds", func() {
				parameters := eventfilter.DataSourceGetterParameters{
					Event: types.Event{
						EventType: types.EventTypeCheck,
						SourceType: types.SourceTypeComponent,
						Component:  "component",
						State:      3,
						Debug:      true,
					},
				}

				entity, err := source.Get(parameters)
				So(err, ShouldBeNil)
				So(enrichmentCenter.HandleCalls, ShouldEqual, 1)

				typedEntity, isEntity := entity.(types.Entity)
				So(isEntity, ShouldBeTrue)
				So(typedEntity.ID, ShouldEqual, "component")
			})
		})
	})
}
