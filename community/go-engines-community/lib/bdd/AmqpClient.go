package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	liblog "git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/cucumber/messages-go/v10"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"net/url"
	"text/template"
	"time"
)

// stepTimeout is used to limit waiting time for wait steps.
const stepTimeout = 10 * time.Second

// AmqpClient represents utility struct which implements AMQP steps to feature context.
type AmqpClient struct {
	amqpConnection    libamqp.Connection
	mongoClient       mongo.DbClient
	mainStreamAckMsgs <-chan amqp.Delivery
	apiURL            *url.URL
	eventLogger       zerolog.Logger
}

// NewAmqpClient creates new AMQP client.
func NewAmqpClient(exchange, key string, eventLogger zerolog.Logger) (*AmqpClient, error) {
	mongoClient, err := mongo.NewClient(0, 0)
	if err != nil {
		return nil, err
	}

	amqpConnection, err := libamqp.NewConnection(liblog.NewLogger(false), 0, 0)
	if err != nil {
		return nil, err
	}

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
		mongoClient:       mongoClient,
		mainStreamAckMsgs: msgs,
		apiURL:            apiURL,
		eventLogger:       eventLogger,
	}, nil
}

func (c *AmqpClient) Reset(_ *messages.Pickle) {
	// Clear channel.
	for {
		select {
		case <-c.mainStreamAckMsgs:
		default:
			return
		}
	}
}

/**
Step example:
	When I wait the end of event processing
*/
func (c *AmqpClient) IWaitTheEndOfEventProcessing() error {
	return c.IWaitTheEndOfEventsProcessing(1)
}

/**
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

			c.eventLogger.Info().RawJSON("event", d.Body).Msg("received")

			if count == msgsCount {
				return nil
			}
		case <-done:
			return errors.New("reached timeout")
		}
	}
}

/**
Step example:
	When I call RPC request to engine-axe with alarm resource/component:
	"""
	{
	  "event_type": "ack"
	}
	"""
*/
func (c *AmqpClient) ICallRPCAxeRequest(eid string, doc *messages.PickleStepArgument_PickleDocString) error {
	alarm, err := c.findAlarm(eid)
	if err != nil {
		return err
	}

	var event types.RPCAxeEvent
	err = json.Unmarshal([]byte(doc.Content), &event)
	if err != nil {
		return err
	}

	event.Alarm = alarm
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	res, err := c.executeRPC(canopsis.AxeRPCQueueServerName, body)
	if err != nil {
		return err
	}

	var resultEvent types.RPCAxeResultEvent
	err = json.Unmarshal(res, &resultEvent)
	if err != nil {
		return err
	}

	if resultEvent.Error != nil {
		return resultEvent.Error.Error
	}

	return nil
}

/**
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
func (c *AmqpClient) ICallRPCWebhookRequest(eid string, doc *messages.PickleStepArgument_PickleDocString) error {
	alarm, err := c.findAlarm(eid)
	if err != nil {
		return err
	}

	content, err := c.executeTemplate(doc.Content)
	if err != nil {
		return err
	}

	var event types.RPCWebhookEvent
	err = json.Unmarshal(content.Bytes(), &event)
	if err != nil {
		return err
	}

	event.Alarm = alarm
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	res, err := c.executeRPC(canopsis.WebhookRPCQueueServerName, body)
	if err != nil {
		return err
	}

	var resultEvent types.RPCWebhookResultEvent
	err = json.Unmarshal(res, &resultEvent)
	if err != nil {
		return err
	}

	if resultEvent.Error != nil {
		return resultEvent.Error.Error
	}

	return nil
}

func (c *AmqpClient) findAlarm(eid string) (*types.Alarm, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res := c.mongoClient.Collection(libalarm.AlarmCollectionName).FindOne(ctx, bson.M{"d": eid})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, fmt.Errorf("couldn't find an alarm for eid = %s", eid)
		}

		return nil, err
	}

	var alarm types.Alarm
	err := res.Decode(&alarm)
	if err != nil {
		return nil, err
	}

	return &alarm, nil
}

func (c *AmqpClient) executeRPC(queue string, body []byte) ([]byte, error) {
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
	err = publishCh.Publish(
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
	t, err := template.New("tpl").Parse(tpl)
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
