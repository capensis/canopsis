package bdd

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/cucumber/godog"
	"github.com/kylelemons/godebug/pretty"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

// stepTimeout is used to limit waiting time for wait steps.
const stepTimeout = 10 * time.Second

const (
	consumePrefetchCount = 1000
	consumePrefetchSize  = 0
)

// AmqpClient represents utility struct which implements AMQP steps to feature context.
type AmqpClient struct {
	amqpConnection libamqp.Connection
	mongoClient    mongo.DbClient
	encoder        encoding.Encoder
	decoder        encoding.Decoder
	eventLogger    zerolog.Logger
	exchange, key  string
	templater      *Templater

	ackConsumersMx   sync.Mutex
	ackConsumers     []<-chan amqp.Delivery
	freeAckConsumers []int
}

// NewAmqpClient creates new AMQP client.
func NewAmqpClient(
	dbClient mongo.DbClient,
	amqpConnection libamqp.Connection,
	exchange, key string,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	eventLogger zerolog.Logger,
	templater *Templater,
) *AmqpClient {
	return &AmqpClient{
		amqpConnection: amqpConnection,
		mongoClient:    dbClient,
		encoder:        encoder,
		decoder:        decoder,
		eventLogger:    eventLogger,
		exchange:       exchange,
		key:            key,
		templater:      templater,

		ackConsumers:     make([]<-chan amqp.Delivery, 0),
		freeAckConsumers: make([]int, 0),
	}
}

func (c *AmqpClient) BeforeScenario(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	c.ackConsumersMx.Lock()
	defer c.ackConsumersMx.Unlock()

	if len(c.freeAckConsumers) > 0 {
		consumer := c.freeAckConsumers[0]
		ctx = setConsumer(ctx, consumer)
		c.freeAckConsumers = c.freeAckConsumers[1:]

		for {
			select {
			case <-c.ackConsumers[consumer]:
			default:
				return ctx, nil
			}
		}
	}

	ch, err := c.amqpConnection.Channel()
	if err != nil {
		return ctx, err
	}

	// Declare queue to detect the end of event processing.
	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return ctx, fmt.Errorf("cannot declare queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,
		c.key,
		c.exchange,
		false,
		nil,
	)
	if err != nil {
		return ctx, fmt.Errorf("cannot bind queue: %v", err)
	}

	err = ch.Qos(consumePrefetchCount, consumePrefetchSize, false)
	if err != nil {
		return ctx, fmt.Errorf("cannot qos: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return ctx, fmt.Errorf("cannot consume queue: %v", err)
	}

	c.ackConsumers = append(c.ackConsumers, msgs)
	ctx = setConsumer(ctx, len(c.ackConsumers)-1)

	return ctx, nil
}

func (c *AmqpClient) AfterScenario(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	c.ackConsumersMx.Lock()
	defer c.ackConsumersMx.Unlock()

	if i, ok := getConsumer(ctx); ok {
		c.freeAckConsumers = append(c.freeAckConsumers, i)
	}

	return ctx, nil
}

// IWaitTheEndOfEventProcessing
// Step example:
//
//	When I wait the end of event processing
func (c *AmqpClient) IWaitTheEndOfEventProcessing(ctx context.Context) error {
	return c.IWaitTheEndOfEventsProcessing(ctx, 1)
}

// IWaitTheEndOfEventsProcessing
// Step example:
//
//	When I wait the end of 2 events processing
func (c *AmqpClient) IWaitTheEndOfEventsProcessing(ctx context.Context, expectedCount int) error {
	msgs, err := c.getMainStreamAckConsumer(ctx)
	if err != nil {
		return err
	}
	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()
	caughtCount := 0

	scName, _ := GetScenarioName(ctx)
	scUri, _ := GetScenarioUri(ctx)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				return errors.New("consume chan is closed")
			}

			caughtCount++

			event := &types.Event{}
			err := c.decoder.Decode(d.Body, event)
			c.eventLogger.Info().Err(err).
				Str("event_type", event.EventType).
				Str("entity", event.GetEID()).
				Str("file", scUri).
				Str("scenario", scName).
				RawJSON("body", d.Body).Msg("received event")

			if expectedCount == caughtCount {
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("reached timeout: caught %d events out of %d", caughtCount, expectedCount)
		}
	}
}

// IWaitTheEndOfEventProcessingWhichContains
//
// Step example:
//
//	When I wait the end of event processing which contains:
//	  """
//	  {
//	    "event_type": "check"
//	  }
//	  """
func (c *AmqpClient) IWaitTheEndOfEventProcessingWhichContains(ctx context.Context, doc string) error {
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	expectedEvent := make(map[string]interface{})
	err = c.decoder.Decode(b.Bytes(), &expectedEvent)
	if err != nil {
		return err
	}

	return c.catchEvents(ctx, []map[string]interface{}{expectedEvent})
}

// IWaitTheEndOfEventsProcessingWhichContain
//
// Step example:
//
//		When I wait the end of events processing which contain:
//	   """
//	   [
//	     {
//	       "event_type": "check"
//	     },
//	     {
//	       "event_type": "check"
//	     }
//	   ]
//	   """
func (c *AmqpClient) IWaitTheEndOfEventsProcessingWhichContain(ctx context.Context, doc string) error {
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	expectedEvents := make([]map[string]interface{}, 0)
	err = c.decoder.Decode(b.Bytes(), &expectedEvents)
	if err != nil {
		return err
	}

	return c.catchEvents(ctx, expectedEvents)
}

// IWaitTheEndOfOneOfEventsProcessingWhichContain
//
// Step example:
//
//		When I wait the end of one of events processing which contain:
//	   """
//	   [
//	     {
//	       "event_type": "check"
//	     },
//	     {
//	       "event_type": "activate"
//	     }
//	   ]
//	   """
func (c *AmqpClient) IWaitTheEndOfOneOfEventsProcessingWhichContain(ctx context.Context, doc string) error {
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	expectedEvents := make([]map[string]interface{}, 0)
	err = c.decoder.Decode(b.Bytes(), &expectedEvents)
	if err != nil {
		return err
	}

	if len(expectedEvents) == 0 {
		return nil
	}

	caughtEvents := make([]map[string]interface{}, 0)
	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()

	msgs, err := c.getMainStreamAckConsumer(ctx)
	if err != nil {
		return err
	}

	scName, _ := GetScenarioName(ctx)
	scUri, _ := GetScenarioUri(ctx)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				return errors.New("consume chan is closed")
			}

			event, eventMap, err := c.decodeEvent(d.Body)
			if err != nil {
				c.eventLogger.Err(err).
					RawJSON("body", d.Body).
					Str("file", scUri).
					Str("scenario", scName).
					Msg("received invalid event")
				continue
			}

			caughtEvents = append(caughtEvents, eventMap)
			foundIndex := c.matchEvent(eventMap, expectedEvents)
			c.eventLogger.Info().
				Str("event_type", event.EventType).
				Str("entity", event.GetEID()).
				Bool("matched", foundIndex >= 0).
				RawJSON("body", d.Body).
				Str("file", scUri).
				Str("scenario", scName).
				Msg("received event")

			if foundIndex >= 0 {
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("reached timeout: caught %d events but none of them matches to expected events\n%s\n",
				len(caughtEvents), pretty.Compare(caughtEvents, expectedEvents))
		}
	}
}

func (c *AmqpClient) IWaitTheEndOfSentEventProcessing(ctx context.Context, doc string) error {
	b, err := c.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	sentEvents := make([]map[string]interface{}, 0)
	err = c.decoder.Decode(b.Bytes(), &sentEvents)
	if err != nil {
		sentEvent := make(map[string]interface{})
		err = c.decoder.Decode(b.Bytes(), &sentEvent)
		if err != nil {
			return err
		}
		sentEvents = append(sentEvents, sentEvent)
	}

	if len(sentEvents) == 0 {
		return nil
	}

	expectedFields := []string{"event_type", "source_type", "connector", "connector_name", "component", "resource"}
	expectedEventsByIndex := make([][]map[string]interface{}, len(sentEvents))

	for i, sentEvent := range sentEvents {
		expectedEvent := make(map[string]interface{}, len(expectedFields))
		for _, field := range expectedFields {
			if v, ok := sentEvent[field]; ok {
				expectedEvent[field] = v
			}
		}
		expectedEvents := []map[string]interface{}{expectedEvent}

		if t, ok := expectedEvent["event_type"].(string); ok && t == types.EventTypeCheck {
			activateEvent := make(map[string]interface{}, len(expectedEvent))
			for k, v := range expectedEvent {
				activateEvent[k] = v
			}
			activateEvent["event_type"] = types.EventTypeActivate
			expectedEvents = append(expectedEvents, activateEvent)
		}

		expectedEventsByIndex[i] = expectedEvents
	}

	caughtEvents := make([]map[string]interface{}, 0)
	matched := make(map[int]bool, len(sentEvents))
	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()

	msgs, err := c.getMainStreamAckConsumer(ctx)
	if err != nil {
		return err
	}

	scName, _ := GetScenarioName(ctx)
	scUri, _ := GetScenarioUri(ctx)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				return errors.New("consume chan is closed")
			}

			event, eventMap, err := c.decodeEvent(d.Body)
			if err != nil {
				c.eventLogger.Err(err).
					RawJSON("body", d.Body).
					Str("file", scUri).
					Str("scenario", scName).
					Msg("received invalid event")
				continue
			}

			caughtEvents = append(caughtEvents, eventMap)
			foundIndex := -1
			for index, expectedEvents := range expectedEventsByIndex {
				if matched[index] {
					continue
				}

				foundIndex = c.matchEvent(eventMap, expectedEvents)
				if foundIndex >= 0 {
					matched[index] = true
					break
				}
			}

			c.eventLogger.Info().
				Str("event_type", event.EventType).
				Str("entity", event.GetEID()).
				Bool("matched", foundIndex >= 0).
				Str("file", scUri).
				Str("scenario", scName).
				RawJSON("body", d.Body).Msg("received event")

			if len(matched) == len(sentEvents) {
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("reached timeout: caught %d events out of %d\n%s\n", len(matched), len(sentEvents),
				pretty.Compare(caughtEvents, sentEvents))
		}
	}
}

// ICallRPCAxeRequest
// Step example:
//
//	When I call RPC request to engine-axe with alarm resource/component:
//	"""
//	{
//	  "event_type": "ack"
//	}
//	"""
func (c *AmqpClient) ICallRPCAxeRequest(ctx context.Context, eid string, doc string) error {
	alarm, err := c.findAlarm(ctx, eid)
	if err != nil {
		return err
	}

	var event rpc.AxeEvent
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

	var resultEvent rpc.AxeResultEvent
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

	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()

	select {
	case <-timer.C:
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

func (c *AmqpClient) catchEvents(ctx context.Context, expectedEvents []map[string]interface{}) error {
	if len(expectedEvents) == 0 {
		return nil
	}

	caughtEvents := make([]map[string]interface{}, 0)
	matched := make(map[int]bool, len(expectedEvents))
	timer := time.NewTimer(stepTimeout)
	defer timer.Stop()

	msgs, err := c.getMainStreamAckConsumer(ctx)
	if err != nil {
		return err
	}

	scName, _ := GetScenarioName(ctx)
	scUri, _ := GetScenarioUri(ctx)

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				return errors.New("consume chan is closed")
			}

			event, eventMap, err := c.decodeEvent(d.Body)
			if err != nil {
				c.eventLogger.
					Err(err).
					RawJSON("body", d.Body).
					Str("file", scUri).
					Str("scenario", scName).
					Msg("received invalid event")
				continue
			}

			caughtEvents = append(caughtEvents, eventMap)
			foundIndex := -1
			for i, expectedEvent := range expectedEvents {
				if matched[i] {
					continue
				}
				m := true
				for k, v := range expectedEvent {
					if v != eventMap[k] {
						m = false
						break
					}
				}
				if m {
					matched[i] = true
					foundIndex = i
					break
				}
			}

			c.eventLogger.Info().
				Str("event_type", event.EventType).
				Str("entity", event.GetEID()).
				Bool("matched", foundIndex >= 0).
				Str("file", scUri).
				Str("scenario", scName).
				RawJSON("body", d.Body).Msg("received event")

			if len(matched) == len(expectedEvents) {
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("reached timeout: caught %d events out of %d\n%s\n", len(matched), len(expectedEvents),
				pretty.Compare(caughtEvents, expectedEvents))
		}
	}
}

func (c *AmqpClient) getMainStreamAckConsumer(ctx context.Context) (<-chan amqp.Delivery, error) {
	c.ackConsumersMx.Lock()
	defer c.ackConsumersMx.Unlock()

	if i, ok := getConsumer(ctx); ok {
		return c.ackConsumers[i], nil
	}

	return nil, fmt.Errorf("scenario isn't started")
}

func (c *AmqpClient) decodeEvent(msg []byte) (types.Event, map[string]interface{}, error) {
	event := types.Event{}
	eventMap := make(map[string]interface{})
	err := c.decoder.Decode(msg, &eventMap)
	if err != nil {
		return event, eventMap, err
	}
	err = c.decoder.Decode(msg, &event)
	return event, eventMap, err
}

func (c *AmqpClient) matchEvent(event map[string]interface{}, expectedEvents []map[string]interface{}) int {
	for i, expectedEvent := range expectedEvents {
		matched := true
		for k, v := range expectedEvent {
			if v != event[k] {
				matched = false
				break
			}
		}
		if matched {
			return i
		}
	}

	return -1
}
