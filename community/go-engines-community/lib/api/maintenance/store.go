package maintenance

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/session"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrEnabled = errors.New("maintenance mode has already been enabled")
var ErrDisabled = errors.New("maintenance mode has already been disabled")

const defaultColor = "#e75e40"

type Store interface {
	Enable(ctx context.Context, message, color, userID string) error
	Disable(ctx context.Context, userID string) error
}

type store struct {
	dbClient            mongo.DbClient
	configCollection    mongo.DbCollection
	broadcastCollection mongo.DbCollection

	userProvider security.UserProvider
	tokenService apisecurity.TokenService
	sessionStore session.Store
}

func NewStore(
	dbClient mongo.DbClient,
	userProvider security.UserProvider,
	tokenService apisecurity.TokenService,
	sessionStore session.Store,
) Store {
	return &store{
		dbClient:            dbClient,
		configCollection:    dbClient.Collection(mongo.ConfigurationMongoCollection),
		broadcastCollection: dbClient.Collection(mongo.BroadcastMessageMongoCollection),

		userProvider: userProvider,
		tokenService: tokenService,
		sessionStore: sessionStore,
	}
}

func (s *store) Enable(ctx context.Context, message, color, userID string) error {
	broadcastID := utils.NewID()

	_, err := s.configCollection.UpdateOne(
		ctx,
		bson.M{"_id": config.MaintenanceKeyName, "enabled": false},
		bson.M{"$set": bson.M{
			"enabled":      true,
			"broadcast_id": broadcastID,
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		if mongodriver.IsDuplicateKeyError(err) {
			return ErrEnabled
		}

		return err
	}

	if color == "" {
		color = defaultColor
	}

	return s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		now := datetime.NewCpsTime()

		users, err := s.userProvider.FindWithoutPermission(ctx, apisecurity.PermMaintenance)
		if err != nil {
			panic(err)
		}

		userIDs := make([]string, len(users))
		for idx := range users {
			userIDs[idx] = users[idx].ID
		}

		err = s.tokenService.DeleteByUserIDs(ctx, userIDs)
		if err != nil {
			return err
		}

		err = s.sessionStore.ExpireSessionsByUserIDs(ctx, userIDs)
		if err != nil {
			return err
		}

		_, err = s.broadcastCollection.InsertOne(ctx, broadcastmessage.CreateRequest{
			ID: broadcastID,
			EditRequest: broadcastmessage.EditRequest{
				Color:   color,
				Message: message,
				Start:   now,
				End:     datetime.NewCpsTime(now.AddDate(1, 0, 0).Unix()),
				Author:  userID,
				Created: &now,
				Updated: &now,
			},
		})

		return err
	})
}

func (s *store) Disable(ctx context.Context, userID string) error {
	return s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		var state config.MaintenanceConf

		err := s.configCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": config.MaintenanceKeyName, "enabled": true},
			bson.M{"$set": bson.M{"enabled": false}},
		).Decode(&state)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return ErrDisabled
			}

			return err
		}

		// required to get the author in action log listener.
		res, err := s.broadcastCollection.UpdateOne(ctx, bson.M{"_id": state.BroadcastID}, bson.M{"$set": bson.M{"author": userID}})
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		_, err = s.broadcastCollection.DeleteOne(ctx, bson.M{"_id": state.BroadcastID})
		return err
	})
}
