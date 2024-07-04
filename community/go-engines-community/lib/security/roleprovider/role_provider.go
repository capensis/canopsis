package roleprovider

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ErrRoleNotFound struct {
	role string
}

func (e ErrRoleNotFound) Error() string {
	return "role " + e.role + " doesn't exist"
}

type roleProvider struct {
	roleCollection mongo.DbCollection
}

func NewRoleProvider(dbClient mongo.DbClient) security.RoleProvider {
	return &roleProvider{
		roleCollection: dbClient.Collection(mongo.RoleCollection),
	}
}

func (v *roleProvider) GetRoleID(ctx context.Context, name string) (string, error) {
	var role struct {
		ID string `bson:"_id"`
	}

	err := v.roleCollection.FindOne(ctx, bson.M{"name": name}, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(&role)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return "", ErrRoleNotFound{role: name}
		}

		return "", err
	}

	return role.ID, nil
}

func (v *roleProvider) GetValidRoleIDs(ctx context.Context, potentialRoles []string, defaultRole string) ([]string, error) {
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
				"found_roles": bson.M{"$addToSet": "$_id"},
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

func (v *roleProvider) getDefaultRole(ctx context.Context, defaultRole string) ([]string, error) {
	id, err := v.GetRoleID(ctx, defaultRole)
	if err != nil {
		return nil, err
	}

	return []string{id}, nil
}
