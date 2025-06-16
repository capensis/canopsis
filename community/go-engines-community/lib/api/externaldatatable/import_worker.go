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

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

const (
	filePerm os.FileMode = 0o770

	maxRetries = 10
)

type ImportWorker interface {
	CreateJob(ctx context.Context, id string, separator rune, f multipart.File) (_ ImportJob, resErr error)
	ProcessJob(ctx context.Context, id string) error
	GetJob(ctx context.Context, id string) (ImportJob, error)
	CompleteJob(ctx context.Context, id string, columnTypes []int) (bool, error)
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
}

func (w *importWorker) CreateJob(ctx context.Context, id string, separator rune, f multipart.File) (_ ImportJob, resErr error) {
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

	g, ctx := errgroup.WithContext(ctx)
	done := make(chan struct{})
	g.Go(func() (resErr error) {
		defer close(done)
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
		var columns []string
		var columnLengths []int
		var failReason string
		switch job.Type {
		case externaldata.TypeMongoDB:
			columns, failReason, err = w.writeToMongo(ctx, job, r)
		case externaldata.TypePostgreSQL:
			columns, columnLengths, failReason, err = w.writeToPostgres(ctx, job, r)
		}

		if err != nil {
			return err
		}

		var update bson.M
		if failReason == "" {
			update = bson.M{
				"status":         ImportStatusSucceeded,
				"columns":        columns,
				"column_lengths": columnLengths,
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

		rmFile = true
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

func (w *importWorker) CompleteJob(ctx context.Context, id string, columnTypes []int) (bool, error) {
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

	if len(columnTypes) != len(job.Columns) {
		return false, common.NewValidationError("column_types", "ColumnTypes must contain "+strconv.Itoa(len(job.Columns))+" items.")
	}

	table := externaldata.Table{}
	update := bson.M{
		"columns":      job.Columns,
		"column_types": columnTypes,
	}
	if len(job.ColumnLengths) > 0 {
		update["column_lengths"] = job.ColumnLengths
	}

	err = w.dbCollection.FindOneAndUpdate(ctx,
		bson.M{"_id": job.ExternalDataTable},
		bson.M{"$set": update},
	).Decode(&table)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, fmt.Errorf("failed to find table: %w", err)
	}

	switch table.Type {
	case externaldata.TypeMongoDB:
		cursor, err := w.dbClient.Collection(job.getDBTableName()).Aggregate(ctx, []bson.M{
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

		existedCols := make(map[string]int, len(table.Columns))
		for i, c := range table.Columns {
			existedCols[c] = table.ColumnLengths[i]
		}

		for i, c := range job.Columns {
			newL := job.ColumnLengths[i]
			sql := ""
			if l, ok := existedCols[c]; !ok {
				sql = "ALTER TABLE " + table.GetDBName() +
					" ADD COLUMN " + pgx.Identifier{c}.Sanitize() + " VARCHAR(" + strconv.Itoa(newL) + ")"
			} else if l != newL {
				sql = "ALTER TABLE " + table.GetDBName() +
					" ALTER COLUMN " + pgx.Identifier{c}.Sanitize() + " TYPE VARCHAR(" + strconv.Itoa(newL) + ")"
			}

			if sql != "" {
				_, err = pgPool.Exec(ctx, sql)
				if err != nil {
					return false, fmt.Errorf("failed to alter postgres table: %w", err)
				}
			}

			delete(existedCols, c)
		}

		for c := range existedCols {
			sql := "ALTER TABLE " + table.GetDBName() +
				" DROP COLUMN " + pgx.Identifier{c}.Sanitize()
			_, err = pgPool.Exec(ctx, sql)
			if err != nil {
				return false, fmt.Errorf("failed to alter postgres table: %w", err)
			}
		}

		sanColsWithID := make([]string, len(job.Columns)+1)
		sanColsWithID[0] = externaldata.IDColumnName
		for j, col := range job.Columns {
			sanColsWithID[j+1] = pgx.Identifier{col}.Sanitize()
		}

		_, err = pgPool.Exec(ctx, "INSERT INTO "+table.GetDBName()+"("+strings.Join(sanColsWithID, ",")+
			") SELECT "+strings.Join(sanColsWithID, ",")+" FROM "+job.getDBTableName())
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

func (w *importWorker) writeToMongo(ctx context.Context, job ImportJob, r *csv.Reader) ([]string, string, error) {
	var columns []string
	i := 0
	docs := make([]any, 0, canopsis.DefaultBulkSize)
	for {
		i++
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err.Error(), nil
		}

		if len(record) == 0 {
			return nil, "empty record " + strconv.Itoa(i), nil
		}

		if len(columns) == 0 {
			columns = make([]string, len(record))
			copy(columns, record)
			err = w.createTable(ctx, job, columns, nil)
			if err != nil {
				return nil, "", err
			}

			continue
		}

		if len(columns) != len(record) {
			return nil, "invalid record " + strconv.Itoa(i) + ": not match fields", nil
		}

		doc := make(map[string]string, len(record)+1)
		doc[externaldata.IDColumnName] = utils.NewID()
		for j, c := range columns {
			doc[c] = record[j]
		}

		docs = append(docs, doc)
		if len(docs) == canopsis.DefaultBulkSize {
			_, err = w.dbClient.Collection(job.getDBTableName()).InsertMany(ctx, docs)
			if err != nil {
				return nil, "", fmt.Errorf("failed to insert into %q collection: %w", job.Table, err)
			}

			docs = docs[:0]
		}
	}

	if len(docs) > 0 {
		_, err := w.dbClient.Collection(job.getDBTableName()).InsertMany(ctx, docs)
		if err != nil {
			return nil, "", fmt.Errorf("failed to insert into %q collection: %w", job.Table, err)
		}
	}

	return columns, "", nil
}

func (w *importWorker) writeToPostgres(ctx context.Context, job ImportJob, r *csv.Reader) ([]string, []int, string, error) {
	pgPool, err := w.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get postgres pool: %w", err)
	}

	var columns []string
	var columnsWithID []string
	var columnLengths []int
	i := 0
	rows := make([][]any, 0, canopsis.DefaultBulkSize)
	for {
		i++
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, nil, err.Error(), nil
		}

		if len(record) == 0 {
			return nil, nil, "empty record " + strconv.Itoa(i), nil
		}

		if len(columns) == 0 {
			columnsWithID = make([]string, len(record)+1)
			columnsWithID[0] = externaldata.IDColumnName
			copy(columnsWithID[1:], record)
			columns = columnsWithID[1:]
			columnLengths = make([]int, len(columns))
			for j := range columns {
				columnLengths[j] = externaldata.PostgresDefaultColumnLen
			}

			err = w.createTable(ctx, job, columns, columnLengths)
			if err != nil {
				return nil, nil, "", err
			}

			continue
		}

		if len(columns) != len(record) {
			return nil, nil, "invalid record " + strconv.Itoa(i) + ": not match fields", nil
		}

		columnLengths, err = w.alterColumnLengths(ctx, job.getDBTableName(), columns, columnLengths, record)
		if err != nil {
			return nil, nil, "", err
		}

		row := make([]any, len(record)+1)
		row[0] = utils.NewID()
		for j, v := range record {
			row[j+1] = v
		}

		rows = append(rows, row)

		if len(rows) == canopsis.DefaultBulkSize {
			_, err = pgPool.CopyFrom(ctx, externaldata.GetPostgresTableIdentifier(job.Table), columnsWithID, pgx.CopyFromRows(rows))
			if err != nil {
				return nil, nil, "", fmt.Errorf("failed to insert into %q table: %w", job.Table, err)
			}

			rows = rows[:0]
		}
	}

	if len(rows) > 0 {
		_, err = pgPool.CopyFrom(ctx, externaldata.GetPostgresTableIdentifier(job.Table), columnsWithID, pgx.CopyFromRows(rows))
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to insert into %q table: %w", job.Table, err)
		}
	}

	return columns, columnLengths, "", nil
}

func (w *importWorker) validateColumns(ctx context.Context, t externaldata.Table, f multipart.File, separator rune) error {
	isLinked := false
	var err error
	if len(t.Columns) > 0 {
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
		for _, c := range t.Columns {
			if !existColumns[c] {
				missingCols = append(missingCols, strconv.Quote(c))
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

func (w *importWorker) createTable(ctx context.Context, job ImportJob, columns []string, columnLens []int) error {
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
			externaldata.IDColumnName + " VARCHAR(" + strconv.Itoa(externaldata.PostgresIDColumnLen) + ") PRIMARY KEY, "
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

func (w *importWorker) alterColumnLengths(
	ctx context.Context,
	tableName string,
	columns []string,
	columnLens []int,
	record []string,
) ([]int, error) {
	updatedLengths := make(map[string]int)
	for i, v := range record {
		if len(v) > columnLens[i] {
			columnLens[i] = len(v)
			updatedLengths[columns[i]] = columnLens[i]
		}
	}

	if len(updatedLengths) == 0 {
		return columnLens, nil
	}

	pgPool, err := w.pgPoolProvider.Get(ctx)
	if err != nil {
		return columnLens, fmt.Errorf("failed to get postgres pool: %w", err)
	}

	for c, l := range updatedLengths {
		sql := "ALTER TABLE " + tableName +
			" ALTER COLUMN " + pgx.Identifier{c}.Sanitize() + " TYPE VARCHAR(" + strconv.Itoa(l) + ")"
		_, err = pgPool.Exec(ctx, sql)
		if err != nil {
			return columnLens, fmt.Errorf("failed to alter postgres table: %w", err)
		}
	}

	return columnLens, nil
}
