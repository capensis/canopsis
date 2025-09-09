package externaldatatable

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	pgErrCodeDuplicateTable     = "42P07"
	mongoErrCodeNamespaceExists = 48

	limitLinkedRules = 11
)

var linkedCollections = []string{
	mongo.EventFilterRuleCollection,
	mongo.LinkRuleMongoCollection,
}

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	FindOne(ctx context.Context, id string) (Table, error)
	Create(ctx context.Context, r EditRequest) (Table, error)
	Update(ctx context.Context, r UpdateRequest) (Table, error)
	Delete(ctx context.Context, id, author string) (bool, error)
	FindData(ctx context.Context, tableName string, tableType int, columnConfigs []externaldata.ColumnConfig, r ListDataRequest) (*AggregationDataResult, error)
	FindPreviewData(ctx context.Context, job ImportJob, r ListPreviewRequest) (*AggregationDataResult, error)
	FindOneData(ctx context.Context, tableID, id string) (map[string]any, error)
	CreateData(ctx context.Context, tableID string, r map[string]any) (map[string]any, error)
	UpdateData(ctx context.Context, tableID, id string, r map[string]any) (map[string]any, error)
	DeleteData(ctx context.Context, table Table, id string) (bool, error)
	Export(ctx context.Context, t export.Task) (export.DataCursor, error)
}

func NewStore(dbClient mongo.DbClient, pgPoolProvider postgres.PoolProvider, dbExportClient mongo.DbClient, exportDecoder encoding.Decoder) Store {
	linkedDbCollections := make([]mongo.DbCollection, len(linkedCollections))
	for i, c := range linkedCollections {
		linkedDbCollections[i] = dbClient.Collection(c)
	}

	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.ExternalDataTableCollection),
		dbWidgetCollection:    dbClient.Collection(mongo.WidgetMongoCollection),
		linkedDbCollections:   linkedDbCollections,
		pgPoolProvider:        pgPoolProvider,
		defaultSearchByFields: []string{"name", "description"},
		defaultSortBy:         "name",
		dbExportClient:        dbExportClient,
		exportDecoder:         exportDecoder,
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	dbWidgetCollection    mongo.DbCollection
	linkedDbCollections   []mongo.DbCollection
	pgPoolProvider        postgres.PoolProvider
	defaultSearchByFields []string
	defaultSortBy         string
	dbExportClient        mongo.DbClient
	exportDecoder         encoding.Decoder
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	if len(r.IDs) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": bson.M{"$in": r.IDs}}})
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	project := make([]bson.M, 0)
	if r.WithFlags {
		project = append(project,
			bson.M{"$lookup": bson.M{
				"from":         s.dbWidgetCollection.Name(),
				"localField":   "_id",
				"foreignField": "parameters.table",
				"as":           "tmp_linked_widgets",
				"pipeline": []bson.M{
					{"$match": bson.M{"type": view.WidgetTypeExternalData}},
					{"$limit": limitLinkedRules},
					{"$project": bson.M{
						"name": "$title",
					}},
				},
			}},
			bson.M{"$unwind": bson.M{"path": "$tmp_linked_widgets", "preserveNullAndEmptyArrays": true}},
			bson.M{"$sort": bson.M{"tmp_linked_widgets.name": 1}},
			bson.M{"$group": bson.M{
				"_id":  "$_id",
				"data": bson.M{"$first": "$$ROOT"},
				"tmp_linked_widgets": bson.M{"$push": bson.M{"$cond": bson.M{
					"if":   "$tmp_linked_widgets",
					"then": "$tmp_linked_widgets",
					"else": "$$REMOVE",
				}}},
			}},
			bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{"linked_rules": bson.M{"$mergeObjects": bson.A{
					"$data.linked_rules",
					bson.M{"widget": "$tmp_linked_widgets"},
				}}},
			}}}},
		)

		for _, c := range s.linkedDbCollections {
			k := strings.ReplaceAll(strings.ToLower(c.Name()), "_", "")
			project = append(project,
				bson.M{"$lookup": bson.M{
					"from":         c.Name(),
					"localField":   "_id",
					"foreignField": "external_data.table",
					"as":           "tmp_linked_rules",
					"pipeline": []bson.M{
						{"$limit": limitLinkedRules},
						{"$project": bson.M{
							"name": bson.M{"$cond": bson.M{
								"if":   "$name",
								"then": "$name",
								"else": "$description",
							}},
						}},
					},
				}},
				bson.M{"$unwind": bson.M{"path": "$tmp_linked_rules", "preserveNullAndEmptyArrays": true}},
				bson.M{"$sort": bson.M{"tmp_linked_rules" + ".name": 1}},
				bson.M{"$group": bson.M{
					"_id":  "$_id",
					"data": bson.M{"$first": "$$ROOT"},
					"tmp_linked_rules": bson.M{"$push": bson.M{"$cond": bson.M{
						"if":   "$tmp_linked_rules",
						"then": "$tmp_linked_rules",
						"else": "$$REMOVE",
					}}},
				}},
				bson.M{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
					"$data",
					bson.M{"linked_rules": bson.M{"$mergeObjects": bson.A{
						"$data.linked_rules",
						bson.M{k: "$tmp_linked_rules"},
					}}},
				}}}},
			)
		}
	}

	sort := common.GetSortQuery(cmp.Or(r.SortBy, s.defaultSortBy), r.Sort)
	if len(project) > 0 {
		project = append(project, sort)
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		sort,
		project,
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

func (s *store) FindOne(ctx context.Context, id string) (Table, error) {
	response := Table{}

	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&response)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return response, fmt.Errorf("failed to get external data table: %w", err)
	}

	return response, nil
}

func (s *store) Create(ctx context.Context, r EditRequest) (Table, error) {
	if r.Type == nil {
		return Table{}, errors.New("type is required")
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
		return Table{}, err
	}

	id := utils.NewID()
	response := Table{}
	now := datetime.NewCpsTime()
	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		response = Table{}
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

func (s *store) Update(ctx context.Context, r UpdateRequest) (Table, error) {
	res := Table{}
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

	if len(oldTable.ColumnConfigs) != len(r.ColumnTags) {
		return res, common.NewValidationError("column_tags", "ColumnTags must contain "+strconv.Itoa(len(oldTable.ColumnConfigs))+" items.")
	}

	updatedColumnConfigs := make([]externaldata.ColumnConfig, 0, len(oldTable.ColumnConfigs))
	for i, cfg := range oldTable.ColumnConfigs {
		cfg.Tag = &r.ColumnTags[i]
		updatedColumnConfigs = append(updatedColumnConfigs, cfg)
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
		res = Table{}

		_, err = s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": externaldata.Table{
			Type:          *r.Type,
			Name:          r.Name,
			Description:   r.Description,
			ColumnConfigs: updatedColumnConfigs,
			Author:        r.Author,
			Updated:       now,
		}})
		if err != nil {
			return fmt.Errorf("failed to update external data table: %w", err)
		}

		res, err = s.FindOne(ctx, r.ID)
		if err != nil {
			return err
		}

		if oldTable.Name != res.Name {
			return s.updateLinkedModels(ctx, res.ID, res.getDBTableName())
		}

		return nil
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

		isLinked, err := isTableLinked(ctx, id, s.dbWidgetCollection, s.linkedDbCollections)
		if err != nil {
			return err
		}

		if isLinked {
			return ErrLinkedNotDeletable
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

func (s *store) FindPreviewData(ctx context.Context, job ImportJob, r ListPreviewRequest) (res *AggregationDataResult, err error) {
	switch job.Type {
	case externaldata.TypeMongoDB:
		res, err = s.findPreviewDataFromMongo(ctx, job, r)
	case externaldata.TypePostgreSQL:
		res, err = s.findPreviewDataFromPostgres(ctx, job, r)
	default:
		err = fmt.Errorf("unknown external data type %d", job.Type)
	}

	return res, err
}

func (s *store) FindData(ctx context.Context, tableName string, tableType int, columnConfigs []externaldata.ColumnConfig, r ListDataRequest) (res *AggregationDataResult, err error) {
	columns := make([]string, len(columnConfigs))
	for idx, cfg := range columnConfigs {
		columns[idx] = cfg.Name
	}

	valErrMsgs := make(map[string]string)
	if len(r.SearchBy) > 0 || r.SortBy != "" {
		hasCol := make(map[string]bool, len(columns))
		for _, c := range columns {
			hasCol[c] = true
		}

		if r.SortBy != "" && !hasCol[r.SortBy] {
			valErrMsgs["sort_by"] = "SortBy must be one of [" + strings.Join(columns, " ") + "]."
		}

		for i, v := range r.SearchBy {
			if !hasCol[v] {
				valErrMsgs["search_by."+strconv.Itoa(i)] = "SearchBy must be one of [" + strings.Join(columns, " ") + "]."
			}
		}
	}

	if len(valErrMsgs) > 0 {
		return nil, common.NewValidationErrors(valErrMsgs)
	}

	searchBy := r.SearchBy
	if len(searchBy) == 0 {
		searchBy = columns
	}

	switch tableType {
	case externaldata.TypeMongoDB:
		res, err = s.findDataFromMongo(ctx, tableName, searchBy, r)
	case externaldata.TypePostgreSQL:
		res, err = s.findDataFromPostgres(ctx, tableName, searchBy, r, columnConfigs)
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

		columns := table.getColumns()

		sql := "SELECT "
		columnsWithID := make([]string, len(columns)+1)
		columnsWithID[0] = externaldata.IDColumnName
		copy(columnsWithID[1:], columns)
		for i, col := range columnsWithID {
			sql += pgx.Identifier{col}.Sanitize()
			if i < len(columnsWithID)-1 {
				sql += ", "
			}
		}

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

			res, err = s.transformPostgresResToData(vals, table.ColumnConfigs)
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

func (s *store) CreateData(ctx context.Context, tableID string, r map[string]any) (map[string]any, error) {
	table, err := s.FindOne(ctx, tableID)
	if err != nil || table.ID == "" {
		return nil, err
	}

	columnsWithID := make([]string, len(table.ColumnConfigs)+1)
	columnsWithID[0] = externaldata.IDColumnName

	doc := make(map[string]any, len(table.ColumnConfigs))
	doc[externaldata.IDColumnName] = utils.NewID()

	row := make([]any, len(columnsWithID))
	row[0] = doc[externaldata.IDColumnName]

	valErrMsgs := make(map[string]string)
	for i, cfg := range table.ColumnConfigs {
		columnName := cfg.Name
		columnsWithID[i+1] = columnName

		rawVal, ok := r[columnName]
		if !ok {
			valErrMsgs[columnName] = columnName + " is missing."
			continue
		}

		var val any

		switch cfg.Type {
		case externaldata.ColumnTypeString:
			strVal, ok := rawVal.(string)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a string."
				continue
			}

			if len(strVal) > MaxStringLen {
				valErrMsgs[columnName] = "string length must be less than " + MaxStringLenStr
				continue
			}

			val = strVal
		case externaldata.ColumnTypeNumber:
			val, ok = rawVal.(float64)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a number."
				continue
			}
		case externaldata.ColumnTypeBoolean:
			val, ok = rawVal.(bool)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a boolean."
				continue
			}
		case externaldata.ColumnTypeStringArray:
			val, ok = utils.IsStringSlice(rawVal)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a string array."
				continue
			}
		case externaldata.ColumnTypeDateTime, externaldata.ColumnTypeTimestamp:
			val, ok = getIntValue(rawVal)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a timestamp."
				continue
			}
		default:
			return nil, fmt.Errorf("unexpected column type: %d", cfg.Type)
		}

		doc[columnName] = val
		row[i+1] = val
	}

	if len(valErrMsgs) > 0 {
		return nil, common.NewValidationErrors(valErrMsgs)
	}

	switch table.Type {
	case externaldata.TypeMongoDB:
		_, err = s.dbClient.Collection(table.getDBTableName()).InsertOne(ctx, doc)

		return doc, err
	case externaldata.TypePostgreSQL:
		pgPool, err := s.pgPoolProvider.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		_, err = pgPool.CopyFrom(ctx, table.getPostgresTableIdentifier(), columnsWithID, pgx.CopyFromRows([][]any{row}))
		if err != nil {
			return nil, err
		}

		return doc, err
	default:
		return nil, fmt.Errorf("unknown external data type %d", table.Type)
	}
}

func (s *store) UpdateData(ctx context.Context, tableID, id string, r map[string]any) (map[string]any, error) {
	table, err := s.FindOne(ctx, tableID)
	if err != nil || table.ID == "" {
		return nil, err
	}

	doc := make(map[string]any, len(table.ColumnConfigs)+1)
	querySql := "UPDATE " + table.getDBTableName() + " SET "
	queryArgs := make([]any, len(table.ColumnConfigs)+1)

	valErrMsgs := make(map[string]string)

	for i, cfg := range table.ColumnConfigs {
		columnName := cfg.Name

		rawVal, ok := r[columnName]
		if !ok {
			valErrMsgs[columnName] = columnName + " is missing."
			continue
		}

		var val any

		switch cfg.Type {
		case externaldata.ColumnTypeString:
			strVal, ok := rawVal.(string)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a string."
				continue
			}

			if len(strVal) > MaxStringLen {
				valErrMsgs[columnName] = "string length must be less than " + MaxStringLenStr
				continue
			}

			val = strVal
		case externaldata.ColumnTypeNumber:
			val, ok = rawVal.(float64)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a number."
				continue
			}
		case externaldata.ColumnTypeBoolean:
			val, ok = rawVal.(bool)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a boolean."
				continue
			}
		case externaldata.ColumnTypeStringArray:
			val, ok = utils.IsStringSlice(rawVal)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a string array."
				continue
			}
		case externaldata.ColumnTypeDateTime, externaldata.ColumnTypeTimestamp:
			val, ok = getIntValue(rawVal)
			if !ok {
				valErrMsgs[columnName] = columnName + " is not a timestamp."
				continue
			}
		default:
			return nil, fmt.Errorf("unexpected column type: %d", cfg.Type)
		}

		doc[columnName] = val

		queryArgs[i] = val
		querySql += pgx.Identifier{columnName}.Sanitize() + " = $" + strconv.Itoa(i+1)
		if i < len(table.ColumnConfigs)-1 {
			querySql += ", "
		}
	}

	if len(valErrMsgs) > 0 {
		return nil, common.NewValidationErrors(valErrMsgs)
	}

	queryArgs[len(queryArgs)-1] = id
	querySql += " WHERE " + externaldata.IDColumnName + " = $" + strconv.Itoa(len(queryArgs))
	switch table.Type {
	case externaldata.TypeMongoDB:
		updateRes, err := s.dbClient.Collection(table.getDBTableName()).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": doc})
		if err != nil || updateRes.MatchedCount == 0 {
			return nil, err
		}
	case externaldata.TypePostgreSQL:
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

	return doc, err
}

func (s *store) DeleteData(ctx context.Context, table Table, id string) (bool, error) {
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

func (s *store) Export(ctx context.Context, t export.Task) (export.DataCursor, error) {
	r := ExportFetchParameters{}
	err := s.exportDecoder.Decode([]byte(t.Parameters), &r)
	if err != nil {
		return nil, err
	}

	table, err := s.FindOne(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	if table.ID == "" {
		return nil, fmt.Errorf("table not found id=%q", r.ID)
	}

	columns := table.getColumns()

	searchBy := r.SearchBy
	if len(searchBy) == 0 {
		searchBy = columns
	}

	var selectColumns []string
	if len(t.Fields) > 0 {
		selectColumns = make([]string, len(t.Fields))
		for i, f := range t.Fields {
			selectColumns[i] = f.Name
		}
	} else {
		selectColumns = columns
	}

	switch table.Type {
	case externaldata.TypeMongoDB:
		pipeline := make([]bson.M, 0)
		filter := common.GetSearchQuery(r.Search, searchBy)
		if len(filter) > 0 {
			pipeline = append(pipeline, bson.M{"$match": filter})
		}

		project := bson.M{}
		for _, c := range selectColumns {
			project[c] = 1
		}

		pipeline = append(pipeline, bson.M{"$project": project})
		pipeline = append(pipeline, common.GetSortQuery(externaldata.IDColumnName, mongo.SortAsc))

		cursor, err := s.dbExportClient.Collection(table.getDBTableName()).Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
		if err != nil {
			return nil, err
		}

		return export.NewMongoCursor(cursor, selectColumns, nil), nil
	case externaldata.TypePostgreSQL:
		pgPool, err := s.pgPoolProvider.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get postgres pool: %w", err)
		}

		sql := "SELECT "
		for i, c := range selectColumns {
			sql += pgx.Identifier{c}.Sanitize()
			if i < len(selectColumns)-1 {
				sql += ", "
			}
		}

		sql += " FROM " + table.getDBTableName()
		queryArgs := make([]any, 0)
		if r.Search != "" {
			sql += " WHERE "
			queryArgs = append(queryArgs, r.Search)
			for i, col := range searchBy {
				sql += pgx.Identifier{col}.Sanitize() + " ~ $" + strconv.Itoa(len(queryArgs))
				if i < len(searchBy)-1 {
					sql += " OR "
				}
			}
		}

		sql += " ORDER BY " + externaldata.IDColumnName

		rows, err := pgPool.Query(ctx, sql, queryArgs...)
		if err != nil {
			return nil, err
		}

		return export.NewPostgresCursor(rows, selectColumns, nil), nil
	default:
		return nil, fmt.Errorf("unknown external data type %d", table.Type)
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
		" ( " + externaldata.IDColumnName + " VARCHAR(" + MaxIDLenStr + ") PRIMARY KEY )"
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

func (s *store) findPreviewDataFromMongo(ctx context.Context, job ImportJob, request ListPreviewRequest) (*AggregationDataResult, error) {
	project := bson.M{}

	for _, cfg := range job.ColumnConfigs {
		name := "$" + cfg.Name
		project[cfg.Name] = bson.M{
			"$ifNull": bson.A{name + ".transformed_value", name},
		}
	}

	cursor, err := s.dbClient.Collection(job.getDBTableName()).Aggregate(ctx, pagination.CreateAggregationPipeline(
		request.Query,
		[]bson.M{},
		bson.M{"$sort": bson.M{externaldata.IDColumnName: 1}},
		[]bson.M{{"$project": project}},
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

func (s *store) findPreviewDataFromPostgres(ctx context.Context, job ImportJob, r ListPreviewRequest) (*AggregationDataResult, error) {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get postgres pool: %w", err)
	}

	limitStmt := ""
	if r.Paginate {
		limitStmt = "OFFSET " + strconv.FormatInt(r.Limit*(r.Page-1), 10) + " LIMIT " + strconv.FormatInt(r.Limit, 10)
	}

	columnsWithID := make([]string, len(job.ColumnConfigs)*sqlColumnsForCsvColumn+1)
	columnsWithID[0] = externaldata.IDColumnName

	sql := "SELECT " + pgx.Identifier{externaldata.IDColumnName}.Sanitize() + ", "
	for i, cfg := range job.ColumnConfigs {
		initialName := cfg.Name + "_initial_value"
		transformedName := cfg.Name + "_transformed_value"
		errorsName := cfg.Name + "_transform_error"

		columnsWithID[i*sqlColumnsForCsvColumn+1] = initialName
		columnsWithID[i*sqlColumnsForCsvColumn+2] = transformedName
		columnsWithID[i*sqlColumnsForCsvColumn+3] = errorsName

		sql += pgx.Identifier{initialName}.Sanitize() + ", " +
			pgx.Identifier{transformedName}.Sanitize() + ", " +
			pgx.Identifier{errorsName}.Sanitize()
		if i < len(job.ColumnConfigs)-1 {
			sql += ", "
		}
	}

	tableName := job.getDBTableName()

	// A hack to avoid the "cached plan must not change result type (SQLSTATE 0A000)" error when column types have been changed.
	// Adding a comment forces Postgres to execute the query without using a cached plan.
	comment := "/*" + strconv.Itoa(int(time.Now().UnixMicro())) + "*/"

	sql += " FROM " + tableName + " " + "ORDER BY " + pgx.Identifier{externaldata.IDColumnName}.Sanitize() + " " + limitStmt + comment
	countSql := "SELECT count(*) FROM " + tableName
	result := &AggregationDataResult{
		Data: make([]map[string]any, 0, r.Limit),
	}

	rows, err := pgPool.Query(ctx, sql)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return result, err
		}

		row, err := s.transformPostgresPreviewResToData(vals, job.ColumnConfigs)
		if err != nil {
			return result, err
		}

		result.Data = append(result.Data, row)
	}

	if err = rows.Err(); err != nil {
		return result, err
	}

	err = pgPool.QueryRow(ctx, countSql).Scan(&result.TotalCount)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *store) findDataFromPostgres(ctx context.Context, tableName string, columns []string, r ListDataRequest, columnConfigs []externaldata.ColumnConfig) (*AggregationDataResult, error) {
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

	// A hack to avoid the "cached plan must not change result type (SQLSTATE 0A000)" error when column types have been changed.
	// Adding a comment forces Postgres to execute the query without using a cached plan.
	comment := "/*" + strconv.Itoa(int(time.Now().UnixMicro())) + "*/"

	sql += " FROM " + tableName + " " + whereStmt + " " + orderStmt + " " + limitStmt + comment
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

		row, err := s.transformPostgresResToData(vals, columnConfigs)
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

func (s *store) transformPostgresResToData(vals []any, columnConfigs []externaldata.ColumnConfig) (map[string]any, error) {
	res := make(map[string]any, len(vals))

	id, ok := vals[0].(string)
	if !ok {
		return res, fmt.Errorf("%q column doesn't contain string", externaldata.IDColumnName)
	}

	res[externaldata.IDColumnName] = id

	for i, cfg := range columnConfigs {
		var value any

		switch cfg.Type {
		case externaldata.ColumnTypeString:
			value, ok = vals[i+1].(string)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a string", cfg.Name)
			}
		case externaldata.ColumnTypeNumber:
			value, ok = vals[i+1].(float64)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a number", cfg.Name)
			}
		case externaldata.ColumnTypeBoolean:
			value, ok = vals[i+1].(bool)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a boolean", cfg.Name)
			}
		case externaldata.ColumnTypeDateTime, externaldata.ColumnTypeTimestamp:
			value, ok = vals[i+1].(int64)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain an int", cfg.Name)
			}
		case externaldata.ColumnTypeStringArray:
			value, ok = utils.IsStringSlice(vals[i+1])
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a string array", cfg.Name)
			}
		default:
			return nil, fmt.Errorf("unsupported column type %q", cfg.Type)
		}

		res[cfg.Name] = value
	}

	return res, nil
}

func (s *store) transformPostgresPreviewResToData(vals []any, columnConfigs []ColumnConfig) (map[string]any, error) {
	res := make(map[string]any, len(vals))

	id, ok := vals[0].(string)
	if !ok {
		return res, fmt.Errorf("%q column doesn't contain string", externaldata.IDColumnName)
	}

	res[externaldata.IDColumnName] = id

	for i, cfg := range columnConfigs {
		initValue, ok := vals[i*sqlColumnsForCsvColumn+1].(string)
		if !ok {
			return nil, fmt.Errorf("initial value for %q column doesn't contain a string", cfg.Name)
		}

		if vals[i*sqlColumnsForCsvColumn+2] == nil {
			if vals[i*sqlColumnsForCsvColumn+3] == nil {
				return nil, fmt.Errorf("error for %q column is nil", cfg.Name)
			}

			errorMessage, ok := vals[i*sqlColumnsForCsvColumn+3].(string)
			if !ok {
				return nil, fmt.Errorf("error for %q column doesn't contain a string", cfg.Name)
			}

			res[cfg.Name] = map[string]string{
				"initial_value":   initValue,
				"transform_error": errorMessage,
			}

			continue
		}

		var transformedValue any

		switch cfg.Type {
		case externaldata.ColumnTypeString:
			transformedValue, ok = vals[i*sqlColumnsForCsvColumn+2].(string)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a string", cfg.Name)
			}
		case externaldata.ColumnTypeNumber:
			transformedValue, ok = vals[i*sqlColumnsForCsvColumn+2].(float64)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a number", cfg.Name)
			}
		case externaldata.ColumnTypeBoolean:
			transformedValue, ok = vals[i*sqlColumnsForCsvColumn+2].(bool)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a boolean", cfg.Name)
			}
		case externaldata.ColumnTypeDateTime, externaldata.ColumnTypeTimestamp:
			transformedValue, ok = vals[i*sqlColumnsForCsvColumn+2].(int64)
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain an int", cfg.Name)
			}
		case externaldata.ColumnTypeStringArray:
			transformedValue, ok = utils.IsStringSlice(vals[i*sqlColumnsForCsvColumn+2])
			if !ok {
				return nil, fmt.Errorf("transformed value for %q column doesn't contain a string array", cfg.Name)
			}
		default:
			return nil, fmt.Errorf("unsupported column type %q", cfg.Type)
		}

		res[cfg.Name] = transformedValue
	}

	return res, nil
}

func (s *store) updateLinkedModels(ctx context.Context, tableID, tableName string) error {
	for _, c := range s.linkedDbCollections {
		_, err := c.UpdateMany(ctx,
			bson.M{"external_data.table": tableID},
			bson.M{"$set": bson.M{"external_data.$[exdata].table_name": tableName}},
			options.UpdateMany().SetArrayFilters([]any{bson.M{"exdata.table": tableID}}),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func TransformRefParameters(ctx context.Context, r []externaldata.RefParameters, dbCollection mongo.DbCollection) ([]externaldata.RefParameters, error) {
	ids := make([]string, 0, len(r))
	for _, params := range r {
		if params.Type == externaldata.RefTypeTable {
			ids = append(ids, params.Table)
		}
	}

	cursor, err := dbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	tables := make(map[string]externaldata.Table, len(ids))
	for cursor.Next(ctx) {
		t := externaldata.Table{}
		err = cursor.Decode(&t)
		if err != nil {
			return nil, err
		}

		tables[t.ID] = t
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	res := make([]externaldata.RefParameters, len(r))
	valErrMsgs := make(map[string]string)
	for i, params := range r {
		if params.Type != externaldata.RefTypeTable {
			res[i] = params
			continue
		}

		t, ok := tables[params.Table]
		if !ok {
			valErrMsgs["external_data."+strconv.Itoa(i)+".table"] = "Table doesn't exist."
			continue
		}

		if len(t.ColumnConfigs) == 0 {
			valErrMsgs["external_data."+strconv.Itoa(i)+".table"] = "Table is empty."
			continue
		}

		columns := make([]string, len(t.ColumnConfigs))
		hasCol := make(map[string]bool, len(t.ColumnConfigs))
		for i, c := range t.ColumnConfigs {
			hasCol[c.Name] = true
			columns[i] = c.Name
		}

		if params.SortBy != "" && !hasCol[params.SortBy] {
			valErrMsgs["external_data."+strconv.Itoa(i)+".sort_by"] = "SortBy must be one of [" + strings.Join(columns, " ") + "]."
		}

		for f := range params.Select {
			if !hasCol[f] {
				valErrMsgs["external_data."+strconv.Itoa(i)+".select."+f] = f + " must be one of [" + strings.Join(columns, " ") + "]."
			}
		}

		for f := range params.Regexp {
			if !hasCol[f] {
				valErrMsgs["external_data."+strconv.Itoa(i)+".regexp."+f] = f + " must be one of [" + strings.Join(columns, " ") + "]."
			}
		}

		if len(valErrMsgs) > 0 {
			continue
		}

		params.TableName = t.GetDBName()
		params.TableType = &t.Type
		params.TableColumns = columns
		res[i] = params
	}

	if len(valErrMsgs) > 0 {
		return nil, common.NewValidationErrors(valErrMsgs)
	}

	return res, nil
}

func GetRefParametersLookups() []bson.M {
	return []bson.M{
		{"$unwind": bson.M{
			"path":                       "$external_data",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "exdata_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.ExternalDataTableCollection,
			"localField":   "external_data.table",
			"foreignField": "_id",
			"as":           "external_data.table",
		}},
		{"$unwind": bson.M{"path": "$external_data.table", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"exdata_index": 1}},
		{"$group": bson.M{
			"_id":  "$_id",
			"data": bson.M{"$first": "$$ROOT"},
			"external_data": bson.M{"$push": bson.M{"$cond": bson.M{
				"if":   "$external_data.reference",
				"then": "$external_data",
				"else": "$$REMOVE",
			}}},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"external_data": "$external_data"},
		}}}},
	}
}

func isTableLinked(ctx context.Context, id string, dbWidgetCollection mongo.DbCollection, linkedDbCollections []mongo.DbCollection) (bool, error) {
	err := dbWidgetCollection.
		FindOne(ctx, bson.M{"type": view.WidgetTypeExternalData, "parameters.table": id}, options.FindOne().SetProjection(bson.M{"_id": 1})).
		Err()
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return false, err
	}

	if err == nil {
		return true, nil
	}

	for _, c := range linkedDbCollections {
		err = c.
			FindOne(ctx, bson.M{"external_data.table": id}, options.FindOne().SetProjection(bson.M{"_id": 1})).
			Err()
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, err
		}

		if err == nil {
			return true, nil
		}
	}

	return false, nil
}

func getIntValue(v interface{}) (int64, bool) {
	switch i := v.(type) {
	case int:
		return int64(i), true
	case int32:
		return int64(i), true
	case int64:
		return i, true
	case uint:
		return int64(i), true
	case uint32:
		return int64(i), true
	case uint64:
		return int64(i), true
	case float32:
		a, b := math.Modf(float64(i))
		if b == 0 {
			return int64(a), true
		}
	case float64:
		a, b := math.Modf(i)
		if b == 0 {
			return int64(a), true
		}
	}

	return 0, false
}
