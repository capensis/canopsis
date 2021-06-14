package main_test

import (
	"testing"
	"time"

	hb "git.canopsis.net/canopsis/go-engines/cmd/engine-heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/streadway/amqp"
)

type testSendAlarm struct {
	Alarms []string
	Status []int
}

func (tsa *testSendAlarm) sendAlarmTest(alarmid string, state int, output string) error {
	tsa.Alarms = append(tsa.Alarms, alarmid)
	tsa.Status = append(tsa.Status, state)
	return nil
}

type TestEngineHeartBeat struct {
	*hb.EngineHeartBeat
}

func (e *TestEngineHeartBeat) ConsumerChan() (<-chan amqp.Delivery, error) {
	return make(chan amqp.Delivery), nil
}

func testNewEngineHeartBeat() (*TestEngineHeartBeat, *testSendAlarm, error) {
	tsa := testSendAlarm{
		Alarms: make([]string, 0),
	}

	depMaker := hb.DependencyMaker{}
	references := depMaker.GetDefaultReferences(log.NewTestLogger())

	engine := TestEngineHeartBeat{
		EngineHeartBeat: hb.NewEngineHeartBeat(references),
	}
	engine.SendAlarmFunc = tsa.sendAlarmTest
	_, err := engine.ConsumerChan()
	return &engine, &tsa, err
}

func TestInitializeHeartBeat(t *testing.T) {
	Convey("Setup", t, func() {

		engine, _, err := testNewEngineHeartBeat()

		So(err, ShouldBeNil)

		li := heartbeat.NewItem(time.Second * 10)
		li.Mappings = map[string]string{"connector": "c1"}

		err = engine.AddHeartBeatItem(li)
		So(err, ShouldBeNil)

		Convey("Initialize will set correct warmups", func() {
			err := engine.Initialize()
			So(err, ShouldBeNil)

			r := engine.References.Redis.Exists("connector:c1")
			So(r.Err(), ShouldBeNil)
			So(r.Val(), ShouldEqual, 1)

			rttl := engine.References.Redis.TTL("connector:c1")
			So(rttl.Err(), ShouldBeNil)
			So(rttl.Val(), ShouldBeBetweenOrEqual, time.Second*9, time.Second*10)
		})
	})
}

func TestWorkHeartBeat(t *testing.T) {
	Convey("Setup", t, func() {
		engine, _, err := testNewEngineHeartBeat()

		So(err, ShouldBeNil)

		li1 := heartbeat.NewItem(time.Second * 10)
		li1.AddMapping("connector", "c1")
		li1.AddMapping("connector_name", "c1")

		li2 := heartbeat.NewItem(time.Second * 10)
		li2.AddMapping("connector", "c2")
		li2.AddMapping("connector_name", "c2")

		li3 := heartbeat.NewItem(time.Second * 10)
		li3.AddMapping("connector", "c3")
		li3.AddMapping("connector_name", "c3")

		engine.AddHeartBeatItem(li1)
		engine.AddHeartBeatItem(li2)
		engine.AddHeartBeatItem(li3)

		Convey("Given an event matching LI1, i must set an ID in Redis", func() {
			engine.References.Redis.FlushDB()
			bevent := []byte(`{"connector": "c2", "connector_name": "c2"}`)
			msg := amqp.Delivery{Body: bevent}
			engine.WorkerProcess(msg)

			Convey("I have only li2.ID() set", func() {
				r := engine.References.Redis.Exists(li2.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 1)
				rttl := engine.References.Redis.TTL(li2.ID())
				So(rttl.Err(), ShouldBeNil)
				So(rttl.Val(), ShouldBeBetweenOrEqual, li1.MaxDuration-time.Second-1, li1.MaxDuration)
				r = engine.References.Redis.Exists(li1.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
				r = engine.References.Redis.Exists(li3.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
			})
		})

		Convey("Given an event with a null field that match LI2 second mapping", func() {
			engine.References.Redis.FlushDB()
			bevent := []byte(`{"connector": null, "connector_name": "c2"}`)
			msg := amqp.Delivery{Body: bevent}
			engine.WorkerProcess(msg)

			Convey("I don't have any id set", func() {
				r := engine.References.Redis.Exists(li2.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
				r = engine.References.Redis.Exists(li2.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
				r = engine.References.Redis.Exists(li3.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
			})
		})

		Convey("Given an event with an absent field", func() {
			engine.References.Redis.FlushDB()
			bevent := []byte(`{"connector_name": "testcn"}`)
			msg := amqp.Delivery{Body: bevent}
			engine.WorkerProcess(msg)

			Convey("I don't have any id set", func() {
				r := engine.References.Redis.Exists(li3.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
				r = engine.References.Redis.Exists(li2.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
				r = engine.References.Redis.Exists(li3.ID())
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
			})
		})
	})
}

func TestHeartBeatPeriodicalProcess(t *testing.T) {
	Convey("Setup", t, func() {
		engine, tsa, err := testNewEngineHeartBeat()
		So(err, ShouldBeNil)
		So(tsa, ShouldNotBeNil)
		engine.References.Redis.FlushDB()

		status := engine.References.Redis.Set("connector:testc.connector_name:testcn", 1, time.Millisecond*500)
		So(status.Err(), ShouldBeNil)

		li1 := heartbeat.NewItem(time.Second)
		li1.AddMapping("connector", "testc")
		li1.AddMapping("connector_name", "testcn")

		li2 := heartbeat.NewItem(time.Second)
		li2.AddMapping("component", "testC")
		li2.AddMapping("resource", "testR")

		So(engine.AddHeartBeatItem(li1), ShouldBeNil)
		So(engine.AddHeartBeatItem(li2), ShouldBeNil)
		Convey("When i run PeriodicalProcess a first time", func() {
			engine.PeriodicalProcess()
			So(len(tsa.Alarms), ShouldEqual, 2)
			So(tsa.Alarms[0], ShouldEqual, "connector:testc.connector_name:testcn")
			So(tsa.Status[0], ShouldEqual, 0)
			So(tsa.Alarms[1], ShouldEqual, "component:testC.resource:testR")
			So(tsa.Status[1], ShouldEqual, types.AlarmStateCritical)

			Convey("When i send an event matching first alarm", func() {
				evt := amqp.Delivery{
					Body: []byte(`{"connector":"testc","connector_name":"testcn"}`),
				}

				r := engine.References.Redis.Exists("alarm:connector:testc.connector_name:testcn")
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)

				engine.WorkerProcess(evt)
				So(len(tsa.Alarms), ShouldEqual, 2)
				So(tsa.Alarms[0], ShouldEqual, "connector:testc.connector_name:testcn")
				So(tsa.Status[0], ShouldEqual, 0)
				So(tsa.Status[1], ShouldEqual, 3)

				r = engine.References.Redis.Exists("alarm:" + tsa.Alarms[0])
				So(r.Err(), ShouldBeNil)
				So(r.Val(), ShouldEqual, 0)
			})
		})
	})
}

func BenchmarkWorkHeartBeat(b *testing.B) {
	engine, _, err := testNewEngineHeartBeat()
	engine.References.Redis.FlushDB()

	if err != nil {
		b.Error(err)
	}

	li1 := heartbeat.NewItem(time.Millisecond * 10)
	li1.AddMapping("connector", "test")
	li1.AddMapping("connector_name", "test")

	li2 := heartbeat.NewItem(time.Millisecond * 10)
	li2.AddMapping("connector", "none")
	li2.AddMapping("component", "none")

	engine.AddHeartBeatItem(li1)
	engine.AddHeartBeatItem(li2)

	for i := 0; i < b.N; i++ {
		sevent := `{
		"connector": "test",
		"connector_name": "test",
		"source_type": "resource",
		"component": "localhost",
		"resource": "benchmarkheartBeatItemid",
		"event_type": "check",
		"status": 0,
		"output": "Mieux, mieux : Un tiers chou, un tiers tulipe, et un tiers blette !"
	}`

		engine.WorkerProcess(amqp.Delivery{
			Body: []byte(sevent),
		})
	}
}
