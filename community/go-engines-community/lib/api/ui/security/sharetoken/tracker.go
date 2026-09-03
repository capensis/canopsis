package sharetoken

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewTracker(client mongo.DbClient, logger zerolog.Logger) token.UsageTracker {
	return &tracker{
		collection: client.Collection(mongo.ShareTokenMongoCollection),
		logger:     logger,
	}
}

type tracker struct {
	collection mongo.DbCollection
	logger     zerolog.Logger
}

func (s *tracker) Exists(ctx context.Context, v string) (bool, error) {
	err := s.collection.FindOne(ctx, bson.M{"value": v}).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *tracker) Access(ctx context.Context, v string) error {
	now := datetime.NewCpsTime()
	_, err := s.collection.UpdateOne(ctx, bson.M{"value": v}, bson.M{"$set": bson.M{"accessed": now}})

	return err
}
