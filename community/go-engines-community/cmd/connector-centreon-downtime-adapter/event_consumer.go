package main

import (
	"context"
	"fmt"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
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
)

type eventConsumer struct {
	amqpCh     libamqp.Channel
	config     Config
	httpClient *http.Client
	logger     zerolog.Logger
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
	connector := string(event.GetStringBytes("connector"))
	connectorName := string(event.GetStringBytes("connector_name"))
	component := string(event.GetStringBytes("component"))
	resource := string(event.GetStringBytes("resource"))
	pbhName := string(event.GetStringBytes("pbehavior_name"))
	downtimeId := string(event.GetStringBytes("downtime_id"))
	start := event.GetInt("start")
	end := event.GetInt("end")

	uniquePbhName := fmt.Sprintf("%s/%s %s %s", connector, connectorName, pbhName, downtimeId)

	switch action {
	case ActionCreate:
		arena := &fastjson.Arena{}

		pbhFilter := arena.NewObject()
		if resource == "" {
			pbhFilter.Set("_id", arena.NewString(resource))
		} else {
			pbhFilter.Set("_id", arena.NewString(fmt.Sprintf("%s/%s", resource, component)))
		}

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
			b, _ := ioutil.ReadAll(response.Body)

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
			b, _ := ioutil.ReadAll(response.Body)

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
				response.StatusCode != http.StatusNoContent

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
