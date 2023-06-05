package httpprovider

import (
	"errors"
	"net/http"
	"net/url"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/rs/zerolog"
)

// cookieProvider implements a Cookie Token Authentication provider.
// It must be used only for file access.
type cookieProvider struct {
	tokenProviders []security.TokenProvider

	cookieName string
	logger     zerolog.Logger
}

func NewCookieProvider(
	tokenProviders []security.TokenProvider,
	cookieName string,
	logger zerolog.Logger,
) security.HttpProvider {
	return &cookieProvider{
		tokenProviders: tokenProviders,

		cookieName: cookieName,
		logger:     logger,
	}
}

func (p *cookieProvider) Auth(r *http.Request) (*security.User, error, bool) {
	cookie, err := r.Cookie(p.cookieName)
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

	for _, provider := range p.tokenProviders {
		user, err := provider.Auth(r.Context(), tokenString)
		if err != nil || user != nil {
			return user, err, true
		}
	}

	return nil, nil, true
}
