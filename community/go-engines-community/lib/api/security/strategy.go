package security

import (
	"context"
	"errors"

	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	OwnershipNone byte = iota
	OwnershipPublic
	OwnershipOwner
	OwnershipNotOwner
	OwnershipNotFound
)

type OwnershipStrategy interface {
	IsOwner(ctx context.Context, id, userID string) (byte, error)
}

type OwnedObjectsProvider interface {
	GetOwnedIDs(ctx context.Context, userID string) ([]string, error)
}

type viewOwnerStrategy struct {
	collection           libmongo.DbCollection
	enforcer             security.Enforcer
	linkedViewPermission string
}

func NewViewOwnerStrategy(client libmongo.DbClient, enforcer security.Enforcer, linkedViewPermission string) OwnershipStrategy {
	return &viewOwnerStrategy{
		collection:           client.Collection(libmongo.ViewMongoCollection),
		enforcer:             enforcer,
		linkedViewPermission: linkedViewPermission,
	}
}

func (v *viewOwnerStrategy) IsOwner(ctx context.Context, id, userID string) (byte, error) {
	var obj PrivacySettings

	err := v.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"author": 1, "is_private": 1})).Decode(&obj)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return OwnershipNotFound, nil
		}

		return OwnershipNone, err
	}

	if !obj.IsPrivate {
		ok, err := v.enforcer.Enforce(userID, id, v.linkedViewPermission)
		if err != nil {
			panic(err)
		}

		if !ok {
			return OwnershipNotOwner, nil
		}

		return OwnershipPublic, nil
	}

	if obj.Author == userID {
		return OwnershipOwner, nil
	}

	return OwnershipNotOwner, nil
}

type viewOwnedObjectsProvider struct {
	collection libmongo.DbCollection
}

func NewViewOwnedObjectsProvider(client libmongo.DbClient) OwnedObjectsProvider {
	return &viewOwnedObjectsProvider{
		collection: client.Collection(libmongo.ViewMongoCollection),
	}
}

func (v *viewOwnedObjectsProvider) GetOwnedIDs(ctx context.Context, userID string) ([]string, error) {
	var doc struct {
		IDs []string `bson:"ids"`
	}

	cursor, err := v.collection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"is_private": true,
				"author":     userID,
			},
		},
		{
			"$group": bson.M{
				"_id": 1,
				"ids": bson.M{
					"$push": "$_id",
				},
			},
		},
		{
			"$project": bson.M{
				"_id": 0,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&doc)
		if err != nil {
			return nil, err
		}
	}

	return doc.IDs, nil
}

type viewGroupOwnerStrategy struct {
	collection libmongo.DbCollection
}

func NewViewGroupOwnershipStrategy(client libmongo.DbClient) OwnershipStrategy {
	return &viewGroupOwnerStrategy{
		collection: client.Collection(libmongo.ViewGroupMongoCollection),
	}
}

func (v *viewGroupOwnerStrategy) IsOwner(ctx context.Context, id, userID string) (byte, error) {
	var obj PrivacySettings

	err := v.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"author": 1, "is_private": 1})).Decode(&obj)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return OwnershipNotFound, nil
		}

		return OwnershipNone, err
	}

	if !obj.IsPrivate {
		return OwnershipPublic, nil
	}

	if obj.Author == userID {
		return OwnershipOwner, nil
	}

	return OwnershipNotOwner, nil
}

type viewTabOwnerStrategy struct {
	collection           libmongo.DbCollection
	enforcer             security.Enforcer
	linkedViewPermission string
}

func NewViewTabOwnershipStrategy(client libmongo.DbClient, enforcer security.Enforcer, linkedViewPermission string) OwnershipStrategy {
	return &viewTabOwnerStrategy{
		collection:           client.Collection(libmongo.ViewTabMongoCollection),
		enforcer:             enforcer,
		linkedViewPermission: linkedViewPermission,
	}
}

func (v *viewTabOwnerStrategy) IsOwner(ctx context.Context, id, userID string) (byte, error) {
	var obj ViewTabPrivacySettings

	err := v.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"author": 1, "is_private": 1, "view": 1})).Decode(&obj)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return OwnershipNotFound, nil
		}

		return OwnershipNone, err
	}

	if !obj.IsPrivate {
		ok, err := v.enforcer.Enforce(userID, obj.View, v.linkedViewPermission)
		if err != nil {
			panic(err)
		}

		if !ok {
			return OwnershipNotOwner, nil
		}

		return OwnershipPublic, nil
	}

	if obj.Author == userID {
		return OwnershipOwner, nil
	}

	return OwnershipNotOwner, nil
}

type widgetOwnershipStrategy struct {
	collection           libmongo.DbCollection
	enforcer             security.Enforcer
	linkedViewPermission string
}

func NewWidgetOwnershipStrategy(client libmongo.DbClient, enforcer security.Enforcer, linkedViewPermission string) OwnershipStrategy {
	return &widgetOwnershipStrategy{
		collection:           client.Collection(libmongo.WidgetMongoCollection),
		enforcer:             enforcer,
		linkedViewPermission: linkedViewPermission,
	}
}

func (v *widgetOwnershipStrategy) IsOwner(ctx context.Context, id, userID string) (byte, error) {
	var obj ViewTabPrivacySettings

	cursor, err := v.collection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"_id": id,
			},
		},
		{
			"$lookup": bson.M{
				"from":         libmongo.ViewTabMongoCollection,
				"localField":   "tab",
				"foreignField": "_id",
				"as":           "tab",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$tab",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"view":       "$tab.view",
				"author":     "$author",
				"is_private": "$is_private",
			},
		},
	})
	if err != nil {
		return OwnershipNone, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err := cursor.Decode(&obj)
		if err != nil {
			return OwnershipNone, err
		}
	} else {
		return OwnershipNotFound, nil
	}

	if !obj.IsPrivate {
		ok, err := v.enforcer.Enforce(userID, obj.View, v.linkedViewPermission)
		if err != nil {
			panic(err)
		}

		if !ok {
			return OwnershipNotOwner, nil
		}

		return OwnershipPublic, nil
	}

	if obj.Author == userID {
		return OwnershipOwner, nil
	}

	return OwnershipNotOwner, nil
}
