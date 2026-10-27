package axe_test

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pausejob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/resumejob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/webhook"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/depprovider/depprovidertest"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/pbehavior"
	mock_postgres "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/postgres"
	mock_techmetrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/techmetrics"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/mock/gomock"
)

// Entities of pause fixtures.
const (
	service1  = "test-service-1"
	service2  = "test-service-2"
	service3  = "test-service-3"
	component = "test-component"
	resource1 = "test-resource-1/test-component"
	resource2 = "test-resource-2/test-component"
)

const (
	pause         = types.EventTypePause
	resume        = types.EventTypeResume
	check         = types.EventTypeCheck
	contextUpdate = types.EventTypeContextUpdate
	entityUpdated = types.EventTypeEntityUpdated
	entityToggled = types.EventTypeEntityToggled

	pauseTestTimeout = 10 * time.Second
)

// Fixture: service1 -> service2 -> resource1, resource2.
// Events: 1. service2 pause, 2. service1 pause, 3. service1 resume.
// service2 receives them in order 1, 2, 3 and resource1 in the reverse order 3, 2, 1.
// resource2 receives them in an undefined order, since the fan-out doesn't sort dependencies,
// but it must end up in the same state as resource1.
// All of them must end up not paused: service1 pause replaced service2 pause and its resume restores nothing.
func TestMessageProcessor_Process_GivenEventsReorderedAtServiceAndAtResource_ShouldApplyTheSamePause(t *testing.T) {
	test := newPauseTest(t, "./testdata/fixtures/pause_nested_services.yml")

	// Each event is processed at its origin, then stops right before it reads service2, then resource1.
	s2Pause := test.startBlockedBefore(t, testEvent{pause, service2, 1}, service2, resource1)
	s1Pause := test.startBlockedBefore(t, testEvent{pause, service1, 2}, service2, resource1)
	s1Resume := test.startBlockedBefore(t, testEvent{resume, service1, 3}, service2, resource1)

	// service2 receives the events in the declared order, then each of them stops before resource1.
	s2Pause.unblockAndWait(t)
	s1Pause.unblockAndWait(t)
	s1Resume.unblockAndWait(t)

	// resource1 receives them in the reverse order, resource2 either before resource1 in the order above or after it in this one.
	s1Resume.unblockAndWait(t)
	s1Pause.unblockAndWait(t)
	s2Pause.unblockAndWait(t)

	test.assertPauseState(t, service1, resumed(service1, 2, 3))
	test.assertPauseState(t, service2, resumed(service1, 2, 3))
	test.assertPauseState(t, resource1, resumed(service1, 2, 3))
	test.assertPauseState(t, resource2, resumed(service1, 2, 3))
}

// Fixture: service1, service2 and component -> resource1, resource2.
// The alarm of resource1 is created by a check event while the service1 pause is in flight,
// so the alarm inherits the pause from service1 and the pause event is ignored when it reaches resource1.
func TestMessageProcessor_Process_GivenServicePauseAndCheckCreatingAlarm_ShouldPauseTheNewAlarmOnce(t *testing.T) {
	test := newPauseTest(t, "./testdata/fixtures/pause_independent_services.yml")
	test.removeAlarm(t, resource1)

	s1Pause := test.startBlockedBefore(t, testEvent{pause, service1, 1}, resource1)
	test.process(t, testEvent{check, resource1, 2})
	s1Pause.unblockAndWait(t)

	test.assertPauseState(t, service1, pausedBy(service1, 1))
	test.assertPauseState(t, resource1, pausedBy(service1, 1))
	test.assertPauseState(t, resource2, pausedBy(service1, 1))
}

// Fixture: service1, service2 and component -> resource1, resource2.
// resource1 is added to service1 while its pause is in flight, the check event which applies
// the services change must pause resource1 and its alarm only once.
func TestMessageProcessor_Process_GivenServicePauseAndCheckOfAddedResource_ShouldPauseItOnce(t *testing.T) {
	for _, withAlarm := range []bool{true, false} {
		t.Run(fmt.Sprintf("with alarm %v", withAlarm), func(t *testing.T) {
			test := newPauseTest(t, "./testdata/fixtures/pause_independent_services.yml")
			if !withAlarm {
				test.removeAlarm(t, resource1)
			}

			test.updateEntity(t, resource1, bson.M{"$set": bson.M{"services_to_add": []string{service1}}})

			s1Pause := test.startBlockedBefore(t, testEvent{pause, service1, 1}, resource1)
			test.process(t, testEvent{check, resource1, 2})
			s1Pause.unblockAndWait(t)

			test.assertPauseState(t, resource1, pausedBy(service1, 1))
			test.assertPauseState(t, resource2, pausedBy(service1, 1))
		})
	}
}

// Fixture: service1, service2 and component -> resource1, resource2.
// resource1 is paused by service1 and removed from it while the service2 pause is in flight,
// so the check event which applies the services change must pause resource1 by service2.
func TestMessageProcessor_Process_GivenServicePauseAndCheckOfRemovedResource_ShouldPauseItByTheOtherService(t *testing.T) {
	for _, withAlarm := range []bool{true, false} {
		t.Run(fmt.Sprintf("with alarm %v", withAlarm), func(t *testing.T) {
			test := newPauseTest(t, "./testdata/fixtures/pause_independent_services.yml")
			if !withAlarm {
				test.removeAlarm(t, resource1)
			}

			test.process(t, testEvent{pause, service1, 1})
			test.assertPauseState(t, resource1, pausedBy(service1, 1))

			s2Pause := test.startBlockedBefore(t, testEvent{pause, service2, 2}, resource1)
			test.updateEntity(t, resource1, bson.M{"$set": bson.M{
				"services":           []string{service2},
				"services_to_remove": []string{service1},
			}})
			test.process(t, testEvent{check, resource1, 3})
			s2Pause.unblockAndWait(t)

			test.assertPauseState(t, resource1, pausedBy(service2, 2))
			test.assertPauseState(t, resource2, pausedBy(service2, 2))
		})
	}
}

// Fixture: service1, service2 and component -> resource1, resource2.
// resource1 is added to service1 while its pause is in flight, the event which applies
// the services change must pause resource1 and its alarm, if any, only once.
func TestMessageProcessor_Process_GivenServicePauseAndEntityEventOfAddedResource_ShouldPauseItOnce(t *testing.T) {
	for _, eventType := range []string{contextUpdate, entityUpdated, entityToggled} {
		for _, withAlarm := range []bool{true, false} {
			t.Run(fmt.Sprintf("%s with alarm %v", eventType, withAlarm), func(t *testing.T) {
				test := newPauseTest(t, "./testdata/fixtures/pause_independent_services.yml")
				if !withAlarm {
					test.removeAlarm(t, resource1)
				}

				test.updateEntity(t, resource1, bson.M{"$set": bson.M{"services_to_add": []string{service1}}})

				s1Pause := test.startBlockedBefore(t, testEvent{pause, service1, 1}, resource1)
				test.process(t, testEvent{eventType, resource1, 2})
				s1Pause.unblockAndWait(t)

				test.assertPauseState(t, resource1, pausedBy(service1, 1))
				test.assertPauseState(t, resource2, pausedBy(service1, 1))
			})
		}
	}
}

// Fixture: service1, service2 and component -> resource1, resource2.
func TestMessageProcessor_Process_GivenPauseEventsOfIndependentParents_ShouldApplyTheSamePauseInAnyOrder(t *testing.T) {
	runPauseCases(t, "./testdata/fixtures/pause_independent_services.yml", []pauseCase{
		{
			name:        "the latest pause wins",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{pause, service2, 2},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 1),
				service2:  pausedBy(service2, 2),
				resource1: pausedBy(service2, 2),
				resource2: pausedBy(service2, 2),
			},
		},
		{
			name:        "pauses at the same time are ordered by source",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{pause, service2, 1},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 1),
				service2:  pausedBy(service2, 1),
				resource1: pausedBy(service2, 1),
				resource2: pausedBy(service2, 1),
			},
		},
		{
			name:        "resume of the latest pause doesn't restore the older pause",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{pause, service2, 2},
				{resume, service2, 3},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 1),
				service2:  resumed(service2, 2, 3),
				resource1: resumed(service2, 2, 3),
				resource2: resumed(service2, 2, 3),
			},
		},
		{
			name:        "resume of a replaced pause is ignored",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{pause, service2, 2},
				{resume, service1, 3},
			},
			expected: map[string]pauseState{
				service1:  resumed(service1, 1, 3),
				service2:  pausedBy(service2, 2),
				resource1: pausedBy(service2, 2),
				resource2: pausedBy(service2, 2),
			},
		},
		{
			name:        "pause after resume of the same service",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{resume, service1, 2},
				{pause, service1, 3},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 3),
				resource1: pausedBy(service1, 3),
				resource2: pausedBy(service1, 3),
			},
		},
		{
			name:        "duplicated pause is ignored",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{pause, service1, 1},
				{resume, service1, 2},
			},
			expected: map[string]pauseState{
				service1:  resumed(service1, 1, 2),
				resource1: resumed(service1, 1, 2),
				resource2: resumed(service1, 1, 2),
			},
		},
		{
			name:        "service resume doesn't restore component pause",
			blockBefore: resource1,
			events: []testEvent{
				{pause, component, 1},
				{pause, service1, 2},
				{resume, service1, 3},
			},
			expected: map[string]pauseState{
				component: pausedBy(component, 1),
				service1:  resumed(service1, 2, 3),
				resource1: resumed(service1, 2, 3),
				resource2: resumed(service1, 2, 3),
			},
		},
		{
			name:        "older service pause replaces resource pause",
			blockBefore: resource1,
			events: []testEvent{
				{pause, resource1, 2},
				{pause, service1, 1},
				{resume, service1, 3},
			},
			expected: map[string]pauseState{
				service1:  resumed(service1, 1, 3),
				resource1: resumed(service1, 1, 3),
				resource2: resumed(service1, 1, 3),
			},
		},
		{
			name:        "service resume doesn't restore resource pause",
			blockBefore: resource1,
			events: []testEvent{
				{pause, resource1, 1},
				{pause, service1, 2},
				{resume, service1, 3},
			},
			expected: map[string]pauseState{
				service1:  resumed(service1, 2, 3),
				resource1: resumed(service1, 2, 3),
				resource2: resumed(service1, 2, 3),
			},
		},
	})
}

// Fixture: service1 -> service2 -> resource1, resource2.
func TestMessageProcessor_Process_GivenPauseEventsOfNestedServices_ShouldApplyTheSamePauseInAnyOrder(t *testing.T) {
	runPauseCases(t, "./testdata/fixtures/pause_nested_services.yml", []pauseCase{
		{
			name:        "newer parent pause replaces child pause",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service2, 1},
				{pause, service1, 2},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 2),
				service2:  pausedBy(service1, 2),
				resource1: pausedBy(service1, 2),
				resource2: pausedBy(service1, 2),
			},
		},
		{
			name:        "own pause of a child paused by its parent is ignored",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{pause, service2, 2},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 1),
				service2:  pausedBy(service1, 1),
				resource1: pausedBy(service1, 1),
				resource2: pausedBy(service1, 1),
			},
		},
		{
			name:        "parent resume doesn't restore child pause, reordered at resource",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service2, 1},
				{pause, service1, 2},
				{resume, service1, 3},
			},
			expected: map[string]pauseState{
				service1:  resumed(service1, 2, 3),
				service2:  resumed(service1, 2, 3),
				resource1: resumed(service1, 2, 3),
				resource2: resumed(service1, 2, 3),
			},
		},
		{
			name:        "parent resume doesn't restore child pause, reordered at child service",
			blockBefore: service2,
			events: []testEvent{
				{pause, service2, 1},
				{pause, service1, 2},
				{resume, service1, 3},
			},
			expected: map[string]pauseState{
				service1:  resumed(service1, 2, 3),
				service2:  resumed(service1, 2, 3),
				resource1: resumed(service1, 2, 3),
				resource2: resumed(service1, 2, 3),
			},
		},
		{
			name:        "child resume while paused by parent is ignored, reordered at resource",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service2, 1},
				{pause, service1, 2},
				{resume, service2, 3},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 2),
				service2:  pausedBy(service1, 2),
				resource1: pausedBy(service1, 2),
				resource2: pausedBy(service1, 2),
			},
		},
		{
			name:        "child resume while paused by parent is ignored, reordered at child service",
			blockBefore: service2,
			events: []testEvent{
				{pause, service2, 1},
				{pause, service1, 2},
				{resume, service2, 3},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 2),
				service2:  pausedBy(service1, 2),
				resource1: pausedBy(service1, 2),
				resource2: pausedBy(service1, 2),
			},
		},
	})
}

// Fixture: service1 -> service2, service3 -> resource1, resource2, so resources receive service1 events twice.
func TestMessageProcessor_Process_GivenPauseEventsOfDiamondServices_ShouldApplyTheSamePauseInAnyOrder(t *testing.T) {
	runPauseCases(t, "./testdata/fixtures/pause_diamond_services.yml", []pauseCase{
		{
			name:        "parent pause and resume received through both children",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service1, 1},
				{resume, service1, 2},
			},
			expected: map[string]pauseState{
				service1:  resumed(service1, 1, 2),
				service2:  resumed(service1, 1, 2),
				service3:  resumed(service1, 1, 2),
				resource1: resumed(service1, 1, 2),
				resource2: resumed(service1, 1, 2),
			},
		},
		{
			name:        "the parent pause replaces the child pause it has passed through",
			blockBefore: resource1,
			events: []testEvent{
				{pause, service2, 2},
				{pause, service1, 1},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 1),
				service2:  pausedBy(service1, 1),
				service3:  pausedBy(service1, 1),
				resource1: pausedBy(service1, 1),
				resource2: pausedBy(service1, 1),
			},
		},
	})
}

// Events are processed at the same time, so concurrent updates of the same entities go through
// transaction conflicts and retries. Events have different origins to not depend on the processing order.
func TestMessageProcessor_Process_GivenConcurrentPauseEvents_ShouldApplyTheLatestPause(t *testing.T) {
	testCases := []struct {
		name        string
		fixturePath string
		events      []testEvent
		expected    map[string]pauseState
	}{
		{
			name:        "independent services",
			fixturePath: "./testdata/fixtures/pause_independent_services.yml",
			events: []testEvent{
				{pause, service1, 1},
				{pause, service2, 2},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 1),
				service2:  pausedBy(service2, 2),
				resource1: pausedBy(service2, 2),
				resource2: pausedBy(service2, 2),
			},
		},
		{
			name:        "nested services",
			fixturePath: "./testdata/fixtures/pause_nested_services.yml",
			events: []testEvent{
				{pause, service2, 1},
				{pause, service1, 2},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 2),
				service2:  pausedBy(service1, 2),
				resource1: pausedBy(service1, 2),
				resource2: pausedBy(service1, 2),
			},
		},
		{
			name:        "diamond services",
			fixturePath: "./testdata/fixtures/pause_diamond_services.yml",
			events: []testEvent{
				{pause, service2, 1},
				{pause, service3, 2},
				{pause, service1, 3},
			},
			expected: map[string]pauseState{
				service1:  pausedBy(service1, 3),
				service2:  pausedBy(service1, 3),
				service3:  pausedBy(service1, 3),
				resource1: pausedBy(service1, 3),
				resource2: pausedBy(service1, 3),
			},
		},
	}

	for _, tc := range testCases {
		for i := range 10 {
			t.Run(fmt.Sprintf("%s %d", tc.name, i+1), func(t *testing.T) {
				test := newPauseTest(t, tc.fixturePath)
				var wg sync.WaitGroup
				for _, e := range tc.events {
					body := test.encode(t, e)
					processor := test.newProcessor(nil)
					wg.Go(func() {
						_, err := processor.Process(t.Context(), amqp.Delivery{Body: body})
						if err != nil {
							t.Errorf("unexpected error %v", err)
						}
					})
				}

				wg.Wait()

				for id, expected := range tc.expected {
					test.assertPauseState(t, id, expected)
				}
			})
		}
	}
}

// pauseCase is run once per arrival order of its events at blockBefore entity.
// Events are processed at their origin in the declared order, only their fan-outs
// reach blockBefore entity in a different order.
// Other entities receive the events in an undefined order, since the fan-out doesn't sort
// dependencies, they must end up in the expected state anyway.
type pauseCase struct {
	name        string
	blockBefore string
	events      []testEvent
	expected    map[string]pauseState
}

type testEvent struct {
	eventType string
	entityID  string
	at        int // seconds from the test start
}

// pauseState is the entity's active pause and latest_pause, times are seconds from the test start.
type pauseState struct {
	PausedBy    string // empty if not paused
	PausedAt    int
	Latest      string // empty if the entity has never received any pause or resume
	LatestAt    int
	LatestEndAt int // 0 if the latest pause is active
}

func pausedBy(source string, at int) pauseState {
	return pauseState{PausedBy: source, PausedAt: at, Latest: source, LatestAt: at}
}

func resumed(source string, at, endAt int) pauseState {
	return pauseState{Latest: source, LatestAt: at, LatestEndAt: endAt}
}

func runPauseCases(t *testing.T, fixturePath string, cases []pauseCase) {
	for _, c := range cases {
		// order contains events indexes, e.g. "arrival 2,1" means that the 2nd event
		// reaches blockBefore entity before the 1st one.
		p := permutations(len(c.events))
		for _, order := range p {
			numbers := make([]string, len(order))
			for j, i := range order {
				numbers[j] = strconv.Itoa(i + 1)
			}

			t.Run(c.name+"/arrival "+strings.Join(numbers, ","), func(t *testing.T) {
				test := newPauseTest(t, fixturePath)

				// Each event is processed at its origin, then stops right before it reads blockBefore entity.
				events := make([]*blockedEvent, len(c.events))
				for i, e := range c.events {
					events[i] = test.startBlockedBefore(t, e, c.blockBefore)
				}

				// Events reach blockBefore entity in the given order, other entities of the fan-out
				// receive them either before or after it.
				for _, i := range order {
					events[i].unblockAndWait(t)
				}

				for id, expected := range c.expected {
					test.assertPauseState(t, id, expected)
				}
			})
		}
	}
}

// permutations returns all orders of 0..n-1.
func permutations(n int) [][]int {
	if n == 0 {
		return [][]int{{}}
	}

	var res [][]int
	for _, p := range permutations(n - 1) {
		for i := 0; i <= len(p); i++ {
			res = append(res, slices.Insert(slices.Clone(p), i, n-1))
		}
	}

	return res
}

type pauseTest struct {
	entities mongo.DbCollection
	alarms   mongo.DbCollection
	// newProcessor returns a processor which calls beforeEntityRead, if set, before each entity read.
	newProcessor func(beforeEntityRead func(filter any)) *pauseTestProcessor
	startTime    time.Time
}

func newPauseTest(t *testing.T, fixturePath string) *pauseTest {
	t.Helper()

	ctrl := gomock.NewController(t)
	ctx := t.Context()
	dbClient, err := mongo.NewClient(ctx)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	t.Cleanup(func() {
		err := dbClient.Disconnect(context.WithoutCancel(ctx))
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
	})

	jobFilter := bson.M{"parent": bson.M{"$in": []string{service1, service2, service3, component, resource1, resource2}}}
	for _, collection := range []string{mongo.ResumeJobCollection, mongo.PauseJobCollection} {
		_, err = dbClient.Collection(collection).DeleteMany(ctx, jobFilter)
		if err != nil {
			t.Fatalf("cannot clear pause test jobs: %v", err)
		}
		t.Cleanup(func() {
			_, err := dbClient.Collection(collection).DeleteMany(context.WithoutCancel(ctx), jobFilter)
			if err != nil {
				t.Errorf("cannot clear pause test jobs: %v", err)
			}
		})
	}

	loader := fixtures.NewLoader(dbClient, []string{fixturePath},
		fixtures.NewParser(fixtures.NewFaker(password.NewBcryptEncoder())), zerolog.Nop())
	err = loader.Load(ctx)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	t.Cleanup(func() {
		err := loader.Clean(context.WithoutCancel(ctx))
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}
	})

	amqpChannel := mock_amqp.NewMockChannel(ctrl)
	amqpChannel.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	techMetricsSender := mock_techmetrics.NewMockSender(ctrl)
	techMetricsSender.EXPECT().SendAxeEvent(gomock.Any()).AnyTimes()
	pbhStore := mock_pbehavior.NewMockStore(ctrl)
	pbhStore.EXPECT().GetSpan(gomock.Any()).Return(timespan.New(time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour)), nil).AnyTimes()
	pbhStore.EXPECT().GetComputed(gomock.Any()).Return(pbehavior.ComputeResult{}, nil).AnyTimes()
	pbhStore.EXPECT().GetInheritedServicesPbhResolveResult(gomock.Any()).Return(pbehavior.InheritedServicesPbhResolveResult{}, nil).AnyTimes()

	logger := zerolog.Nop()
	cfg := config.CanopsisConf{}
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	userInterfaceConfigProvider := config.NewUserInterfaceConfigProvider(config.UserInterfaceConf{}, logger)
	templateConfigProvider := config.NewTemplateConfigProvider(cfg, logger)
	tzConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, logger)
	eventsSender := entitycounters.NewEventSender(json.NewEncoder(), amqpChannel, canopsis.DefaultExchangeName, canopsis.FIFOQueueName, "", alarmConfigProvider)
	alarmStatusService := alarmstatus.NewService(dbClient, flappingrule.NewAdapter(dbClient), alarmConfigProvider, logger)
	metricsSender := metrics.NewTimescaleDBSender(mock_postgres.NewMockPoolProvider(ctrl), metricsConfigProvider, logger)
	metaAlarmStatesService := correlation.NewMetaAlarmStateService(dbClient)
	metaAlarmPostProcessor := event.NewMetaAlarmPostProcessor(dbClient, alarm.NewAdapter(dbClient), correlation.NewRuleAdapter(dbClient),
		alarmStatusService, metaAlarmStatesService, json.NewEncoder(), libevent.NewGenerator(canopsis.AxeConnector, canopsis.AxeConnector), nil, metricsSender, logger)

	return &pauseTest{
		entities:  dbClient.Collection(mongo.EntityMongoCollection),
		alarms:    dbClient.Collection(mongo.AlarmMongoCollection),
		startTime: time.Now().Truncate(time.Second),
		newProcessor: func(beforeEntityRead func(filter any)) *pauseTestProcessor {
			eventDbClient := dbClient
			if beforeEntityRead != nil {
				eventDbClient = depprovidertest.NewMongoClient(dbClient, depprovidertest.MongoClientOptions{
					WrapCollection: func(coll mongo.DbCollection) mongo.DbCollection {
						if coll.Name() != mongo.EntityMongoCollection {
							return coll
						}

						return depprovidertest.NewMongoCollection(coll, depprovidertest.MongoCollOptions{
							BeforeFindOne: func(_ context.Context, filter any) {
								beforeEntityRead(filter)
							},
						})
					},
				})
			}

			return &pauseTestProcessor{client: dbClient, MessageProcessor: &axe.MessageProcessor{
				EventProcessor: axe.DependencyMaker{}.EventProcessor(
					eventDbClient,
					alarmConfigProvider,
					userInterfaceConfigProvider,
					alarmStatusService,
					pbehavior.NewEntityTypeResolver(pbhStore, logger),
					event.NewNullAutoInstructionMatcher(),
					calculator.NewEntityServiceCountersCalculator(dbClient, template.NewExecutor(templateConfigProvider, tzConfigProvider), eventsSender, logger),
					calculator.NewComponentCountersCalculator(dbClient, eventsSender),
					eventsSender,
					metaAlarmPostProcessor,
					metaAlarmStatesService,
					metricsSender,
					statistics.NewEventStatisticsSender(dbClient, logger, tzConfigProvider),
					nil,
					alarmtag.NewExternalUpdater(dbClient),
					alarmtag.NewInternalTagAlarmMatcher(dbClient),
					amqpChannel,
					libevent.NewGenerator(canopsis.AxeConnector, canopsis.AxeConnector),
					template.NewExecutor(templateConfigProvider, tzConfigProvider),
					webhook.NewNullJobStatusService(),
					logger,
				),
				Encoder:           json.NewEncoder(),
				Decoder:           json.NewDecoder(),
				TechMetricsSender: techMetricsSender,
				AlarmCollection:   dbClient.Collection(mongo.AlarmMongoCollection),
				Logger:            logger,
			}}
		},
	}
}

// The production engine starts pause/resume workers separately. These processor-only
// tests drain its jobs synchronously through the same processor so their existing
// read barriers continue to control delivery to each dependency.
type pauseTestProcessor struct {
	*axe.MessageProcessor
	client mongo.DbClient
}

func (p *pauseTestProcessor) Process(ctx context.Context, delivery amqp.Delivery) ([]byte, error) {
	body, err := p.MessageProcessor.Process(ctx, delivery)
	if err != nil {
		return body, err
	}

	var origin types.Event
	if err := json.NewDecoder().Decode(delivery.Body, &origin); err != nil {
		return nil, err
	}

	source := origin.PauseSource
	if source == "" && origin.Entity != nil {
		source = origin.Entity.ID
		if origin.EventType != types.EventTypePause && origin.Entity.LatestPause != nil {
			source = origin.Entity.LatestPause.Source
		}
	}

	pauseSource := source
	pauseTimestamp := origin.Timestamp
	if origin.EventType != types.EventTypePause {
		var current types.Entity
		err := p.client.Collection(mongo.EntityMongoCollection).FindOne(ctx, bson.M{"_id": origin.GetEID()}).Decode(&current)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, err
		}
		pauseSource = ""
		if current.LatestPause != nil && current.Pause != nil {
			pauseSource = current.LatestPause.Source
			pauseTimestamp = current.LatestPause.Timestamp
		}
	}

	if pauseSource != "" {
		cursor, err := p.client.Collection(mongo.PauseJobCollection).Find(ctx, bson.M{
			"parent": origin.GetEID(), "event.timestamp": pauseTimestamp, "event.pause_source": pauseSource,
		})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			var job pausejob.Job
			if err := cursor.Decode(&job); err != nil {
				return nil, err
			}
			if err := axe.ProcessPauseJobForTest(ctx, p.client, p, job.ID); err != nil {
				return nil, err
			}
		}
		if err := cursor.Err(); err != nil {
			return nil, err
		}
	}

	cursor, err := p.client.Collection(mongo.ResumeJobCollection).Find(ctx, bson.M{
		"parent": origin.GetEID(), "event.timestamp": origin.Timestamp, "event.pause_source": source,
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var job resumejob.Job
		if err := cursor.Decode(&job); err != nil {
			return nil, err
		}

		if err := axe.ProcessResumeJobForTest(ctx, p.client, p, job.ID); err != nil {
			return nil, err
		}
	}

	return body, cursor.Err()
}

func (p *pauseTestProcessor) PublishWithContext(ctx context.Context, _, _ string, _, _ bool, msg amqp.Publishing) error {
	var event types.Event
	if err := json.NewDecoder().Decode(msg.Body, &event); err != nil {
		return err
	}

	var entity types.Entity
	if err := p.client.Collection(mongo.EntityMongoCollection).FindOne(ctx, bson.M{"_id": event.GetEID()}).Decode(&entity); err != nil {
		return err
	}

	event.Entity = &entity
	body, err := json.NewEncoder().Encode(event)
	if err != nil {
		return err
	}

	_, err = p.Process(ctx, amqp.Delivery{Body: body})
	return err
}

// process processes e completely.
func (test *pauseTest) process(t *testing.T, e testEvent) {
	t.Helper()

	_, err := test.newProcessor(nil).Process(t.Context(), amqp.Delivery{Body: test.encode(t, e)})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

// startBlockedBefore starts processing of e in background and waits until it stops right before it reads
// the first of entityIDs, or until it's processed if it's ignored before reading any of them.
func (test *pauseTest) startBlockedBefore(t *testing.T, e testEvent, entityIDs ...string) *blockedEvent {
	t.Helper()

	b := &blockedEvent{
		stopped: make(chan struct{}),
		unblock: make(chan struct{}),
		done:    make(chan error, 1),
	}
	// Don't leave the processing blocked if the test fails.
	t.Cleanup(func() { close(b.unblock) })

	stopBefore := slices.Clone(entityIDs)
	processor := test.newProcessor(func(filter any) {
		f, ok := filter.(bson.M)
		if !ok || len(f) != 1 {
			return
		}

		id, _ := f["_id"].(string)
		idx := slices.Index(stopBefore, id)
		if idx < 0 {
			return
		}

		// Stop only once before each entity.
		stopBefore = slices.Delete(stopBefore, idx, idx+1)
		select {
		case b.stopped <- struct{}{}:
			<-b.unblock
		case <-b.unblock: // closed on test cleanup
		}
	})
	body := test.encode(t, e)
	go func() {
		_, err := processor.Process(t.Context(), amqp.Delivery{Body: body})
		b.done <- err
	}()

	b.wait(t)

	return b
}

func (test *pauseTest) encode(t *testing.T, e testEvent) []byte {
	t.Helper()

	var entity types.Entity
	err := test.entities.FindOne(t.Context(), bson.M{"_id": e.entityID}).Decode(&entity)
	if err != nil {
		t.Fatalf("cannot find entity %q: %v", e.entityID, err)
	}

	evt := types.Event{
		EventType:     e.eventType,
		Connector:     "api",
		ConnectorName: "api",
		Component:     entity.ID,
		SourceType:    entity.Type, // source types have the same values as entity types
		Output:        e.eventType + " " + e.entityID,
		Author:        "test-author",
		UserID:        "test-user",
		Initiator:     types.InitiatorUser,
		Timestamp:     datetime.CpsTime{Time: test.startTime.Add(time.Duration(e.at) * time.Second)},
		Inherited:     true,
		Entity:        &entity,
	}
	if entity.Type == types.EntityTypeResource {
		evt.Component = entity.Component
		evt.Resource = entity.Name
	}

	switch e.eventType {
	case types.EventTypeCheck:
		evt.State = types.AlarmStateCritical
	case types.EventTypeEntityToggled:
		evt.Initiator = types.InitiatorSystem
	}

	body, err := json.NewEncoder().Encode(evt)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	return body
}

// assertPauseState checks the entity's pause and its alarm, if any: the alarm is paused as the entity
// and pause and resume steps alternate, e.g. a pause source change doesn't add a pause step.
func (test *pauseTest) assertPauseState(t *testing.T, entityID string, expected pauseState) {
	t.Helper()

	var entity types.Entity
	err := test.entities.FindOne(t.Context(), bson.M{"_id": entityID}).Decode(&entity)
	if err != nil {
		t.Fatalf("cannot find entity %q: %v", entityID, err)
	}

	secondsFromStart := func(t datetime.CpsTime) int {
		return int(t.Unix() - test.startTime.Unix())
	}
	var state pauseState
	if p := entity.Pause; p != nil {
		// The entity's own pause has no source.
		state.PausedBy = cmp.Or(p.Source, entityID)
		state.PausedAt = secondsFromStart(p.Timestamp)
	}
	if p := entity.LatestPause; p != nil {
		state.Latest = p.Source
		state.LatestAt = secondsFromStart(p.Timestamp)
		if p.End != nil {
			state.LatestEndAt = secondsFromStart(*p.End)
		}
	}

	if state != expected {
		t.Errorf("%s: expected %+v, got %+v", entityID, expected, state)
	}

	var a types.Alarm
	err = test.alarms.FindOne(t.Context(), bson.M{"d": entityID}).Decode(&a)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return
	}

	if err != nil {
		t.Fatalf("cannot find alarm of %q: %v", entityID, err)
	}

	// A pause step can follow another one, e.g. when the pause source changes, but a resume step
	// always follows a pause step.
	pausedByLastStep := false
	for _, step := range a.Value.Steps {
		switch step.Type {
		case types.AlarmStepPause:
			pausedByLastStep = true
		case types.AlarmStepResume:
			if !pausedByLastStep {
				t.Errorf("%s: alarm has a resume step which follows no pause step", entityID)
			}

			pausedByLastStep = false
		}
	}

	expectedPaused := expected.PausedBy != ""
	if pausedByLastStep != expectedPaused {
		t.Errorf("%s: expected the last alarm pause or resume step to leave the alarm paused=%v, got paused=%v",
			entityID, expectedPaused, pausedByLastStep)
	}

	if isPaused := a.Value.Pause != nil; isPaused != expectedPaused {
		t.Errorf("%s: expected alarm v.pause to be set=%v, got set=%v", entityID, expectedPaused, isPaused)
	}
}

// blockedEvent is an event processing which stops right before it reads each of the given entities.
type blockedEvent struct {
	stopped  chan struct{}
	unblock  chan struct{}
	done     chan error
	finished bool
}

// unblockAndWait lets the processing continue and waits until it stops right before the next entity,
// or until it's processed.
func (b *blockedEvent) unblockAndWait(t *testing.T) {
	t.Helper()

	if b.finished {
		return
	}

	b.unblock <- struct{}{}
	b.wait(t)
}

func (b *blockedEvent) wait(t *testing.T) {
	t.Helper()

	select {
	case <-b.stopped:
	case err := <-b.done:
		b.finished = true
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	case <-time.After(pauseTestTimeout):
		t.Fatal("timeout waiting for event processing")
	}
}

// updateEntity updates the entity in the database, e.g. to change its services.
func (test *pauseTest) updateEntity(t *testing.T, entityID string, update bson.M) {
	t.Helper()

	_, err := test.entities.UpdateOne(t.Context(), bson.M{"_id": entityID}, update)
	if err != nil {
		t.Fatalf("cannot update entity %q: %v", entityID, err)
	}
}

// removeAlarm removes the entity alarm, so the next check event creates a new one.
func (test *pauseTest) removeAlarm(t *testing.T, entityID string) {
	t.Helper()

	_, err := test.alarms.DeleteMany(t.Context(), bson.M{"d": entityID})
	if err != nil {
		t.Fatalf("cannot remove alarm of %q: %v", entityID, err)
	}
}
