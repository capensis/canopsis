package providers

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrDefaultRoleNotFound = errors.New("default role not found")

type RoleProvider interface {
	// GetValidRoles checks if potentialRoles slice contains valid roles and returns at least one valid role.
	// If no roles found, then it check if default role is valid and returns it. Return ErrDefaultRoleNotFound error if default role not found.
	GetValidRoles(ctx context.Context, potentialRoles []string, defaultRole string) ([]string, error)
}

type roleProvider struct {
	roleCollection mongo.DbCollection
}

func NewRoleProvider(dbClient mongo.DbClient) RoleProvider {
	return &roleProvider{
		roleCollection: dbClient.Collection(mongo.RoleCollection),
	}
}

func (v *roleProvider) getDefaultRole(ctx context.Context, defaultRole string) ([]string, error) {
	err := v.roleCollection.FindOne(ctx, bson.M{"name": defaultRole}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, ErrDefaultRoleNotFound
		}

		return nil, err
	}

	return []string{defaultRole}, nil
}

func (v *roleProvider) GetValidRoles(ctx context.Context, potentialRoles []string, defaultRole string) ([]string, error) {
	if len(potentialRoles) == 0 {
		return v.getDefaultRole(ctx, defaultRole)
	}

	cursor, err := v.roleCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"name": bson.M{
					"$in": potentialRoles,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":         nil,
				"found_roles": bson.M{"$addToSet": "$name"},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var res struct {
		FoundRoles []string `bson:"found_roles"`
	}

	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("failed to decode get roles query result: %w", err)
		}
	}

	if len(res.FoundRoles) == 0 {
		return v.getDefaultRole(ctx, defaultRole)
	}

	return res.FoundRoles, nil
}
