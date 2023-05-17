package alarmaction

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libauthor "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Ack(ctx context.Context, id string, r AckRequest, userId string) (bool, error)
	AckRemove(ctx context.Context, id string, r Request, userId string) (bool, error)
	Snooze(ctx context.Context, id string, r SnoozeRequest, userId string) (bool, error)
	Cancel(ctx context.Context, id string, r Request, userId string) (bool, error)
	Uncancel(ctx context.Context, id string, r Request, userId string) (bool, error)
	AssocTicket(ctx context.Context, id string, r AssocTicketRequest, userId string) (bool, error)
	Comment(ctx context.Context, id string, r CommentRequest, userId string) (bool, error)
	ChangeState(ctx context.Context, id string, r ChangeStateRequest, userId string) (bool, error)
}

func NewStore(
	dbClient mongo.DbClient,
	amqpPublisher libamqp.Publisher,
	exchange, queue string,
	encoder encoding.Encoder,
	contentType string,
	authorProvider libauthor.Provider,
	logger zerolog.Logger,
) Store {
	return &store{
		dbClient:         dbClient,
		dbCollection:     dbClient.Collection(mongo.AlarmMongoCollection),
		userDbCollection: dbClient.Collection(mongo.RightsMongoCollection),
		amqpPublisher:    amqpPublisher,
		exchange:         exchange,
		queue:            queue,
		encoder:          encoder,
		contentType:      contentType,
		authorProvider:   authorProvider,
		logger:           logger,
	}
}

type store struct {
	dbClient         mongo.DbClient
	dbCollection     mongo.DbCollection
	userDbCollection mongo.DbCollection
	amqpPublisher    libamqp.Publisher
	exchange, queue  string
	encoder          encoding.Encoder
	contentType      string
	authorProvider   libauthor.Provider
	logger           zerolog.Logger
}

func (s *store) Ack(ctx context.Context, id string, r AckRequest, userId string) (bool, error) {
	// Double ack can be enabled. Check in engine-axe.
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	author, err := s.getAuthor(ctx, userId)
	if err != nil {
		return false, err
	}

	event := types.Event{
		EventType: types.EventTypeAck,
		Output:    r.Comment,
		Component: alarm.Alarm.Value.Component,
		Resource:  alarm.Alarm.Value.Resource,
		UserID:    userId,
		Author:    author,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	if r.AckResources && alarm.Entity.Type == types.EntityTypeComponent {
		go func() {
			err = s.ackResources(context.Background(), alarm.Entity.ID, r.Comment, userId, author)
			if err != nil {
				s.logger.Err(err).Msg("cannot ack resources")
			}
		}()
	}

	return true, nil
}

func (s *store) AckRemove(ctx context.Context, id string, r Request, userId string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.ack": bson.M{"$ne": nil}})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType: types.EventTypeAckremove,
		Output:    r.Comment,
		Component: alarm.Alarm.Value.Component,
		Resource:  alarm.Alarm.Value.Resource,
		UserID:    userId,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Snooze(ctx context.Context, id string, r SnoozeRequest, userId string) (bool, error) {
	d, err := r.Duration.To("s")
	if err != nil {
		return false, common.NewValidationError("duration", "Duration is invalid.")
	}

	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.snooze": nil})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType: types.EventTypeSnooze,
		Output:    r.Comment,
		Duration:  types.CpsNumber(d.Value),
		Component: alarm.Alarm.Value.Component,
		Resource:  alarm.Alarm.Value.Resource,
		UserID:    userId,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Cancel(ctx context.Context, id string, r Request, userId string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.canceled": nil})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType: types.EventTypeCancel,
		Output:    r.Comment,
		Component: alarm.Alarm.Value.Component,
		Resource:  alarm.Alarm.Value.Resource,
		UserID:    userId,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Uncancel(ctx context.Context, id string, r Request, userId string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.canceled": bson.M{"$ne": nil}})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType: types.EventTypeUncancel,
		Output:    r.Comment,
		Component: alarm.Alarm.Value.Component,
		Resource:  alarm.Alarm.Value.Resource,
		UserID:    userId,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) AssocTicket(ctx context.Context, id string, r AssocTicketRequest, userId string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	author, err := s.getAuthor(ctx, userId)
	if err != nil {
		return false, err
	}

	ticketInfo := types.TicketInfo{
		Ticket:           r.Ticket,
		TicketURL:        r.Url,
		TicketComment:    r.Comment,
		TicketSystemName: r.SystemName,
		TicketData:       r.Data,
	}
	event := types.Event{
		EventType:  types.EventTypeAssocTicket,
		TicketInfo: ticketInfo,
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		UserID:     userId,
		Author:     author,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	if r.TicketResources && alarm.Entity.Type == types.EntityTypeComponent {
		go func() {
			err = s.ticketResources(context.Background(), alarm.Entity.ID, ticketInfo, userId, author)
			if err != nil {
				s.logger.Err(err).Msg("cannot ticket resources")
			}
		}()
	}

	return true, nil
}

func (s *store) Comment(ctx context.Context, id string, r CommentRequest, userId string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType: types.EventTypeComment,
		Output:    r.Comment,
		Component: alarm.Alarm.Value.Component,
		Resource:  alarm.Alarm.Value.Resource,
		UserID:    userId,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) ChangeState(ctx context.Context, id string, r ChangeStateRequest, userId string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType: types.EventTypeChangestate,
		State:     types.CpsNumber(*r.State),
		Output:    r.Comment,
		Component: alarm.Alarm.Value.Component,
		Resource:  alarm.Alarm.Value.Resource,
		UserID:    userId,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) findAlarm(ctx context.Context, match bson.M) (types.AlarmWithEntity, error) {
	if match == nil {
		match = bson.M{}
	}
	match["v.resolved"] = nil
	alarm := types.AlarmWithEntity{}
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": match},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$project": bson.M{
			"alarm._id":         "$_id",
			"alarm.v.component": "$v.component",
			"alarm.v.resource":  "$v.resource",
			"entity._id":        1,
			"entity.type":       1,
		}},
	})
	if err != nil {
		return alarm, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&alarm)
		if err != nil {
			return alarm, err
		}
	}

	return alarm, nil
}

func (s *store) sendEvent(ctx context.Context, event types.Event) error {
	if event.Author == "" && event.UserID != "" {
		var err error
		event.Author, err = s.getAuthor(ctx, event.UserID)
		if err != nil {
			return err
		}
	}

	event.Connector = canopsis.ApiName
	event.ConnectorName = canopsis.ApiName
	event.Initiator = types.InitiatorUser
	event.Timestamp = types.NewCpsTime()
	event.SourceType = event.DetectSourceType()
	body, err := s.encoder.Encode(event)
	if err != nil {
		return err
	}

	return s.amqpPublisher.PublishWithContext(
		ctx,
		s.exchange,
		s.queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  s.contentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (s *store) ackResources(
	ctx context.Context,
	component string,
	comment string,
	userId, author string,
) error {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"v.component": component,
			"v.meta":      bson.M{"$exists": false},
			"v.resolved":  nil,
			"v.ack":       nil,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": bson.M{
			"entity.type": types.EntityTypeResource,
		}},
		{"$project": bson.M{
			"alarm._id":         "$_id",
			"alarm.v.component": "$v.component",
			"alarm.v.resource":  "$v.resource",
			"entity._id":        1,
			"entity.type":       1,
		}},
	})
	if err != nil {
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		alarm := types.AlarmWithEntity{}
		err = cursor.Decode(&alarm)
		if err != nil {
			return fmt.Errorf("cannot decode alarm: %w", err)
		}

		event := types.Event{
			EventType: types.EventTypeAck,
			Output:    comment,
			Component: alarm.Alarm.Value.Component,
			Resource:  alarm.Alarm.Value.Resource,
			UserID:    userId,
			Author:    author,
		}
		err = s.sendEvent(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) ticketResources(
	ctx context.Context,
	component string,
	ticketInfo types.TicketInfo,
	userId, author string,
) error {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"v.component": component,
			"v.meta":      bson.M{"$exists": false},
			"v.resolved":  nil,
			"v.ticket":    nil,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": bson.M{
			"entity.type": types.EntityTypeResource,
		}},
		{"$project": bson.M{
			"alarm._id":         "$_id",
			"alarm.v.component": "$v.component",
			"alarm.v.resource":  "$v.resource",
			"entity._id":        1,
			"entity.type":       1,
		}},
	})
	if err != nil {
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		alarm := types.AlarmWithEntity{}
		err = cursor.Decode(&alarm)
		if err != nil {
			return fmt.Errorf("cannot decode alarm: %w", err)
		}

		event := types.Event{
			EventType:  types.EventTypeAssocTicket,
			TicketInfo: ticketInfo,
			Component:  alarm.Alarm.Value.Component,
			Resource:   alarm.Alarm.Value.Resource,
			UserID:     userId,
			Author:     author,
		}
		err = s.sendEvent(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) getAuthor(ctx context.Context, id string) (string, error) {
	author, err := s.authorProvider.Find(ctx, id)
	return author.DisplayName, err
}
