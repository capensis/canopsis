package tokenprovider

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"github.com/rs/zerolog"
)

type TokenStore interface {
	Exists(ctx context.Context, id string) (bool, error)
	Access(ctx context.Context, id string) error
}

func NewTokenProvider(
	tokenGenerator token.Generator,
	tokenStore TokenStore,
	userProvider security.UserProvider,
	logger zerolog.Logger,
) security.TokenProvider {
	return &tokenProvider{
		tokenGenerator: tokenGenerator,
		tokenStore:     tokenStore,
		userProvider:   userProvider,
		logger:         logger,
	}
}

type tokenProvider struct {
	tokenGenerator token.Generator
	tokenStore     TokenStore
	userProvider   security.UserProvider
	logger         zerolog.Logger
}

func (p *tokenProvider) Auth(ctx context.Context, token string) (*security.User, error) {
	ok, err := p.tokenStore.Exists(ctx, token)
	if err != nil || !ok {
		return nil, err
	}

	userID, err := p.tokenGenerator.Validate(token)
	if err != nil {
		p.logger.Debug().Err(err).Msg("invalid token")
		return nil, nil
	}

	user, err := p.userProvider.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %w", err)
	}

	if user == nil || !user.IsEnabled {
		return nil, nil
	}

	err = p.tokenStore.Access(ctx, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}
