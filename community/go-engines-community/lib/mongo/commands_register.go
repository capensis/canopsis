package mongo

//go:generate mockgen -destination=../../mocks/lib/mongo/commands_register.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo CommandsRegister

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandsRegister interface {
	Clear()
	RegisterInsert(i any)
	RegisterUpdate(id string, set bson.M)
	Commit(ctx context.Context) error
}

type commandsRegister struct {
	dbCollection DbCollection
	bulkSize     int

	models []mongo.WriteModel
}

func NewCommandsRegister(collection DbCollection, bulkSize int) CommandsRegister {
	return &commandsRegister{
		dbCollection: collection,
		bulkSize:     bulkSize,
		models:       nil,
	}
}

func (s *commandsRegister) Clear() {
	// todo: shrink capacity dynamically
	s.models = s.models[:0]
}

func (s *commandsRegister) RegisterInsert(i any) {
	s.models = append(s.models, mongo.NewInsertOneModel().SetDocument(i))
}

func (s *commandsRegister) RegisterUpdate(id string, set bson.M) {
	s.models = append(s.models, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": id}).SetUpdate(bson.M{"$set": set}))
}

func (s *commandsRegister) Commit(ctx context.Context) error {
	// todo: check bulk size in bytes

	modelsLen := len(s.models)
	if modelsLen == 0 {
		return nil
	}

	defer s.Clear()

	from := 0

	for to := s.bulkSize; to <= modelsLen; to += s.bulkSize {
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
