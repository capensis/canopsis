package externaldata

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
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	filePerm os.FileMode = 0o770

	maxRetries = 10
)

type ImportWorker interface {
	CreateJob(ctx context.Context, id string, delimiter rune, f multipart.File, fh *multipart.FileHeader) (_ ImportJob, resErr error)
	ProcessJob(ctx context.Context, id string) error
	GetJob(ctx context.Context, id string) (ImportJob, error)
	CompleteJob(ctx context.Context, id string, columnTypes map[string]int) (bool, error)
	ProcessAbandonedJobs(ctx context.Context)
	DeleteOldJobs(ctx context.Context)
}

func NewImportWorker(
	dbClient mongo.DbClient,
	pgPoolProvider postgres.PoolProvider,
	tmpImportDir string,
	jobPublisher workers.JobPublisher,
	logger zerolog.Logger,
) ImportWorker {
	return &importWorker{
		dbClient:                dbClient,
		dbCollection:            dbClient.Collection(mongo.ExternalDataTableCollection),
		dbImportCollection:      dbClient.Collection(mongo.ExternalDataImportWorkerCollection),
		pgPoolProvider:          pgPoolProvider,
		tmpImportDir:            tmpImportDir,
		jobPublisher:            jobPublisher,
		logger:                  logger,
		abandonedTickerInterval: time.Minute,
		pingInterval:            time.Second,
		deleteTickerInterval:    time.Hour,
		deleteInterval:          24 * time.Hour,
	}
}

type importWorker struct {
	dbClient                mongo.DbClient
	dbCollection            mongo.DbCollection
	dbImportCollection      mongo.DbCollection
	pgPoolProvider          postgres.PoolProvider
	tmpImportDir            string
	jobPublisher            workers.JobPublisher
	abandonedTickerInterval time.Duration
	pingInterval            time.Duration
	deleteTickerInterval    time.Duration
	deleteInterval          time.Duration
	logger                  zerolog.Logger
}

func (w *importWorker) CreateJob(ctx context.Context, id string, delimiter rune, f multipart.File, fh *multipart.FileHeader) (_ ImportJob, resErr error) {
	defer func() {
		err := f.Close()
		if err != nil && resErr == nil {
			resErr = fmt.Errorf("cannot close file: %w", err)
		}
	}()

	job := ImportJob{}
	externalDataTable := Document{}
	err := w.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&externalDataTable)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return job, common.NewValidationError("_id", "ID doesn't exist.")
		}

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
		Delimiter:         delimiter,
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

	_, err = io.Copy(df, f)
	if err != nil {
		return job, fmt.Errorf("failed to copy to file %q: %w", job.Filepath, err)
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

func (w *importWorker) ProcessJob(ctx context.Context, id string) error {
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

	g, ctx := errgroup.WithContext(ctx)
	done := make(chan struct{})
	g.Go(func() error {
		defer close(done)
		f, err := os.Open(job.Filepath)
		if err != nil {
			return fmt.Errorf("failed to open file %q: %w", job.Filepath, err)
		}

		defer f.Close()
		r := csv.NewReader(f)
		r.Comma = job.Delimiter
		r.ReuseRecord = true
		var fields []string
		var columns []string
		var maxColumnValLens []int
		var failReason string
		i := 0
		docs := make([]any, 0, canopsis.DefaultBulkSize)
		rows := make([][]any, 0, canopsis.DefaultBulkSize)
		for {
			i++
			record, err := r.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				failReason = err.Error()
				break
			}

			if len(record) == 0 {
				failReason = "empty record " + strconv.Itoa(i)
				break
			}

			if len(fields) == 0 {
				fields = make([]string, len(record))
				copy(fields, record)
				columns = make([]string, len(record)+1)
				columns[0] = idColumnName
				copy(columns[1:], record)
				err = w.validateColumns(fields)
				if err != nil {
					failReason = err.Error()
					break
				}

				maxColumnValLens = make([]int, len(fields))
				for j := range fields {
					maxColumnValLens[j] = postgresDefaultColumnLen
				}

				err = w.createTable(ctx, job, fields, maxColumnValLens)
				if err != nil {
					return err
				}

				continue
			}

			if len(fields) != len(record) {
				failReason = "invalid record " + strconv.Itoa(i) + ": not match fields"
				break
			}

			maxColumnValLens, err = w.alterColumnLengths(ctx, job, fields, maxColumnValLens, record)
			if err != nil {
				return err
			}

			switch job.Type {
			case TypeMongoDB:
				doc := make(map[string]string, len(record)+1)
				doc[idColumnName] = utils.NewID()
				for j, field := range fields {
					doc[field] = record[j]
				}

				docs = append(docs, doc)
			case TypePostgreSQL:
				row := make([]any, len(record)+1)
				row[0] = utils.NewID()
				for j, v := range record {
					row[j+1] = v
				}

				rows = append(rows, row)
			}

			if len(docs) == canopsis.DefaultBulkSize || len(rows) == canopsis.DefaultBulkSize {
				err = w.insertIntoTable(ctx, job, docs, columns, rows)
				if err != nil {
					return err
				}

				docs = docs[:0]
				rows = rows[:0]
			}
		}

		if len(docs) > 0 || len(rows) > 0 {
			err = w.insertIntoTable(ctx, job, docs, columns, rows)
			if err != nil {
				return err
			}
		}

		var update bson.M
		if failReason == "" {
			update = bson.M{
				"status":  ImportStatusSucceeded,
				"columns": fields,
			}
		} else {
			w.logger.Err(errors.New(failReason)).Str("job", job.ID).Msg("failed to import external data")
			update = bson.M{
				"status":      ImportStatusFailed,
				"fail_reason": failReason,
			}
		}

		updateRes, err := w.dbImportCollection.UpdateOne(ctx,
			bson.M{
				"_id":    job.ID,
				"status": ImportStatusRunning,
			},
			bson.M{"$set": update},
		)
		if err != nil {
			return fmt.Errorf("failed to update import job: %w", err)
		}

		if updateRes.ModifiedCount == 0 {
			return errors.New("import job is processing by another worker")
		}

		err = w.deleteJobFile(job)
		if err != nil {
			return err
		}

		if failReason != "" {
			err = w.deleteTable(ctx, job)
			if err != nil {
				return err
			}
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

func (w *importWorker) CompleteJob(ctx context.Context, id string, columnTypes map[string]int) (bool, error) {
	job := ImportJob{}
	err := w.dbImportCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, fmt.Errorf("failed to find import job: %w", err)
	}

	if job.Status != ImportStatusSucceeded {
		return false, nil
	}

	for _, col := range job.Columns {
		t, ok := columnTypes[col]
		if !ok {
			return false, common.NewValidationError("columns."+col, "Column is required.")
		}

		switch t {
		case ColumnTypeNoType, ColumnTypeFilter, ColumnTypeContext:
		default:
			return false, common.NewValidationError("columns."+col, "Column must be one of ["+
				strconv.Itoa(ColumnTypeNoType)+" "+
				strconv.Itoa(ColumnTypeFilter)+" "+
				strconv.Itoa(ColumnTypeContext)+
				"].")
		}
	}

	if len(columnTypes) > len(job.Columns) {
		return false, common.NewValidationError("columns", "Columns is invalid.")
	}

	table := Document{}
	err = w.dbCollection.FindOneAndUpdate(ctx,
		bson.M{"_id": job.ExternalDataTable},
		bson.M{"$set": bson.M{"columns": columnTypes}},
	).Decode(&table)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, fmt.Errorf("failed to find table: %w", err)
	}

	switch table.Type {
	case TypeMongoDB:
		err = w.dbClient.RunAdminCommand(ctx, bson.D{
			{Key: "renameCollection", Value: w.dbClient.Name() + "." + getCollectionName(job.Table)},
			{Key: "to", Value: w.dbClient.Name() + "." + getCollectionName(table.Name)},
			{Key: "dropTarget", Value: true},
		}).Err()
		if err != nil {
			return false, fmt.Errorf("failed to rename mongo collection: %w", err)
		}
	case TypePostgreSQL:
		pgPool, err := w.pgPoolProvider.Get(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.Exec(ctx, "DROP TABLE IF EXISTS "+getPostgresTableName(table.Name))
		if err != nil {
			return false, fmt.Errorf("failed to drop postgres table: %w", err)
		}

		_, err = pgPool.Exec(ctx, "ALTER TABLE "+getPostgresTableName(job.Table)+" RENAME TO "+pgx.Identifier{table.Name}.Sanitize())
		if err != nil {
			return false, fmt.Errorf("failed to rename postgres table: %w", err)
		}
	default:
		return false, fmt.Errorf("invalid table type: %q", table.Type)
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

func (w *importWorker) deleteJobFile(job ImportJob) error {
	err := os.Remove(job.Filepath)
	if err != nil {
		return fmt.Errorf("failed to remove file %q: %w", job.Filepath, err)
	}

	return nil
}

func (w *importWorker) validateColumns(columns []string) error {
	for _, c := range columns {
		if !common.IsTableName(c) || c == idColumnName {
			return fmt.Errorf("invalid field name: %q", c)
		}
	}

	return nil
}

func (w *importWorker) createTable(ctx context.Context, job ImportJob, columns []string, columnLens []int) error {
	switch job.Type {
	case TypeMongoDB:
		// clean collection if job is retried
		_, err := w.dbClient.Collection(getCollectionName(job.Table)).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("failed to clean collection: %w", err)
		}

		return nil
	case TypePostgreSQL:
		sql := "CREATE TABLE IF NOT EXISTS " + getPostgresTableName(job.Table) + " ( " +
			idColumnName + " VARCHAR(" + strconv.Itoa(uuidLen) + ") PRIMARY KEY, "
		for i, field := range columns {
			sql += pgx.Identifier{field}.Sanitize() + " VARCHAR(" + strconv.Itoa(columnLens[i]) + ") "
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
		_, err = pgPool.Exec(ctx, "TRUNCATE TABLE "+getPostgresTableName(job.Table))
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
	case TypeMongoDB:
		err := w.dbClient.Collection(getCollectionName(job.Table)).Drop(ctx)
		if err != nil {
			return fmt.Errorf("failed to drop mongo collection: %w", err)
		}

		return nil
	case TypePostgreSQL:
		pgPool, err := w.pgPoolProvider.Get(ctx)
		if err != nil {
			return fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.Exec(ctx, "DROP TABLE IF EXISTS "+getPostgresTableName(job.Table))
		if err != nil {
			return fmt.Errorf("failed to create postgres table: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("invalid job type: %q", job.Type)
	}
}

func (w *importWorker) alterColumnLengths(
	ctx context.Context,
	job ImportJob,
	columns []string,
	columnLens []int,
	record []string,
) ([]int, error) {
	updatedLengths := make(map[string]int, len(columns))
	for i, v := range record {
		if len(v) > columnLens[i] {
			columnLens[i] = len(v)
			updatedLengths[columns[i]] = columnLens[i]
		}
	}

	if len(updatedLengths) == 0 {
		return columnLens, nil
	}

	switch job.Type {
	case TypeMongoDB:
		return columnLens, nil
	case TypePostgreSQL:
		pgPool, err := w.pgPoolProvider.Get(ctx)
		if err != nil {
			return columnLens, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		for c, l := range updatedLengths {
			sql := "ALTER TABLE " + getPostgresTableName(job.Table) +
				" ALTER COLUMN " + pgx.Identifier{c}.Sanitize() + " TYPE VARCHAR(" + strconv.Itoa(l) + ")"
			_, err = pgPool.Exec(ctx, sql)
			if err != nil {
				return columnLens, fmt.Errorf("failed to alter postgres table: %w", err)
			}
		}

		return columnLens, nil
	default:
		return columnLens, fmt.Errorf("invalid job type: %q", job.Type)
	}
}

func (w *importWorker) insertIntoTable(
	ctx context.Context,
	job ImportJob,
	docs []any,
	columns []string,
	rows [][]any,
) error {
	switch job.Type {
	case TypeMongoDB:
		_, err := w.dbClient.Collection(getCollectionName(job.Table)).InsertMany(ctx, docs)
		if err != nil {
			return fmt.Errorf("failed to insert into %q collection: %w", job.Table, err)
		}

		return nil
	case TypePostgreSQL:
		pgPool, err := w.pgPoolProvider.Get(ctx)
		if err != nil {
			return fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.CopyFrom(ctx, getPostgresTableIdentifier(job.Table), columns, pgx.CopyFromRows(rows))
		if err != nil {
			return fmt.Errorf("failed to insert into %q table: %w", job.Table, err)
		}

		return nil
	default:
		return fmt.Errorf("invalid job type: %q", job.Type)
	}
}
