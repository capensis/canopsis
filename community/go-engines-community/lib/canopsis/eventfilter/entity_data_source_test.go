package eventfilter_test

import (
	"context"
	mock_context "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/context"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"testing"

	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEntityDataSourceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given an entity data source factory", t, func() {
		enrichmentCenter := mock_context.NewMockEnrichmentCenter(ctrl)
		enrichFields := libcontext.NewEnrichFields("", "")

		dbClient := mock_mongo.NewMockDbClient(ctrl)

		factory := eventfilter.NewEntityDataSourceFactory(enrichmentCenter, enrichFields)

		Convey("Creating an entity data source with parameters returns an error", func() {
			_, err := factory.Create(dbClient, map[string]interface{}{
				"unexpected_parameters": "test",
			})
			So(err, ShouldNotBeNil)
		})
	})
}

func TestEntityDataSourceGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	Convey("Given an entity data source factory", t, func() {
		entity := types.NewEntity("component", "name", types.EntityTypeComponent, nil, nil, nil)
		enrichmentCenter := mock_context.NewMockEnrichmentCenter(ctrl)
		enrichmentCenter.EXPECT().Handle(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity, libcontext.UpdatedEntityServices{}, nil)
		enrichFields := libcontext.NewEnrichFields("", "")

		factory := eventfilter.NewEntityDataSourceFactory(enrichmentCenter, enrichFields)

		Convey("Creating an entity data source without parameters succeeds", func() {
			dbClient := mock_mongo.NewMockDbClient(ctrl)

			source, err := factory.Create(dbClient, map[string]interface{}{})
			So(err, ShouldBeNil)
			So(source, ShouldNotBeNil)

			Convey("Getting an event's entity succeeds", func() {
				parameters := eventfilter.DataSourceGetterParameters{
					Event: types.Event{
						EventType:  types.EventTypeCheck,
						SourceType: types.SourceTypeComponent,
						Component:  "component",
						State:      3,
						Debug:      true,
					},
				}

				entity, err := source.Get(context.Background(), parameters, nil)
				So(err, ShouldBeNil)

				typedEntity, isEntity := entity.(types.Entity)
				So(isEntity, ShouldBeTrue)
				So(typedEntity.ID, ShouldEqual, "component")
			})
		})
	})
}
