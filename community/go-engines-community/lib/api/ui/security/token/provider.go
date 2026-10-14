package token

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/rs/zerolog"
)

type Validator interface {
	Validate(token string) (userID string, err error)
}

type UsageTracker interface {
	Exists(ctx context.Context, id string) (bool, error)
	Access(ctx context.Context, id string) error
}

func NewProvider(
	validator Validator,
	tracker UsageTracker,
	userProvider security.UserProvider,
	logger zerolog.Logger,
) security.TokenProvider {
	return &provider{
		validator:    validator,
		tracker:      tracker,
		userProvider: userProvider,
		logger:       logger,
	}
}

type provider struct {
	validator    Validator
	tracker      UsageTracker
	userProvider security.UserProvider
	logger       zerolog.Logger
}

func (p *provider) Auth(ctx context.Context, token string) (*security.User, error) {
	ok, err := p.tracker.Exists(ctx, token)
	if err != nil || !ok {
		return nil, err
	}

	userID, err := p.validator.Validate(token)
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

	err = p.tracker.Access(ctx, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}
