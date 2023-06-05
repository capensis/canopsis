package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
	"golang.org/x/sync/errgroup"
)

const (
	EventTypePbehavior = "pbehavior"

	ActionCreate = "create"
	ActionDelete = "delete"

	ApiEndpoint = "/api/v4/pbehaviors"

	Queue      = "legacy_pbehavior"
	RoutingKey = "*.*.pbehavior.*.#"

	httpRetryCount    = 3
	httpRetryInterval = 500 * time.Millisecond

	workersCount = 20
)

type eventConsumer struct {
	amqpCh     libamqp.Channel
	config     Config
	httpClient *http.Client
	logger     zerolog.Logger

	locksMx sync.Mutex
	locks   map[string]*pbhLock
}

type pbhLock struct {
	Mx     sync.Mutex
	Count  int64
	Action string
}

func NewEventConsumer(
	amqpCh libamqp.Channel,
	config Config,
	httpClient *http.Client,
	logger zerolog.Logger,
) *eventConsumer {
	return &eventConsumer{
		amqpCh:     amqpCh,
		httpClient: httpClient,
		config:     config,
		logger:     logger,

		locks: make(map[string]*pbhLock),
	}
}

func (c *eventConsumer) Start(ctx context.Context) error {
	_, err := c.amqpCh.QueueDeclare(
		Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("cannot declare rmq queue: %w", err)
	}

	err = c.amqpCh.QueueBind(
		Queue,
		RoutingKey,
		canopsis.CanopsisEventsExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("cannot bind rmq queue: %w", err)
	}

	msgs, err := c.amqpCh.Consume(
		Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("cannot consume from rmq channel: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < workersCount; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case msg, ok := <-msgs:
					if !ok {
						return nil
					}

					err := c.processMsg(ctx, msg)
					if err != nil {
						nackErr := c.amqpCh.Nack(msg.DeliveryTag, false, true)
						if nackErr != nil {
							c.logger.Err(nackErr).Msg("cannot nack amqp delivery")
						}

						return err
					}

					err = c.amqpCh.Ack(msg.DeliveryTag, false)
					if err != nil {
						c.logger.Err(err).Msg("cannot ack amqp delivery")
					}
				}
			}
		})
	}

	return g.Wait()
}

func (c *eventConsumer) processMsg(ctx context.Context, msg amqp.Delivery) error {
	c.logger.Debug().Str("event", string(msg.Body)).Msg("received event")

	event, err := fastjson.ParseBytes(msg.Body)
	if err != nil {
		c.logger.Err(err).Str("event", string(msg.Body)).Msg("cannot parse event")
		return nil
	}

	eventType := string(event.GetStringBytes("event_type"))
	if eventType != EventTypePbehavior {
		c.logger.Err(err).
			Str("event_type", eventType).
			Str("event", string(msg.Body)).
			Msg("unknown event type")
		return nil
	}

	action := string(event.GetStringBytes("action"))
	component := string(event.GetStringBytes("component"))
	resource := string(event.GetStringBytes("resource"))
	start := event.GetInt("start")
	end := event.GetInt("end")

	uniquePbhName := strings.Join([]string{
		string(event.GetStringBytes("connector")),
		string(event.GetStringBytes("pbehavior_name")),
		event.Get("downtime_id").String(),
		event.Get("entry").String(),
		event.Get("timestamp").String(),
	}, "-")

	skip := c.lock(uniquePbhName, action)
	defer func() {
		c.unlock(uniquePbhName, action, skip)
	}()

	if skip {
		c.logger.Debug().Str("event", string(msg.Body)).Msg("skip event")
		return nil
	}

	hgValues := event.GetArray("hostgroups")
	if hgValues != nil {
		hostgroups := make([]string, 0, len(hgValues))
		for _, v := range hgValues {
			if s := v.GetStringBytes(); s != nil {
				hostgroups = append(hostgroups, string(s))
			}
		}

		skip = c.inactiveMatch(event.GetFloat64("timestamp"), hostgroups)

		if skip {
			c.logger.Debug().Str("event", string(msg.Body)).Msg("inactive matches: skip event")
			return nil
		}
	}

	entityID := ""
	if resource == "" {
		entityID = component
	} else {
		entityID = fmt.Sprintf("%s/%s", resource, component)
	}

	switch action {
	case ActionCreate:
		arena := &fastjson.Arena{}

		pbhFilter := arena.NewObject()
		pbhFilter.Set("_id", arena.NewString(entityID))

		requestBody := arena.NewObject()
		requestBody.Set("name", arena.NewString(uniquePbhName))
		requestBody.Set("filter", pbhFilter)
		requestBody.Set("tstart", arena.NewNumberInt(start))
		requestBody.Set("tstop", arena.NewNumberInt(end))
		requestBody.Set("type", arena.NewString(c.config.Pbehavior.Type))
		requestBody.Set("reason", arena.NewString(c.config.Pbehavior.Reason))
		requestBody.Set("enabled", arena.NewTrue())

		request, err := c.config.Api.CreateRequest(ctx, http.MethodPost, ApiEndpoint, requestBody.MarshalTo(nil), nil)
		if err != nil {
			return err
		}

		response, err := c.doRequest(request, "create pbehavior in canopsis api")
		if err != nil {
			return fmt.Errorf("cannot execute api request: %w", err)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusCreated {
			b, err := io.ReadAll(response.Body)
			if err == nil && response.StatusCode == http.StatusBadRequest {
				responseBody, err := fastjson.ParseBytes(b)
				if err == nil {
					valErr := string(responseBody.GetStringBytes("errors", "name"))
					if strings.Contains(valErr, "Name already exists") {
						c.logger.Debug().Msg("pbehavior already created")
						return nil
					}
				}
			}

			c.logger.Error().
				Int("response_code", response.StatusCode).
				Str("response", string(b)).
				Msg("cannot create pbehavior")
		}
	case ActionDelete:
		query := make(url.Values)
		query.Set("name", uniquePbhName)
		request, err := c.config.Api.CreateRequest(ctx, http.MethodDelete, ApiEndpoint, nil, query)
		if err != nil {
			return err
		}

		response, err := c.doRequest(request, "delete pbehavior from canopsis api")
		if err != nil {
			return err
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusNoContent {
			b, _ := io.ReadAll(response.Body)

			c.logger.Error().
				Int("response_code", response.StatusCode).
				Str("response", string(b)).
				Msg("cannot delete pbehavior")
		}
	default:
		c.logger.Err(err).
			Str("action", action).
			Str("event", string(msg.Body)).
			Msg("unknown event action")
	}

	return nil
}

func (c *eventConsumer) doRequest(request *http.Request, logMsg string) (*http.Response, error) {
	var reqDump, resDump []byte
	if c.logger.GetLevel() == zerolog.DebugLevel {
		reqDump, _ = httputil.DumpRequest(request, true)
	}

	retryCount := httpRetryCount + 1
	timeout := httpRetryInterval
	var response *http.Response
	var err error

	for i := 0; i < retryCount; i++ {
		response, err = c.httpClient.Do(request)
		if err == nil {
			retry := response.StatusCode != http.StatusOK &&
				response.StatusCode != http.StatusCreated &&
				response.StatusCode != http.StatusNoContent &&
				response.StatusCode != http.StatusBadRequest

			if c.logger.GetLevel() == zerolog.DebugLevel {
				resDump, _ = httputil.DumpResponse(response, true)
			}

			c.logger.Debug().
				Str("request", string(reqDump)).
				Str("response", string(resDump)).
				Int("retry", i).
				Msg(logMsg)

			if !retry {
				break
			}
		} else {
			c.logger.Debug().
				Err(err).
				Str("request", string(reqDump)).
				Int("retry", i).
				Msg(logMsg)
		}

		if i < retryCount-1 {
			time.Sleep(timeout)
			timeout *= 2
		}
	}

	return response, err
}

// lock prevents to process events for the same pbehavior simultaneously.
// It returns skip=true if events aren't processed in order.
func (c *eventConsumer) lock(pbhName, action string) (skip bool) {
	lock := c.getLockOnLock(pbhName)
	lock.Mx.Lock()

	return lock.Action == ActionDelete && action == ActionCreate
}

func (c *eventConsumer) unlock(pbhName, action string, skip bool) {
	lock := c.getLockOnUnlock(pbhName, action, skip)
	if lock != nil {
		lock.Mx.Unlock()
	}
}

// getLockOnLock gets a lock from locks or creates new one on absence.
func (c *eventConsumer) getLockOnLock(pbhName string) *pbhLock {
	c.locksMx.Lock()
	defer c.locksMx.Unlock()

	var lock *pbhLock
	var ok bool
	if lock, ok = c.locks[pbhName]; !ok {
		lock = &pbhLock{}
		c.locks[pbhName] = lock
	}

	lock.Count++

	return lock
}

// getLockOnUnlock gets a lock from locks and deletes it if there is not any waiting goroutines.
// It updates last event data in lock if an event wasn't skipped.
func (c *eventConsumer) getLockOnUnlock(pbhName, action string, skip bool) *pbhLock {
	c.locksMx.Lock()
	defer c.locksMx.Unlock()

	if lock, ok := c.locks[pbhName]; ok {
		lock.Count--

		if lock.Count == 0 {
			delete(c.locks, pbhName)
		} else if !skip {
			lock.Action = action
		}

		return lock
	}

	return nil
}

func (c *eventConsumer) inactiveMatch(eventTs float64, eventHostgroups []string) (match bool) {
	if eventTs == 0 || len(eventHostgroups) == 0 {
		return match
	}
	eventTime := time.Unix(int64(eventTs), 0).In(c.config.location)
	weekdayIndex := int(eventTime.Weekday()) - 1
	if weekdayIndex < 0 {
		weekdayIndex = 6
	}
	h, err := strconv.Atoi(eventTime.Format("1504"))
	if err != nil {
		return match
	}
	for _, inactive := range c.config.Inactive {
		timeMatches := false
		for _, hr := range inactive.hourRanges {
			if !hr.weekdays[weekdayIndex] {
				break
			}
			timeMatches = hr.hhmm[0] <= h && h < hr.hhmm[1] && hr.hhmm[0] < hr.hhmm[1] ||
				hr.hhmm[0] > hr.hhmm[1] && (hr.hhmm[0] <= h || h < hr.hhmm[1])
			if timeMatches {
				break
			}
		}
		if !timeMatches {
			continue
		}
		// check hostgroups
		for _, v := range inactive.Hostgroups {
			for _, ehg := range eventHostgroups {
				if match = v == ehg; match {
					return match
				}
			}
		}
	}

	return match
}
