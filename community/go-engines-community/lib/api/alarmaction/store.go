package alarmaction

import (
	"context"
	"errors"
	"fmt"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Ack(ctx context.Context, id string, r AckRequest, userId, username string) (bool, error)
	AckRemove(ctx context.Context, id string, r Request, userId, username string) (bool, error)
	Snooze(ctx context.Context, id string, r SnoozeRequest, userId, username string) (bool, error)
	Cancel(ctx context.Context, id string, r Request, userId, username string) (bool, error)
	Uncancel(ctx context.Context, id string, r Request, userId, username string) (bool, error)
	AssocTicket(ctx context.Context, id string, r AssocTicketRequest, userId, username string) (bool, error)
	Comment(ctx context.Context, id string, r CommentRequest, userId, username string) (bool, error)
	ChangeState(ctx context.Context, id string, r ChangeStateRequest, userId, username string) (bool, error)
	AddBookmark(ctx context.Context, alarmID, userID string) (bool, error)
	RemoveBookmark(ctx context.Context, alarmID, userID string) (bool, error)
}

func NewStore(
	dbClient libmongo.DbClient,
	amqpPublisher libamqp.Publisher,
	exchange, queue string,
	encoder encoding.Encoder,
	contentType string,
	logger zerolog.Logger,
) Store {
	return &store{
		dbClient:             dbClient,
		dbCollection:         dbClient.Collection(libmongo.AlarmMongoCollection),
		resolvedDbCollection: dbClient.Collection(libmongo.ResolvedAlarmMongoCollection),
		amqpPublisher:        amqpPublisher,
		exchange:             exchange,
		queue:                queue,
		encoder:              encoder,
		contentType:          contentType,
		logger:               logger,
	}
}

type store struct {
	dbClient             libmongo.DbClient
	dbCollection         libmongo.DbCollection
	resolvedDbCollection libmongo.DbCollection
	amqpPublisher        libamqp.Publisher
	exchange, queue      string
	encoder              encoding.Encoder
	contentType          string
	logger               zerolog.Logger
}

func (s *store) Ack(ctx context.Context, id string, r AckRequest, userId, username string) (bool, error) {
	// Double ack can be enabled. Check in engine-axe.
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType:  types.EventTypeAck,
		Output:     r.Comment,
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	if r.AckResources && alarm.Entity.Type == types.EntityTypeComponent {
		go func() {
			err = s.ackResources(context.Background(), alarm.Entity.ID, r.Comment, userId, username)
			if err != nil {
				s.logger.Err(err).Msg("cannot ack resources")
			}
		}()
	}

	return true, nil
}

func (s *store) AckRemove(ctx context.Context, id string, r Request, userId, username string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.ack": bson.M{"$ne": nil}})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType:  types.EventTypeAckremove,
		Output:     r.Comment,
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Snooze(ctx context.Context, id string, r SnoozeRequest, userId, username string) (bool, error) {
	d, err := r.Duration.To("s")
	if err != nil {
		return false, common.NewValidationError("duration", "Duration is invalid.")
	}

	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.snooze": nil})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType:  types.EventTypeSnooze,
		Output:     r.Comment,
		Duration:   types.CpsNumber(d.Value),
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Cancel(ctx context.Context, id string, r Request, userId, username string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.canceled": nil})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType:  types.EventTypeCancel,
		Output:     r.Comment,
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Uncancel(ctx context.Context, id string, r Request, userId, username string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id, "v.canceled": bson.M{"$ne": nil}})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType:  types.EventTypeUncancel,
		Output:     r.Comment,
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) AssocTicket(ctx context.Context, id string, r AssocTicketRequest, userId, username string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
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
		Output:     ticketInfo.GetStepMessage(),
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	if r.TicketResources && alarm.Entity.Type == types.EntityTypeComponent {
		go func() {
			err = s.ticketResources(context.Background(), alarm.Entity.ID, ticketInfo, userId, username)
			if err != nil {
				s.logger.Err(err).Msg("cannot ticket resources")
			}
		}()
	}

	return true, nil
}

func (s *store) Comment(ctx context.Context, id string, r CommentRequest, userId, username string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType:  types.EventTypeComment,
		Output:     r.Comment,
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
	}
	err = s.sendEvent(ctx, event)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) ChangeState(ctx context.Context, id string, r ChangeStateRequest, userId, username string) (bool, error) {
	alarm, err := s.findAlarm(ctx, bson.M{"_id": id})
	if err != nil || alarm.Alarm.ID == "" {
		return false, err
	}

	event := types.Event{
		EventType:  types.EventTypeChangestate,
		State:      types.CpsNumber(*r.State),
		Output:     r.Comment,
		Component:  alarm.Alarm.Value.Component,
		Resource:   alarm.Alarm.Value.Resource,
		SourceType: alarm.Entity.Type,
		UserID:     userId,
		Author:     username,
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
			"from":         libmongo.EntityMongoCollection,
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
	event.Connector = canopsis.ApiConnector
	event.ConnectorName = canopsis.ApiConnector
	event.Initiator = types.InitiatorUser
	event.Timestamp = datetime.NewCpsTime()
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
	userId, username string,
) error {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"v.component": component,
			"v.meta":      bson.M{"$exists": false},
			"v.resolved":  nil,
			"v.ack":       nil,
		}},
		{"$lookup": bson.M{
			"from":         libmongo.EntityMongoCollection,
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

	outputBuilder := strings.Builder{}
	if comment != "" {
		outputBuilder.WriteString(comment)
		outputBuilder.WriteString("\n")
	}

	outputBuilder.WriteString(types.OutputComponentPrefix)
	outputBuilder.WriteString(component)

	output := outputBuilder.String()
	for cursor.Next(ctx) {
		alarm := types.AlarmWithEntity{}
		err = cursor.Decode(&alarm)
		if err != nil {
			return fmt.Errorf("cannot decode alarm: %w", err)
		}

		event := types.Event{
			EventType:  types.EventTypeAck,
			Output:     output,
			Component:  alarm.Alarm.Value.Component,
			Resource:   alarm.Alarm.Value.Resource,
			SourceType: alarm.Entity.Type,
			UserID:     userId,
			Author:     username,
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
	userId, username string,
) error {
	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"v.component": component,
			"v.meta":      bson.M{"$exists": false},
			"v.resolved":  nil,
			"v.ticket":    nil,
		}},
		{"$lookup": bson.M{
			"from":         libmongo.EntityMongoCollection,
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

	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(ticketInfo.GetStepMessage())
	outputBuilder.WriteString(" ")
	outputBuilder.WriteString(types.OutputComponentPrefix)
	outputBuilder.WriteString(component)
	outputBuilder.WriteRune('.')
	output := outputBuilder.String()
	for cursor.Next(ctx) {
		alarm := types.AlarmWithEntity{}
		err = cursor.Decode(&alarm)
		if err != nil {
			return fmt.Errorf("cannot decode alarm: %w", err)
		}

		event := types.Event{
			EventType:  types.EventTypeAssocTicket,
			TicketInfo: ticketInfo,
			Output:     output,
			Component:  alarm.Alarm.Value.Component,
			Resource:   alarm.Alarm.Value.Resource,
			SourceType: alarm.Entity.Type,
			UserID:     userId,
			Author:     username,
		}
		err = s.sendEvent(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) AddBookmark(ctx context.Context, alarmID, userID string) (bool, error) {
	found := false

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		found = false

		var doc alarmResolvedField

		err := s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": alarmID},
			bson.M{"$addToSet": bson.M{"bookmarks": userID}},
			options.FindOneAndUpdate().SetProjection(bson.M{"resolved": "$v.resolved"}).SetReturnDocument(options.After),
		).Decode(&doc)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		if doc.ID != "" && doc.Resolved == nil {
			found = true

			return nil
		}

		resolvedRes, err := s.resolvedDbCollection.UpdateOne(
			ctx,
			bson.M{"_id": alarmID},
			bson.M{"$addToSet": bson.M{"bookmarks": userID}},
		)
		if err != nil {
			return err
		}

		found = resolvedRes.MatchedCount != 0

		return nil
	})

	return found, err
}

func (s *store) RemoveBookmark(ctx context.Context, alarmID, userID string) (bool, error) {
	found := false

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		found = false

		var doc alarmResolvedField

		err := s.dbCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": alarmID},
			bson.M{"$pull": bson.M{"bookmarks": userID}},
			options.FindOneAndUpdate().SetProjection(bson.M{"resolved": "$v.resolved"}).SetReturnDocument(options.After),
		).Decode(&doc)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		if doc.ID != "" && doc.Resolved == nil {
			found = true

			return nil
		}

		resolvedRes, err := s.resolvedDbCollection.UpdateOne(
			ctx,
			bson.M{"_id": alarmID},
			bson.M{"$pull": bson.M{"bookmarks": userID}},
		)
		if err != nil {
			return err
		}

		found = resolvedRes.MatchedCount != 0

		return nil
	})

	return found, err
}
