package mongo

//go:generate mockgen -destination=../../mocks/lib/mongo/commands_register.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo CommandsRegister

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandsRegister interface {
	Clear()
	RegisterInsert(entity *types.Entity)
	RegisterUpdate(id string, set bson.M)
	Commit(ctx context.Context) error
}

type commandsRegister struct {
	dbCollection DbCollection

	models []mongo.WriteModel
}

func NewCommandsRegister(collection DbCollection) CommandsRegister {
	return &commandsRegister{
		dbCollection: collection,
		models:       nil,
	}
}

func (s *commandsRegister) Clear() {
	// todo: shrink capacity dynamically
	s.models = s.models[:0]
}

func (s *commandsRegister) RegisterInsert(entity *types.Entity) {
	s.models = append(s.models, mongo.NewInsertOneModel().SetDocument(entity))
}

func (s *commandsRegister) RegisterUpdate(id string, set bson.M) {
	s.models = append(s.models, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": id}).SetUpdate(bson.M{"$set": set}))
}

func (s *commandsRegister) Commit(ctx context.Context) error {
	// todo: check bulk size in bytes

	if len(s.models) == 0 {
		return nil
	}

	defer s.Clear()

	modelsLen := len(s.models)
	from := 0

	for to := canopsis.DefaultBulkSize; to <= modelsLen; to += canopsis.DefaultBulkSize {
		_, err := s.dbCollection.BulkWrite(ctx, s.models[from:to])
		if err != nil {
			return err
		}

		from = to
	}

	if from < modelsLen {
		_, err := s.dbCollection.BulkWrite(ctx, s.models[from:modelsLen])
		if err != nil {
			return err
		}
	}

	return nil
}
