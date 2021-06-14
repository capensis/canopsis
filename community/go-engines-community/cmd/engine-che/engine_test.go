package main_test

import (
	"math/rand"
	"strconv"
	"testing"

	che "git.canopsis.net/canopsis/go-engines/cmd/engine-che"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/depmake"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/streadway/amqp"
)

type EngineTestChe struct {
	*che.EngineChe
}

func (e *EngineTestChe) ConsumerChan() (<-chan amqp.Delivery, error) {
	ch := make(chan amqp.Delivery)
	return ch, nil
}

func testNewEngineChe() (*EngineTestChe, error) {
	options := che.Options{
		FeatureContextCreation: true,
		FeatureContextEnrich:   true,
		DataSourceDirectory:    ".",
	}
	depMaker := che.DependencyMaker{}
	references := depMaker.GetDefaultReferences(options, log.NewTestLogger())

	engine := EngineTestChe{
		EngineChe: che.NewEngineCHE(options, references),
	}
	_, err := engine.ConsumerChan()
	return &engine, err
}

func TestInitializeChe(t *testing.T) {
	Convey("Setup", t, func() {
		engine, err := testNewEngineChe()
		So(err, ShouldBeNil)

		Convey("Initialize will set correct warmups", func() {
			err := engine.Initialize()
			So(err, ShouldBeNil)
		})
	})
}

func TestWorkChe(t *testing.T) {
	Convey("Setup", t, func() {
		engine, err := testNewEngineChe()
		So(err, ShouldBeNil)
		err = engine.Initialize()
		So(err, ShouldBeNil)

		maker := depmake.DependencyMaker{}
		mongoSession := maker.DepMongoSession()
		entityCollection := entity.DefaultCollection(mongoSession)
		entityService := entity.NewAdapter(entityCollection)

		err = entityService.Insert(types.Entity{ID: "force_collection_create"})
		if err != nil {
			So(err.Error(), ShouldStartWith, "E11000 duplicate key error collection")
		}

		Convey("Given a new event, it must create a context behind it", func() {
			So(err, ShouldBeNil)
			bevent := []byte(`{
					"event_type": "check",
					"source_type": "component",
					"state": 1,
					"connector": "air",
					"connector_name": "strike",
					"component": "electron",
					"spin": "1/2",
					"mass": "9.10938356(11)×10−31 kg"
				}`)
			event, err := types.NewEventFromJSON(bevent)
			So(err, ShouldBeNil)
			newEvent, err := engine.EnrichContextFromEvent(event)
			So(err, ShouldBeNil)
			So(newEvent.Connector, ShouldEqual, "air")
			So(newEvent.Component, ShouldEqual, "electron")
			engine.WorkerProcess(amqp.Delivery{
				Body: bevent,
			})
			err = engine.References.EnrichmentCenter.Flush()
			So(err, ShouldBeNil)

			entity, found := entityService.Get("air/strike")
			So(entity, ShouldNotBeNil)
			So(entity.ID, ShouldNotEqual, "")
			So(entity.Type, ShouldEqual, types.EntityTypeConnector)
			So(found, ShouldBeTrue)

			entity, found = entityService.Get("electron")
			So(entity, ShouldNotBeNil)
			So(entity.Type, ShouldEqual, types.EntityTypeComponent)
			So(found, ShouldBeTrue)

			So(entity.Infos, ShouldContainKey, "spin")
			So(entity.Infos["spin"].Value, ShouldEqual, "1/2")
			So(entity.Infos["spin"].RealValue, ShouldEqual, "1/2")
			So(entity.Infos, ShouldNotContainKey, "parity")

			Convey("Given a richer event on the same component", func() {
				bevent := []byte(`{
					"event_type": "check",
					"source_type": "component",
					"state": 1,
					"connector": "air",
					"connector_name": "strike",
					"component": "electron",
					"spin": "−1/2",
					"charge": "−1,602 × 10^−19 C"
				}`)
				event, err := types.NewEventFromJSON(bevent)
				So(err, ShouldBeNil)
				newEvent, err := engine.EnrichContextFromEvent(event)
				So(err, ShouldBeNil)
				So(newEvent.Component, ShouldEqual, "electron")
				engine.WorkerProcess(amqp.Delivery{
					Body: bevent,
				})
				err = engine.References.EnrichmentCenter.Flush()
				So(err, ShouldBeNil)

				Convey("Then it must enrich the context without destroying old informations", func() {
					entity2, found := entityService.Get("electron")
					So(entity2, ShouldNotBeNil)
					So(entity2.Type, ShouldEqual, types.EntityTypeComponent)
					So(found, ShouldBeTrue)

					So(entity2.Infos, ShouldContainKey, "charge")
					charge := entity2.Infos["charge"]
					So(charge, ShouldNotBeNil)
					So(charge.Value, ShouldEqual, "−1,602 × 10^−19 C")

					So(entity2.Infos, ShouldContainKey, "spin")
					spin := entity2.Infos["spin"]
					So(spin, ShouldNotBeNil)
					So(spin.Value, ShouldEqual, "−1/2")

					So(entity2.Infos, ShouldContainKey, "mass")
					mass := entity2.Infos["mass"]
					So(mass, ShouldNotBeNil)
					So(mass.Value, ShouldEqual, "9.10938356(11)×10−31 kg")
				})
			})
		})

		Convey("Given a bad event, it must not create a context", func() {
			//sessions.Redis.FlushDB()
			bevent := []byte(`{"stevie": "wonder"}`)
			event, err := types.NewEventFromJSON(bevent)
			So(err, ShouldBeNil)
			event.EventType = ""

			newEvent, err := engine.EnrichContextFromEvent(event)
			So(err, ShouldBeNil)
			So(newEvent.Entity, ShouldBeNil)
		})
	})
}

func TestEmptyEntityInfos(t *testing.T) {
	Convey("Given an entity without infos", t, func() {
		info := types.NewInfo("talking", "", "heads")
		entity := types.NewEntity("the", "great", "curve", nil, nil, nil)
		So(entity.Infos, ShouldBeEmpty)

		Convey("Then we should initialise infos to populate it", func() {
			_, exists := entity.Infos[info.Name]
			So(exists, ShouldBeFalse)

			entity.Infos = make(map[string]types.Info, 0)
			_, exists = entity.Infos[info.Name]
			So(exists, ShouldBeFalse)
			entity.Infos[info.Name] = info
			_, exists = entity.Infos[info.Name]
			So(exists, ShouldBeTrue)
		})
	})
}

func BenchmarkWorkChe(b *testing.B) {
	engine, err := testNewEngineChe()
	if err != nil {
		b.Error(err)
	}

	engine.Initialize()

	conn := 1
	comp := 1
	for i := 1; i < b.N; i++ {
		if i%10 == 0 {
			if i%100 == 0 {
				conn = conn + 1
			} else {
				comp = comp + 1
			}
		}

		sevent := `{
			"event_type": "check",
			"source_type": "resource",
			"connector": "perceval_` + strconv.Itoa(conn) + `",
			"connector_name": "pellinor",
			"component": "gallois_` + strconv.Itoa(comp) + `",
			"resource": "benchmark_work_che_` + strconv.Itoa(i) + `",
			"state": ` + strconv.Itoa(rand.Intn(3)) + `,
			"output": "Selon comme on est tourné, ça change tout !"
		}`

		engine.WorkerProcess(amqp.Delivery{
			Body: []byte(sevent),
		})
	}
}
