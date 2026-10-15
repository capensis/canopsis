package crud

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
)

type Store[ListRequest, CreateRequest, UpdateRequest, Response any] interface {
	Find(ctx context.Context, r ListRequest) (pagination.ListResponse[Response], error)
	GetByID(ctx context.Context, id string) (Response, error)
	Insert(ctx context.Context, r CreateRequest) (Response, error)
	Update(ctx context.Context, id string, r UpdateRequest) (Response, error)
	Delete(ctx context.Context, id, userID string) error
}
