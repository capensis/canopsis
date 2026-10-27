package datastorage

import (
	"context"

	storage "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	timescaleJobsTable     = "timescaledb_information.jobs"
	timescaleJobStatsTable = "timescaledb_information.job_stats"
	healthCheckHistoryKey  = "health_check"
	healthCheckSchema      = "public"
)

const getHealthCheckHistoryQuery = `
SELECT j.job_id, j.hypertable_schema, j.hypertable_name,
       CASE WHEN isfinite(s.last_run_started_at)
            THEN floor(extract(epoch FROM s.last_run_started_at))::bigint
       END AS last_execution_unix
FROM ` + timescaleJobsTable + ` AS j
LEFT JOIN ` + timescaleJobStatsTable + ` AS s ON s.job_id = j.job_id
WHERE j.proc_name = 'policy_retention'
  AND j.hypertable_schema = $1
  AND j.hypertable_name = $2`

const healthCheckRetentionPolicyQuery = `
SELECT EXISTS (
  SELECT 1 FROM ` + timescaleJobsTable + `
  WHERE proc_name = 'policy_retention' AND hypertable_schema = $1 AND hypertable_name = $2
    AND scheduled
    AND (config ->> 'drop_after')::interval = $3::interval
)`

const healthCheckRetentionPolicyExistsQuery = `
SELECT EXISTS (
  SELECT 1 FROM ` + timescaleJobsTable + `
  WHERE proc_name = 'policy_retention' AND hypertable_schema = $1 AND hypertable_name = $2
)`

func (s *store) getHealthCheckHistory(ctx context.Context) *storage.HistoryWithCount {
	if s.pgPoolProvider == nil {
		return nil
	}
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		s.logger.Err(err).Msg("cannot get PostgreSQL pool for health-check retention history")
		return nil
	}
	rows, err := pgPool.Query(ctx, getHealthCheckHistoryQuery, healthCheckSchema, metrics.MessageRateHourly)
	if err != nil {
		s.logger.Err(err).Msg("cannot query PostgreSQL health-check retention history")
		return nil
	}
	defer rows.Close()

	var latest *storage.HistoryWithCount
	for rows.Next() {
		var jobID int64
		var schema, table string
		var unixTimestamp *int64
		if err := rows.Scan(&jobID, &schema, &table, &unixTimestamp); err != nil {
			s.logger.Err(err).Msg("cannot scan PostgreSQL health-check retention history")
			return nil
		}
		if unixTimestamp == nil {
			continue
		}
		observation := storage.HistoryWithCount{Time: datetime.NewCpsTime(*unixTimestamp)}
		if latest == nil || latest.Time.Before(observation.Time) {
			latest = &observation
		}
	}
	if err := rows.Err(); err != nil {
		s.logger.Err(err).Msg("cannot iterate PostgreSQL health-check retention history")
		return nil
	}
	return latest
}

func (s *store) mergeHealthCheckHistory(data *storage.DataStorage, observation *storage.HistoryWithCount) {
	if observation == nil {
		return
	}
	existing, ok := data.History[healthCheckHistoryKey]
	if ok && !existing.Time.Before(observation.Time) {
		return
	}
	if data.History == nil {
		data.History = make(map[string]storage.HistoryWithCount)
	}
	data.History[healthCheckHistoryKey] = *observation
}

func (s *store) updateConfigWithHealthCheckHistory(conf storage.Config, observation storage.HistoryWithCount) bson.A {
	return bson.A{bson.M{"$set": bson.M{
		"config": conf,
		"history": bson.M{"$mergeObjects": bson.A{
			bson.M{"$ifNull": bson.A{"$history", bson.M{}}},
			bson.M{healthCheckHistoryKey: bson.M{"$cond": bson.M{
				"if":   bson.M{"$gt": bson.A{observation.Time, bson.M{"$ifNull": bson.A{"$history." + healthCheckHistoryKey + ".time", 0}}}},
				"then": observation,
				"else": "$history." + healthCheckHistoryKey,
			}}},
		}},
	}}}
}

func healthCheckRetentionPolicyMatches(ctx context.Context, pgPool postgres.Pool, deleteAfter *datetime.DurationWithEnabled) (bool, error) {
	var current bool
	if !datetime.IsDurationEnabledAndValid(deleteAfter) {
		err := pgPool.QueryRow(ctx, healthCheckRetentionPolicyExistsQuery, healthCheckSchema, metrics.MessageRateHourly).Scan(&current)
		return !current, err
	}
	err := pgPool.QueryRow(ctx, healthCheckRetentionPolicyQuery, healthCheckSchema, metrics.MessageRateHourly, deleteAfter.String()).Scan(&current)
	return current, err
}
