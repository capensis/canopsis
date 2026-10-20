package entity

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

type Store interface {
	Find(ctx context.Context, r PaginatedListRequest, userID string) (pagination.ListResponse[Response], error)
	Search(ctx context.Context, r PaginatedSearchRequest, userID string) (pagination.ListResponse[Response], error)
	Pause(ctx context.Context, r PauseRequest, userID, username string) error
	Resume(ctx context.Context, r ResumeRequest, userID, username string) error
}

type store struct {
	dbClient       mongo.DbClient
	dbCollection   mongo.DbCollection
	transformer    patternfields.Transformer
	generator      link.Generator
	userProvider   liblink.UserProvider
	amqpPublisher  amqp.Publisher
	eventGenerator event.Generator
	encoder        encoding.Encoder
}

func NewStore(
	dbClient mongo.DbClient,
	transformer patternfields.Transformer,
	generator link.Generator,
	userProvider liblink.UserProvider,
	amqpPublisher amqp.Publisher,
	eventGenerator event.Generator,
	encoder encoding.Encoder,
) Store {
	return &store{
		dbClient:       dbClient,
		dbCollection:   dbClient.Collection(mongo.EntityMongoCollection),
		transformer:    transformer,
		generator:      generator,
		userProvider:   userProvider,
		amqpPublisher:  amqpPublisher,
		eventGenerator: eventGenerator,
		encoder:        encoder,
	}
}

func (s *store) Find(ctx context.Context, r PaginatedListRequest, userID string) (res pagination.ListResponse[Response], err error) {
	beforeLimit, afterLimit, sort, err := ListPipeline(r.ListRequest)
	if err != nil {
		return res, err
	}

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateAll[Response](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	user, err := s.userProvider.User(ctx, userID)
	if err != nil {
		return res, err
	}

	for i, item := range ar {
		ar[i].SetDefaults()
		ar[i].Links = s.generator.GenerateForEntity(ctx, item.Raw, user)
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) Search(ctx context.Context, r PaginatedSearchRequest, userID string) (res pagination.ListResponse[Response], err error) {
	if len(r.EntityPattern) > 0 {
		patternAliases := patternfields.GetAliases(r.EntityPattern)
		if len(patternAliases) != 0 {
			aliases, err := s.transformer.FetchAliases(ctx, patternAliases)
			if err != nil {
				return res, err
			}

			var valErrs validator.ValidationErrors
			r.EntityPattern, _, valErrs = s.transformer.ApplyAliases(r.EntityPattern, aliases, "EntityPattern")
			if len(valErrs) > 0 {
				return res, validation.NewError(valErrs, r)
			}
		}
	}

	beforeLimit, afterLimit, sort, err := SearchPipeline(r.SearchRequest)
	if err != nil {
		return res, err
	}

	pipeline := pagination.CreateAggregationPipeline(
		r.Query,
		beforeLimit,
		sort,
		afterLimit,
	)
	ar, err := db.AggregateAll[Response](ctx, s.dbCollection, pipeline)
	if err != nil {
		return res, err
	}

	user, err := s.userProvider.User(ctx, userID)
	if err != nil {
		return res, err
	}

	for i, item := range ar {
		ar[i].SetDefaults()
		ar[i].Links = s.generator.GenerateForEntity(ctx, item.Raw, user)
	}

	res = pagination.NewResponse[Response](r.Query, ar)

	return res, nil
}

func (s *store) Pause(ctx context.Context, r PauseRequest, userID, username string) error {
	entity, err := s.findEntity(ctx, bson.M{"_id": r.ID, "pause": nil})
	if err != nil {
		return err
	}

	e, err := s.genEvent(types.EventTypePause, entity, r.Comment, username, userID)
	if err != nil {
		return err
	}

	e.Inherited = r.Inherited
	e.PauseType = r.Type

	return s.sendEvent(ctx, e)
}

func (s *store) Resume(ctx context.Context, r ResumeRequest, userID, username string) error {
	entity, err := s.findEntity(ctx, bson.M{"_id": r.ID, "pause": bson.M{"$ne": nil}})
	if err != nil {
		return err
	}

	e, err := s.genEvent(types.EventTypeResume, entity, r.Comment, username, userID)
	if err != nil {
		return err
	}

	return s.sendEvent(ctx, e)
}

func (s *store) findEntity(ctx context.Context, match bson.M) (types.Entity, error) {
	if match == nil {
		match = bson.M{}
	}

	match["v.resolved"] = nil
	entity := types.Entity{}
	err := s.dbCollection.FindOne(ctx, match).Decode(&entity)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return entity, httperror.ErrNotFound
		}

		return entity, err
	}

	return entity, nil
}

func (s *store) genEvent(
	eventType string,
	entity types.Entity,
	output string,
	username, userID string,
) (types.Event, error) {
	e, err := s.eventGenerator.Generate(entity)
	if err != nil {
		return e, fmt.Errorf("cannot generate event: %w", err)
	}

	e.EventType = eventType
	e.Timestamp = datetime.NewCpsTime()
	e.Output = output
	e.Author = username
	e.UserID = userID
	e.Initiator = types.InitiatorUser

	return e, nil
}

func (s *store) sendEvent(ctx context.Context, event types.Event) error {
	body, err := s.encoder.Encode(event)
	if err != nil {
		return err
	}

	return s.amqpPublisher.PublishWithContext(
		ctx,
		canopsis.DefaultExchangeName,
		canopsis.FIFOQueueName,
		false,
		false,
		amqp091.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
}
