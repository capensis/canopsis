package statesetting

import (
	"context"
	"sync"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

const workers = 10

type Listener interface {
	Work(ctx context.Context)
	Listen(ctx context.Context, ch <-chan RuleUpdatedMessage)
}

type RuleUpdatedMessage struct {
	ID      string
	Updated datetime.CpsTime

	OldPattern *pattern.Entity // nil if a state setting is new
	OldType    string

	NewPattern *pattern.Entity
	NewType    string
}

type listener struct {
	entityCollection mongo.DbCollection

	amqpPublisher libamqp.Publisher
	encoder       encoding.Encoder

	periodicalInterval, delay time.Duration

	processedEntityIDsMx sync.Mutex
	processedEntityIDs   map[string]struct{}

	updatedRulesMx sync.Mutex
	updatedRules   map[string]RuleUpdatedMessage

	connector string

	logger zerolog.Logger
}

type entityDoc struct {
	ID   string `bson:"_id"`
	Type string `bson:"type"`
}

func NewListener(
	dbClient mongo.DbClient,
	publisher libamqp.Publisher,
	connector string,
	periodicalInterval, delay time.Duration,
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Listener {
	return &listener{
		entityCollection:   dbClient.Collection(mongo.EntityMongoCollection),
		amqpPublisher:      publisher,
		connector:          connector,
		periodicalInterval: periodicalInterval,
		delay:              delay,
		encoder:            encoder,
		logger:             logger,

		processedEntityIDsMx: sync.Mutex{},
		processedEntityIDs:   make(map[string]struct{}),
		updatedRulesMx:       sync.Mutex{},
		updatedRules:         make(map[string]RuleUpdatedMessage),
	}
}

func (l *listener) Work(ctx context.Context) {
	ticker := time.NewTicker(l.periodicalInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updatedRules := l.flushRules()
			if len(updatedRules) == 0 {
				continue
			}

			err := l.processRules(ctx, updatedRules)
			if err != nil {
				l.logger.Err(err).Msg("failed to process changes in state settings, please use engine-axe to recompute data")
			}
		}
	}
}

func (l *listener) Listen(ctx context.Context, ch <-chan RuleUpdatedMessage) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			l.addRule(msg)
		}
	}
}

func (l *listener) processRules(ctx context.Context, updatedRules map[string]RuleUpdatedMessage) error {
	defer l.clearProcessedEntityIDs()
	g, ctx := errgroup.WithContext(ctx)
	ch := make(chan RuleUpdatedMessage)
	g.Go(func() error {
		defer close(ch)
		for k := range updatedRules {
			select {
			case <-ctx.Done():
				return nil
			case ch <- updatedRules[k]:
			}
		}

		return nil
	})

	for i := 0; i < workers; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case msg, ok := <-ch:
					if !ok {
						return nil
					}

					if msg.NewPattern != nil {
						err := l.processPattern(ctx, *msg.NewPattern, msg.NewType)
						if err != nil {
							return err
						}
					}

					if msg.OldPattern != nil {
						err := l.processPattern(ctx, *msg.OldPattern, msg.OldType)
						if err != nil {
							return err
						}
					}
				}
			}
		})
	}

	return g.Wait()
}

func (l *listener) processPattern(ctx context.Context, pattern pattern.Entity, t string) error {
	if len(pattern) == 0 {
		return nil
	}

	newQuery, err := db.EntityPatternToMongoQuery(pattern, "")
	if err != nil {
		return err
	}

	cursor, err := l.entityCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"type": t,
			},
		},
		{
			"$match": newQuery,
		},
		{
			"$project": bson.M{
				"_id":  1,
				"type": 1,
			},
		},
	})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc entityDoc
		err = cursor.Decode(&doc)
		if err != nil {
			return err
		}

		if l.isProcessedEntityID(doc.ID) {
			continue
		}

		err = l.publishEvent(ctx, doc.ID, doc.Type)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *listener) publishEvent(ctx context.Context, entID, entType string) error {
	var event types.Event

	switch entType {
	case types.EntityTypeComponent:
		event = types.Event{
			EventType:     types.EventTypeEntityUpdated,
			Connector:     l.connector,
			ConnectorName: l.connector,
			SourceType:    types.SourceTypeComponent,
			Component:     entID,
		}
	case types.EntityTypeService:
		event = types.Event{
			EventType:     types.EventTypeRecomputeEntityService,
			Connector:     l.connector,
			ConnectorName: l.connector,
			SourceType:    types.SourceTypeService,
			Component:     entID,
		}
	}

	event.Timestamp = datetime.NewCpsTime()
	event.Author = canopsis.DefaultEventAuthor
	event.Initiator = types.InitiatorSystem

	body, err := l.encoder.Encode(event)
	if err != nil {
		return err
	}

	return l.amqpPublisher.PublishWithContext(
		ctx,
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (l *listener) addRule(msg RuleUpdatedMessage) {
	l.updatedRulesMx.Lock()
	defer l.updatedRulesMx.Unlock()

	if prev, ok := l.updatedRules[msg.ID]; ok {
		msg.OldType = prev.OldType
		msg.OldPattern = prev.OldPattern
		l.updatedRules[msg.ID] = msg

		return
	}

	l.updatedRules[msg.ID] = msg
}

func (l *listener) flushRules() map[string]RuleUpdatedMessage {
	l.updatedRulesMx.Lock()
	defer l.updatedRulesMx.Unlock()

	updatedRules := l.updatedRules
	l.updatedRules = make(map[string]RuleUpdatedMessage, len(l.updatedRules))
	now := datetime.NewCpsTime()
	for key, msg := range l.updatedRules {
		if now.Sub(msg.Updated.Time) < l.delay {
			l.updatedRules[key] = msg
			delete(updatedRules, key)
		}
	}

	return updatedRules
}

func (l *listener) clearProcessedEntityIDs() {
	l.processedEntityIDsMx.Lock()
	defer l.processedEntityIDsMx.Unlock()

	clear(l.processedEntityIDs)
}

func (l *listener) isProcessedEntityID(id string) bool {
	l.processedEntityIDsMx.Lock()
	defer l.processedEntityIDsMx.Unlock()

	if _, ok := l.processedEntityIDs[id]; ok {
		return true
	}

	l.processedEntityIDs[id] = struct{}{}

	return false
}
