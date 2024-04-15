package statesetting

import (
	"context"
	"sync"

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

type Listener interface {
	Listen(ctx context.Context, ch <-chan RuleUpdatedMessage)
}

type RuleUpdatedMessage struct {
	OldPattern *pattern.Entity // nil if a state setting is new
	OldType    string

	NewPattern *pattern.Entity
	NewType    string
}

type listener struct {
	entityCollection mongo.DbCollection

	amqpPublisher libamqp.Publisher
	encoder       encoding.Encoder

	processedIDsMapMx sync.Mutex
	processedIDsMap   map[string]bool

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
	encoder encoding.Encoder,
	logger zerolog.Logger,
) Listener {
	return &listener{
		entityCollection: dbClient.Collection(mongo.EntityMongoCollection),
		amqpPublisher:    publisher,
		connector:        connector,
		encoder:          encoder,
		logger:           logger,

		processedIDsMapMx: sync.Mutex{},
		processedIDsMap:   make(map[string]bool),
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

			err := l.processMessage(ctx, msg)
			if err != nil {
				l.logger.Err(err).Msg("failed to process changes in state settings, please use engine-axe to recompute data")
			}
		}
	}
}

func (l *listener) processMessage(ctx context.Context, msg RuleUpdatedMessage) error {
	g, ctx := errgroup.WithContext(ctx)
	defer clear(l.processedIDsMap)

	if msg.NewPattern != nil {
		g.Go(func() error {
			return l.processPattern(ctx, *msg.NewPattern, msg.NewType)
		})
	}

	if msg.OldPattern != nil {
		g.Go(func() error {
			return l.processPattern(ctx, *msg.OldPattern, msg.OldType)
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

		l.processedIDsMapMx.Lock()

		if l.processedIDsMap[doc.ID] {
			l.processedIDsMapMx.Unlock()
			continue
		}

		l.processedIDsMap[doc.ID] = true

		l.processedIDsMapMx.Unlock()

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
