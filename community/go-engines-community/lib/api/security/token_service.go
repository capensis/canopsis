package security

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"go.mongodb.org/mongo-driver/bson"
)

type TokenStore interface {
	Save(ctx context.Context, token token.Token) error
	Count(ctx context.Context) (int64, error)
	Delete(ctx context.Context, id string) (bool, error)
	DeleteBy(ctx context.Context, user, provider string) error
}

type TokenService interface {
	Create(ctx context.Context, user security.User, provider string) (string, error)
	CreateWithExpiration(ctx context.Context, user security.User, provider string, expiredAt time.Time) (string, error)
	Delete(ctx context.Context, token string) (bool, error)
	DeleteBy(ctx context.Context, user, provider string) error
	Count(ctx context.Context) (int64, error)
}

type AuthMethodConf struct {
	ExpirationInterval *types.DurationWithUnit `bson:"expiration_interval" json:"expiration_interval"`
	InactivityInterval *types.DurationWithUnit `bson:"inactivity_interval" json:"inactivity_interval"`
}

func NewTokenService(
	client mongo.DbClient,
	generator token.Generator,
	store TokenStore,
) TokenService {
	return &tokenService{
		dbConfigCollection: client.Collection(mongo.ConfigurationMongoCollection),
		dbRoleCollection:   client.Collection(mongo.RightsMongoCollection),
		tokenGenerator:     generator,
		tokenStore:         store,
	}
}

type tokenService struct {
	dbConfigCollection mongo.DbCollection
	dbRoleCollection   mongo.DbCollection

	tokenGenerator token.Generator
	tokenStore     TokenStore
}

func (s *tokenService) Create(ctx context.Context, user security.User, provider string) (string, error) {
	expirationInterval, inactivityInterval, err := s.getIntervals(ctx, user, provider)
	if err != nil {
		return "", err
	}

	now := types.NewCpsTime()
	var expired types.CpsTime
	if expirationInterval.Value > 0 {
		expired = expirationInterval.AddTo(now)
	}
	accessToken, err := s.tokenGenerator.Generate(user.ID, expired.Time)
	if err != nil {
		return "", err
	}

	t := token.Token{
		ID:       accessToken,
		User:     user.ID,
		Provider: provider,
		Created:  now,
	}
	if inactivityInterval.Value > 0 {
		t.MaxInactiveInterval = &inactivityInterval
	}
	if !expired.IsZero() {
		t.Expired = &expired
	}
	err = s.tokenStore.Save(ctx, t)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *tokenService) CreateWithExpiration(ctx context.Context, user security.User, provider string, expired time.Time) (string, error) {
	expirationInterval, inactivityInterval, err := s.getIntervals(ctx, user, provider)
	if err != nil {
		return "", err
	}

	now := types.NewCpsTime()
	minExpired := types.CpsTime{Time: expired}
	if expirationInterval.Value > 0 {
		expiredByInterval := expirationInterval.AddTo(now)
		if expiredByInterval.Before(minExpired) {
			minExpired = expiredByInterval
		}
	}
	accessToken, err := s.tokenGenerator.Generate(user.ID, minExpired.Time)
	if err != nil {
		return "", err
	}

	t := token.Token{
		ID:       accessToken,
		User:     user.ID,
		Provider: provider,
		Created:  now,
	}
	if inactivityInterval.Value > 0 {
		t.MaxInactiveInterval = &inactivityInterval
	}
	if !minExpired.IsZero() {
		t.Expired = &minExpired
	}
	err = s.tokenStore.Save(ctx, t)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *tokenService) Delete(ctx context.Context, token string) (bool, error) {
	return s.tokenStore.Delete(ctx, token)
}

func (s *tokenService) DeleteBy(ctx context.Context, user, provider string) error {
	return s.tokenStore.DeleteBy(ctx, user, provider)
}

func (s *tokenService) Count(ctx context.Context) (int64, error) {
	return s.tokenStore.Count(ctx)
}

func (s *tokenService) getIntervals(ctx context.Context, user security.User, provider string) (types.DurationWithUnit, types.DurationWithUnit, error) {
	var expirationInterval, inactivityInterval types.DurationWithUnit
	roleConf := AuthMethodConf{}
	cursor, err := s.dbRoleCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": user.Role}},
		{"$project": bson.M{
			"expiration_interval": "$auth_config.expiration_interval",
			"inactivity_interval": "$auth_config.inactivity_interval",
		}},
	})
	if err != nil {
		return expirationInterval, inactivityInterval, err
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		err = cursor.Decode(&roleConf)
		if err != nil {
			return expirationInterval, inactivityInterval, err
		}
	}
	if roleConf.ExpirationInterval != nil {
		expirationInterval = *roleConf.ExpirationInterval
	}
	if roleConf.InactivityInterval != nil {
		inactivityInterval = *roleConf.InactivityInterval
	}

	if expirationInterval.Value == 0 || inactivityInterval.Value == 0 {
		conf := AuthMethodConf{}
		if provider == "" {
			provider = security.AuthMethodBasic
		}
		cursor, err := s.dbConfigCollection.Aggregate(ctx, []bson.M{
			{"$match": bson.M{"_id": config.ApiSecurityKeyName}},
			{"$project": bson.M{
				"expiration_interval": fmt.Sprintf("$%s.expiration_interval", provider),
				"inactivity_interval": fmt.Sprintf("$%s.inactivity_interval", provider),
			}},
		})
		if err != nil {
			return expirationInterval, inactivityInterval, err
		}
		defer cursor.Close(ctx)
		if cursor.Next(ctx) {
			err = cursor.Decode(&conf)
			if err != nil {
				return expirationInterval, inactivityInterval, err
			}
		}

		if expirationInterval.Value == 0 && conf.ExpirationInterval != nil {
			expirationInterval = *conf.ExpirationInterval
		}
		if inactivityInterval.Value == 0 && conf.InactivityInterval != nil {
			inactivityInterval = *conf.InactivityInterval
		}
	}

	if inactivityInterval.Value == 0 {
		inactivityInterval.Value = config.ApiSecurityInactivityInterval
		inactivityInterval.Unit = "h"
	}

	return expirationInterval, inactivityInterval, nil
}
