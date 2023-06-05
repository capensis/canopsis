// userprovider contains user storages.
package userprovider

import (
	"context"

	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoProvider decorates request to mongo db.
type mongoProvider struct {
	collection libmongo.DbCollection
}

// NewMongoProvider creates new provider.
func NewMongoProvider(db libmongo.DbClient) security.UserProvider {
	return &mongoProvider{
		collection: db.Collection(libmongo.RightsMongoCollection),
	}
}

func (p *mongoProvider) FindByUsername(ctx context.Context, username string) (*security.User, error) {
	return p.findByFilter(ctx, bson.M{
		"crecord_type": model.LineTypeSubject,
		"_id":          username,
		"source":       bson.M{"$in": bson.A{"", nil}},
	})
}

func (p *mongoProvider) FindByAuthApiKey(ctx context.Context, apiKey string) (*security.User, error) {
	return p.findByFilter(ctx, bson.M{
		"crecord_type": model.LineTypeSubject,
		"authkey":      apiKey,
	})
}

func (p *mongoProvider) FindByID(ctx context.Context, id string) (*security.User, error) {
	var objID interface{}
	var err error
	objID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		objID = id
	}

	return p.findByFilter(ctx, bson.M{
		"crecord_type": model.LineTypeSubject,
		"_id":          objID,
	})
}

func (p *mongoProvider) FindByExternalSource(
	ctx context.Context,
	externalID string,
	source security.Source,
) (*security.User, error) {
	return p.findByFilter(ctx, bson.M{
		"crecord_type": model.LineTypeSubject,
		"external_id":  externalID,
		"source":       source,
	})
}

func (p *mongoProvider) Save(ctx context.Context, u *security.User) error {
	if u.ID == "" {
		u.ID = utils.NewID()
		u.AuthApiKey = utils.NewID()
	}

	m := transformUserToDbModel(u)
	_, err := p.collection.UpdateOne(
		ctx,
		bson.M{"_id": m.ID},
		bson.M{"$set": *m},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		return err
	}

	return nil
}

// findByFilter returns User or nil if no user matches filter.
func (p *mongoProvider) findByFilter(ctx context.Context, f interface{}) (*security.User, error) {
	cursor, err := p.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		var line model.Rbac
		err := cursor.Decode(&line)

		if err != nil {
			return nil, err
		}

		return transformDbModelToUser(&line), nil
	}

	return nil, nil
}

// transformUserToDbModel transforms User model to mongo document.
func transformUserToDbModel(u *security.User) *model.Rbac {
	var m model.Rbac
	m.Type = model.LineTypeSubject
	m.ID = u.ID
	m.Name = u.Name
	m.Email = u.Email
	m.Firstname = u.Firstname
	m.Lastname = u.Lastname
	m.Role = u.Role
	m.HashedPassword = u.HashedPassword
	m.AuthApiKey = u.AuthApiKey
	m.IsEnabled = u.IsEnabled
	m.ExternalID = u.ExternalID
	m.Source = string(u.Source)
	m.Contact.Name = u.Contact.Name
	m.Contact.Address = u.Contact.Address

	return &m
}

// transformDbModelToUser transforms mongo document to User model.
func transformDbModelToUser(m *model.Rbac) *security.User {
	var u security.User
	u.ID = m.ID
	u.Name = m.Name
	u.Email = m.Email
	u.Firstname = m.Firstname
	u.Lastname = m.Lastname
	u.Role = m.Role
	u.HashedPassword = m.HashedPassword
	u.AuthApiKey = m.AuthApiKey
	u.IsEnabled = m.IsEnabled
	u.ExternalID = m.ExternalID
	u.Source = security.Source(m.Source)
	u.Contact.Name = m.Contact.Name
	u.Contact.Address = m.Contact.Address

	return &u
}
