package userprovider

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	mockmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"go.uber.org/mock/gomock"
)

// setRoleNames fills the aggregation result the way cursor.Decode does.
func setRoleNames(val any, names ...string) error {
	result, ok := val.(*roleNamesAggregationResult)
	if !ok {
		return fmt.Errorf("unexpected decode target %T", val)
	}
	for _, name := range names {
		result.Roles = append(result.Roles, struct {
			Name string `bson:"name"`
		}{Name: name})
	}
	return nil
}

func TestRoleNamesReaderGetRoleNames(t *testing.T) {
	t.Run("empty user ID skips the database", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		collection := mockmongo.NewMockDbCollection(ctrl)
		got, err := (&roleNamesReader{collection: collection}).GetRoleNames(context.Background(), "")
		if err != nil {
			t.Fatal(err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("expected a non-nil empty slice, got %#v", got)
		}
	})

	t.Run("reports aggregate error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		collection := mockmongo.NewMockDbCollection(ctrl)
		wantErr := errors.New("mongo unavailable")
		collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(nil, wantErr)
		_, err := (&roleNamesReader{collection: collection}).GetRoleNames(context.Background(), "user-1")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected wrapped aggregate error, got %v", err)
		}
	})

	t.Run("reports cursor and decode errors", func(t *testing.T) {
		decodeError := errors.New("invalid document")
		for _, test := range []struct {
			name       string
			decodeErr  error
			cursorErr  error
			wantErr    error
			nextResult bool
		}{
			{name: "cursor", cursorErr: context.Canceled, wantErr: context.Canceled},
			{name: "decode", decodeErr: decodeError, wantErr: decodeError, nextResult: true},
		} {
			t.Run(test.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				collection := mockmongo.NewMockDbCollection(ctrl)
				cursor := mockmongo.NewMockCursor(ctrl)
				collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil)
				cursor.EXPECT().Next(gomock.Any()).Return(test.nextResult)
				if test.nextResult {
					cursor.EXPECT().Decode(gomock.Any()).Return(test.decodeErr)
				}
				if test.decodeErr == nil {
					cursor.EXPECT().Err().Return(test.cursorErr)
				}
				cursor.EXPECT().Close(gomock.Any()).Return(nil)
				_, err := (&roleNamesReader{collection: collection}).GetRoleNames(context.Background(), "user-1")
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("expected wrapped %v error, got %v", test.wantErr, err)
				}
			})
		}
	})

	t.Run("returns empty slice when no document is found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		collection := mockmongo.NewMockDbCollection(ctrl)
		cursor := mockmongo.NewMockCursor(ctrl)
		collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil)
		cursor.EXPECT().Next(gomock.Any()).Return(false)
		cursor.EXPECT().Err().Return(nil)
		cursor.EXPECT().Close(gomock.Any()).Return(nil)
		got, err := (&roleNamesReader{collection: collection}).GetRoleNames(context.Background(), "user-1")
		if err != nil {
			t.Fatal(err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("expected a non-nil empty slice, got %#v", got)
		}
	})

	t.Run("reports context errors when no document is found", func(t *testing.T) {
		for _, test := range []struct {
			name    string
			newCtx  func() (context.Context, context.CancelFunc)
			wantErr error
		}{
			{
				name: "cancelled",
				newCtx: func() (context.Context, context.CancelFunc) {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx, cancel
				},
				wantErr: context.Canceled,
			},
			{
				name: "deadline exceeded",
				newCtx: func() (context.Context, context.CancelFunc) {
					return context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
				},
				wantErr: context.DeadlineExceeded,
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				collection := mockmongo.NewMockDbCollection(ctrl)
				cursor := mockmongo.NewMockCursor(ctrl)
				ctx, cancel := test.newCtx()
				defer cancel()
				collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil)
				cursor.EXPECT().Next(gomock.Any()).Return(false)
				cursor.EXPECT().Err().Return(nil)
				cursor.EXPECT().Close(gomock.Any()).Return(nil)
				_, err := (&roleNamesReader{collection: collection}).GetRoleNames(ctx, "user-1")
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("expected wrapped %v error, got %v", test.wantErr, err)
				}
			})
		}
	})

	t.Run("reports cursor error after decoding a document", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		collection := mockmongo.NewMockDbCollection(ctrl)
		cursor := mockmongo.NewMockCursor(ctrl)
		wantErr := errors.New("cursor failed")
		collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil)
		cursor.EXPECT().Next(gomock.Any()).Return(true)
		cursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val any) error {
			return setRoleNames(val, "admin")
		})
		cursor.EXPECT().Err().Return(wantErr)
		cursor.EXPECT().Close(gomock.Any()).Return(nil)
		_, err := (&roleNamesReader{collection: collection}).GetRoleNames(context.Background(), "user-1")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected wrapped cursor error, got %v", err)
		}
	})

	t.Run("reports context cancellation after decoding a document", func(t *testing.T) {
		for _, test := range []struct {
			name         string
			cancelInNext bool
		}{
			{name: "cancelled during next", cancelInNext: true},
			{name: "cancelled during decode"},
		} {
			t.Run(test.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				collection := mockmongo.NewMockDbCollection(ctrl)
				cursor := mockmongo.NewMockCursor(ctrl)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil)
				cursor.EXPECT().Next(gomock.Any()).DoAndReturn(func(context.Context) bool {
					if test.cancelInNext {
						cancel()
					}
					return true
				})
				cursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val any) error {
					if !test.cancelInNext {
						cancel()
					}
					return setRoleNames(val, "admin")
				})
				cursor.EXPECT().Err().Return(nil)
				cursor.EXPECT().Close(gomock.Any()).Return(nil)
				_, err := (&roleNamesReader{collection: collection}).GetRoleNames(ctx, "user-1")
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("expected wrapped context.Canceled error, got %v", err)
				}
			})
		}
	})

	t.Run("returns sorted non-empty role names", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		collection := mockmongo.NewMockDbCollection(ctrl)
		cursor := mockmongo.NewMockCursor(ctrl)
		collection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(cursor, nil)
		cursor.EXPECT().Next(gomock.Any()).Return(true)
		cursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(val any) error {
			return setRoleNames(val, "operator", "", "admin", "user")
		})
		cursor.EXPECT().Err().Return(nil)
		cursor.EXPECT().Close(gomock.Any()).Return(nil)
		got, err := (&roleNamesReader{collection: collection}).GetRoleNames(context.Background(), "user-1")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"admin", "operator", "user"}
		if !slices.Equal(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
}
