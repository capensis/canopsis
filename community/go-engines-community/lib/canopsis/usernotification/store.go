package usernotification

import (
	"context"
	"errors"
	"fmt"
	"sync"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Flush(ctx context.Context) error
	AddForInstructionApprove(time datetime.CpsTime, id, name, author, comment, userID, roleID string)
	AddForInstructionRate(time datetime.CpsTime, id, name, userID string)
	AddForEventFilterFailure(time datetime.CpsTime, id, description string, roleIDs []string)

	DeleteForInstructionApprove(ctx context.Context, id string) error
	DeleteForInstructionRate(ctx context.Context, id, userID string) error
	DeleteForEventFilterFailure(ctx context.Context, id string) error
}

type store struct {
	collection               mongo.DbCollection
	publishCh                libamqp.Publisher
	exchangeName, routingKey string
	msgContentType           string
	encoder                  encoding.Encoder

	dataMX      sync.Mutex
	writeModels []mongodriver.WriteModel
	userIDs     map[string]struct{}
	roleIDs     map[string]struct{}
}

func NewStore(dbClient mongo.DbClient, publishCh libamqp.Publisher, encoder encoding.Encoder, exchangeName, routingKey, msgContentType string) Store {
	return &store{
		collection:     dbClient.Collection(mongo.UserNotificationCollection),
		publishCh:      publishCh,
		encoder:        encoder,
		exchangeName:   exchangeName,
		routingKey:     routingKey,
		msgContentType: msgContentType,
		writeModels:    make([]mongodriver.WriteModel, 0),
		userIDs:        make(map[string]struct{}),
		roleIDs:        make(map[string]struct{}),
	}
}

func (s *store) Flush(ctx context.Context) error {
	writeModels, mapUserIDs, mapRoleIDs := s.flushData()
	modelsLen := len(writeModels)
	if modelsLen == 0 {
		return nil
	}

	bulkSize := canopsis.DefaultBulkSize
	updated := false
	from := 0
	for to := min(bulkSize, modelsLen); from < modelsLen; to = min(to+bulkSize, modelsLen) {
		bulkRes, err := s.collection.BulkWrite(ctx, writeModels[from:to])
		if err != nil {
			return fmt.Errorf("cannot bulk write user notifications: %w", err)
		}

		from = to
		if bulkRes.InsertedCount > 0 || bulkRes.ModifiedCount > 0 || bulkRes.UpsertedCount > 0 || bulkRes.DeletedCount > 0 {
			updated = true
		}
	}

	if !updated {
		return nil
	}

	userIDs := make([]string, len(mapUserIDs))
	i := 0
	for id := range mapUserIDs {
		userIDs[i] = id
		i++
	}

	roleIDs := make([]string, len(mapRoleIDs))
	i = 0
	for id := range mapRoleIDs {
		roleIDs[i] = id
		i++
	}

	return s.sendEvent(ctx, userIDs, roleIDs)
}

func (s *store) AddForInstructionApprove(time datetime.CpsTime, id, name, author, comment, userID, roleID string) {
	s.dataMX.Lock()
	defer s.dataMX.Unlock()
	t := TypeInstructionApprove
	if userID != "" {
		s.userIDs[userID] = struct{}{}
	}

	roleIDs := make([]string, 0)
	if roleID != "" {
		roleIDs = append(roleIDs, roleID)
		s.roleIDs[roleID] = struct{}{}
	}

	s.writeModels = append(s.writeModels, mongodriver.NewUpdateOneModel().
		SetFilter(bson.M{
			"type":     t,
			"rule._id": id,
		}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"time":    time,
				"user":    userID,
				"roles":   roleIDs,
				"comment": comment,
				"rule": bson.M{
					"_id":  id,
					"name": name,
				},
				"author": author,
			},
			"$setOnInsert": bson.M{
				"_id":  utils.NewID(),
				"type": t,
			},
		}).
		SetUpsert(true),
	)
}

func (s *store) AddForInstructionRate(time datetime.CpsTime, id, name, userID string) {
	s.dataMX.Lock()
	defer s.dataMX.Unlock()

	t := TypeInstructionRate
	s.userIDs[userID] = struct{}{}
	s.writeModels = append(s.writeModels, mongodriver.NewUpdateOneModel().
		SetFilter(bson.M{
			"type":     t,
			"rule._id": id,
			"user":     userID,
		}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"time": time,
				"rule": bson.M{
					"_id":  id,
					"name": name,
				},
			},
			"$setOnInsert": bson.M{
				"_id":   utils.NewID(),
				"type":  t,
				"user":  userID,
				"roles": []string{},
			},
		}).
		SetUpsert(true),
	)
}

func (s *store) AddForEventFilterFailure(time datetime.CpsTime, id, description string, roleIDs []string) {
	s.dataMX.Lock()
	defer s.dataMX.Unlock()

	for _, v := range roleIDs {
		s.roleIDs[v] = struct{}{}
	}

	t := TypeEventFilterFailure
	s.writeModels = append(s.writeModels, mongodriver.NewUpdateOneModel().
		SetFilter(bson.M{
			"type":     t,
			"rule._id": id,
		}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"time":  time,
				"roles": roleIDs,
				"rule": bson.M{
					"_id":  id,
					"name": description,
				},
			},
			"$setOnInsert": bson.M{
				"_id":  utils.NewID(),
				"type": t,
				"user": "",
			},
		}).
		SetUpsert(true),
	)
}

func (s *store) DeleteForInstructionApprove(ctx context.Context, id string) error {
	return s.delete(ctx, bson.M{"type": TypeInstructionApprove, "rule._id": id})
}

func (s *store) DeleteForInstructionRate(ctx context.Context, id, userID string) error {
	if userID != "" {
		return s.delete(ctx, bson.M{
			"type":     TypeInstructionRate,
			"user":     userID,
			"rule._id": id,
		})
	}

	cursor, err := s.collection.Find(ctx, bson.M{"type": TypeInstructionRate, "rule._id": id})
	if err != nil {
		return fmt.Errorf("cannot find notifications: %w", err)
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, canopsis.DefaultBulkSize)
	userIDs := make([]string, 0, canopsis.DefaultBulkSize)
	for cursor.Next(ctx) {
		n := Notification{}
		err = cursor.Decode(&n)
		if err != nil {
			return fmt.Errorf("cannot decode notification: %w", err)
		}

		ids = append(ids, n.ID)
		userIDs = append(userIDs, n.User)
		if len(ids) == canopsis.DefaultBulkSize {
			_, err = s.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				return fmt.Errorf("cannot delete notificaions: %w", err)
			}

			ids = ids[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return fmt.Errorf("cannot fetch notifications: %w", err)
	}

	if len(ids) > 0 {
		_, err = s.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return fmt.Errorf("cannot delete notificaions: %w", err)
		}
	}

	return s.sendEvent(ctx, userIDs, nil)
}

func (s *store) DeleteForEventFilterFailure(ctx context.Context, id string) error {
	return s.delete(ctx, bson.M{"type": TypeEventFilterFailure, "rule._id": id})
}

func (s *store) flushData() ([]mongodriver.WriteModel, map[string]struct{}, map[string]struct{}) {
	s.dataMX.Lock()
	defer s.dataMX.Unlock()

	writeModels := s.writeModels
	userIDs := s.userIDs
	roleIDs := s.roleIDs

	s.writeModels = make([]mongodriver.WriteModel, 0, len(writeModels))
	s.userIDs = make(map[string]struct{}, len(userIDs))
	s.roleIDs = make(map[string]struct{}, len(roleIDs))

	return writeModels, userIDs, roleIDs
}

func (s *store) sendEvent(ctx context.Context, userIDs, roleIDs []string) error {
	b, err := s.encoder.Encode(rpc.ApiNotificationEvent{
		Users: userIDs,
		Roles: roleIDs,
	})
	if err != nil {
		return fmt.Errorf("cannot encode notification event: %w", err)
	}

	err = s.publishCh.PublishWithContext(
		ctx,
		s.exchangeName,
		s.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: s.msgContentType,
			Body:        b,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot send notification event: %w", err)
	}

	return nil
}

func (s *store) delete(ctx context.Context, f bson.M) error {
	n := Notification{}
	err := s.collection.FindOneAndDelete(ctx, f).Decode(&n)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}

		return fmt.Errorf("cannot delete notification: %w", err)
	}

	var userIDs []string
	if n.User != "" {
		userIDs = append(userIDs, n.User)
	}

	return s.sendEvent(ctx, userIDs, n.Roles)
}
