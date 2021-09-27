package httpprovider

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	"github.com/rs/zerolog"
)

const cookieToken = "token"

// cookieProvider implements a Cookie Token Authentication provider.
// It must be used only for file access.
type cookieProvider struct {
	tokenService token.Service
	tokenStore   token.Store
	userProvider security.UserProvider
	logger       zerolog.Logger
}

func NewCookieProvider(
	tokenService token.Service,
	tokenStore token.Store,
	userProvider security.UserProvider,
	logger zerolog.Logger,
) security.HttpProvider {
	return &cookieProvider{
		tokenService: tokenService,
		tokenStore:   tokenStore,
		userProvider: userProvider,
		logger:       logger,
	}
}

func (p *cookieProvider) Auth(r *http.Request) (*security.User, error, bool) {
	cookie, err := r.Cookie(cookieToken)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil, false
		}
		return nil, err, false
	}

	tokenString, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, err, false
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
