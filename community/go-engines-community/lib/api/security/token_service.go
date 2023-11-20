package security

import (
	"context"
	"time"

	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"go.mongodb.org/mongo-driver/bson"
)

type TokenStore interface {
	Save(ctx context.Context, token token.Token) error
	Delete(ctx context.Context, id string) (bool, error)
	DeleteBy(ctx context.Context, user, provider string) error
	DeleteByUserIDs(ctx context.Context, ids []string) error
}

type TokenService interface {
	Create(ctx context.Context, user security.User, provider string) (string, error)
	CreateWithExpiration(ctx context.Context, user security.User, provider string, expiredAt time.Time) (string, error)
	Delete(ctx context.Context, token string) (bool, error)
	DeleteBy(ctx context.Context, user, provider string) error
	DeleteByUserIDs(ctx context.Context, ids []string) error
}

type AuthMethodConf struct {
	ExpirationInterval *libtime.DurationWithUnit `bson:"expiration_interval" json:"expiration_interval"`
	InactivityInterval *libtime.DurationWithUnit `bson:"inactivity_interval" json:"inactivity_interval"`
}

func NewTokenService(
	config security.Config,
	client mongo.DbClient,
	generator token.Generator,
	store TokenStore,
) TokenService {
	return &tokenService{
		config:           config,
		dbRoleCollection: client.Collection(mongo.RoleCollection),
		tokenGenerator:   generator,
		tokenStore:       store,
	}
}

type tokenService struct {
	config security.Config

	dbRoleCollection mongo.DbCollection

	tokenGenerator token.Generator
	tokenStore     TokenStore
}

func (s *tokenService) Create(ctx context.Context, user security.User, provider string) (string, error) {
	expirationInterval, inactivityInterval, err := s.getIntervals(ctx, user, provider)
	if err != nil {
		return "", err
	}

	now := libtime.NewCpsTime()
	var expired libtime.CpsTime
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

	now := libtime.NewCpsTime()
	minExpired := libtime.CpsTime{Time: expired}
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

func (s *tokenService) DeleteByUserIDs(ctx context.Context, ids []string) error {
	return s.tokenStore.DeleteByUserIDs(ctx, ids)
}

func (s *tokenService) getIntervals(ctx context.Context, user security.User, provider string) (libtime.DurationWithUnit, libtime.DurationWithUnit, error) {
	var expirationInterval, inactivityInterval libtime.DurationWithUnit
	roleConf := AuthMethodConf{}
	if len(user.Roles) > 0 {
		cursor, err := s.dbRoleCollection.Aggregate(ctx, []bson.M{
			{"$match": bson.M{"_id": user.Roles[0]}},
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
	}

	if roleConf.ExpirationInterval != nil {
		expirationInterval = *roleConf.ExpirationInterval
	}
	if roleConf.InactivityInterval != nil {
		inactivityInterval = *roleConf.InactivityInterval
	}

	if expirationInterval.Value == 0 || inactivityInterval.Value == 0 {
		if provider == "" {
			provider = security.AuthMethodBasic
		}
		var expirationIntervalStr, inactivityIntervalStr string
		switch provider {
		case security.AuthMethodBasic:
			expirationIntervalStr = s.config.Security.Basic.ExpirationInterval
			inactivityIntervalStr = s.config.Security.Basic.InactivityInterval
		case security.AuthMethodLdap:
			expirationIntervalStr = s.config.Security.Ldap.ExpirationInterval
			inactivityIntervalStr = s.config.Security.Ldap.InactivityInterval
		case security.AuthMethodCas:
			expirationIntervalStr = s.config.Security.Cas.ExpirationInterval
			inactivityIntervalStr = s.config.Security.Cas.InactivityInterval
		case security.AuthMethodSaml:
			expirationIntervalStr = s.config.Security.Saml.ExpirationInterval
			inactivityIntervalStr = s.config.Security.Saml.InactivityInterval

		}

		if expirationInterval.Value == 0 && expirationIntervalStr != "" {
			expirationInterval, _ = libtime.ParseDurationWithUnit(expirationIntervalStr)
		}
		if inactivityInterval.Value == 0 && inactivityIntervalStr != "" {
			inactivityInterval, _ = libtime.ParseDurationWithUnit(inactivityIntervalStr)
		}
	}

	if inactivityInterval.Value == 0 {
		inactivityInterval.Value = security.DefaultInactivityInterval
		inactivityInterval.Unit = "h"
	}

	return expirationInterval, inactivityInterval, nil
}
