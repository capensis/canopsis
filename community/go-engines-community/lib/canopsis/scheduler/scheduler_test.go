package scheduler

import (
	"context"
	"testing"
	"time"

	amqpLib "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	. "github.com/smartystreets/goconvey/convey"
)

func testNewSchedulerService() (Scheduler, QueueLock) {
	ctx := context.Background()
	logger := log.NewTestLogger()

	redisLockStorage, err := redis.NewSession(ctx, redis.LockStorage, logger, 0, 0)
	if err != nil {
		panic(err)
	}

	redisQueueStorage, err := redis.NewSession(ctx, redis.QueueStorage, logger, 0, 0)
	if err != nil {
		panic(err)
	}

	queue := NewQueueLock(redisLockStorage, time.Second, redisQueueStorage, true, logger)

	amqpSession, err := amqpLib.NewConnection(log.NewLogger(false), 0, 0)
	if err != nil {
		panic(err)
	}

	pubChannel, err := amqpSession.Channel()
	if err != nil {
		panic(err)
	}

	shd := scheduler{
		redisConn:      redisLockStorage,
		channelPub:     pubChannel,
		publishToQueue: "",
		queueLock:      queue,
		logger:         logger,
		jsonDecoder:    json.NewDecoder(),
	}
	shd.subscribe(ctx)

	return &shd, queue
}

func TestScheduler(t *testing.T) {
	ctx := context.Background()
	shd, queue := testNewSchedulerService()

	bytesEvent := []byte(`{"_id":"testEvent","component":"testschedulerComponent","connector":"testschedulerConnector"}`)
	event := types.Event{}

	amqpSession, err := amqpLib.NewConnection(log.NewLogger(false), 0, 0)
	if err != nil {
		panic(err)
	}

	pubChannel, err := amqpSession.Channel()
	if err != nil {
		panic(err)
	}

	Convey("scheduler should process event without errors", t, func() {
		decoder := json.NewDecoder()
		err := decoder.Decode(bytesEvent, &event)
		So(err, ShouldBeNil)

		lockID := event.GetLockID()
		err = shd.ProcessEvent(ctx, pubChannel, lockID, bytesEvent)
		So(err, ShouldBeNil)

		Convey("Then event should be locked and should not be queued", func() {
			locked := queue.IsLocked(ctx, event.GetLockID())
			So(locked, ShouldBeTrue)

			queued := !queue.IsEmpty(ctx, event.GetLockID())
			So(queued, ShouldBeFalse)

			Convey("Then event should be acked without errors and should not be locked", func() {
				err = shd.AckEvent(ctx, pubChannel, event)
				So(err, ShouldBeNil)

				//sleep one second, because unlock processed in goroutine
				time.Sleep(time.Second * 1)

				locked = queue.IsLocked(ctx, event.GetLockID())
				So(locked, ShouldBeFalse)
			})
		})

	})

	encoder := json.NewEncoder()

	Convey("scheduler should must process two messages related to the same alarm", t, func() {
		event1 := &types.Event{
			ID:        strPtr("testEvent1"),
			Component: "testschedulerComponent1",
			Connector: "testschedulerConnector1",
		}
		bytesEvent1, err := encoder.Encode(&event1)
		event2 := &types.Event{
			ID:        strPtr("testEvent2"),
			Component: "testschedulerComponent1",
			Connector: "testschedulerConnector1",
		}
		bytesEvent2, err := encoder.Encode(&event2)
		lockID1 := event1.GetLockID()
		lockID2 := event2.GetLockID()

		err = shd.ProcessEvent(ctx, pubChannel, lockID1, bytesEvent1)
		So(err, ShouldBeNil)

		err = shd.ProcessEvent(ctx, pubChannel, lockID2, bytesEvent2)
		So(err, ShouldBeNil)

		Convey("Then event should be locked and one event queued", func() {
			locked := queue.IsLocked(ctx, event1.GetLockID())
			So(locked, ShouldBeTrue)

			queued := !queue.IsEmpty(ctx, event1.GetLockID())
			So(queued, ShouldBeTrue)

			Convey("Then ack message, event should be locked and queue should be empty", func() {
				err = shd.AckEvent(ctx, pubChannel, *event1)
				So(err, ShouldBeNil)

				locked = queue.IsLocked(ctx, event1.GetLockID())
				So(locked, ShouldBeTrue)

				queued = !queue.IsEmpty(ctx, event1.GetLockID())
				So(queued, ShouldBeFalse)
			})
		})

	})
}

func strPtr(v string) *string {
	return &v
}
