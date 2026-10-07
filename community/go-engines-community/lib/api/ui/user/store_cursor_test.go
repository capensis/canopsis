package user

import (
	"context"
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/author"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"
)

func TestGetOneByReturnsCursorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	collection := mock_mongo.NewMockDbCollection(ctrl)
	cursor := mock_mongo.NewMockCursor(ctrl)
	ctx := context.Background()
	wantErr := errors.New("cursor failed")

	collection.EXPECT().Aggregate(ctx, gomock.Any()).Return(cursor, nil)
	cursor.EXPECT().Next(ctx).Return(false)
	cursor.EXPECT().Err().Return(wantErr)
	cursor.EXPECT().Close(ctx).Return(nil)

	store := &store{
		userCollection: collection,
		authorProvider: cursorTestAuthorProvider{},
	}
	got, err := store.GetOneBy(ctx, "user-1")
	if got != nil || !errors.Is(err, wantErr) {
		t.Fatalf("GetOneBy() = (%#v, %v), want (nil, cursor error)", got, err)
	}
}

type cursorTestAuthorProvider struct{}

func (cursorTestAuthorProvider) Pipeline() []bson.M { return nil }

func (cursorTestAuthorProvider) PipelineForField(string) []bson.M { return nil }

func (cursorTestAuthorProvider) GetDisplayNameQuery(string) bson.M { return bson.M{} }

func (cursorTestAuthorProvider) AddAuthorStr(string) bson.M { return bson.M{} }

var _ author.Provider = cursorTestAuthorProvider{}
var _ libmongo.DbCollection = (*mock_mongo.MockDbCollection)(nil)
