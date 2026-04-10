package dbvalidation

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ValidateExist(
	ctx context.Context,
	collection mongo.DbCollection,
	request any,
	field string,
	value any,
) error {
	var q bson.M
	var expectedCount int64
	switch v := value.(type) {
	case string:
		if v == "" {
			return nil
		}

		q = bson.M{"_id": value}
		expectedCount = 1
	case *string:
		if v == nil || *v == "" {
			return nil
		}

		q = bson.M{"_id": value}
		expectedCount = 1
	case []string:
		if len(v) == 0 {
			return nil
		}

		q = bson.M{"_id": bson.M{"$in": value}}
		expectedCount = int64(len(v))
	default:
		return fmt.Errorf("unsupported type: %T, collection %q, field %q, request %+v", value, collection.Name(), field, request)
	}

	count, err := collection.CountDocuments(ctx, q)
	if err != nil {
		return err
	}

	if count != expectedCount {
		return validation.NewSingleError("not_exist", field, field, request)
	}

	return nil
}

func ValidateLinkedReference(
	ctx context.Context,
	collection mongo.DbCollection,
	filter bson.M,
	entityName string,
	referencedBy string,
) error {
	err := collection.
		FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"_id": 1})).
		Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil
		}

		return err
	}

	return httperror.NewConflictError("The " + entityName + " cannot be deleted because it is referenced by " + referencedBy + ".")
}
