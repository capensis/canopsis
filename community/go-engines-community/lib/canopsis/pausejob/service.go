package pausejob

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Service interface {
	Create(ctx context.Context, entity types.Entity, pause types.LatestPause, parameters rpc.AxeParameters) error
	CreateForAdded(ctx context.Context, parent string, entities []types.Entity, pause *types.LatestPause) error
	Cancel(ctx context.Context, entity types.Entity, pause types.Pause) error
	IsCurrent(ctx context.Context, parameters rpc.AxeParameters) (bool, error)
	Complete(ctx context.Context, jobID, entityID string) error
	IsPending(ctx context.Context, jobID, entityID string) (bool, error)
}

func NewService(client mongo.DbClient, eventGenerator libevent.Generator, logger zerolog.Logger) Service {
	return &service{
		collection:       client.Collection(mongo.PauseJobCollection),
		entityCollection: client.Collection(mongo.EntityMongoCollection),
		eventGenerator:   eventGenerator,
		logger:           logger,
	}
}

type service struct {
	collection       mongo.DbCollection
	entityCollection mongo.DbCollection
	eventGenerator   libevent.Generator
	logger           zerolog.Logger
}

// Create must run in the transaction which accepts the parent's pause.
func (s *service) Create(ctx context.Context, entity types.Entity, pause types.LatestPause,
	parameters rpc.AxeParameters,
) error {
	if !pause.Inherited || (entity.Type != types.EntityTypeService && entity.Type != types.EntityTypeComponent) {
		return nil
	}

	inherited, err := entity.InheritedPause()
	if err != nil {
		s.logger.Err(err).Str("entity", entity.ID).Msg("cannot inherit pause")
		return nil
	}

	if inherited == nil {
		return nil
	}

	_, err = s.collection.InsertOne(ctx, Job{
		ID:         utils.NewID(),
		Parent:     entity.ID,
		ParentType: entity.Type,
		Event: types.Event{
			EventType:       types.EventTypePause,
			Timestamp:       inherited.Timestamp,
			Output:          inherited.Message,
			Author:          parameters.Author,
			UserID:          parameters.User,
			Role:            parameters.Role,
			Initiator:       parameters.Initiator,
			Inherited:       true,
			PauseSource:     inherited.Source,
			PauseSourceType: inherited.SourceType,
			PauseType:       inherited.Type,
			PausePath:       slices.Clone(inherited.Path),
		},
		Pending: []string{},
		RetryAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("cannot create pause propagation job: %w", err)
	}

	return nil
}

// CreateForAdded persists the explicit recipients of a service recompute in the
// same transaction as the membership changes.
func (s *service) CreateForAdded(ctx context.Context, parent string,
	entities []types.Entity, pause *types.LatestPause,
) error {
	if len(entities) == 0 || pause == nil {
		return nil
	}

	groups := make(map[string]int)
	jobs := make([]Job, 0)
	for _, entity := range entities {
		event, err := s.eventGenerator.Generate(entity)
		if err != nil {
			return fmt.Errorf("cannot generate pause event of %s: %w", entity.ID, err)
		}

		key, err := bson.Marshal(struct {
			Connector     string
			ConnectorName string
		}{event.Connector, event.ConnectorName})
		if err != nil {
			return fmt.Errorf("cannot group pause propagation jobs: %w", err)
		}

		if i, ok := groups[string(key)]; ok {
			jobs[i].Pending = append(jobs[i].Pending, entity.ID)
			continue
		}

		groups[string(key)] = len(jobs)
		jobs = append(jobs, Job{
			ID: utils.NewID(), Parent: parent, ParentType: types.EntityTypeService,
			Initialized: true,
			Pending:     []string{entity.ID},
			RetryAt:     time.Now(),
			Event: types.Event{
				EventType:       types.EventTypePause,
				Timestamp:       pause.Timestamp,
				Connector:       event.Connector,
				ConnectorName:   event.ConnectorName,
				Output:          pause.Message,
				Author:          pause.Author,
				UserID:          pause.UserID,
				Initiator:       pause.Initiator,
				Inherited:       true,
				PauseType:       pause.Type,
				PauseSource:     pause.Source,
				PauseSourceType: pause.SourceType,
				PausePath:       slices.Clone(pause.Path),
			},
		})
	}

	for _, job := range jobs {
		if _, err := s.collection.InsertOne(ctx, job); err != nil {
			return fmt.Errorf("cannot create added-member pause propagation job: %w", err)
		}
	}

	return nil
}

// Cancel runs with the state transition which ends or replaces a pause. Removing
// the receipt also conflicts with an in-flight recipient's state transaction.
func (s *service) Cancel(ctx context.Context, entity types.Entity, pause types.Pause) error {
	source := pause.Source
	if source == "" {
		source = entity.ID
	}

	filter := bson.M{"event.pause_source": source, "event.timestamp": pause.Timestamp}
	if source != entity.ID {
		filter["$or"] = bson.A{bson.M{"parent": entity.ID}, bson.M{"event.pause_path": entity.ID}}
	}

	if _, err := s.collection.DeleteMany(ctx, filter); err != nil {
		return fmt.Errorf("cannot cancel pause propagation of %s: %w", entity.ID, err)
	}

	return nil
}

// IsCurrent checks the source and immediate parent before a queued pause is
// accepted. Membership alone does not prove that its originating pause is active.
func (s *service) IsCurrent(ctx context.Context, parameters rpc.AxeParameters) (bool, error) {
	parent := parameters.PauseSource
	if len(parameters.PausePath) > 0 {
		parent = parameters.PausePath[len(parameters.PausePath)-1]
	}

	for _, id := range []string{parameters.PauseSource, parent} {
		var entity types.Entity
		err := s.entityCollection.FindOne(ctx, bson.M{
			"_id": id, "enabled": true, "pause.inherited": true,
			"latest_pause.source":    parameters.PauseSource,
			"latest_pause.timestamp": parameters.Timestamp,
			"latest_pause.end":       bson.M{"$exists": false},
		}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return false, nil
			}

			return false, fmt.Errorf("cannot find pause propagation parent %s: %w", id, err)
		}

		if parent == parameters.PauseSource {
			break
		}
	}

	return true, nil
}

// Complete is part of the recipient's state transaction. A deleted job is harmless:
// queued duplicates are ignored when their pending receipt no longer exists.
func (s *service) Complete(ctx context.Context, jobID, entityID string) error {
	if jobID == "" {
		return nil
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": jobID}, bson.M{"$pull": bson.M{"pending": entityID}})
	if err != nil {
		return fmt.Errorf("cannot complete pause propagation for %s: %w", entityID, err)
	}

	return nil
}

func (s *service) IsPending(ctx context.Context, jobID, entityID string) (bool, error) {
	var receipt struct {
		ID string `bson:"_id"`
	}

	err := s.collection.FindOne(ctx, bson.M{
		"_id": jobID, "pending": entityID,
	}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&receipt)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
