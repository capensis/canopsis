package token

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewProvider(dbClient mongo.DbClient, userProvider security.UserProvider) security.TokenProvider {
	return &provider{
		dbCollection: dbClient.Collection(mongo.PublicAPITokenCollection),
		userProvider: userProvider,
	}
}

type provider struct {
	dbCollection mongo.DbCollection
	userProvider security.UserProvider
}

// Auth returns the user which owns the given token.
// It returns a nil user and no error when the token is unknown or expired, or when its owner doesn't
// exist anymore or is disabled.
func (p *provider) Auth(ctx context.Context, token string) (*security.User, error) {
	hash := Hash(token)
	userID, err := p.findTokenUserID(ctx, hash)
	if err != nil {
		return nil, err
	}

	if userID == "" {
		return nil, nil
	}

	user, err := p.userProvider.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %w", err)
	}

	if user == nil || !user.IsEnabled {
		return nil, nil
	}

	err = p.markAccessed(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("cannot update token: %w", err)
	}

	return user, nil
}

func (p *provider) findTokenUserID(ctx context.Context, hash string) (string, error) {
	now := datetime.NewCpsTime()
	doc := struct {
		User string `bson:"user"`
	}{}
	opts := options.FindOne().SetProjection(bson.M{"user": 1})
	err := p.dbCollection.
		FindOne(
			ctx,
			bson.M{
				"hash": hash,
				// A token without an expiration date never expires.
				"$or": []bson.M{
					{"expired": nil},
					{"expired": bson.M{"$gt": now}},
				},
			},
			opts,
		).
		Decode(&doc)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return "", nil
		}

		return "", fmt.Errorf("cannot find token: %w", err)
	}

	return doc.User, nil
}

// markAccessed saves the last time the token was used to authenticate.
func (p *provider) markAccessed(ctx context.Context, hash string) error {
	_, err := p.dbCollection.UpdateOne(ctx, bson.M{"hash": hash}, bson.M{"$set": bson.M{
		"accessed": datetime.NewCpsTime(),
	}})

	return err
}
