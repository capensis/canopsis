package providers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleValidator interface {
	// AreRolesValid checks whether provided slice of roles are valid roles, returns an error with a list of roles, which are not valid.
	AreRolesValid(ctx context.Context, roles []string) error
}

type notFoundRolesRes struct {
	NotFoundRoles []string `bson:"not_found_roles"`
}

type roleValidator struct {
	roleCollection mongo.DbCollection
}

func NewRoleValidator(dbClient mongo.DbClient) RoleValidator {
	return &roleValidator{
		roleCollection: dbClient.Collection(mongo.RoleCollection),
	}
}

func (v *roleValidator) AreRolesValid(ctx context.Context, roles []string) error {
	if len(roles) == 0 {
		return errors.New("empty roles")
	}

	cursor, err := v.roleCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"name": bson.M{
					"$in": roles,
				},
			},
		},
		{
			"$group": bson.M{
				"_id":         nil,
				"found_roles": bson.M{"$addToSet": "$name"},
			},
		},
		{
			"$project": bson.M{
				"not_found_roles": bson.M{
					"$setDifference": bson.A{roles, "$found_roles"},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	// First, assume that no roles was found. If the cursor doesn't have anything, then no roles was found => all of them are invalid.
	notFoundRoles := roles

	if cursor.Next(ctx) {
		var res notFoundRolesRes

		err = cursor.Decode(&res)
		if err != nil {
			panic(fmt.Errorf("failed to decode not found roles query result: %w", err))
		}

		notFoundRoles = res.NotFoundRoles
	}

	if len(notFoundRoles) != 0 {
		return fmt.Errorf("roles %s doesn't exist", strings.Join(notFoundRoles, ","))
	}

	return nil
}
