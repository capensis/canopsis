package icon

import (
	"context"
	"errors"
	"mime/multipart"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libfile "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Create(context.Context, EditRequest) (*Response, error)
	Update(context.Context, EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
	List(ctx context.Context, query pagination.FilteredQuery) (*AggregationResult, error)
	Get(ctx context.Context, id string) (*Response, error)
	GetFilepath(model Response) string
}

type store struct {
	dbCollection          mongo.DbCollection
	storage               libfile.Storage
	defaultSortBy         string
	defaultSearchByFields []string
}

func NewStore(dbClient mongo.DbClient, storage libfile.Storage) Store {
	return &store{
		dbCollection:          dbClient.Collection(mongo.IconCollection),
		storage:               storage,
		defaultSortBy:         "title",
		defaultSearchByFields: []string{"_id", "title"},
	}
}

func (s *store) Create(ctx context.Context, r EditRequest) (*Response, error) {
	id := utils.NewID()
	res, err := s.storeFile(id, r.File)
	if err != nil {
		return nil, err
	}

	if res != nil {
		now := types.NewCpsTime()
		res.ID = id
		res.Title = r.Title
		res.MimeType = r.MimeType
		res.Created = now
		res.Updated = now
		_, err := s.dbCollection.InsertOne(ctx, res)
		if err != nil {
			return nil, err
		}
	}

	return res, err
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	id := r.ID
	old, err := s.Get(ctx, id)
	if err != nil || old == nil {
		return nil, err
	}

	res, err := s.storeFile(id, r.File)
	if err != nil {
		return nil, err
	}

	err = s.storage.Delete(id, old.Storage)
	if err != nil {
		return nil, err
	}

	if res != nil {
		now := types.NewCpsTime()
		res.ID = id
		res.Title = r.Title
		res.MimeType = r.MimeType
		res.Updated = now
		updateRes, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": res})
		if err != nil || updateRes.MatchedCount == 0 {
			return nil, err
		}
	}

	return res, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	f := Response{}
	err := s.dbCollection.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&f)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	err = s.storage.Delete(f.ID, f.Storage)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) List(ctx context.Context, query pagination.FilteredQuery) (*AggregationResult, error) {
	var pipeline []bson.M
	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(s.defaultSortBy, common.SortAsc),
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

func (s *store) Get(ctx context.Context, id string) (*Response, error) {
	res := Response{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &res, nil
}

func (s *store) GetFilepath(model Response) string {
	return s.storage.GetFilepath(model.ID, model.Storage)
}

func (s *store) storeFile(id string, file *multipart.FileHeader) (_ *Response, resErr error) {
	tmp, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		err = tmp.Close()
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	storage, err := s.storage.CopyReader(id, tmp)
	if err != nil {
		return nil, err
	}

	etag, err := s.storage.GetEtag(id, storage)
	if err != nil {
		return nil, err
	}

	return &Response{
		Storage: storage,
		Etag:    etag,
	}, nil
}
