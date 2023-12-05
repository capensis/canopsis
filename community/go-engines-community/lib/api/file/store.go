package file

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/file"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	contentType = "Content-Type"
)

// Store interface consists of helpers for api handlers
type Store interface {
	Create(context.Context, bool, *multipart.Form) ([]File, error)
	Delete(ctx context.Context, id string) (bool, error)
	List(ctx context.Context, ids []string) ([]File, error)
	Get(ctx context.Context, id string) (*File, error)
	GetFilepath(model File) string
}

type store struct {
	dbCollection mongo.DbCollection
	storage      file.Storage
	maxSize      int64
}

// NewStore initializes Store implementation.
func NewStore(dbClient mongo.DbClient, storage file.Storage, maxSize int64) Store {
	return &store{
		dbCollection: dbClient.Collection(mongo.FileMongoCollection),
		storage:      storage,
		maxSize:      maxSize,
	}
}

// Create parses form data from request and stores files and linked database records
func (s *store) Create(ctx context.Context, isPublic bool, form *multipart.Form) ([]File, error) {
	files, err := s.validateFormRequest(form)
	if err != nil {
		return nil, err
	}
	models, err := s.storeFiles(isPublic, files)
	if err != nil {
		return nil, err
	}

	if len(models) > 0 {
		docs := make([]interface{}, len(models))
		for i := range models {
			docs[i] = models[i]
		}

		_, err := s.dbCollection.InsertMany(ctx, docs)
		if err != nil {
			return nil, err
		}
	}

	return models, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	f := File{}
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

// List files
func (s *store) List(ctx context.Context, ids []string) ([]File, error) {
	cursor, err := s.dbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	files := make([]File, 0)
	err = cursor.All(ctx, &files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (s *store) Get(ctx context.Context, id string) (*File, error) {
	f := File{}
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&f)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &f, nil
}

func (s *store) GetFilepath(model File) string {
	return s.storage.GetFilepath(model.ID, model.Storage)
}

func (s *store) storeFiles(isPublic bool, files []*multipart.FileHeader) ([]File, error) {
	models := make([]File, len(files))

	for i, f := range files {
		id := utils.NewID()
		tmp, err := f.Open()
		if err != nil {
			return nil, err
		}

		storage, err := s.storage.CopyReader(id, tmp)
		if err != nil {
			return nil, err
		}

		err = tmp.Close()
		if err != nil {
			return nil, err
		}

		etag, err := s.storage.GetEtag(id, storage)
		if err != nil {
			return nil, err
		}

		models[i] = File{
			ID:        id,
			FileName:  f.Filename,
			MediaType: f.Header.Get(contentType),
			Created:   datetime.NewCpsTime(),
			Storage:   storage,
			Etag:      etag,
			IsPublic:  isPublic,
		}
	}

	return models, nil
}

func (s *store) validateFormRequest(form *multipart.Form) ([]*multipart.FileHeader, error) {
	files := make([]*multipart.FileHeader, 0)

	for field, headers := range form.File {
		for i, header := range headers {
			if s.maxSize > 0 && header.Size > s.maxSize {
				return nil, ValidationError{
					field: fmt.Sprintf("%s[%d]", field, i),
					error: fmt.Sprintf("file size %d exceeds limit %d", header.Size, s.maxSize),
				}
			}

			files = append(files, header)
		}
	}

	return files, nil
}
