package bdd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"text/template"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/cucumber/godog"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

// stepTimeout is used to limit waiting time for wait steps.
const stepTimeout = 10 * time.Second

// AmqpClient represents utility struct which implements AMQP steps to feature context.
type AmqpClient struct {
	amqpConnection    libamqp.Connection
	mongoClient       mongo.DbClient
	mainStreamAckMsgs <-chan amqp.Delivery
	apiURL            *url.URL
	encoder           encoding.Encoder
	decoder           encoding.Decoder
	eventLogger       zerolog.Logger
}

// NewAmqpClient creates new AMQP client.
func NewAmqpClient(
	dbClient mongo.DbClient,
	amqpConnection libamqp.Connection,
	exchange, key string,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	eventLogger zerolog.Logger,
) (*AmqpClient, error) {
	ch, err := amqpConnection.Channel()
	if err != nil {
		return nil, err
	}

	// Declare queue to detect the end of event processing.
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("cannot declare queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("cannot consume queue: %v", err)
	}

	apiURL, err := GetApiURL()
	if err != nil {
		return nil, err
	}

	return &AmqpClient{
		amqpConnection:    amqpConnection,
		mongoClient:       dbClient,
		mainStreamAckMsgs: msgs,
		encoder:           encoder,
		decoder:           decoder,
		apiURL:            apiURL,
		eventLogger:       eventLogger,
	}, nil
}

func (c *AmqpClient) Reset(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	// Clear channel.
	for {
		select {
		case <-c.mainStreamAckMsgs:
		default:
			return ctx, nil
		}
	}
}

/*
IWaitTheEndOfEventProcessing
Step example:
	When I wait the end of event processing
*/
func (c *AmqpClient) IWaitTheEndOfEventProcessing() error {
	return c.IWaitTheEndOfEventsProcessing(1)
}

/*
IWaitTheEndOfEventsProcessing
Step example:
	When I wait the end of 2 events processing
*/
func (c *AmqpClient) IWaitTheEndOfEventsProcessing(count int) error {
	done := time.After(stepTimeout)
	msgsCount := 0

	for {
		select {
		case d, ok := <-c.mainStreamAckMsgs:
			if !ok {
				return errors.New("consume chan is closed")
			}

			msgsCount++

			event := &types.Event{}
			err := c.decoder.Decode(d.Body, event)
			c.eventLogger.Info().Err(err).
				Str("event_type", event.EventType).
				Str("entity", event.GetEID()).
				RawJSON("body", d.Body).Msg("received event")

			if count == msgsCount {
				return nil
			}
		case <-done:
			return fmt.Errorf("reached timeout: waiting for %d events, but got %d", count, msgsCount)
		}
	}
}

/*
ICallRPCAxeRequest
Step example:
	When I call RPC request to engine-axe with alarm resource/component:
	"""
	{
	  "event_type": "ack"
	}
	"""
*/
func (c *AmqpClient) ICallRPCAxeRequest(ctx context.Context, eid string, doc string) error {
	alarm, err := c.findAlarm(ctx, eid)
	if err != nil {
		return err
	}

	var event types.RPCAxeEvent
	err = c.decoder.Decode([]byte(doc), &event)
	if err != nil {
		return err
	}

	event.Alarm = &alarm.Alarm
	event.Entity = &alarm.Entity
	body, err := c.encoder.Encode(event)
	if err != nil {
		return err
	}

	res, err := c.executeRPC(ctx, canopsis.AxeRPCQueueServerName, body)
	if err != nil {
		return err
	}

	var resultEvent types.RPCAxeResultEvent
	err = c.decoder.Decode(res, &resultEvent)
	if err != nil {
		return err
	}

	if resultEvent.Error != nil {
		return resultEvent.Error.Error
	}

	return nil
}

/*
ICallRPCWebhookRequest
Step example:
	When I call RPC request to engine-webhook with alarm resource/component:
	"""
	{
	  "request": {
		  "method": "POST",
		  "url": "http://test-url.com"
		}
	}
	"""
*/
func (c *AmqpClient) ICallRPCWebhookRequest(ctx context.Context, eid string, doc string) error {
	alarm, err := c.findAlarm(ctx, eid)
	if err != nil {
		return err
	}

	content, err := c.executeTemplate(doc)
	if err != nil {
		return err
	}

	var event types.RPCWebhookEvent
	err = c.decoder.Decode(content.Bytes(), &event)
	if err != nil {
		return err
	}

	event.Alarm = &alarm.Alarm
	event.Entity = &alarm.Entity
	body, err := c.encoder.Encode(event)
	if err != nil {
		return err
	}

	res, err := c.executeRPC(ctx, canopsis.WebhookRPCQueueServerName, body)
	if err != nil {
		return err
	}

	var resultEvent types.RPCWebhookResultEvent
	err = c.decoder.Decode(res, &resultEvent)
	if err != nil {
		return err
	}

	if resultEvent.Error != nil {
		return resultEvent.Error.Error
	}

	return nil
}

func (c *AmqpClient) findAlarm(ctx context.Context, eid string) (*types.AlarmWithEntity, error) {
	cursor, err := c.mongoClient.Collection(mongo.AlarmMongoCollection).Aggregate(ctx, []bson.M{
		{"$match": bson.M{"d": eid}},
		{"$sort": bson.M{"v.creation_date": -1}},
		{"$limit": 1},
		{"$project": bson.M{
			"alarm": "$$ROOT",
			"_id":   0,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		var alarm types.AlarmWithEntity
		err := cursor.Decode(&alarm)
		if err != nil {
			return nil, err
		}

		return &alarm, nil
	}

	return nil, fmt.Errorf("couldn't find an alarm for eid = %s", eid)
}

func (c *AmqpClient) executeRPC(ctx context.Context, queue string, body []byte) ([]byte, error) {
	publishCh, err := c.amqpConnection.Channel()
	if err != nil {
		return nil, err
	}

	defer publishCh.Close()

	consumeCh, err := c.amqpConnection.Channel()
	if err != nil {
		return nil, err
	}

	defer consumeCh.Close()

	// Declare queue to receive RPC response.
	q, err := consumeCh.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("cannot declare queue: %v", err)
	}

	msgs, err := consumeCh.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, fmt.Errorf("cannot consume queue: %v", err)
	}

	corrID := fmt.Sprintf("test-%d", rand.Int())
	err = publishCh.PublishWithContext(
		ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			Body:          body,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot publish to queue: %v", err)
	}

	select {
	case <-time.After(stepTimeout):
		return nil, errors.New("reached timeout")
	case d, ok := <-msgs:
		if !ok {
			return nil, errors.New("channel is closed")
		}

		if d.CorrelationId != corrID {
			return nil, fmt.Errorf("expected correlation id %v but got %v", corrID, d.CorrelationId)
		}

		return d.Body, nil
	}
}

func (c *AmqpClient) executeTemplate(tpl string) (*bytes.Buffer, error) {
	t, err := template.New("tpl").Option("missingkey=error").Parse(tpl)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"apiUrl": c.apiURL,
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
