package watcher_test

import (
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProcessState(t *testing.T) {
	Convey("Given a CountersCache", t, func() {
		redisClient, err := redis.NewSession(
			redis.CacheWatcher,
			log.NewTestLogger(),
			0,
			0,
		)
		So(err, ShouldBeNil)

		// FIXME: This removes all data from the redis cache, and is necessary
		// to make sure that running the tests twice leads to the same result.
		// This prevents from running multiple tests in parallel. A better
		// option might be to used https://github.com/alicebob/miniredis (but
		// this will only fix this problem for redis)
		err = redisClient.FlushDB().Err()
		So(err, ShouldBeNil)

		countersCache := watcher.NewCountersCache(redisClient, log.NewTestLogger())

		Convey("With two watchers and four entities", func() {
			counters, err := countersCache.ProcessState(watcher.DependencyState{
				EntityID:         "cpu/server1",
				ImpactedWatchers: []string{"watcher1"},
				HasAlarm:         false,
				LastUpdateDate:   time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
				IsEntityActive:   true,
			})
			So(err, ShouldBeNil)
			So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
				"watcher1": watcher.AlarmCounters{
					State: watcher.StateCounters{
						Info: 1,
					},
					PbehaviorCounters: make(map[string]int64),
				},
			})

			counters, err = countersCache.ProcessState(watcher.DependencyState{
				EntityID:          "cpu/server2",
				ImpactedWatchers:  []string{"watcher2"},
				HasAlarm:          true,
				AlarmState:        3,
				AlarmAcknowledged: true,
				LastUpdateDate:    time.Date(2019, 1, 1, 0, 0, 0, 1, time.UTC),
				IsEntityActive:    true,
			})
			So(err, ShouldBeNil)
			So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
				"watcher2": watcher.AlarmCounters{
					All:    1,
					Alarms: 1,
					State: watcher.StateCounters{
						Critical: 1,
					},
					Acknowledged:      1,
					PbehaviorCounters: make(map[string]int64),
				},
			})

			counters, err = countersCache.ProcessState(watcher.DependencyState{
				EntityID:          "mem/server2",
				ImpactedWatchers:  []string{"watcher2"},
				HasAlarm:          true,
				AlarmState:        2,
				AlarmAcknowledged: false,
				IsEntityActive:    false,
				LastUpdateDate:    time.Date(2019, 1, 1, 0, 0, 0, 1, time.UTC),
			})
			So(err, ShouldBeNil)
			So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
				"watcher2": watcher.AlarmCounters{
					All:    2,
					Alarms: 1,
					State: watcher.StateCounters{
						Critical: 1,
						Info:     1,
					},
					Acknowledged:      1,
					PbehaviorCounters: make(map[string]int64),
				},
			})

			counters, err = countersCache.ProcessState(watcher.DependencyState{
				EntityID:          "cpu/server3",
				ImpactedWatchers:  []string{},
				HasAlarm:          true,
				AlarmState:        3,
				AlarmAcknowledged: true,
				LastUpdateDate:    time.Date(2019, 1, 1, 0, 0, 0, 2, time.UTC),
				IsEntityActive:    true,
			})
			So(err, ShouldBeNil)
			So(counters, ShouldBeEmpty)

			Convey("Changing the state of an entity without impacts does not update the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "cpu/server3",
					ImpactedWatchers:  []string{},
					HasAlarm:          true,
					AlarmState:        2,
					AlarmAcknowledged: false,
					LastUpdateDate:    time.Date(2019, 1, 1, 1, 0, 0, 0, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldBeEmpty)
			})

			Convey("Keeping the state of an entity identical does not update the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:         "cpu/server1",
					ImpactedWatchers: []string{"watcher1"},
					HasAlarm:         false,
					LastUpdateDate:   time.Date(2019, 1, 1, 1, 0, 0, 0, time.UTC),
					IsEntityActive:   true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldBeEmpty)
			})

			Convey("Changing the state of an entity updates the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "cpu/server2",
					ImpactedWatchers:  []string{"watcher2"},
					HasAlarm:          true,
					AlarmState:        2,
					AlarmAcknowledged: false,
					LastUpdateDate:    time.Date(2019, 1, 1, 1, 0, 0, 0, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher2": watcher.AlarmCounters{
						All:    2,
						Alarms: 1,
						State: watcher.StateCounters{
							Major: 1,
							Info:  1,
						},
						NotAcknowledged:   1,
						PbehaviorCounters: make(map[string]int64),
					},
				})
			})

			Convey("Processing an outdated state does not update the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "cpu/server2",
					ImpactedWatchers:  []string{"watcher2"},
					HasAlarm:          true,
					AlarmState:        2,
					AlarmAcknowledged: false,
					LastUpdateDate:    time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldBeEmpty)
			})

			Convey("Resolving the alarm of an entity updates the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:         "cpu/server2",
					ImpactedWatchers: []string{"watcher2"},
					HasAlarm:         false,
					LastUpdateDate:   time.Date(2019, 1, 1, 1, 0, 0, 0, time.UTC),
					IsEntityActive:   true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher2": watcher.AlarmCounters{
						All: 1,
						State: watcher.StateCounters{
							Info: 2,
						},
						PbehaviorCounters: make(map[string]int64),
					},
				})
			})

			Convey("Adding a new entity updates the watchers' counters", func() {
				counters, err := countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "mem/server1",
					ImpactedWatchers:  []string{"watcher1"},
					HasAlarm:          true,
					AlarmState:        1,
					AlarmAcknowledged: false,
					LastUpdateDate:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher1": watcher.AlarmCounters{
						All:    1,
						Alarms: 1,
						State: watcher.StateCounters{
							Minor: 1,
							Info:  1,
						},
						NotAcknowledged:   1,
						PbehaviorCounters: make(map[string]int64),
					},
				})
			})

			Convey("Adding an existing entity to a watcher updates the watchers' counters", func() {
				counters, err := countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "cpu/server1",
					ImpactedWatchers:  []string{"watcher1"},
					HasAlarm:          true,
					AlarmState:        2,
					AlarmAcknowledged: true,
					LastUpdateDate:    time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher1": watcher.AlarmCounters{
						All:    1,
						Alarms: 1,
						State: watcher.StateCounters{
							Major: 1,
						},
						Acknowledged:      1,
						PbehaviorCounters: make(map[string]int64),
					},
				})

				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "cpu/server1",
					ImpactedWatchers:  []string{"watcher1", "watcher2"},
					HasAlarm:          true,
					AlarmState:        3,
					AlarmAcknowledged: false,
					LastUpdateDate:    time.Date(2019, 1, 1, 0, 0, 0, 1, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher1": watcher.AlarmCounters{
						All:    1,
						Alarms: 1,
						State: watcher.StateCounters{
							Critical: 1,
						},
						NotAcknowledged:   1,
						PbehaviorCounters: make(map[string]int64),
					},
					"watcher2": watcher.AlarmCounters{
						All:    3,
						Alarms: 2,
						State: watcher.StateCounters{
							Critical: 2,
							Info:     1,
						},
						Acknowledged:      1,
						NotAcknowledged:   1,
						PbehaviorCounters: make(map[string]int64),
					},
				})
			})

			Convey("Removing an entity from a watcher updates the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "cpu/server2",
					ImpactedWatchers:  []string{},
					HasAlarm:          true,
					AlarmState:        3,
					AlarmAcknowledged: true,
					LastUpdateDate:    time.Date(2019, 1, 1, 1, 0, 0, 0, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher2": watcher.AlarmCounters{
						All: 1,
						State: watcher.StateCounters{
							Info: 1,
						},
						PbehaviorCounters: make(map[string]int64),
					},
				})
			})

			Convey("Adding an active pbehavior updates the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "cpu/server2",
					ImpactedWatchers:  []string{"watcher2"},
					HasAlarm:          true,
					AlarmState:        3,
					AlarmAcknowledged: true,
					IsEntityActive:    false,
					LastUpdateDate:    time.Date(2019, 1, 1, 1, 0, 0, 0, time.UTC),
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher2": watcher.AlarmCounters{
						All: 2,
						State: watcher.StateCounters{
							Info: 2,
						},
						PbehaviorCounters: make(map[string]int64),
					},
				})
			})

			Convey("Removing an active pbehavior updates the watchers' counters", func() {
				counters, err = countersCache.ProcessState(watcher.DependencyState{
					EntityID:          "mem/server2",
					ImpactedWatchers:  []string{"watcher2"},
					HasAlarm:          true,
					AlarmState:        2,
					AlarmAcknowledged: false,
					LastUpdateDate:    time.Date(2019, 1, 1, 1, 0, 0, 0, time.UTC),
					IsEntityActive:    true,
				})
				So(err, ShouldBeNil)
				So(counters, ShouldResemble, map[string]watcher.AlarmCounters{
					"watcher2": watcher.AlarmCounters{
						All:    2,
						Alarms: 2,
						State: watcher.StateCounters{
							Critical: 1,
							Major:    1,
						},
						Acknowledged:      1,
						NotAcknowledged:   1,
						PbehaviorCounters: make(map[string]int64),
					},
				})
			})
		})
	})
}
