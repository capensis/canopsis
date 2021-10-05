package httpprovider

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"github.com/rs/zerolog"
	"net/http"
)

// queryTokenProvider implements a Token HTTP Authentication provider.
// It uses token from request query and retrieves user using user provider.
type queryTokenProvider struct {
	tokenService token.Service
	tokenStore   token.Store
	userProvider security.UserProvider
	logger       zerolog.Logger
}

// NewQueryTokenProvider creates new provider.
func NewQueryTokenProvider(
	tokenService token.Service,
	tokenStore token.Store,
	userProvider security.UserProvider,
	logger zerolog.Logger,
) security.HttpProvider {
	return &queryTokenProvider{
		tokenService: tokenService,
		tokenStore:   tokenStore,
		userProvider: userProvider,
		logger:       logger,
	}
}

func (p *queryTokenProvider) Auth(r *http.Request) (*security.User, error, bool) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		return nil, nil, false
	}

	ok, err := p.tokenStore.Exists(r.Context(), tokenString)
	if err != nil || !ok {
		return nil, err, true
	}

	userID, err := p.tokenService.ValidateToken(tokenString)
	if err != nil {
		p.logger.Debug().Err(err).Msg("invalid token")
		return nil, nil, true
	}

	user, err := p.userProvider.FindByID(r.Context(), userID)
	if err != nil {
		return nil, fmt.Errorf("cannot find user: %w", err), true
	}

	if user == nil || !user.IsEnabled {
		return nil, nil, true
	}

	return user, nil, true
}
