package externaldata

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	pgErrCodeDuplicateTable     = "42P07"
	mongoErrCodeNamespaceExists = 48
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindOne(ctx context.Context, id string) (Response, error)
	Create(ctx context.Context, r EditRequest) (Response, error)
	Update(ctx context.Context, r EditRequest) (Response, error)
	Delete(ctx context.Context, id, author string) (bool, error)
	FindData(ctx context.Context, tableName string, tableType int, columns []string, r ListDataRequest) (*AggregationDataResult, error)
	FindOneData(ctx context.Context, tableID, id string) (map[string]any, error)
	CreateData(ctx context.Context, tableID string, r map[string]string) (map[string]string, error)
	UpdateData(ctx context.Context, tableID, id string, r map[string]string) (map[string]string, error)
	DeleteData(ctx context.Context, table Response, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, pgPoolProvider postgres.PoolProvider) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.ExternalDataTableCollection),
		pgPoolProvider:        pgPoolProvider,
		defaultSearchByFields: []string{"name", "description"},
		defaultSortBy:         "name",
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	pgPoolProvider        postgres.PoolProvider
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) FindOne(ctx context.Context, id string) (Response, error) {
	response := Response{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&response)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return response, fmt.Errorf("failed to get external data table: %w", err)
	}

	return response, nil
}

func (s *store) Create(ctx context.Context, r EditRequest) (Response, error) {
	if r.Type == nil {
		return Response{}, errors.New("type is required")
	}

	var err error
	switch *r.Type {
	case externaldata.TypeMongoDB:
		err = s.createMongoCollection(ctx, r.Name)
	case externaldata.TypePostgreSQL:
		err = s.createPostgresTable(ctx, r.Name)
	default:
		err = fmt.Errorf("unknown external data type %d", r.Type)
	}

	if err != nil {
		return Response{}, err
	}

	id := utils.NewID()
	response := Response{}
	now := datetime.NewCpsTime()
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = Response{}
		_, err = s.dbCollection.InsertOne(ctx, externaldata.Table{
			ID:          id,
			Type:        *r.Type,
			Name:        r.Name,
			Description: r.Description,
			Author:      r.Author,
			Created:     now,
			Updated:     now,
		})
		if err != nil {
			return fmt.Errorf("failed to insert external data table: %w", err)
		}

		response, err = s.FindOne(ctx, id)

		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r EditRequest) (Response, error) {
	res := Response{}
	if r.Type == nil {
		return res, errors.New("type is required")
	}

	oldTable, err := s.FindOne(ctx, r.ID)
	if err != nil || oldTable.ID == "" {
		return res, err
	}

	if oldTable.Type != *r.Type {
		return res, common.NewValidationError("type", "Type cannot be changed.")
	}

	if oldTable.Name != r.Name {
		if oldTable.FromConfig {
			return res, common.NewValidationError("name", "Name cannot be changed.")
		}

		switch oldTable.Type {
		case externaldata.TypeMongoDB:
			err = s.dbClient.RunAdminCommand(ctx, bson.D{
				{Key: "renameCollection", Value: s.dbClient.Name() + "." + oldTable.getDBTableName()},
				{Key: "to", Value: s.dbClient.Name() + "." + externaldata.GetMongoCollectionName(r.Name, false)},
			}).Err()
			if err != nil {
				commErr := mongodriver.CommandError{}
				if errors.As(err, &commErr) && commErr.Code == mongoErrCodeNamespaceExists {
					return res, common.NewValidationError("name", "MongoDB collection already exists.")
				}

				return res, fmt.Errorf("failed to rename mongo collection: %w", err)
			}
		case externaldata.TypePostgreSQL:
			pgPool, err := s.pgPoolProvider.Get(ctx)
			if err != nil {
				return res, fmt.Errorf("failed to get postgres pool: %w", err)
			}

			_, err = pgPool.Exec(ctx, "ALTER TABLE "+oldTable.getDBTableName()+" RENAME TO "+pgx.Identifier{r.Name}.Sanitize())
			if err != nil {
				pgErr := &pgconn.PgError{}
				if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeDuplicateTable {
					return res, common.NewValidationError("name", "PostgreSQL table already exists.")
				}

				return res, fmt.Errorf("failed to rename postgres table: %w", err)
			}
		default:
			return res, fmt.Errorf("invalid table type: %q", oldTable.Type)
		}
	}

	now := datetime.NewCpsTime()
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = Response{}

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": externaldata.Table{
			Type:        *r.Type,
			Name:        r.Name,
			Description: r.Description,
			Author:      r.Author,
			Updated:     now,
		}})
		if err != nil {
			return fmt.Errorf("failed to update external data table: %w", err)
		}

		res, err = s.FindOne(ctx, r.ID)

		return err
	})

	return res, err
}

func (s *store) Delete(ctx context.Context, id, author string) (bool, error) {
	var res externaldata.Table
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = externaldata.Table{}

		table := externaldata.Table{}
		err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&table)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		if table.FromConfig {
			return ErrConfigNotDeletable
		}

		// required to get the author in action log listener.
		ur, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id, "from_config": bson.M{"$ne": true}}, bson.M{"$set": bson.M{"author": author}})
		if err != nil || ur.MatchedCount == 0 {
			return err
		}

		d, err := s.dbCollection.DeleteOne(ctx, bson.M{"_id": id, "from_config": bson.M{"$ne": true}})
		if err != nil || d == 0 {
			return err
		}

		res = table

		return nil
	})
	if err != nil || res.ID == "" {
		return false, err
	}

	switch res.Type {
	case externaldata.TypeMongoDB:
		err = s.deleteMongoCollection(ctx, res.Name)
	case externaldata.TypePostgreSQL:
		err = s.deletePostgresTable(ctx, res.Name)
	default:
		err = fmt.Errorf("unknown external data type %d", res.Type)
	}

	return err == nil, err
}

func (s *store) FindData(ctx context.Context, tableName string, tableType int, columns []string, r ListDataRequest) (*AggregationDataResult, error) {
	foundSort := false
	for _, col := range columns {
		if r.SortBy == col {
			foundSort = true
			break
		}
	}

	if r.SortBy != "" && !foundSort {
		return nil, common.NewValidationError("sort_by", "SortBy must be one of ["+strings.Join(columns, " ")+"].")
	}

	var res *AggregationDataResult
	var err error
	switch tableType {
	case externaldata.TypeMongoDB:
		res, err = s.findDataFromMongo(ctx, tableName, columns, r)
	case externaldata.TypePostgreSQL:
		res, err = s.findDataFromPostgres(ctx, tableName, columns, r)
	default:
		err = fmt.Errorf("unknown external data type %d", tableType)
	}

	return res, err
}

func (s *store) FindOneData(ctx context.Context, tableID, id string) (map[string]any, error) {
	table, err := s.FindOne(ctx, tableID)
	if err != nil || table.ID == "" {
		return nil, err
	}

	var res map[string]any
	switch table.Type {
	case externaldata.TypeMongoDB:
		err = s.dbClient.Collection(table.getDBTableName()).FindOne(ctx, bson.M{"_id": id}).Decode(&res)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, err
		}

		return res, nil
	case externaldata.TypePostgreSQL:
		pgPool, err := s.pgPoolProvider.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		sql := "SELECT "
		columnsWithID := make([]string, len(table.ColumnTypes)+1)
		i := 0
		for col := range table.ColumnTypes {
			columnsWithID[i] = col
			i++
			sql += pgx.Identifier{col}.Sanitize() + ", "
		}

		columnsWithID[i] = externaldata.IDColumnName
		sql += externaldata.IDColumnName
		sql += " FROM " + table.getDBTableName() + " WHERE " + externaldata.IDColumnName + " = $1"
		rows, err := pgPool.Query(ctx, sql, id)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		if rows.Next() {
			vals, err := rows.Values()
			if err != nil {
				return nil, err
			}

			res, err = s.transformPostgresResToData(vals, columnsWithID)
			if err != nil {
				return nil, err
			}
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return res, nil
	default:
		return nil, fmt.Errorf("unknown external data type %d", table.Type)
	}
}

func (s *store) CreateData(ctx context.Context, tableID string, r map[string]string) (map[string]string, error) {
	table, err := s.FindOne(ctx, tableID)
	if err != nil || table.ID == "" {
		return nil, err
	}

	doc := make(map[string]string, len(table.ColumnTypes)+1)
	row := make([]any, len(table.ColumnTypes)+1)
	columnsWithID := make([]string, len(table.ColumnTypes)+1)
	i := 0
	doc[externaldata.IDColumnName] = utils.NewID()
	row[i] = doc[externaldata.IDColumnName]
	columnsWithID[i] = externaldata.IDColumnName
	i++
	var ok bool
	updatedLengths := make(map[string]int)
	valErrMsgs := make(map[string]string)
	for col := range table.ColumnTypes {
		doc[col], ok = r[col]
		if !ok {
			valErrMsgs[col] = col + " is missing."
			continue
		}

		row[i] = doc[col]
		columnsWithID[i] = col
		i++

		if l, ok := table.ColumnLengths[col]; ok && len(doc[col]) > l {
			updatedLengths[col] = len(doc[col])
		}
	}

	if len(valErrMsgs) > 0 {
		return nil, common.NewValidationErrors(valErrMsgs)
	}

	switch table.Type {
	case externaldata.TypeMongoDB:
		_, err = s.dbClient.Collection(table.getDBTableName()).InsertOne(ctx, doc)

		return doc, err
	case externaldata.TypePostgreSQL:
		err = s.alterPostgresColumns(ctx, table.getDBTableName(), updatedLengths)
		if err != nil {
			return nil, err
		}

		pgPool, err := s.pgPoolProvider.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.CopyFrom(ctx, table.getPostgresTableIdentifier(), columnsWithID, pgx.CopyFromRows([][]any{row}))
		if err != nil {
			return nil, err
		}

		err = s.updateColumnLengths(ctx, table, updatedLengths)

		return doc, err
	default:
		return nil, fmt.Errorf("unknown external data type %d", table.Type)
	}
}

func (s *store) UpdateData(ctx context.Context, tableID, id string, r map[string]string) (map[string]string, error) {
	table, err := s.FindOne(ctx, tableID)
	if err != nil || table.ID == "" {
		return nil, err
	}

	doc := make(map[string]string, len(table.ColumnTypes)+1)
	querySql := "UPDATE " + table.getDBTableName() + " SET "
	queryArgs := make([]any, len(table.ColumnTypes)+1)
	i := 0
	var ok bool
	updatedLengths := make(map[string]int)
	valErrMsgs := make(map[string]string)
	for col := range table.ColumnTypes {
		doc[col], ok = r[col]
		if !ok {
			valErrMsgs[col] = col + " is missing."
			continue
		}

		queryArgs[i] = doc[col]
		querySql += pgx.Identifier{col}.Sanitize() + " = $" + strconv.Itoa(i+1)
		if i < len(table.ColumnTypes)-1 {
			querySql += ", "
		}

		i++
		if l, ok := table.ColumnLengths[col]; ok && len(doc[col]) > l {
			updatedLengths[col] = len(doc[col])
		}
	}

	if len(valErrMsgs) > 0 {
		return nil, common.NewValidationErrors(valErrMsgs)
	}

	queryArgs[i] = id
	querySql += " WHERE " + externaldata.IDColumnName + " = $" + strconv.Itoa(len(queryArgs))
	switch table.Type {
	case externaldata.TypeMongoDB:
		updateRes, err := s.dbClient.Collection(table.getDBTableName()).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": doc})
		if err != nil || updateRes.MatchedCount == 0 {
			return nil, err
		}
	case externaldata.TypePostgreSQL:
		err = s.alterPostgresColumns(ctx, table.getDBTableName(), updatedLengths)
		if err != nil {
			return nil, err
		}

		pgPool, err := s.pgPoolProvider.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		execRes, err := pgPool.Exec(ctx, querySql, queryArgs...)
		if err != nil || execRes.RowsAffected() == 0 {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown external data type %d", table.Type)
	}

	doc[externaldata.IDColumnName] = id
	err = s.updateColumnLengths(ctx, table, updatedLengths)

	return doc, err
}

func (s *store) DeleteData(ctx context.Context, table Response, id string) (bool, error) {
	switch table.Type {
	case externaldata.TypeMongoDB:
		deleted, err := s.dbClient.Collection(table.getDBTableName()).DeleteOne(ctx, bson.M{"_id": id})

		return deleted > 0, err
	case externaldata.TypePostgreSQL:
		pgPool, err := s.pgPoolProvider.Get(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		sql := "DELETE FROM " + table.getDBTableName() + " WHERE " + externaldata.IDColumnName + " = $1"
		execRes, err := pgPool.Exec(ctx, sql, id)
		if err != nil {
			return false, err
		}

		return execRes.RowsAffected() > 0, nil
	default:
		return false, fmt.Errorf("unknown external data type %d", table.Type)
	}
}

func (s *store) createMongoCollection(ctx context.Context, name string) error {
	collName := externaldata.GetMongoCollectionName(name, false)
	collections, err := s.dbClient.ListCollectionNames(ctx, bson.M{"name": collName})
	if err != nil {
		return fmt.Errorf("failed to get collections: %w", err)
	}

	if len(collections) == 1 {
		return common.NewValidationError("name", "MongoDB collection already exists.")
	}

	err = s.dbClient.CreateCollection(ctx, collName)
	if err != nil {
		return fmt.Errorf("failed to create mongo collection: %w", err)
	}

	return nil
}

func (s *store) deleteMongoCollection(ctx context.Context, name string) error {
	err := s.dbClient.Collection(externaldata.GetMongoCollectionName(name, false)).Drop(ctx)
	if err != nil {
		return fmt.Errorf("failed to drop mongo collection: %w", err)
	}

	return nil
}

func (s *store) createPostgresTable(ctx context.Context, name string) error {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get postgres pool: %w", err)
	}

	sql := "CREATE TABLE " + externaldata.GetPostgresTableName(name) +
		" ( " + externaldata.IDColumnName + " VARCHAR(" + strconv.Itoa(externaldata.PostgresIDColumnLen) + ") PRIMARY KEY )"
	_, err = pgPool.Exec(ctx, sql)
	if err != nil {
		pgErr := &pgconn.PgError{}
		if errors.As(err, &pgErr) && pgErr.Code == pgErrCodeDuplicateTable {
			return common.NewValidationError("name", "PostgreSQL table already exists.")
		}

		return fmt.Errorf("failed to create postgres table: %w", err)
	}

	return nil
}

func (s *store) deletePostgresTable(ctx context.Context, name string) error {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get postgres pool: %w", err)
	}

	sql := "DROP TABLE IF EXISTS " + externaldata.GetPostgresTableName(name)
	_, err = pgPool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed to delete postgres table: %w", err)
	}

	return nil
}

func (s *store) alterPostgresColumns(ctx context.Context, name string, columnLengths map[string]int) error {
	if len(columnLengths) == 0 {
		return nil
	}

	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get postgres pool: %w", err)
	}

	for col, l := range columnLengths {
		sql := "ALTER TABLE " + name +
			" ALTER COLUMN " + pgx.Identifier{col}.Sanitize() + " TYPE VARCHAR(" + strconv.Itoa(l) + ")"
		_, err = pgPool.Exec(ctx, sql)
		if err != nil {
			return fmt.Errorf("failed to alter postgres table: %w", err)
		}
	}

	return nil
}

func (s *store) updateColumnLengths(ctx context.Context, table Response, columnLengths map[string]int) error {
	if len(columnLengths) == 0 {
		return nil
	}

	update := make(bson.M, len(columnLengths))
	for k, v := range columnLengths {
		update["column_lengths."+k] = v
	}

	_, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": table.ID}, bson.M{"$set": update})

	return err
}

func (s *store) findDataFromMongo(ctx context.Context, collectionName string, columns []string, request ListDataRequest) (*AggregationDataResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(request.Search, columns)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.dbClient.Collection(collectionName).Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(request.SortBy, externaldata.IDColumnName), request.Sort),
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var result AggregationDataResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *store) findDataFromPostgres(ctx context.Context, tableName string, columns []string, r ListDataRequest) (*AggregationDataResult, error) {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres pool: %w", err)
	}

	whereStmts := make([]string, 0)
	queryArgs := make([]any, 0)
	if r.Search != "" {
		orStmts := make([]string, len(columns))
		queryArgs = append(queryArgs, r.Search)
		for i, col := range columns {
			orStmts[i] = pgx.Identifier{col}.Sanitize() + " ~ $" + strconv.Itoa(len(queryArgs))
		}

		whereStmts = append(whereStmts, "("+strings.Join(orStmts, " OR ")+")")
	}

	whereStmt := ""
	if len(whereStmts) > 0 {
		whereStmt = "WHERE " + strings.Join(whereStmts, " AND ")
	}

	limitStmt := ""
	if r.Paginate {
		limitStmt = "OFFSET " + strconv.FormatInt(r.Limit*(r.Page-1), 10) +
			" LIMIT " + strconv.FormatInt(r.Limit, 10)
	}

	orderStmt := "ORDER BY " + pgx.Identifier{cmp.Or(r.SortBy, externaldata.IDColumnName)}.Sanitize()
	if r.Sort == mongo.SortDesc {
		orderStmt = orderStmt + " DESC"
	}

	columnsWithID := make([]string, len(columns)+1)
	columnsWithID[0] = externaldata.IDColumnName
	copy(columnsWithID[1:], columns)
	sql := "SELECT "
	for i, col := range columnsWithID {
		sql += pgx.Identifier{col}.Sanitize()
		if i < len(columnsWithID)-1 {
			sql += ", "
		}
	}

	sql += " FROM " + tableName + " " + whereStmt + " " + orderStmt + " " + limitStmt
	countSql := "SELECT count(*) FROM " + tableName + " " + whereStmt
	result := &AggregationDataResult{
		Data: make([]map[string]any, 0, r.Limit),
	}

	rows, err := pgPool.Query(ctx, sql, queryArgs...)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return result, err
		}

		row, err := s.transformPostgresResToData(vals, columnsWithID)
		if err != nil {
			return result, err
		}

		result.Data = append(result.Data, row)
	}

	if err = rows.Err(); err != nil {
		return result, err
	}

	err = pgPool.QueryRow(ctx, countSql, queryArgs...).Scan(&result.TotalCount)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *store) transformPostgresResToData(vals []any, columns []string) (map[string]any, error) {
	res := make(map[string]any, len(vals))
	var ok bool
	for i, val := range vals {
		res[columns[i]], ok = val.(string)
		if !ok {
			return res, fmt.Errorf("%q column doesn't contain string", columns[i])
		}
	}

	return res, nil
}
