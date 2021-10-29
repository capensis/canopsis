package httpprovider

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
)

const bearerPrefix = "Bearer"

// bearerProvider implements a Bearer Token Authentication provider.
type bearerProvider struct {
	tokenService token.Service
	tokenStore   token.Store
	userProvider security.UserProvider
	logger       zerolog.Logger
}

// NewBearerProvider creates new provider.
func NewBearerProvider(
	tokenService token.Service,
	tokenStore token.Store,
	userProvider security.UserProvider,
	logger zerolog.Logger,
) security.HttpProvider {
	return &bearerProvider{
		tokenService: tokenService,
		tokenStore:   tokenStore,
		userProvider: userProvider,
		logger:       logger,
	}
}

func (p *bearerProvider) Auth(r *http.Request) (*security.User, error, bool) {
	header := r.Header.Get(headerAuthorization)
	if len(header) < len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return nil, nil, false
	}

	tokenString := strings.TrimSpace(header[len(bearerPrefix):])

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
