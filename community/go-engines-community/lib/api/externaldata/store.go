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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const pgErrCodeDuplicateTable = "42P07"

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindOne(ctx context.Context, id string) (Response, error)
	Create(ctx context.Context, r CreateRequest) (Response, error)
	FindData(ctx context.Context, tableName string, tableType int, columns []string, r ListDataRequest) (*AggregationDataResult, error)
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

func (s *store) Find(ctx context.Context, request ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(request.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(request.SortBy, s.defaultSortBy), request.Sort),
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

func (s *store) Create(ctx context.Context, r CreateRequest) (Response, error) {
	if r.Type == nil {
		return Response{}, errors.New("type is required")
	}

	var err error
	switch *r.Type {
	case TypeMongoDB:
		err = s.createMongoCollection(ctx, r.Name)
	case TypePostgreSQL:
		err = s.createPostgresTable(ctx, r.Name)
	default:
		err = fmt.Errorf("unknown external data type %d", r.Type)
	}

	if err != nil {
		return Response{}, err
	}

	id := utils.NewID()
	response := Response{}
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = Response{}
		_, err = s.dbCollection.InsertOne(ctx, Document{
			ID:          id,
			Type:        *r.Type,
			Name:        r.Name,
			Description: r.Description,
			Author:      r.Author,
		})
		if err != nil {
			return fmt.Errorf("failed to insert external data table: %w", err)
		}

		response, err = s.FindOne(ctx, id)

		return err
	})

	return response, err
}

func (s *store) FindData(ctx context.Context, tableName string, tableType int, columns []string, request ListDataRequest) (*AggregationDataResult, error) {
	foundSort := false
	for _, col := range columns {
		if request.SortBy == col {
			foundSort = true
			break
		}
	}

	if request.SortBy != "" && !foundSort {
		return nil, common.NewValidationError("sort_by", "SortBy must be one of ["+strings.Join(columns, " ")+"].")
	}

	var res *AggregationDataResult
	var err error
	switch tableType {
	case TypeMongoDB:
		res, err = s.findDataFromMongo(ctx, tableName, columns, request)
	case TypePostgreSQL:
		res, err = s.findDataFromPostgres(ctx, tableName, columns, request)
	default:
		err = fmt.Errorf("unknown external data type %q", tableType)
	}

	return res, err
}

func (s *store) createMongoCollection(ctx context.Context, name string) error {
	collName := getCollectionName(name)
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

func (s *store) createPostgresTable(ctx context.Context, name string) error {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get postgres pool: %w", err)
	}

	sql := "CREATE TABLE " + getPostgresTableName(name) + " ( " + idColumnName + " VARCHAR(" + strconv.Itoa(uuidLen) + ") PRIMARY KEY )"
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

func (s *store) findDataFromMongo(ctx context.Context, collectionName string, columns []string, request ListDataRequest) (*AggregationDataResult, error) {
	pipeline := make([]bson.M, 0)
	filter := common.GetSearchQuery(request.Search, columns)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.dbClient.Collection(getCollectionName(collectionName)).Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		pipeline,
		common.GetSortQuery(cmp.Or(request.SortBy, idColumnName), request.Sort),
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

	return &result, nil
}

func (s *store) findDataFromPostgres(ctx context.Context, tableNameWOPrefix string, columns []string, request ListDataRequest) (*AggregationDataResult, error) {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres pool: %w", err)
	}

	tableName := getPostgresTableName(tableNameWOPrefix)
	whereStmts := make([]string, 0)
	queryArgs := make([]any, 0)
	if request.Search != "" {
		queryArgs = append(queryArgs, request.Search)
		for _, col := range columns {
			whereStmts = append(whereStmts, pgx.Identifier{col}.Sanitize()+" ~ $"+strconv.Itoa(len(queryArgs)))
		}
	}

	whereStmt := ""
	if len(whereStmts) > 0 {
		whereStmt = "WHERE " + strings.Join(whereStmts, " AND ")
	}

	limitStmt := ""
	if request.Paginate {
		limitStmt = "OFFSET " + strconv.FormatInt(request.Limit*(request.Page-1), 10) +
			" LIMIT " + strconv.FormatInt(request.Limit, 10)
	}

	orderStmt := "ORDER BY " + pgx.Identifier{cmp.Or(request.SortBy, idColumnName)}.Sanitize()
	if request.Sort == mongo.SortDesc {
		orderStmt = orderStmt + " DESC"
	}

	columnsWithID := []string{idColumnName}
	columnsWithID = append(columnsWithID, columns...)
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
		Data: make([]map[string]string, 0),
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

		row := make(map[string]string, len(vals))
		var ok bool
		for i, val := range vals {
			row[columnsWithID[i]], ok = val.(string)
			if !ok {
				return result, fmt.Errorf("%q column doesn't contain string", columnsWithID[i])
			}
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
