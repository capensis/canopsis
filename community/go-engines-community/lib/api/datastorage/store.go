package datastorage

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const InstructionExecutionHourly = "instruction_execution_hourly"
const InstructionExecutionByModifiedOn = "instruction_execution_by_modified_on"

type Store interface {
	Get(ctx context.Context) (datastorage.DataStorage, error)
	Update(context.Context, datastorage.Config) (datastorage.DataStorage, error)
}

func NewStore(client mongo.DbClient, pgPoolProvider postgres.PoolProvider, logger zerolog.Logger) Store {
	return &store{
		collection:     client.Collection(mongo.ConfigurationMongoCollection),
		pgPoolProvider: pgPoolProvider,
		logger:         logger,
	}
}

type store struct {
	collection     mongo.DbCollection
	pgPoolProvider postgres.PoolProvider
	logger         zerolog.Logger
}

func (s *store) Get(ctx context.Context) (datastorage.DataStorage, error) {
	data := datastorage.DataStorage{}
	err := s.collection.FindOne(ctx, bson.M{"_id": datastorage.ID}).Decode(&data)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return data, nil
		}

		return data, err
	}

	return data, nil
}

func (s *store) Update(ctx context.Context, conf datastorage.Config) (datastorage.DataStorage, error) {
	data := datastorage.DataStorage{}
	err := s.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": datastorage.ID},
		bson.M{"$set": bson.M{
			"config": conf,
		}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&data)
	if err != nil {
		return data, err
	}

	err = s.updateRetentionPolicy(ctx, data)
	if err != nil {
		s.logger.Err(err).Msg("cannot update retention policy for remediation metrics")
	}

	return data, nil
}

func (s *store) updateRetentionPolicy(ctx context.Context, data datastorage.DataStorage) error {
	if s.pgPoolProvider == nil {
		return nil
	}
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return err
	}
	err = s.deleteRetentionPolicy(ctx, pgPool, InstructionExecutionHourly)
	if err != nil {
		return err
	}
	deleteStatsAfter := data.Config.Remediation.DeleteStatsAfter
	if deleteStatsAfter != nil && deleteStatsAfter.Enabled != nil && *deleteStatsAfter.Enabled && deleteStatsAfter.Value > 0 {
		err = s.addRetentionPolicy(ctx, pgPool, InstructionExecutionHourly, deleteStatsAfter.String())
		if err != nil {
			return err
		}
	}

	err = s.deleteRetentionPolicy(ctx, pgPool, InstructionExecutionByModifiedOn)
	if err != nil {
		return err
	}
	deleteModStatsAfter := data.Config.Remediation.DeleteModStatsAfter
	if deleteModStatsAfter != nil && deleteModStatsAfter.Enabled != nil && *deleteModStatsAfter.Enabled && deleteModStatsAfter.Value > 0 {
		err = s.addRetentionPolicy(ctx, pgPool, InstructionExecutionByModifiedOn, deleteStatsAfter.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) deleteRetentionPolicy(ctx context.Context, pgPool postgres.Pool, table string) error {
	_, err := pgPool.Exec(ctx, "SELECT remove_retention_policy('"+table+"', if_exists => true)")
	return err
}

func (s *store) addRetentionPolicy(ctx context.Context, pgPool postgres.Pool, table, interval string) error {
	_, err := pgPool.Exec(ctx, "SELECT add_retention_policy('"+table+"', INTERVAL '"+interval+"')")
	return err
}
