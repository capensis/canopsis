package contextgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"os"
	"time"
)

const (
	DefaultThdWarnMinPerImport = 30 * time.Minute
	DefaultThdCritMinPerImport = 60 * time.Minute

	BulkMaxSize = 1000
)

type worker struct {
	importQueue         JobQueue
	reporter            StatusReporter
	publisher           EventPublisher
	entityCollection    libmongo.DbCollection
	categoryCollection  libmongo.DbCollection
	logger              zerolog.Logger
	filePattern         string
	thdWarnMinPerImport time.Duration
	thdCritMinPerImport time.Duration
}

func NewImportWorker(
	conf config.CanopsisConf,
	dbClient libmongo.DbClient,
	publisher EventPublisher,
	reporter StatusReporter,
	queue JobQueue,
	logger zerolog.Logger,
) ImportWorker {
	w := &worker{
		importQueue:        queue,
		publisher:          publisher,
		reporter:           reporter,
		entityCollection:   dbClient.Collection(libmongo.EntityMongoCollection),
		categoryCollection: dbClient.Collection(libmongo.EntityCategoryMongoCollection),
		filePattern:        conf.ImportCtx.FilePattern,
	}

	thdWarnMinPerImport, err := time.ParseDuration(conf.ImportCtx.ThdWarnMinPerImport)
	if err != nil {
		logger.Warn().Err(err).Msg("Can't parse thdWarnMinPerImport value, use default")
		thdWarnMinPerImport = DefaultThdWarnMinPerImport
	}

	thdCritMinPerImport, err := time.ParseDuration(conf.ImportCtx.ThdCritMinPerImport)
	if err != nil {
		logger.Warn().Err(err).Msg("Can't parse thdCritMinPerImport value, use default")
		thdCritMinPerImport = DefaultThdCritMinPerImport
	}

	w.thdWarnMinPerImport = thdWarnMinPerImport
	w.thdCritMinPerImport = thdCritMinPerImport

	return w
}

func (w *worker) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			job := w.importQueue.Pop()

			if job.ID != "" {
				err := w.reporter.ReportOngoing(job)
				if err != nil {
					w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")

					err = w.publisher.SendImportResultEvent(job.ID, 0, types.AlarmStateCritical)
					if err != nil {
						w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send import result event")
					}

					break
				}

				startTime := time.Now()
				stats, err := w.doJob(ctx, job)
				stats.ExecTime = time.Since(startTime)

				resultState := types.AlarmStateOK
				if err != nil {
					w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Error during the import.")

					resultState = types.AlarmStateCritical
					if err == ErrNotImplemented {
						resultState = types.AlarmStateMinor
					}

					err = w.reporter.ReportError(job, stats.ExecTime, err)
					if err != nil {
						w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")
					}
				} else {
					w.logger.Info().Str("job_id", job.ID).Msg("Import-ctx: import done")

					err = w.reporter.ReportDone(job, stats)
					if err != nil {
						w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to update import info")
					}
				}

				perfDataState := types.AlarmStateOK
				if stats.ExecTime > w.thdCritMinPerImport {
					perfDataState = types.AlarmStateMajor
				} else if stats.ExecTime > w.thdWarnMinPerImport {
					perfDataState = types.AlarmStateMinor
				}

				err = w.publisher.SendPerfDataEvent(job.ID, stats, types.CpsNumber(perfDataState))
				if err != nil {
					w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send perf data")
				}

				err = w.publisher.SendImportResultEvent(job.ID, stats.ExecTime, types.CpsNumber(resultState))
				if err != nil {
					w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to send import result event")
				}
			}
		}
	}
}

func (w *worker) doJob(ctx context.Context, job ImportJob) (JobStats, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		writeModels []mongo.WriteModel
		jobStats    JobStats
	)

	w.logger.Info().Str("job_id", job.ID).Msg("Import-ctx: Processing import")

	filename := fmt.Sprintf(w.filePattern, job.ID)
	file, err := os.Open(filename)
	if err != nil {
		w.logger.Err(err).Str("job_id", job.ID).Msg("Import-ctx: Failed to open import file")

		return jobStats, err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			w.logger.Err(err).Msg("Import-ctx: Failed to close file")
		} else {
			err = os.Remove(filename)
			if err != nil {
				w.logger.Err(err).Msg("Import-ctx: Failed to remove file")
			}
		}
	}()

	dec := json.NewDecoder(file)

	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return jobStats, err
			}
		}

		if t == "cis" {
			t, err := dec.Token()
			if err != nil {
				return jobStats, fmt.Errorf("failed to parse cis: %v", err)
			}

			if t != json.Delim('[') {
				return jobStats, fmt.Errorf("cis should be an array")
			}

			for dec.More() {
				var ci ConfigurationItem
				err := dec.Decode(&ci)
				if err != nil {
					return jobStats, fmt.Errorf("failed to decode cis item: %v", err)
				}

				w.fillDefaultFields(&ci)
				err = w.validate(ci)
				if err != nil {
					return jobStats, fmt.Errorf("ci = %s, validation error: %s", ci.ID, err.Error())
				}

				switch ci.Action {
				case actionCreate:
					writeModels = append(writeModels, w.createEntity(ci))
				case actionSet:
					res := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID})
					if res.Err() != nil && res.Err() != mongo.ErrNoDocuments {
						return jobStats, res.Err()
					}

					if res.Err() == mongo.ErrNoDocuments {
						writeModels = append(writeModels, w.createEntity(ci))

						break
					}

					var oldEntity ConfigurationItem
					err := res.Decode(&oldEntity)
					if err != nil {
						return jobStats, err
					}

					writeModels = append(writeModels, w.updateEntity(ci, oldEntity, true))
				case actionUpdate:
					res := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID})
					if res.Err() != nil {
						if res.Err() == mongo.ErrNoDocuments {
							return jobStats, fmt.Errorf("failed to update an entity with _id = %s", ci.ID)
						}

						return jobStats, res.Err()
					}

					var oldEntity ConfigurationItem
					err := res.Decode(&oldEntity)
					if err != nil {
						return jobStats, err
					}

					writeModels = append(writeModels, w.updateEntity(ci, oldEntity, false))
				case actionDelete:
					res := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID})
					if res.Err() != nil {
						if res.Err() == mongo.ErrNoDocuments {
							return jobStats, fmt.Errorf("failed to delete an entity with _id = %s", ci.ID)
						}

						return jobStats, res.Err()
					}

					writeModels = append(writeModels, w.deleteEntity(ci)...)
				case actionEnable:
					res := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID})
					if res.Err() != nil {
						if res.Err() == mongo.ErrNoDocuments {
							return jobStats, fmt.Errorf("failed to enable an entity with _id = %s", ci.ID)
						}

						return jobStats, res.Err()
					}

					writeModels = append(writeModels, w.changeState(ci.ID, true))
				case actionDisable:
					res := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID})
					if res.Err() != nil {
						if res.Err() == mongo.ErrNoDocuments {
							return jobStats, fmt.Errorf("failed to disable an entity with _id = %s", ci.ID)
						}

						return jobStats, res.Err()
					}

					writeModels = append(writeModels, w.changeState(ci.ID, false))
				default:
					return jobStats, fmt.Errorf("the action %s is not recognized", ci.Action)
				}
			}
		}

		if t == "links" {
			t, err := dec.Token()
			if err != nil {
				return jobStats, fmt.Errorf("failed to parse links: %v", err)
			}

			if t != json.Delim('[') {
				return jobStats, fmt.Errorf("links should be an array")
			}

			for dec.More() {
				var link Link
				err := dec.Decode(&link)
				if err != nil {
					return jobStats, fmt.Errorf("failed to decode links item: %v", err)
				}

				switch link.Action {
				case actionCreate:
					writeModels = append(writeModels, w.createLink(link)...)
				case actionDelete:
					writeModels = append(writeModels, w.deleteLink(link)...)
				case actionUpdate:
					//wasn't implemented in python code
					return jobStats, ErrNotImplemented
				case actionEnable:
					//wasn't implemented in python code
					return jobStats, ErrNotImplemented
				case actionDisable:
					//wasn't implemented in python code
					return jobStats, ErrNotImplemented
				default:
					return jobStats, fmt.Errorf("the action %s is not recognized", link.Action)
				}
			}
		}
	}

	if len(writeModels) == 0 {
		return jobStats, fmt.Errorf("empty import")
	}

	// execute queries
	for i := 0; i < len(writeModels); i += BulkMaxSize {
		end := i + BulkMaxSize
		if i+BulkMaxSize > len(writeModels) {
			end = len(writeModels)
		}

		p := writeModels[i:end]
		result, err := w.entityCollection.BulkWrite(ctx, p)
		if err != nil {
			return jobStats, err
		}

		jobStats.Updated += result.UpsertedCount + result.ModifiedCount
		jobStats.Deleted += result.DeletedCount
	}

	//send events for services
	cursor, err := w.entityCollection.Find(ctx, bson.M{"type": types.EntityTypeService})
	if err != nil {
		return jobStats, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var service types.Entity

		err := cursor.Decode(&service)
		if err != nil {
			return jobStats, err
		}

		err = w.publisher.SendUpdateEntityServiceEvent(service.ID)
		if err != nil {
			return jobStats, err
		}
	}

	return jobStats, err
}

func (w *worker) validate(ci ConfigurationItem) error {
	if ci.ID == "" {
		return fmt.Errorf("_id is required")
	}

	if ci.Type == nil {
		return fmt.Errorf("type is required")
	}

	switch *ci.Type {
	case types.EntityTypeService:
	case types.EntityTypeResource:
	case types.EntityTypeComponent:
	case types.EntityTypeConnector:
	default:
		return fmt.Errorf("type is not valid")
	}

	if *ci.Type != types.EntityTypeService && ci.EntityPatterns != nil {
		return fmt.Errorf("contains entity patterns, but ci is not a service")
	}

	return nil
}

func (w *worker) fillDefaultFields(ci *ConfigurationItem) {
	if ci.Name == nil {
		ci.Name = &ci.ID
	}

	if ci.Category == nil {
		def := new(string)
		*def = entity.DefaultCategory

		ci.Category = def
	}

	if ci.ImpactLevel == nil {
		def := new(int64)
		*def = 1

		ci.ImpactLevel = def
	}
}

func (w *worker) createLink(link Link) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": link.To}).
			SetUpdate(bson.M{"$addToSet": bson.M{"depends": bson.M{"$each": link.From}}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": bson.M{"$in": link.From}}).
			SetUpdate(bson.M{"$addToSet": bson.M{"impact": link.To}}),
	}
}

func (w *worker) deleteLink(link Link) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": link.To}).
			SetUpdate(bson.M{"$pull": bson.M{"depends": bson.M{"$in": link.From}}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": bson.M{"$in": link.From}}).
			SetUpdate(bson.M{"$pull": bson.M{"impact": link.To}}),
	}
}

func (w *worker) createEntity(ci ConfigurationItem) mongo.WriteModel {
	ci.Depends = []string{}
	ci.Impact = []string{}
	ci.EnableHistory = make([]string, 0)

	if ci.Infos == nil {
		ci.Infos = make(map[string]interface{})
	}

	if ci.Measurements == nil {
		ci.Measurements = make([]interface{}, 0)
	}

	now := types.CpsTime{Time: time.Now()}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ci.ID}).
		SetUpdate(bson.M{"$set": ci, "$setOnInsert": bson.M{"created": now}}).
		SetUpsert(true)
}

func (w *worker) updateEntity(ci ConfigurationItem, oldEntity ConfigurationItem, mergeInfos bool) mongo.WriteModel {
	ci.Depends = oldEntity.Depends
	ci.Impact = oldEntity.Impact
	ci.EnableHistory = oldEntity.EnableHistory

	if ci.Infos == nil {
		ci.Infos = make(map[string]interface{})
	}

	if mergeInfos {
		for k, v := range ci.Infos {
			oldEntity.Infos[k] = v
		}

		ci.Infos = oldEntity.Infos
	}

	now := types.CpsTime{Time: time.Now()}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ci.ID}).
		SetUpdate(bson.M{"$set": ci, "$setOnInsert": bson.M{"created": now}}).
		SetUpsert(true)
}

func (w *worker) changeState(id string, enabled bool) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": id}).
		SetUpdate(bson.M{"$set": bson.M{"enabled": enabled}})
}

func (w *worker) deleteEntity(ci ConfigurationItem) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"impact": ci.ID}).
			SetUpdate(bson.M{"$pull": bson.M{"impact": ci.ID}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"depends": ci.ID}).
			SetUpdate(bson.M{"$pull": bson.M{"depends": ci.ID}}),
		mongo.NewDeleteOneModel().
			SetFilter(bson.M{"_id": ci.ID}),
	}
}
