package storage

import (
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	redisV7 "github.com/go-redis/redis/v7"
	"testing"
	"time"

	//"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStorageSetGet(t *testing.T) {
	storage := NewRedisGroupingStorage()

	redisClient, err := redis.NewSession(redis.AlarmGroupStorage, log.NewLogger(true), 0, 0)
	if err != nil {
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

		_ = redisClient.Watch(func(tx *redisV7.Tx) error {
			length, err := storage.GetGroupLen(tx, "test_rule39")
			So(err, ShouldBeNil)

			So(length, ShouldEqual, 0)

			err = storage.Push(tx, testRule, testAlarm, "")
			So(err, ShouldBeNil)

			length, err = storage.GetGroupLen(tx, "test_rule")
			So(err, ShouldBeNil)

			So(length, ShouldEqual, 1)

			err = storage.Push(tx, testRule2, testAlarm, "")
			So(err, ShouldBeNil)
			err = storage.Push(tx, testRule2, testAlarm2, "")
			So(err, ShouldBeNil)

			length, err = storage.GetGroupLen(tx, "test_rule_2")
			So(err, ShouldBeNil)

			So(length, ShouldEqual, 2)

			err = storage.Push(tx, testRule2, testAlarm2, "")
			So(err, ShouldBeNil)

			length, err = storage.GetGroupLen(tx, "test_rule_2")
			So(err, ShouldBeNil)

			So(length, ShouldEqual, 2)

			err = storage.CleanPush(tx, testRule2, testAlarm, "")
			So(err, ShouldBeNil)

			length, err = storage.GetGroupLen(tx, "test_rule_2")
			So(err, ShouldBeNil)

			So(length, ShouldEqual, 1)

			return nil
		}, "test_key")
	})
}

func TestStorageShiftTimeInterval(t *testing.T) {
	storage := NewRedisGroupingStorage()

	redisClient, err := redis.NewSession(redis.AlarmGroupStorage, log.NewLogger(true), 0, 0)
	if err != nil {
		panic(err)
	}

	Convey("Test time-interval shifting: basic grouping logic", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		_ = redisClient.Watch(func(tx *redisV7.Tx) error {
			//every new alarm has bigger timestamp as previous so the map is sorted by default
			err = storage.CleanPush(tx, minuteRule, types.Alarm{
				ID: "test_alarm",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_2",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 10).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_3",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 20).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_4",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 30).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_5",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 40).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_6",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 50).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err := storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 6)
			So(group.OpenTime.Equal(now.Time), ShouldBeTrue)

			//This call should shift time interval => so the storage should delete the first alarm in the Group
			// and update 'create' time, since the first alarm won't be in the minute time window anymore.
			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_7",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 65).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err = storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 6)
			So(group.OpenTime.Equal(now.Add(time.Second*10)), ShouldBeTrue)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_8",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 300).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err = storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 1)
			So(group.OpenTime.Equal(now.Add(time.Second*300)), ShouldBeTrue)

			return nil
		}, "test_key")
	})

	Convey("Test time-interval shifting: check that open time is changed if the next event is older than previous", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		_ = redisClient.Watch(func(tx *redisV7.Tx) error {
			err = storage.CleanPush(tx, minuteRule, types.Alarm{
				ID: "test_alarm",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_2",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * -10).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err := storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 2)
			So(group.OpenTime.Equal(now.Add(time.Second*-10)), ShouldBeTrue)

			return nil
		}, "test_key")
	})

	Convey("Test time-interval shifting: check that time-interval is shifted properly, if the new alarm is late and no alarm should be missed", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		_ = redisClient.Watch(func(tx *redisV7.Tx) error {
			err = storage.CleanPush(tx, minuteRule, types.Alarm{
				ID: "test_alarm",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_2",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 10).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_3",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 20).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_4",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 40).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err := storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 4)
			So(group.OpenTime.Equal(now.Time), ShouldBeTrue)

			//test_alarm_4 will be missed, so there shouldn't be any interval shifting
			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_5",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * -30).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err = storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 4)
			So(group.OpenTime.Equal(now.Time), ShouldBeTrue)

			//Interval can be shifted, since none alarm will be lost
			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_5",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * -5).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err = storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 5)
			So(group.OpenTime.Equal(now.Add(time.Second*-5)), ShouldBeTrue)

			return nil
		}, "test_key")
	})

	Convey("Test time-interval shifting: check that grouping is calculated properly with alarm updates", t, func() {
		minuteRule := metaalarm.Rule{
			ID: "minute_rule",
			Config: metaalarm.RuleConfig{
				TimeInterval: 60,
			},
		}

		now := types.NewCpsTime(time.Now().Unix())

		_ = redisClient.Watch(func(tx *redisV7.Tx) error {
			//every new alarm has bigger timestamp as previous so the map is sorted by default
			err = storage.CleanPush(tx, minuteRule, types.Alarm{
				ID: "test_alarm",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_2",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 5).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_3",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 10).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_4",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 15).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_5",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 20).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_6",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 25).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err := storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 6)
			So(group.OpenTime.Equal(now.Time), ShouldBeTrue)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 40).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_2",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 55).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_3",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 45).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_4",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 30).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_5",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 35).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_6",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 40).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err = storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 6)
			So(group.OpenTime.Equal(now.Add(time.Second*30)), ShouldBeTrue)

			//This call should shift time interval, but no alarms should be deleted, since they belong to the new time interval.
			err = storage.Push(tx, minuteRule, types.Alarm{
				ID: "test_alarm_7",
				Value: types.AlarmValue{
					LastUpdateDate: types.NewCpsTime(now.Add(time.Second * 65).Unix()),
				},
			}, "")
			So(err, ShouldBeNil)

			group, err = storage.Get(tx, minuteRule.ID)
			So(err, ShouldBeNil)

			So(len(group.Group), ShouldEqual, 7)
			//The new start time should be equal the minimum alarm's time.
			So(group.OpenTime.Equal(now.Add(time.Second*30)), ShouldBeTrue)

			return nil
		}, "test_key")
	})
}
