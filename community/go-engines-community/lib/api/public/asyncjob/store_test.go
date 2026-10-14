package asyncjob

import (
	"context"
	"errors"
	"testing"
	"time"

	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/mock/gomock"
)

type storeTestJob struct {
	Base `bson:",inline"`
}

func TestStoreClaim_ReturnsPostUpdateLeaseGeneration(t *testing.T) {
	ctrl := gomock.NewController(t)
	collection := mock_mongo.NewMockDbCollection(ctrl)
	result := mock_mongo.NewMockSingleResultHelper(ctrl)
	store := Store[storeTestJob]{collection: collection}

	collection.EXPECT().
		FindOneAndUpdate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(
			_ context.Context,
			filter, update any,
			_ ...options.Lister[options.FindOneAndUpdateOptions],
		) libmongo.SingleResultHelper {
			if got := filter.(bson.M); got["_id"] != "job-id" || got["status"] != StatusPending {
				t.Errorf("unexpected claim filter: %#v", got)
			}

			inc := update.(bson.M)["$inc"].(bson.M)
			if inc["lease_generation"] != 1 {
				t.Errorf("unexpected claim increment: %#v", inc)
			}

			return result
		})
	result.EXPECT().Decode(gomock.Any()).DoAndReturn(func(v any) error {
		*v.(*storeTestJob) = storeTestJob{
			ID:              "job-id",
			Status:          StatusOngoing,
			LeaseGeneration: 3,
		}

		return nil
	})

	job, claimed, err := store.Claim(t.Context(), "job-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !claimed {
		t.Fatal("expected the pending job to be claimed")
	}
	if job.LeaseGeneration != 3 {
		t.Fatalf("expected lease generation 3, got %d", job.LeaseGeneration)
	}
}

func TestStorePing_RejectsAnOldLease(t *testing.T) {
	ctrl := gomock.NewController(t)
	collection := mock_mongo.NewMockDbCollection(ctrl)
	store := Store[storeTestJob]{collection: collection}

	collection.EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, filter, _ any, _ ...options.Lister[options.UpdateOneOptions]) (*mongodriver.UpdateResult, error) {
			got := filter.(bson.M)
			if got["_id"] != "job-id" || got["status"] != StatusOngoing || got["lease_generation"] != int64(2) {
				t.Errorf("unexpected fenced update filter: %#v", got)
			}

			return &mongodriver.UpdateResult{}, nil
		})

	err := store.Ping(t.Context(), "job-id", 2)
	if !errors.Is(err, ErrLeaseLost) {
		t.Fatalf("expected ErrLeaseLost, got %v", err)
	}
}

func TestStoreSetDone_UpdatesTheClaimedLease(t *testing.T) {
	ctrl := gomock.NewController(t)
	collection := mock_mongo.NewMockDbCollection(ctrl)
	store := Store[storeTestJob]{collection: collection}

	collection.EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, filter, update any, _ ...options.Lister[options.UpdateOneOptions]) (*mongodriver.UpdateResult, error) {
			gotFilter := filter.(bson.M)
			if gotFilter["lease_generation"] != int64(2) {
				t.Errorf("expected fenced completion, got %#v", gotFilter)
			}

			set := update.(bson.M)["$set"].(bson.M)
			if set["status"] != StatusDone {
				t.Errorf("expected done status update, got %#v", set)
			}

			return &mongodriver.UpdateResult{MatchedCount: 1}, nil
		})

	err := store.SetDone(t.Context(), "job-id", 2, time.Second, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
