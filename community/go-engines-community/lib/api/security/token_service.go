package security

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
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
	ExpirationInterval *datetime.DurationWithUnit `bson:"expiration_interval" json:"expiration_interval"`
	InactivityInterval *datetime.DurationWithUnit `bson:"inactivity_interval" json:"inactivity_interval"`
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

	now := datetime.NewCpsTime()
	var expired datetime.CpsTime
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

func (s *tokenService) CreateWithExpiration(ctx context.Context, user security.User, provider string, expiredAt time.Time) (string, error) {
	_, inactivityInterval, err := s.getIntervals(ctx, user, provider)
	if err != nil {
		return "", err
	}

	accessToken, err := s.tokenGenerator.Generate(user.ID, expiredAt)
	if err != nil {
		return "", err
	}

	expired := datetime.NewCpsTime(expiredAt.Unix())
	t := token.Token{
		ID:       accessToken,
		User:     user.ID,
		Provider: provider,
		Created:  datetime.NewCpsTime(),
		Expired:  &expired,
	}

	if inactivityInterval.Value > 0 {
		t.MaxInactiveInterval = &inactivityInterval
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

func (s *tokenService) getIntervals(ctx context.Context, user security.User, provider string) (datetime.DurationWithUnit, datetime.DurationWithUnit, error) {
	var expirationInterval, inactivityInterval datetime.DurationWithUnit
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
		default:
			if config, ok := s.config.Security.OAuth2.Providers[provider]; ok {
				expirationIntervalStr = config.ExpirationInterval
				inactivityIntervalStr = config.InactivityInterval
			} else {
				return expirationInterval, inactivityInterval, fmt.Errorf("provider %s is not defined", provider)
			}
		}

		if expirationInterval.Value == 0 && expirationIntervalStr != "" {
			expirationInterval, _ = datetime.ParseDurationWithUnit(expirationIntervalStr)
		}
		if inactivityInterval.Value == 0 && inactivityIntervalStr != "" {
			inactivityInterval, _ = datetime.ParseDurationWithUnit(inactivityIntervalStr)
		}
	}

	if inactivityInterval.Value == 0 {
		inactivityInterval.Value = security.DefaultInactivityInterval
		inactivityInterval.Unit = "h"
	}

	return expirationInterval, inactivityInterval, nil
}
