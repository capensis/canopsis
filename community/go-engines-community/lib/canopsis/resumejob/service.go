package resumejob

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
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
	Create(ctx context.Context, entity types.Entity, ended types.LatestPause, parameters rpc.AxeParameters) error
	CreateForRemoved(ctx context.Context, parent string, entities []types.Entity, end datetime.CpsTime) error
	Complete(ctx context.Context, jobID, entityID string) error
	IsPending(ctx context.Context, jobID, entityID string) (bool, error)
}

func NewService(client mongo.DbClient, eventGenerator libevent.Generator, logger zerolog.Logger) Service {
	return &service{
		collection:     client.Collection(mongo.ResumeJobCollection),
		eventGenerator: eventGenerator,
		logger:         logger,
	}
}

type service struct {
	collection     mongo.DbCollection
	eventGenerator libevent.Generator
	logger         zerolog.Logger
}

// Create must run in the transaction which accepts the parent's resume.
func (s *service) Create(ctx context.Context, entity types.Entity, ended types.LatestPause,
	parameters rpc.AxeParameters,
) error {
	if !ended.Inherited {
		return nil
	}

	switch entity.Type {
	case types.EntityTypeService, types.EntityTypeComponent:
	case "": // A deleted intermediary can still have dependents retaining its pause path.
	default:
		return nil
	}

	path := slices.Clone(ended.Path)
	if ended.Source != entity.ID {
		if len(path) >= types.MaxPausePathLength || slices.Contains(path, entity.ID) {
			s.logger.Error().Str("entity", entity.ID).Int("path_length", len(path)).
				Msg("the pause path is too long or cyclic, the resume is not propagated to the dependencies")

			return nil
		}

		path = append(path, entity.ID)
	}

	output := parameters.Output
	if ended.Source == entity.ID {
		output = ended.GetInheritedMessage(entity.Type, entity.Name)
	}

	_, err := s.collection.InsertOne(ctx, Job{
		ID:         utils.NewID(),
		Parent:     entity.ID,
		ParentType: entity.Type,
		Event: types.Event{
			EventType:       types.EventTypeResume,
			Timestamp:       parameters.Timestamp,
			Output:          output,
			Author:          parameters.Author,
			UserID:          parameters.User,
			Role:            parameters.Role,
			Initiator:       parameters.Initiator,
			Inherited:       true,
			PauseSource:     ended.Source,
			PauseSourceType: ended.SourceType,
			PauseStart:      new(ended.Timestamp),
			PausePath:       path,
		},
		Pending: []string{},
		RetryAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("cannot create resume propagation job: %w", err)
	}

	return nil
}

// CreateForRemoved persists the explicit recipients discovered by a service
// recompute in the same transaction as the membership changes. Different pause
// routes remain separate; one document holds all recipients of the same route.
func (s *service) CreateForRemoved(ctx context.Context, parent string,
	entities []types.Entity, end datetime.CpsTime,
) error {
	groups := make(map[string]int)
	jobs := make([]Job, 0)
	for _, entity := range entities {
		if entity.Pause == nil {
			return fmt.Errorf("cannot create resume propagation job for %s without its pause", entity.ID)
		}

		event, err := s.eventGenerator.Generate(entity)
		if err != nil {
			return fmt.Errorf("cannot generate resume event of %s: %w", entity.ID, err)
		}

		pause := entity.Pause
		key, err := bson.Marshal(struct {
			Source        string
			SourceType    string
			Start         datetime.CpsTime
			Path          []string
			Output        string
			Connector     string
			ConnectorName string
			Author        string
		}{pause.Source, pause.SourceType, pause.Timestamp, pause.Path, pause.Message,
			event.Connector, event.ConnectorName, event.Author})
		if err != nil {
			return fmt.Errorf("cannot group resume propagation jobs: %w", err)
		}

		if i, ok := groups[string(key)]; ok {
			jobs[i].Pending = append(jobs[i].Pending, entity.ID)
			continue
		}

		resumeTimestamp := end
		if resumeTimestamp.Before(pause.Timestamp) {
			// Membership cleanup must also end pauses with a future timestamp.
			resumeTimestamp = pause.Timestamp
		}

		groups[string(key)] = len(jobs)
		jobs = append(jobs, Job{
			ID: utils.NewID(), Parent: parent, ParentType: types.EntityTypeService,
			Initialized: true,
			Pending:     []string{entity.ID},
			RetryAt:     time.Now(),
			Event: types.Event{
				EventType:       types.EventTypeResume,
				Timestamp:       resumeTimestamp,
				Author:          event.Author,
				Connector:       event.Connector,
				ConnectorName:   event.ConnectorName,
				Output:          pause.Message,
				Initiator:       types.InitiatorSystem,
				Inherited:       true,
				PauseSource:     pause.Source,
				PauseSourceType: pause.SourceType,
				PauseStart:      new(pause.Timestamp),
				PausePath:       slices.Clone(pause.Path),
			},
		})
	}

	for _, job := range jobs {
		if _, err := s.collection.InsertOne(ctx, job); err != nil {
			return fmt.Errorf("cannot create removed-member resume propagation job: %w", err)
		}
	}

	return nil
}

// Complete is part of the recipient's state transaction. A deleted job is harmless:
// queued duplicates are ignored when their pending receipt no longer exists.
func (s *service) Complete(ctx context.Context, jobID, entityID string) error {
	if jobID == "" {
		return nil
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": jobID}, bson.M{"$pull": bson.M{"pending": entityID}})
	if err != nil {
		return fmt.Errorf("cannot complete resume propagation for %s: %w", entityID, err)
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
