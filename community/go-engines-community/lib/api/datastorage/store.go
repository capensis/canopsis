package datastorage

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		if errors.Is(err, mongodriver.ErrNoDocuments) {
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

	err = s.updateInstructionRetentionPolicy(ctx, data, pgPool)
	if err != nil {
		return err
	}

	err = s.updateMetricsRetentionPolicy(ctx, data, pgPool)
	if err != nil {
		return err
	}

	return s.updatePerfDataRetentionPolicy(ctx, data, pgPool)
}

func (s *store) updateInstructionRetentionPolicy(ctx context.Context, data datastorage.DataStorage, pgPool postgres.Pool) error {
	err := s.deleteRetentionPolicy(ctx, pgPool, metrics.InstructionExecutionHourly)
	if err != nil {
		return err
	}
	deleteStatsAfter := data.Config.Remediation.DeleteStatsAfter
	if datetime.IsDurationEnabledAndValid(deleteStatsAfter) {
		err = s.addRetentionPolicy(ctx, pgPool, metrics.InstructionExecutionHourly, deleteStatsAfter.String())
		if err != nil {
			return err
		}
	}

	err = s.deleteRetentionPolicy(ctx, pgPool, metrics.InstructionExecutionByModifiedOn)
	if err != nil {
		return err
	}
	deleteModStatsAfter := data.Config.Remediation.DeleteModStatsAfter
	if datetime.IsDurationEnabledAndValid(deleteModStatsAfter) {
		err = s.addRetentionPolicy(ctx, pgPool, metrics.InstructionExecutionByModifiedOn, deleteStatsAfter.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) updateMetricsRetentionPolicy(ctx context.Context, data datastorage.DataStorage, pgPool postgres.Pool) error {
	tables := []string{
		metrics.TotalAlarmNumber,
		metrics.NonDisplayedAlarmNumber,
		metrics.PbhAlarmNumber,
		metrics.InstructionAlarmNumber,
		metrics.TicketAlarmNumber,
		metrics.CorrelationAlarmNumber,
		metrics.AckAlarmNumber,
		metrics.CancelAckAlarmNumber,
		metrics.AckDuration,
		metrics.ResolveDuration,
		metrics.UserActivity,
		metrics.UserSessions,
		metrics.TicketNumber,
		metrics.SliDuration,
		metrics.ManualInstructionAssignedAlarms,
		metrics.ManualInstructionExecutedAlarms,
		metrics.InstructionAssignedInstructions,
		metrics.InstructionExecutedInstructions,
		metrics.NotAckedInHourAlarms,
		metrics.NotAckedInFourHoursAlarms,
		metrics.NotAckedInDayAlarms,
	}

	deleteAfter := data.Config.Metrics.DeleteAfter
	interval := ""
	if datetime.IsDurationEnabledAndValid(deleteAfter) {
		interval = deleteAfter.String()
	}
	for _, table := range tables {
		err := s.deleteRetentionPolicy(ctx, pgPool, table)
		if err != nil {
			return err
		}
		if interval != "" {
			err = s.addRetentionPolicy(ctx, pgPool, table, interval)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *store) updatePerfDataRetentionPolicy(ctx context.Context, data datastorage.DataStorage, pgPool postgres.Pool) error {
	err := s.deleteRetentionPolicy(ctx, pgPool, metrics.PerfData)
	if err != nil {
		return err
	}
	deleteAfter := data.Config.PerfDataMetrics.DeleteAfter
	if datetime.IsDurationEnabledAndValid(deleteAfter) {
		err = s.addRetentionPolicy(ctx, pgPool, metrics.PerfData, deleteAfter.String())
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
