package storage

import (
	"context"
	"fmt"
	"math"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	redisV8 "github.com/go-redis/redis/v8"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStorageSetGet(t *testing.T) {
	storage := NewRedisGroupingStorage()
	ctx := context.Background()

	redisClient, err := redis.NewSession(ctx, redis.AlarmGroupStorage, log.NewLogger(true), 0, 0)
	if err != nil {
		panic(err)
	}
	res := redisClient.FlushDB(ctx)
	if res.Err() != nil {
		panic(err)
	}

	Convey("Test basic manipulations with storage", t, func() {
		testRule := metaalarm.Rule{
			ID: "test_rule",
		}

		testRule2 := metaalarm.Rule{
			ID: "test_rule_2",
		}

		testAlarm := types.Alarm{
			ID: "test_alarm",
		}

		testAlarm2 := types.Alarm{
			ID: "test_alarm_2",
		}

		testAlarm3 := types.Alarm{
			ID: "test_alarm_2",
		}

		testAlarm4 := types.Alarm{
			ID: "test_alarm_4",
		}

		testAlarm5 := types.Alarm{
			ID: "test_alarm_5",
		}

		testAlarm6 := types.Alarm{
			ID: "test_alarm_6",
		}

		_ = redisClient.Watch(ctx, func(tx *redisV8.Tx) error {
			alarmGroup, err := storage.Get(ctx, tx, "test_rule39")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 0)
			So(alarmGroup.GetKey(), ShouldEqual, "test_rule39")

			alarmGroup = NewAlarmGroup(testRule.ID)
			err = storage.Set(ctx, tx, alarmGroup, 60)
			So(err, ShouldBeNil)
			alarmGroup, err = storage.Get(ctx, tx, "test_rule")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 0)

			alarmGroup.Push(testAlarm, 60)
			err = storage.Set(ctx, tx, alarmGroup, 60)
			So(err, ShouldBeNil)

			alarmGroup, err = storage.Get(ctx, tx, "test_rule")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 1)

			alarmGroup = NewAlarmGroup(testRule2.ID)
			alarmGroup.Push(testAlarm, 60)
			alarmGroup.Push(testAlarm2, 60)
			err = storage.Set(ctx, tx, alarmGroup, 60)
			So(err, ShouldBeNil)

			alarmGroup, err = storage.Get(ctx, tx, "test_rule_2")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 2)

			alarmGroup.Push(testAlarm2, 60)
			err = storage.Set(ctx, tx, alarmGroup, 60)
			So(err, ShouldBeNil)

			alarmGroup, err = storage.Get(ctx, tx, "test_rule_2")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 2)

			alarmGroup = NewAlarmGroup(testRule2.ID)
			alarmGroup.Push(testAlarm3, 60)
			err = storage.Set(ctx, tx, alarmGroup, 60)
			So(err, ShouldBeNil)

			alarmGroup, err = storage.Get(ctx, tx, "test_rule_2")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 1)

			alarmGroup1 := NewAlarmGroup("test_set_many_1")
			alarmGroup1.Push(testAlarm4, 60)
			alarmGroup1.Push(testAlarm5, 60)
			alarmGroup2 := NewAlarmGroup("test_set_many_2")
			alarmGroup2.Push(testAlarm6, 60)

			err = storage.SetMany(ctx, tx, 60, alarmGroup1, alarmGroup2)
			So(err, ShouldBeNil)

			alarmGroup, err = storage.Get(ctx, tx, "test_set_many_1")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 2)

			alarmGroup, err = storage.Get(ctx, tx, "test_set_many_2")
			So(err, ShouldBeNil)
			So(alarmGroup.GetGroupLength(), ShouldEqual, 1)

			return nil
		}, "test_key")
	})
}

func TestStorageShiftTimeInterval(t *testing.T) {
	Convey("Test time-interval shifting: basic grouping logic", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := time.Now()

		times := make(map[int]types.CpsTime)
		times[0] = types.NewCpsTime(now.Unix())
		times[1] = types.NewCpsTime(now.Add(time.Second * 10).Unix())
		times[2] = types.NewCpsTime(now.Add(time.Second * 20).Unix())
		times[3] = types.NewCpsTime(now.Add(time.Second * 30).Unix())
		times[4] = types.NewCpsTime(now.Add(time.Second * 40).Unix())
		times[5] = types.NewCpsTime(now.Add(time.Second * 50).Unix())

		//fill alarm group, should be sorted by time
		alarmGroup := NewAlarmGroup("test")
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_3",
			Value: types.AlarmValue{
				LastUpdateDate: times[3],
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_0",
			Value: types.AlarmValue{
				LastUpdateDate: times[0],
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_5",
			Value: types.AlarmValue{
				LastUpdateDate: times[5],
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_2",
			Value: types.AlarmValue{
				LastUpdateDate: times[2],
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_4",
			Value: types.AlarmValue{
				LastUpdateDate: times[4],
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_1",
			Value: types.AlarmValue{
				LastUpdateDate: times[1],
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetTimes()), ShouldEqual, 6)
		So(alarmGroup.GetOpenTime() == now.Unix(), ShouldBeTrue)

		for idx := 0; idx < 6; idx++ {
			So(alarmGroup.GetTimes()[idx] == times[idx].Unix(), ShouldBeTrue)
			So(alarmGroup.GetAlarmIds()[idx] == fmt.Sprintf("test_alarm_%d", idx), ShouldBeTrue)
		}

		//This call should shift time interval => so the storage should delete the first alarm in the Group
		// and update 'create' time, since the first alarm won't be in the minute time window anymore.
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_6",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 65).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 6)
		So(alarmGroup.GetOpenTime() == now.Add(time.Second*10).Unix(), ShouldBeTrue)

		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_7",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 300).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 1)
		So(alarmGroup.GetOpenTime() == now.Add(time.Second*300).Unix(), ShouldBeTrue)
	})

	Convey("Test time-interval shifting: check that open time is changed if the next event is older than previous", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		alarmGroup := NewAlarmGroup("test")
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_2",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * -10).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 2)
		So(alarmGroup.GetOpenTime() == now.Add(time.Second*-10).Unix(), ShouldBeTrue)
	})

	Convey("Test time-interval shifting: check that time-interval is shifted properly, if the new alarm is late and no alarm should be missed", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		alarmGroup := NewAlarmGroup("test")
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_2",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 10).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_3",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 20).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_4",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 40).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 4)
		So(alarmGroup.GetOpenTime() == now.Unix(), ShouldBeTrue)

		//test_alarm_4 will be missed, so there shouldn't be any interval shifting
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_5",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * -30).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 4)
		So(alarmGroup.GetOpenTime() == now.Unix(), ShouldBeTrue)

		//Interval can be shifted, since none alarm will be lost
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_5",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * -5).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 5)
		So(alarmGroup.GetOpenTime() == now.Add(time.Second*-5).Unix(), ShouldBeTrue)
	})

	Convey("Test time-interval shifting: check that grouping is calculated properly with alarm updates", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		//every new alarm has bigger timestamp as previous so the map is sorted by default
		alarmGroup := NewAlarmGroup("test")
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_2",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 5).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_3",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 10).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_4",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 15).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_5",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 20).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_6",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 25).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 6)
		So(alarmGroup.GetOpenTime() == now.Unix(), ShouldBeTrue)

		alarmGroup.Push(types.Alarm{
			ID: "test_alarm",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 40).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_2",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 55).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_3",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 45).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_4",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 30).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_5",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 35).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_6",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 40).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 6)
		So(alarmGroup.GetOpenTime() == now.Add(time.Second*30).Unix(), ShouldBeTrue)

		//This call should shift time interval, but no alarms should be deleted, since they belong to the new time interval.
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_7",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 65).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 7)
		//The new start time should be equal the minimum alarm's time.
		So(alarmGroup.GetOpenTime() == now.Add(time.Second*30).Unix(), ShouldBeTrue)
	})

	Convey("Test remove before: should be able to remove outdated group elements", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		//every new alarm has bigger timestamp as previous so the map is sorted by default
		alarmGroup := NewAlarmGroup("test")
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_2",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 5).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_3",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 10).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_4",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 15).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_5",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 20).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))
		alarmGroup.Push(types.Alarm{
			ID: "test_alarm_6",
			Value: types.AlarmValue{
				LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 25).Unix()),
			},
		}, int64(minuteRule.Config.TimeInterval))

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 6)
		So(alarmGroup.GetOpenTime() == now.Unix(), ShouldBeTrue)

		alarmGroup.RemoveBefore(now.Add(time.Second * 20).Unix())

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 2)
		So(alarmGroup.GetOpenTime() == now.Add(time.Second*20).Unix(), ShouldBeTrue)

		alarmGroup = NewAlarmGroup("test")
		alarmGroup.RemoveBefore(1)

		So(len(alarmGroup.GetAlarmIds()), ShouldEqual, 0)
		So(alarmGroup.GetOpenTime() == math.MaxInt64, ShouldBeTrue)
	})
}
