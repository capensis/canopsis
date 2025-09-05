package externaldatatable

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	filePerm os.FileMode = 0o770

	maxRetries = 10

	sqlColumnsForCsvColumn = 3
	sqlUndefinedColumnCode = "42703"
)

type ImportWorker interface {
	CreateImportJob(ctx context.Context, id string, separator rune, f multipart.File) (_ ImportJob, resErr error)
	ProcessJob(ctx context.Context, id string) error
	GetJob(ctx context.Context, id string) (ImportJob, error)
	CompleteJob(ctx context.Context, id string, columnTags []int) (bool, error)
	ProcessAbandonedJobs(ctx context.Context)
	DeleteOldJobs(ctx context.Context)
	CreatePreviewJob(ctx context.Context, id string, req PreviewRequest) (ImportJob, error)
}

func NewImportWorker(
	dbClient mongo.DbClient,
	pgPoolProvider postgres.PoolProvider,
	tmpImportDir string,
	jobPublisher workers.JobPublisher,
	logger zerolog.Logger,
) ImportWorker {
	linkedDbCollections := make([]mongo.DbCollection, len(linkedCollections))
	for i, c := range linkedCollections {
		linkedDbCollections[i] = dbClient.Collection(c)
	}

	return &importWorker{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(mongo.ExternalDataTableCollection),
		dbImportCollection:      dbClient.Collection(mongo.ExternalDataImportWorkerCollection),
		dbWidgetCollection:      dbClient.Collection(mongo.WidgetMongoCollection),
		linkedDbCollections:     linkedDbCollections,
		pgPoolProvider:          pgPoolProvider,
		tmpImportDir:            tmpImportDir,
		jobPublisher:            jobPublisher,
		logger:                  logger,
		abandonedTickerInterval: time.Minute,
		pingInterval:            time.Second,
		deleteTickerInterval:    time.Hour,
		deleteInterval:          24 * time.Hour,

		parser: NewParser(),
	}
}

type importWorker struct {
	dbClient                mongo.DbClient
	dbCollection            mongo.DbCollection
	dbImportCollection      mongo.DbCollection
	dbWidgetCollection      mongo.DbCollection
	linkedDbCollections     []mongo.DbCollection
	pgPoolProvider          postgres.PoolProvider
	tmpImportDir            string
	jobPublisher            workers.JobPublisher
	abandonedTickerInterval time.Duration
	pingInterval            time.Duration
	deleteTickerInterval    time.Duration
	deleteInterval          time.Duration
	logger                  zerolog.Logger

	parser Parser
}

type Row struct {
	ID      string            `bson:"_id"`
	Columns map[string]Column `bson:",inline"`
}
type Column struct {
	InitialValue     string `bson:"initial_value"`
	TransformedValue any    `bson:"transformed_value,omitempty"`
	TransformError   string `bson:"transform_error,omitempty"`
}

func (w *importWorker) CreateImportJob(ctx context.Context, id string, separator rune, f multipart.File) (_ ImportJob, resErr error) {
	defer func() {
		err := f.Close()
		if err != nil && resErr == nil {
			resErr = fmt.Errorf("cannot close file: %w", err)
		}
	}()

	job := ImportJob{}
	externalDataTable := externaldata.Table{}
	err := w.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&externalDataTable)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return job, common.NewValidationError("_id", "ID doesn't exist.")
		}

		return job, err
	}

	err = w.validateColumns(ctx, externalDataTable, f, separator)
	if err != nil {
		return job, err
	}

	now := datetime.NewCpsTime()
	jobID := utils.NewID()
	job = ImportJob{
		ID:                jobID,
		Type:              externalDataTable.Type,
		Table:             tmpTablePrefix + jobID,
		ExternalDataTable: id,
		Status:            ImportStatusCreated,
		Separator:         separator,
		Filepath:          filepath.Join(w.tmpImportDir, jobID),
		Created:           now,
	}

	err = os.MkdirAll(w.tmpImportDir, os.ModeDir|filePerm)
	if err != nil {
		return job, fmt.Errorf("failed to create dir %q: %w", w.tmpImportDir, err)
	}

	df, err := os.OpenFile(job.Filepath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, filePerm)
	if err != nil {
		return job, fmt.Errorf("failed to open file %q: %w", job.Filepath, err)
	}

	defer func() {
		err = df.Close()
		if err != nil && resErr == nil {
			resErr = fmt.Errorf("failed to close file %q: %w", job.Filepath, err)
		}
	}()

	n, err := io.Copy(df, f)
	if err != nil {
		return job, fmt.Errorf("failed to copy to file %q: %w", job.Filepath, err)
	}

	if n == 0 {
		return job, fmt.Errorf("failed to copy to file %q: source file is empty", job.Filepath)
	}

	_, err = w.dbImportCollection.InsertOne(ctx, job)
	if err != nil {
		return job, fmt.Errorf("failed to insert document: %w", err)
	}

	err = w.jobPublisher.Publish(ctx, job.ID)
	if err != nil {
		return job, fmt.Errorf("failed to publish job: %w", err)
	}

	return job, nil
}

func (w *importWorker) CreatePreviewJob(ctx context.Context, id string, req PreviewRequest) (ImportJob, error) {
	var job ImportJob

	err := w.dbImportCollection.FindOneAndUpdate(ctx,
		bson.M{
			"_id":    id,
			"status": bson.M{"$in": bson.A{ImportStatusSucceeded, ImportStatusPreviewSucceeded, ImportStatusPreviewFailed}},
		},
		[]bson.M{
			{
				"$set": bson.M{
					"job_type":            JobTypePreview,
					"status":              ImportStatusCreated,
					"prev_column_configs": "$column_configs",
					"column_configs":      req.ColumnConfigs,
				},
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return job, nil
		}

		return job, fmt.Errorf("failed to update import job: %w", err)
	}

	err = w.jobPublisher.Publish(ctx, job.ID)
	if err != nil {
		return job, fmt.Errorf("failed to publish job: %w", err)
	}

	return job, nil
}

func (w *importWorker) ProcessJob(ctx context.Context, id string) (resErr error) {
	job := ImportJob{}
	err := w.dbImportCollection.FindOneAndUpdate(ctx,
		bson.M{
			"_id": id,
			"$or": []bson.M{
				{"status": ImportStatusCreated},
				{
					"status":    ImportStatusRunning,
					"last_ping": bson.M{"$lt": time.Now().Add(-2 * w.pingInterval).Unix()},
				},
			},
		},
		bson.M{"$set": bson.M{
			"status":    ImportStatusRunning,
			"last_ping": datetime.NewCpsTime(),
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}

		return err
	}

	var errorInfo map[string]ErrorInfo

	g, ctx := errgroup.WithContext(ctx)
	done := make(chan struct{})
	g.Go(func() (resErr error) {
		defer close(done)

		var update bson.M

		switch job.JobType {
		case JobTypeImport:
			rmFile := false
			f, err := os.Open(job.Filepath)
			if err != nil {
				return fmt.Errorf("failed to open file %q: %w", job.Filepath, err)
			}

			defer func() {
				err = f.Close()
				if err != nil && resErr == nil {
					resErr = fmt.Errorf("failed to close file %q: %w", job.Filepath, err)
				}

				if err == nil && rmFile {
					err = os.Remove(job.Filepath)
					if err != nil && resErr == nil {
						resErr = fmt.Errorf("failed to remove file %q for job %q: %w", job.Filepath, job.ID, err)
					}
				}
			}()
			r := csv.NewReader(f)
			if job.Separator != 0 {
				r.Comma = job.Separator
			}

			r.ReuseRecord = true
			var columnsConfig []externaldata.ColumnConfig

			switch job.Type {
			case externaldata.TypeMongoDB:
				columnsConfig, err = w.writeToMongo(ctx, job, r)
			case externaldata.TypePostgreSQL:
				columnsConfig, err = w.writeToPostgres(ctx, job, r)
			}

			var failReason string

			if err != nil {
				if mongo.IsConnectionError(err) || postgres.IsConnectionError(err) {
					return err
				}

				failReason = err.Error()
				w.logger.Err(errors.New(failReason)).Str("job", job.ID).Msg("failed to import external data")

				update = bson.M{
					"$set": bson.M{
						"status":      ImportStatusFailed,
						"fail_reason": failReason,
					},
				}
			} else {
				update = bson.M{
					"$set": bson.M{
						"status":         ImportStatusSucceeded,
						"column_configs": columnsConfig,
					},
				}
			}

			rmFile = true
			if failReason != "" {
				err = w.deleteTable(ctx, job)
				if err != nil {
					return err
				}
			}
		case JobTypePreview:
			switch job.Type {
			case externaldata.TypeMongoDB:
				errorInfo, err = w.previewMongo(ctx, job)
			case externaldata.TypePostgreSQL:
				errorInfo, err = w.previewPostgres(ctx, job)
			}

			if err != nil {
				if mongo.IsConnectionError(err) || postgres.IsConnectionError(err) {
					return err
				}

				update = bson.M{
					"$set": bson.M{
						"status":      ImportStatusPreviewFailed,
						"fail_reason": err.Error(),
						// replace it with prev configs in order not to save invalid ones.
						"column_configs": job.PrevColumnConfigs,
					},
					"$unset": bson.M{
						"error_info": "",
					},
				}
			} else if len(errorInfo) != 0 {
				update = bson.M{
					"$set": bson.M{
						"status":      ImportStatusPreviewFailed,
						"error_info":  errorInfo,
						"fail_reason": "preview failed",
					},
				}
			} else {
				update = bson.M{
					"$set": bson.M{
						"status": ImportStatusPreviewSucceeded,
					},
					"$unset": bson.M{
						"error_info":  "",
						"fail_reason": "",
					},
				}
			}
		}

		updateRes, err := w.dbImportCollection.UpdateOne(ctx,
			bson.M{
				"_id":    job.ID,
				"status": ImportStatusRunning,
			}, update,
		)
		if err != nil {
			return fmt.Errorf("failed to update import job: %w", err)
		}

		if updateRes.ModifiedCount == 0 {
			return errors.New("import job is processing by another worker")
		}

		return nil
	})
	g.Go(func() error {
		ticker := time.NewTicker(w.pingInterval)
		lastPing := *job.LastPing
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			case <-ticker.C:
				newLastPing := datetime.NewCpsTime()
				updateRes, err := w.dbImportCollection.UpdateOne(ctx,
					bson.M{
						"_id":       job.ID,
						"status":    ImportStatusRunning,
						"last_ping": lastPing,
					},
					bson.M{"$set": bson.M{
						"last_ping": newLastPing,
					}},
				)
				if err != nil {
					return fmt.Errorf("failed to update import job status: %w", err)
				}

				if updateRes.ModifiedCount == 0 {
					return errors.New("import job is processing by another worker")
				}

				lastPing = newLastPing
			}
		}
	})

	return g.Wait()
}

func (w *importWorker) GetJob(ctx context.Context, id string) (ImportJob, error) {
	job := ImportJob{}
	err := w.dbImportCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return job, nil
		}

		return job, err
	}

	return job, nil
}

func (w *importWorker) CompleteJob(ctx context.Context, id string, columnTags []int) (bool, error) {
	job := ImportJob{}

	err := w.dbImportCollection.FindOne(ctx, bson.M{"_id": id, "status": ImportStatusPreviewSucceeded}).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, fmt.Errorf("failed to find import job: %w", err)
	}

	if len(columnTags) != len(job.ColumnConfigs) {
		return false, common.NewValidationError("column_tags", "ColumnTags must contain "+strconv.Itoa(len(job.ColumnConfigs))+" items.")
	}

	for i, columnTag := range columnTags {
		job.ColumnConfigs[i].Tag = &columnTag
	}

	table := externaldata.Table{}

	err = w.dbCollection.FindOneAndUpdate(ctx,
		bson.M{"_id": job.ExternalDataTable},
		bson.M{"$set": bson.M{"column_configs": job.ColumnConfigs}}).Decode(&table)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, fmt.Errorf("failed to find table: %w", err)
	}

	switch table.Type {
	case externaldata.TypeMongoDB:
		project := bson.M{}
		for _, c := range job.ColumnConfigs {
			project[c.Name] = "$" + c.Name + ".transformed_value"
		}

		cursor, err := w.dbClient.Collection(job.getDBTableName()).Aggregate(ctx, []bson.M{
			{"$project": project},
			{"$out": table.GetDBName()},
		})
		if err != nil {
			return false, fmt.Errorf("failed to copy to mongo collection: %w", err)
		}

		if err = cursor.Err(); err != nil {
			return false, fmt.Errorf("failed to copy to mongo collection, cursor err: %w", err)
		}

		err = cursor.Close(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to copy to mongo collection, cursor close: %w", err)
		}
	case externaldata.TypePostgreSQL:
		pgPool, err := w.pgPoolProvider.Get(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.Exec(ctx, "TRUNCATE "+table.GetDBName())
		if err != nil {
			return false, fmt.Errorf("failed to truncate postgres table: %w", err)
		}

		existedCols := make(map[string]bool, len(table.ColumnConfigs))
		for _, c := range table.ColumnConfigs {
			existedCols[c.Name] = true
		}

		for _, c := range job.ColumnConfigs {
			var columnType string

			switch c.Type {
			case externaldata.ColumnTypeString:
				columnType = "VARCHAR(" + externaldata.MaxStringLenStr + ")"
			case externaldata.ColumnTypeNumber:
				columnType = "DOUBLE PRECISION"
			case externaldata.ColumnTypeBoolean:
				columnType = "BOOLEAN"
			case externaldata.ColumnTypeDateTime, externaldata.ColumnTypeTimestamp:
				columnType = "BIGINT"
			case externaldata.ColumnTypeStringArray:
				columnType = "VARCHAR(" + externaldata.MaxStringLenStr + ")[]"
			default:
				return false, fmt.Errorf("unsupported column type %q", c.Type)
			}

			sql := ""
			if existedCols[c.Name] {
				sql = "ALTER TABLE " + table.GetDBName() +
					" ALTER COLUMN " + pgx.Identifier{c.Name}.Sanitize() + " TYPE " + columnType + " USING NULL"
			} else {
				sql = "ALTER TABLE " + table.GetDBName() +
					" ADD COLUMN " + pgx.Identifier{c.Name}.Sanitize() + " " + columnType
			}

			if sql != "" {
				_, err = pgPool.Exec(ctx, sql)
				if err != nil {
					return false, fmt.Errorf("failed to alter postgres table: %w", err)
				}
			}

			delete(existedCols, c.Name)
		}

		for c := range existedCols {
			sql := "ALTER TABLE " + table.GetDBName() +
				" DROP COLUMN " + pgx.Identifier{c}.Sanitize()
			_, err = pgPool.Exec(ctx, sql)
			if err != nil {
				return false, fmt.Errorf("failed to alter postgres table: %w", err)
			}
		}

		sanColsWithID := make([]string, len(job.ColumnConfigs)+1)
		sanColsWithID[0] = externaldata.IDColumnName
		sanTransformedColsWithID := make([]string, len(job.ColumnConfigs)+1)
		sanTransformedColsWithID[0] = externaldata.IDColumnName
		for j, col := range job.ColumnConfigs {
			sanColsWithID[j+1] = pgx.Identifier{col.Name}.Sanitize()
			sanTransformedColsWithID[j+1] = pgx.Identifier{col.Name + "_transformed_value"}.Sanitize()
		}

		_, err = pgPool.Exec(ctx, "INSERT INTO "+table.GetDBName()+"("+strings.Join(sanColsWithID, ",")+
			") SELECT "+strings.Join(sanTransformedColsWithID, ",")+" FROM "+job.getDBTableName())
		if err != nil {
			return false, fmt.Errorf("failed to copy to postgres table: %w", err)
		}
	default:
		return false, fmt.Errorf("invalid table type: %q", table.Type)
	}

	err = w.deleteTable(ctx, job)
	if err != nil {
		return false, err
	}

	_, err = w.dbImportCollection.DeleteOne(ctx, bson.M{"_id": job.ID})
	if err != nil {
		return false, fmt.Errorf("failed to delete import job: %w", err)
	}

	return true, nil
}

func (w *importWorker) ProcessAbandonedJobs(ctx context.Context) {
	ticker := time.NewTicker(w.abandonedTickerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := datetime.NewCpsTime()
			cursor, err := w.dbImportCollection.Find(ctx, bson.M{
				"status":    ImportStatusRunning,
				"last_ping": bson.M{"$lt": now.Time.Add(-2 * w.pingInterval).Unix()},
			})
			if err != nil {
				w.logger.Err(err).Msg("failed to find abandoned import jobs")
				continue
			}

			for cursor.Next(ctx) {
				job := ImportJob{}
				if err = cursor.Decode(&job); err != nil {
					w.logger.Err(err).Msg("failed to decode abandoned import job")
					continue
				}

				if job.Retries == maxRetries {
					_, err = w.dbImportCollection.UpdateOne(ctx, bson.M{"_id": job.ID}, bson.M{"$set": bson.M{
						"status":      ImportStatusFailed,
						"fail_reason": "max retries exceeded",
					}})
					if err != nil {
						w.logger.Err(err).Msg("failed to update abandoned failed import job")
					}

					continue
				}

				job.Retries++
				_, err = w.dbImportCollection.UpdateOne(ctx, bson.M{"_id": job.ID}, bson.M{"$set": bson.M{
					"retries":   job.Retries,
					"last_ping": now,
				}})
				if err != nil {
					w.logger.Err(err).Msg("failed to update abandoned import job")
					continue
				}

				err = w.jobPublisher.Publish(ctx, job.ID)
				if err != nil {
					w.logger.Err(err).Msg("failed to publish abandoned import job")
					continue
				}
			}

			if err = cursor.Err(); err != nil {
				w.logger.Err(err).Msg("failed to fetch abandoned import jobs")
			}

			err = cursor.Close(ctx)
			if err != nil {
				w.logger.Err(err).Msg("failed to close cursor")
			}
		}
	}
}

func (w *importWorker) DeleteOldJobs(ctx context.Context) {
	ticker := time.NewTicker(w.deleteTickerInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cursor, err := w.dbImportCollection.Find(ctx, bson.M{
				"status":    bson.M{"$in": bson.A{ImportStatusSucceeded, ImportStatusFailed}},
				"last_ping": bson.M{"$lt": time.Now().Add(-w.deleteInterval).Unix()},
			})
			if err != nil {
				w.logger.Err(err).Msg("failed to find old import jobs")
				continue
			}

			ids := make([]string, 0, canopsis.DefaultBulkSize)
			for cursor.Next(ctx) {
				job := ImportJob{}
				if err := cursor.Decode(&job); err != nil {
					w.logger.Err(err).Msg("failed to decode old import job")
					continue
				}

				if job.Status == ImportStatusSucceeded {
					err = w.deleteTable(ctx, job)
					if err != nil {
						w.logger.Err(err).Msg("failed to delete table of old import job")
						continue
					}
				}

				ids = append(ids, job.ID)
				if len(ids) == canopsis.DefaultBulkSize {
					_, err = w.dbImportCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
					if err != nil {
						w.logger.Err(err).Msg("failed to delete old import jobs")
						continue
					}

					ids = ids[:0]
				}
			}

			if len(ids) > 0 {
				_, err = w.dbImportCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
				if err != nil {
					w.logger.Err(err).Msg("failed to delete old import jobs")
				}
			}

			if err = cursor.Err(); err != nil {
				w.logger.Err(err).Msg("failed to fetch old import jobs")
			}

			err = cursor.Close(ctx)
			if err != nil {
				w.logger.Err(err).Msg("failed to close cursor")
			}
		}
	}
}

func (w *importWorker) writeToMongo(ctx context.Context, job ImportJob, r *csv.Reader) ([]externaldata.ColumnConfig, error) {
	tableName := job.getDBTableName()

	var columns []externaldata.ColumnConfig
	i := 0
	docs := make([]any, 0, canopsis.DefaultBulkSize)
	for {
		i++
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		if len(record) == 0 {
			return nil, fmt.Errorf("empty record %d", i)
		}

		if len(columns) == 0 {
			columns = make([]externaldata.ColumnConfig, len(record))
			for idx, r := range record {
				columns[idx] = externaldata.ColumnConfig{
					BaseColumnConfig: externaldata.BaseColumnConfig{
						Name: r,
						Type: externaldata.ColumnTypeString,
					},
				}
			}

			err = w.createTable(ctx, job, nil)
			if err != nil {
				return nil, err
			}

			continue
		}

		if len(columns) != len(record) {
			return nil, fmt.Errorf("invalid record %d: not match fields", i)
		}

		doc := make(map[string]any, len(record)+1)
		doc[externaldata.IDColumnName] = utils.NewID()
		for j, c := range columns {
			if utf8.RuneCountInString(record[j]) > externaldata.MaxStringLen {
				return nil, fmt.Errorf("invalid record %d: string length must be less than "+externaldata.MaxStringLenStr, i)
			}

			doc[c.Name] = Column{InitialValue: record[j], TransformedValue: record[j]}
		}

		docs = append(docs, doc)
		if len(docs) == canopsis.DefaultBulkSize {
			_, err = w.dbClient.Collection(tableName).InsertMany(ctx, docs)
			if err != nil {
				return nil, fmt.Errorf("failed to insert into %q collection: %w", job.Table, err)
			}

			docs = docs[:0]
		}
	}

	if len(docs) > 0 {
		_, err := w.dbClient.Collection(tableName).InsertMany(ctx, docs)
		if err != nil {
			return nil, fmt.Errorf("failed to insert into %q collection: %w", job.Table, err)
		}
	}

	return columns, nil
}

func (w *importWorker) writeToPostgres(ctx context.Context, job ImportJob, r *csv.Reader) ([]externaldata.ColumnConfig, error) {
	pgPool, err := w.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres pool: %w", err)
	}

	var columnConfigs []externaldata.ColumnConfig
	var columnsWithID []string

	i := 0
	rows := make([][]any, 0, canopsis.DefaultBulkSize)
	for {
		i++
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		if len(record) == 0 {
			return nil, fmt.Errorf("empty record %d", i)
		}

		if len(columnConfigs) == 0 {
			columnsWithID = make([]string, len(record)*sqlColumnsForCsvColumn+1)
			columnsWithID[0] = externaldata.IDColumnName

			columnConfigs = make([]externaldata.ColumnConfig, len(record))
			for idx, r := range record {
				columnConfigs[idx] = externaldata.ColumnConfig{
					BaseColumnConfig: externaldata.BaseColumnConfig{
						Name: r,
						Type: externaldata.ColumnTypeString,
					},
				}

				columnsWithID[idx*sqlColumnsForCsvColumn+1] = r + "_initial_value"
				columnsWithID[idx*sqlColumnsForCsvColumn+2] = r + "_transformed_value"
				columnsWithID[idx*sqlColumnsForCsvColumn+3] = r + "_transform_error"
			}

			err = w.createTable(ctx, job, columnsWithID[1:])
			if err != nil {
				return nil, err
			}

			continue
		}

		if len(columnConfigs) != len(record) {
			return nil, fmt.Errorf("invalid record %d: not match fields", i)
		}

		row := make([]any, len(record)*sqlColumnsForCsvColumn+1)
		row[0] = utils.NewID()
		for j, v := range record {
			if utf8.RuneCountInString(v) > externaldata.MaxStringLen {
				return nil, fmt.Errorf("invalid record %d: string length must be less than "+externaldata.MaxStringLenStr, i)
			}

			row[sqlColumnsForCsvColumn*j+1] = v
			row[sqlColumnsForCsvColumn*j+2] = v
		}

		rows = append(rows, row)

		if len(rows) == canopsis.DefaultBulkSize {
			_, err = pgPool.CopyFrom(ctx, externaldata.GetPostgresTableIdentifier(job.Table), columnsWithID, pgx.CopyFromRows(rows))
			if err != nil {
				return nil, fmt.Errorf("failed to insert into %q table: %w", job.Table, err)
			}

			rows = rows[:0]
		}
	}

	if len(rows) > 0 {
		_, err = pgPool.CopyFrom(ctx, externaldata.GetPostgresTableIdentifier(job.Table), columnsWithID, pgx.CopyFromRows(rows))
		if err != nil {
			return nil, fmt.Errorf("failed to insert into %q table: %w", job.Table, err)
		}
	}

	return columnConfigs, nil
}

func (w *importWorker) previewMongo(ctx context.Context, job ImportJob) (map[string]ErrorInfo, error) {
	prevColumnConfigs := job.PrevColumnConfigs
	if len(prevColumnConfigs) != len(job.ColumnConfigs) {
		return nil, errors.New("column config count mismatch")
	}

	errorInfo := make(map[string]ErrorInfo)

	columnConfigs := make([]externaldata.ColumnConfig, 0, len(job.ColumnConfigs))
	for i, cfg := range job.ColumnConfigs {
		if cfg.BaseColumnConfig != prevColumnConfigs[i].BaseColumnConfig {
			columnConfigs = append(columnConfigs, job.ColumnConfigs[i])
		} else {
			// if configs are not changed, check if there was error info, if yes, keep it
			if errInfo, ok := job.ErrorInfo[cfg.Name]; ok {
				errorInfo[cfg.Name] = errInfo
			}
		}
	}

	if len(columnConfigs) == 0 {
		return errorInfo, nil
	}

	columnErrRows := make(map[string][]int, len(columnConfigs))
	columnErrMessages := make(map[string][]string, len(columnConfigs))
	uniqueColumnErrMessages := make(map[string]map[string]bool, len(columnConfigs))

	for _, columnConfig := range columnConfigs {
		columnErrRows[columnConfig.Name] = make([]int, 0)
		columnErrMessages[columnConfig.Name] = make([]string, 0)
		uniqueColumnErrMessages[columnConfig.Name] = make(map[string]bool)
	}

	collectionName := job.getDBTableName()

	cursor, err := w.dbClient.Collection(collectionName).Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find temporary external data: %w", err)
	}

	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)

	for i := 0; cursor.Next(ctx); i++ {
		var row Row

		if err = cursor.Decode(&row); err != nil {
			return nil, fmt.Errorf("failed to decode preview data: %w", err)
		}

		if row.ID == "" {
			return nil, errors.New("failed to get id for preview data")
		}

		set := bson.M{}
		unset := bson.M{}

		for _, columnConfig := range columnConfigs {
			columnName := columnConfig.Name

			column, ok := row.Columns[columnName]
			if !ok {
				return nil, fmt.Errorf("no such column %q", columnName)
			}

			transformedValue, err := w.parser.Parse(columnConfig, column.InitialValue)
			if err != nil {
				errorMsg := err.Error()

				set[columnName+".transform_error"] = errorMsg
				unset[columnName+".transformed_value"] = ""

				columnErrRows[columnName] = append(columnErrRows[columnName], i)
				if !uniqueColumnErrMessages[columnName][errorMsg] {
					uniqueColumnErrMessages[columnName][errorMsg] = true
					columnErrMessages[columnName] = append(columnErrMessages[columnName], errorMsg)
				}
			} else {
				set[columnName+".transformed_value"] = transformedValue
				unset[columnName+".transform_error"] = ""
			}
		}

		writeModels = append(writeModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{externaldata.IDColumnName: row.ID}).
			SetUpdate(bson.M{"$set": set, "$unset": unset}))
		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err = w.dbClient.Collection(collectionName).BulkWrite(ctx, writeModels)
			if err != nil {
				return nil, fmt.Errorf("failed to insert into %q collection: %w", job.Table, err)
			}

			writeModels = writeModels[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to read temporary external data: %w", err)
	}

	if len(writeModels) > 0 {
		_, err := w.dbClient.Collection(collectionName).BulkWrite(ctx, writeModels)
		if err != nil {
			return nil, fmt.Errorf("failed to insert into %q collection: %w", job.Table, err)
		}
	}

	for columnName := range columnErrRows {
		if len(columnErrRows[columnName]) == 0 {
			continue
		}

		errorInfo[columnName] = ErrorInfo{
			Rows:     columnErrRows[columnName],
			Messages: columnErrMessages[columnName],
		}
	}

	return errorInfo, nil
}

func (w *importWorker) previewPostgres(ctx context.Context, job ImportJob) (map[string]ErrorInfo, error) {
	prevColumnConfigs := job.PrevColumnConfigs
	if len(prevColumnConfigs) != len(job.ColumnConfigs) {
		return nil, errors.New("column config count mismatch")
	}

	errorInfo := make(map[string]ErrorInfo)

	columnConfigs := make([]externaldata.ColumnConfig, 0, len(job.ColumnConfigs))
	for i, cfg := range job.ColumnConfigs {
		if cfg.BaseColumnConfig != prevColumnConfigs[i].BaseColumnConfig {
			columnConfigs = append(columnConfigs, job.ColumnConfigs[i])
		} else {
			// if configs are not changed, check if there was error info, if yes, keep it
			if errInfo, ok := job.ErrorInfo[cfg.Name]; ok {
				errorInfo[cfg.Name] = errInfo
			}
		}
	}

	if len(columnConfigs) == 0 {
		return errorInfo, nil
	}

	pgPool, err := w.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres pool: %w", err)
	}

	tableName := job.getDBTableName()

	selectCols := make([]string, 0, len(columnConfigs)+1)
	selectCols = append(selectCols, pgx.Identifier{externaldata.IDColumnName}.Sanitize())

	for _, cfg := range columnConfigs {
		selectCols = append(selectCols, pgx.Identifier{cfg.Name + "_initial_value"}.Sanitize())

		var columnType string
		switch cfg.Type {
		case externaldata.ColumnTypeString:
			columnType = "VARCHAR(" + externaldata.MaxStringLenStr + ")"
		case externaldata.ColumnTypeNumber:
			columnType = "DOUBLE PRECISION"
		case externaldata.ColumnTypeBoolean:
			columnType = "BOOLEAN"
		case externaldata.ColumnTypeDateTime, externaldata.ColumnTypeTimestamp:
			columnType = "BIGINT"
		case externaldata.ColumnTypeStringArray:
			columnType = "VARCHAR(" + externaldata.MaxStringLenStr + ")[]"
		default:
			return nil, fmt.Errorf("unexpected column type: %d", cfg.Type)
		}

		if _, err := pgPool.Exec(ctx, "ALTER TABLE "+tableName+" ALTER COLUMN "+pgx.Identifier{cfg.Name + "_transformed_value"}.Sanitize()+" TYPE "+columnType+" USING NULL"); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == sqlUndefinedColumnCode {
					return nil, fmt.Errorf("no such column %q", cfg.Name)
				}
			}

			return nil, fmt.Errorf("failed to alter column in temprorary table: %w", err)
		}
	}

	columnErrRows := make(map[string][]int, len(columnConfigs))
	columnErrMessages := make(map[string][]string, len(columnConfigs))
	uniqueColumnErrMessages := make(map[string]map[string]bool, len(columnConfigs))

	for _, columnConfig := range columnConfigs {
		columnErrRows[columnConfig.Name] = make([]int, 0)
		columnErrMessages[columnConfig.Name] = make([]string, 0)
		uniqueColumnErrMessages[columnConfig.Name] = make(map[string]bool)
	}

	batch := &pgx.Batch{}

	rows, err := pgPool.Query(ctx, "SELECT "+strings.Join(selectCols, ", ")+" FROM "+tableName+" ORDER BY "+pgx.Identifier{externaldata.IDColumnName}.Sanitize())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rowIdx := 0; rows.Next(); rowIdx++ {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}

		id, ok := vals[0].(string)
		if !ok {
			return nil, fmt.Errorf("%q column doesn't contain string", externaldata.IDColumnName)
		}

		setStmts := make([]string, 0, len(columnConfigs)*2)
		args := make([]any, 0, len(columnConfigs)*2+1)

		for i, columnConfig := range columnConfigs {
			columnName := columnConfig.Name

			initialValue := ""
			if vals[i+1] != nil {
				if initialValue, ok = vals[i+1].(string); !ok {
					return nil, fmt.Errorf("initial value for %q column doesn't contain a string", columnName)
				}
			}

			transformedValue, err := w.parser.Parse(columnConfig, initialValue)
			if err != nil {
				errorMsg := err.Error()

				args = append(args, errorMsg, nil)

				setStmts = append(setStmts, pgx.Identifier{columnName + "_transform_error"}.Sanitize()+" = $"+strconv.Itoa(len(args)-1))
				setStmts = append(setStmts, pgx.Identifier{columnName + "_transformed_value"}.Sanitize()+" = $"+strconv.Itoa(len(args)))

				columnErrRows[columnName] = append(columnErrRows[columnName], rowIdx)
				if !uniqueColumnErrMessages[columnName][errorMsg] {
					uniqueColumnErrMessages[columnName][errorMsg] = true
					columnErrMessages[columnName] = append(columnErrMessages[columnName], errorMsg)
				}
			} else {
				args = append(args, transformedValue, nil)

				setStmts = append(setStmts, pgx.Identifier{columnName + "_transformed_value"}.Sanitize()+" = $"+strconv.Itoa(len(args)-1))
				setStmts = append(setStmts, pgx.Identifier{columnName + "_transform_error"}.Sanitize()+" = $"+strconv.Itoa(len(args)))
			}
		}

		args = append(args, id)
		updateSQL := "UPDATE " + tableName + " SET " + strings.Join(setStmts, ", ") +
			" WHERE " + pgx.Identifier{externaldata.IDColumnName}.Sanitize() + " = $" + strconv.Itoa(len(args))
		batch.Queue(updateSQL, args...)

		if batch.Len() == canopsis.DefaultBulkSize {
			if err := pgPool.SendBatch(ctx, batch); err != nil {
				return nil, fmt.Errorf("failed to update preview rows: %w", err)
			}
			batch = &pgx.Batch{}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if batch.Len() > 0 {
		if err := pgPool.SendBatch(ctx, batch); err != nil {
			return nil, fmt.Errorf("failed to update preview rows: %w", err)
		}
	}

	for columnName := range columnErrRows {
		if len(columnErrRows[columnName]) == 0 {
			continue
		}

		errorInfo[columnName] = ErrorInfo{
			Rows:     columnErrRows[columnName],
			Messages: columnErrMessages[columnName],
		}
	}

	return errorInfo, nil
}

func (w *importWorker) validateColumns(ctx context.Context, t externaldata.Table, f multipart.File, separator rune) error {
	isLinked := false
	var err error
	if len(t.ColumnConfigs) > 0 {
		isLinked, err = isTableLinked(ctx, t.ID, w.dbWidgetCollection, w.linkedDbCollections)
		if err != nil {
			return err
		}
	}

	r := csv.NewReader(f)
	if separator != 0 {
		r.Comma = separator
	}

	columns, err := r.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return common.NewValidationError("file", "File is empty.")
		}

		return common.NewValidationError("file", "File is invalid.")
	}

	invalidCols := make([]string, 0)
	existColumns := make(map[string]bool, len(columns))
	for _, c := range columns {
		existColumns[c] = true
		if !common.IsTableName(c) || c == externaldata.IDColumnName {
			invalidCols = append(invalidCols, strconv.Quote(c))
		}
	}

	if len(invalidCols) > 0 {
		return common.NewValidationError("file", "Fields ["+strings.Join(invalidCols, ",")+"] in file are invalid.")
	}

	missingCols := make([]string, 0)
	if isLinked {
		for _, c := range t.ColumnConfigs {
			if !existColumns[c.Name] {
				missingCols = append(missingCols, strconv.Quote(c.Name))
			}
		}
	}

	if len(missingCols) > 0 {
		return common.NewValidationError("file", "Fields ["+strings.Join(missingCols, ",")+"] in file are missing.")
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek file to start: %w", err)
	}

	return nil
}

func (w *importWorker) createTable(ctx context.Context, job ImportJob, columns []string) error {
	switch job.Type {
	case externaldata.TypeMongoDB:
		// clean collection if job is retried
		_, err := w.dbClient.Collection(job.getDBTableName()).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("failed to clean collection: %w", err)
		}

		return nil
	case externaldata.TypePostgreSQL:
		sql := "CREATE TABLE IF NOT EXISTS " + job.getDBTableName() + " ( " +
			externaldata.IDColumnName + " VARCHAR(" + externaldata.MaxIDLenStr + ") PRIMARY KEY, "
		for i, field := range columns {
			sql += pgx.Identifier{field}.Sanitize() + " VARCHAR(" + externaldata.MaxStringLenStr + ") "
			if i != len(columns)-1 {
				sql += ","
			}
		}

		sql += ")"
		pgPool, err := w.pgPoolProvider.Get(ctx)
		if err != nil {
			return fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.Exec(ctx, sql)
		if err != nil {
			return fmt.Errorf("failed to create postgres table: %w", err)
		}

		// clean table if job is retried
		_, err = pgPool.Exec(ctx, "TRUNCATE TABLE "+job.getDBTableName())
		if err != nil {
			return fmt.Errorf("failed to truncate postgres table: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("invalid job type: %q", job.Type)
	}
}

func (w *importWorker) deleteTable(ctx context.Context, job ImportJob) error {
	switch job.Type {
	case externaldata.TypeMongoDB:
		err := w.dbClient.Collection(job.getDBTableName()).Drop(ctx)
		if err != nil {
			return fmt.Errorf("failed to drop mongo collection: %w", err)
		}

		return nil
	case externaldata.TypePostgreSQL:
		pgPool, err := w.pgPoolProvider.Get(ctx)
		if err != nil {
			return fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.Exec(ctx, "DROP TABLE IF EXISTS "+job.getDBTableName())
		if err != nil {
			return fmt.Errorf("failed to create postgres table: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("invalid job type: %q", job.Type)
	}
}
